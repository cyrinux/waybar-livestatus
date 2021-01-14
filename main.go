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
}

func formatData(hAlerts *lql.AlertStruct, sAlerts *lql.AlertStruct, popup bool, debug bool) (wOutput waybarOutput) {

	// format
	globalClass := "ok"

	var text, tooltip, icon string

	// test and format
	if int(hAlerts.Count) > 0 {
		icon = ""
		tooltip += fmt.Sprintf("<b>Hosts:</b>\n\n%s", hAlerts.Items)
		text += fmt.Sprintf("%s %d", icon, hAlerts.Count)
		globalClass = hAlerts.Class
	}

	if int(sAlerts.Count) > 0 {
		if int(hAlerts.Count) > 0 {
			tooltip += "\n\n"
		}

		icon = ""

		tooltip += fmt.Sprintf("<b>Services:</b>\n\n%s", sAlerts.Items)

		if len(text) > 0 {
			text += "|"
		}
		text += fmt.Sprintf("%s %d", icon, sAlerts.Count)
		globalClass += sAlerts.Class
	}

	if len(text) == 0 {
		text = ""
	}

	// Trim Right
	tooltip = strings.TrimRight(tooltip, "\n")
	text = strings.TrimRight(text, "\n")

	log.Debugf("class: %s", globalClass)

	// waybar output
	wOutput = waybarOutput{Text: text, Tooltip: tooltip, Class: globalClass}

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

	log.Debugf("Refresh rate: %d seconds", config.Refresh)

	// create channels and start goroutines
	hostAlerts := make(chan lql.AlertStruct, 100)
	go lql.GetItems("hosts", config.Server, config.Warnings, config.Debug, hostAlerts, config.Refresh, config.HostsPatternString)
	hAlerts := new(lql.AlertStruct)

	serviceAlerts := make(chan lql.AlertStruct, 100)
	go lql.GetItems("services", config.Server, config.Warnings, config.Debug, serviceAlerts, config.Refresh, config.HostsPatternString)
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
		}

		// save previous waybar output
		previousWOutput = wOutput

		// new waybar content
		wOutput = formatData(hAlerts, sAlerts, config.Popup, config.Debug)

		if helpers.Pause {
			wOutput.Class = "pause"
			wOutput.Text = "  "
			wOutput.Tooltip = "Paused, click to resume "
		}

		// convert in JSON
		jsonOutput, err := json.Marshal(wOutput)
		if err != nil {
			log.Error(err)
		}

		// Finally print the expected waybar JSON
		fmt.Println(string(jsonOutput))

		// popup notification
		if config.Popup && !helpers.Pause && wOutput != previousWOutput && wOutput.Class != "ok" {
			if popup, err := helpers.SendNotification("Hero, check the alerts", wOutput.Text, ""); err != nil {
				log.Errorf("Error sending notification: %v", popup)
			}
		}

		if helpers.Pause {
			// backoff
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}
}
