package main

import (
	"encoding/json"
	"fmt"
	"github.com/cyrinux/waybar-livestatus/helpers"
	"github.com/cyrinux/waybar-livestatus/lql"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var version string

type waybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
	Count   int    `json:"count"`
}

func formatData(hAlerts *lql.AlertStruct, sAlerts *lql.AlertStruct, config *helpers.CONFIG) (wOutput waybarOutput) {

	// format
	globalClass := "ok"

	var text, tooltip, icon string

	// test and format
	var hostAlertsCount = int(hAlerts.Count)
	if hostAlertsCount > 0 {
		icon = " " + config.HostPrefix + " "
		text += fmt.Sprintf("%s %d", icon, hAlerts.Count)
		tooltip += fmt.Sprintf("<b>Hosts: %d</b>\n\n%s", hAlerts.Count, hAlerts.Items)
		globalClass = hAlerts.Class
	}

	var serviceAlertsCount = int(sAlerts.Count)
	if serviceAlertsCount > 0 {
		if hostAlertsCount > 0 {
			tooltip += "\n\n"
		}

		tooltip += fmt.Sprintf("<b>Services: %d</b>\n\n%s", sAlerts.Count, sAlerts.Items)
		if len(text) > 0 {
			text += " | "
		}
		icon = " " + config.ServicePrefix + " "
		text += fmt.Sprintf("%s %d", icon, sAlerts.Count)
		globalClass += sAlerts.Class
	}

	if len(text) == 0 {
		icon = " " + config.OkPrefix + " "
		text = icon
	}

	// on error
	if hAlerts.Error != nil || sAlerts.Error != nil {
		icon = " " + config.ErrorPrefix + " "
		text = icon
		tooltip = "Can't connect to the livestatus server"
		globalClass = "error"
	}

	tooltip = strings.TrimRight(tooltip, "\n")
	text = strings.TrimRight(text, "\n")

	// waybar output
	wOutput = waybarOutput{Text: text, Tooltip: tooltip, Class: globalClass, Count: hostAlertsCount + serviceAlertsCount}

	return
}

func main() {

	// get config
	var config = helpers.GetConfig()

	if config.Version {
		fmt.Println("Waybar Livestatus version:", version)
		os.Exit(0)
	}

	if config.Refresh < 15 {
		log.Info("Refresh rate can't be under 15 seconds ! Fallback to default: 60 seconds")
		config.Refresh = 60
	}

	log.Debugf("Refresh rate: %d seconds, long refresh: %d seconds", config.Refresh, config.LongRefresh)

	// create channels and start goroutines
	hostAlerts := make(chan lql.AlertStruct)
	go lql.GetItems("hosts", config, hostAlerts)
	hAlerts := new(lql.AlertStruct)

	serviceAlerts := make(chan lql.AlertStruct)
	go lql.GetItems("services", config, serviceAlerts)
	sAlerts := new(lql.AlertStruct)

	// toggle pause on SIGUSR1
	go helpers.PauseHandler()

	// keep version of previous output
	var wOutput, previousWOutput waybarOutput

	// main loop
	for {
		select {
		case *sAlerts = <-serviceAlerts:
			log.Debugf("received %d service alerts", sAlerts.Count)
		case *hAlerts = <-hostAlerts:
			log.Debugf("received %d hosts alerts", hAlerts.Count)
		default:
			// nothing receive, sleep a little
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// save previous waybar output
		previousWOutput = wOutput

		// new waybar content
		wOutput = formatData(hAlerts, sAlerts, config)

		if helpers.Pause {
			wOutput.Class = "pause"
			wOutput.Text = fmt.Sprintf(" %s ", config.PausePrefix)
			wOutput.Tooltip = "Paused, click to resume "
			wOutput.Count = 0
		}

		// convert in JSON
		jsonOutput, err := json.Marshal(wOutput)
		if err != nil {
			log.Error(err)
		}

		// Finally print the expected waybar JSON
		if wOutput != previousWOutput {
			fmt.Println(string(jsonOutput))
		}

		// popup notification
		if config.Popup && !helpers.Pause && wOutput != previousWOutput && wOutput.Class != "ok" && wOutput.Class != "error" && wOutput.Count > previousWOutput.Count {
			if popup, err := helpers.SendNotification(fmt.Sprintf("Hero, check the %d alerts", wOutput.Count), wOutput.Text, ""); err != nil {
				log.Errorf("Error sending notification: %v", popup)
			}
		}
	}
}
