//go:generate protoc --proto_path=. --descriptor_set_out=./alert/alert.protoset --include_imports ./alert/alert.proto
//go:generate protoc --proto_path=. -I/usr/include -I. -Ivendor --go_out=plugins=grpc:. ./alert/alert.proto
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cyrinux/waybar-livestatus/client"
	"github.com/cyrinux/waybar-livestatus/helpers"
	"github.com/cyrinux/waybar-livestatus/lql"
	"github.com/cyrinux/waybar-livestatus/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func formatData(
	hAlerts *lql.AlertStruct,
	sAlerts *lql.AlertStruct,
	config *helpers.CONFIG) (wOutput helpers.WaybarOutput) {

	// format
	globalClass := "ok"

	var text, tooltip, icon string

	// test and format
	var hostAlertsCount = int(hAlerts.Count)
	if hostAlertsCount > 0 {
		icon = fmt.Sprintf(" %s ", config.HostPrefix)
		text += fmt.Sprintf("%s%d", icon, hAlerts.Count)
		tooltip += fmt.Sprintf("<b>Hosts: %d</b>\n\n%s", hAlerts.Count, hAlerts.Items)
		globalClass = hAlerts.Class
	}

	var serviceAlertsCount = int(sAlerts.Count)
	if serviceAlertsCount > 0 {
		if hostAlertsCount > 0 {
			tooltip += "\n\n"
		}

		tooltip += fmt.Sprintf("<b>Services: %d</b>\n\n%s", sAlerts.Count, sAlerts.Items)
		if len(strings.TrimSpace(text)) > 0 {
			text += " | "
		}
		icon = fmt.Sprintf(" %s ", strings.TrimSpace(config.ServicePrefix))
		text += fmt.Sprintf("%s%d", icon, sAlerts.Count)
		globalClass += sAlerts.Class
	}

	if len(strings.TrimSpace(text)) == 0 {
		icon = fmt.Sprintf(" %s ", strings.TrimSpace(config.OkPrefix))
		text = icon
	}

	// on error
	if hAlerts.Error != nil || sAlerts.Error != nil {
		icon = fmt.Sprintf(" %s ", strings.TrimSpace(config.ErrorPrefix))
		text = icon
		tooltip = "Can't connect to the livestatus server"
		globalClass = "error"
	}

	tooltip = strings.TrimRight(tooltip, "\n")

	if len(strings.TrimSpace(text)) == 0 {
		text = strings.TrimSpace(text)
	} else {
		text = strings.TrimRight(text, "\n")
	}

	// waybar output
	wOutput = helpers.WaybarOutput{
		Text:    text,
		Tooltip: tooltip,
		Class:   globalClass,
		Count:   hostAlertsCount + serviceAlertsCount,
	}

	return
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// get config
	var config = helpers.GetConfig()

	// start the client and exit
	if config.Client {
		err := client.Start(config)
		if err != nil {
			log.Error().Msgf("%v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// start cron based auto-pause system
	helpers.AutoPause(config)

	// create channels and start goroutines
	hostAlerts, serviceAlerts := make(chan lql.AlertStruct), make(chan lql.AlertStruct)
	notificationsChannel, serverChannel := make(chan *helpers.Alert, 10), make(chan []*helpers.Alert)
	hAlerts, sAlerts := new(lql.AlertStruct), new(lql.AlertStruct)

	if config.HostsOnly && !config.ServicesOnly || (!config.ServicesOnly && !config.HostsOnly) {
		go lql.GetItems("hosts", config, hostAlerts, notificationsChannel, serverChannel)
	}
	if (config.ServicesOnly && !config.HostsOnly) || (!config.ServicesOnly && !config.HostsOnly) {
		go lql.GetItems("services", config, serviceAlerts, notificationsChannel, serverChannel)
	}

	// Handle Signal
	go helpers.SignalHandler()

	// notification channel
	if config.Popup {
		go helpers.SendNotification(notificationsChannel, config)
	}

	// keep version of previous output
	var wOutput, previousWOutput helpers.WaybarOutput

	// Start gRPC server
	go server.GRPCListen(serverChannel, config)

	log.Debug().Msgf("Refresh rate: %d seconds, long refresh: %d seconds", config.Refresh, config.LongRefresh)

	// main loop
	for {
		select {
		case *sAlerts = <-serviceAlerts:
			log.Debug().Msgf("received %d service alerts", sAlerts.Count)
		case *hAlerts = <-hostAlerts:
			log.Debug().Msgf("received %d hosts alerts", hAlerts.Count)
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
			log.Error().Msgf("%v", err)
		}

		// Finally print the expected waybar JSON
		if wOutput != previousWOutput {
			fmt.Println(string(jsonOutput))
		}
	}
}
