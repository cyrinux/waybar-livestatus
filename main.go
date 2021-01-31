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

func formatData(hAlerts *lql.AlertStruct, sAlerts *lql.AlertStruct, config *helpers.CONFIG) (wOutput helpers.WaybarOutput) {

	// format
	globalClass := "ok"

	var text, tooltip, icon string

	// test and format
	var hostAlertsCount = int(hAlerts.Count)
	if hostAlertsCount > 0 {
		icon = config.HostPrefix + " "
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
		icon = config.ServicePrefix + " "
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
	wOutput = helpers.WaybarOutput{Text: text, Tooltip: tooltip, Class: globalClass, Count: hostAlertsCount + serviceAlertsCount}

	return
}

func main() {

	// get config
	var config = helpers.GetConfig()

	helpers.AutoPause(config)

	if config.Version {
		fmt.Println("Waybar Livestatus version:", version)
		os.Exit(0)
	}

	log.Debugf("Refresh rate: %d seconds, long refresh: %d seconds", config.Refresh, config.LongRefresh)

	// create channels and start goroutines
	hostAlerts, serviceAlerts, notificationsChannel := make(chan lql.AlertStruct), make(chan lql.AlertStruct), make(chan *helpers.Alert, 10)
	hAlerts, sAlerts := new(lql.AlertStruct), new(lql.AlertStruct)

	if config.HostsOnly && !config.ServicesOnly || (!config.ServicesOnly && !config.HostsOnly) {
		go lql.GetItems("hosts", config, hostAlerts, notificationsChannel)
	}
	if (config.ServicesOnly && !config.HostsOnly) || (!config.ServicesOnly && !config.HostsOnly) {
		go lql.GetItems("services", config, serviceAlerts, notificationsChannel)
	}

	// toggle pause on SIGUSR1
	go helpers.PauseHandler()

	// notification channel
	if config.Popup {
		go helpers.SendNotification(notificationsChannel, config)
	}

	// keep version of previous output
	var wOutput, previousWOutput helpers.WaybarOutput

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
	}
}
