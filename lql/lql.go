package lql

import (
	"fmt"
	"github.com/cyrinux/waybar-livestatus/helpers"
	"github.com/hyperjumptech/jiffy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	livestatus "github.com/vbatoufflet/go-livestatus"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	statusOK       = 0
	statusWarning  = 1
	statusCritical = 2
	statusUnknown  = 3
)

// AlertStruct define an alert
type AlertStruct struct {
	Count int
	Items string
	Class string
	Error error
}

func sanitizeObjectType(objectType string) {
	if !(objectType == "hosts" || objectType == "services") {
		log.Fatal().Msg("Bad objectType, valid are 'services' or 'hosts'")
	}
}

func getQuery(objectType string, config *helpers.CONFIG) (q *livestatus.Query) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	sanitizeObjectType(objectType)

	q = livestatus.NewQuery(objectType)
	columns := []string{"host_name", "state", "description", "is_flapping", "notes_url"}
	if config.GetDuration {
		columns = append(columns, "last_hard_state_change")
	}
	q.Columns(columns...)

	q.Filter("is_problem = 1")
	q.Filter(fmt.Sprintf("state = %v", statusCritical))
	if config.Warnings {
		q.Filter(fmt.Sprintf("state = %v", statusWarning))
		q.Or(2)
	}
	q.Filter(fmt.Sprintf("acknowledged = %v", int(config.Acknowledged)))
	q.Filter(fmt.Sprintf("notifications_enabled = %v", int(config.NotificationsEnabled)))
	q.Filter(fmt.Sprintf("in_notification_period = %v", int(config.InNotificationPeriod)))
	q.Filter(fmt.Sprintf("scheduled_downtime_depth = %v", int(config.ScheduledDowntimeDepth)))
	if objectType == "services" {
		q.Filter(fmt.Sprintf("service_scheduled_downtime_depth = %v", int(config.ServiceScheduledDowntimeDepth)))
	} else if objectType == "hosts" {
		q.Filter(fmt.Sprintf("host_scheduled_downtime_depth = %v", int(config.HostScheduledDowntimeDepth)))
	}

	countHosts := len(config.HostsPattern)
	if countHosts >= 1 {
		for _, hostName := range config.HostsPattern {
			log.Debug().Msgf("filter host_name ~ '%s'", hostName)
			hostFilter := fmt.Sprintf("host_name ~ %s", hostName)
			q.Filter(hostFilter)
		}
		q.Or(countHosts)
	}
	q.Limit(config.Limit)

	log.Debug().Msgf("%v", q)

	return
}

// getResponse parse the response from the livestatus client
// objectType: hosts or services
func getResponse(objectType string, c *livestatus.Client, q *livestatus.Query, config *helpers.CONFIG) (resp *livestatus.Response, localRefresh int, err error) {

	sanitizeObjectType(objectType)

	startTime := time.Now()
	log.Debug().Msgf("start of LQL %s query", objectType)
	resp, err = c.Exec(q)
	if err != nil {
		log.Error().Err(err)
		return
	}
	if resp.Status != 200 {
		log.Error().Msgf("bad LQL query status response code %v", resp.Status)
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	// if the query is two slow, use backoff
	if int(duration.Seconds()) >= config.Refresh {
		localRefresh = config.LongRefresh
		log.Debug().Msgf("end of LQL %s query, took %.3f seconds, too slow, next refresh in %d seconds",
			objectType, duration.Seconds(), localRefresh)
	} else {
		localRefresh = config.Refresh
		log.Debug().Msgf("end of LQL %s query, took %.3f seconds",
			objectType, duration.Seconds())
	}

	return
}

// GetItems open livestatus connection then query it
func GetItems(objectType string, config *helpers.CONFIG, alertChannel chan AlertStruct, notificationsChannel chan *helpers.Alert, serverChannel chan []*helpers.Alert) {
	sanitizeObjectType(objectType)

	// services alerts
	client := livestatus.NewClient("tcp", config.Server)
	defer client.Close()

	// LQL query
	query := getQuery(objectType, config)

	for {
		var class string
		now := time.Now()
		var lastHardStateChangeDuration time.Duration

		if helpers.Pause {
			time.Sleep(5 * time.Second)
			continue
		}

		resp, localRefresh, err := getResponse(objectType, client, query, config)

		var alert AlertStruct
		var items string
		var itemsMenu []*helpers.Alert

		if err != nil {
			class = "error"
			count := 0
			alert = AlertStruct{
				Count: count,
				Items: items,
				Class: class,
				Error: err,
			}
			time.Sleep(10 * time.Second)
		} else {
			class = "ok"

			// get count
			count := len(resp.Records)
			for _, r := range resp.Records {
				var item string

				if config.GetDuration {
					lastHardStateChange, err := r.GetTime("last_hard_state_change")
					if err != nil {
						log.Warn().Err(err)
					}
					lastHardStateChangeDuration = now.Sub(lastHardStateChange)
				}

				host, err := r.GetString("host_name")
				if err != nil {
					log.Warn().Err(err)
				}

				state, err := r.GetInt("state")
				if err != nil {
					log.Warn().Err(err)
				}

				desc, err := r.GetString("description")
				if err != nil {
					log.Warn().Err(err)
				}

				// keep the worse state
				if state == statusCritical {
					class = "critical"
				} else if state == statusWarning && class != "critical" {
					class = "warning"
				}

				// flapping handle
				isFlapping, err := r.GetInt("is_flapping")
				if err != nil {
					log.Warn().Err(err)
				}
				isFlappingStr := strconv.FormatInt(isFlapping, 10)
				if isFlappingStr == "1" {
					isFlappingStr = config.FlappingIcon
				} else {
					isFlappingStr = ""
				}

				notesURL, err := r.GetString("notes_url")
				if err != nil {
					log.Warn().Err(err)
				} else {
					if len(notesURL) > 0 {
						log.Info().Msgf("%s: %s: %s", host, desc, notesURL)
					}
				}

				if objectType == "services" {
					item = fmt.Sprintf("* %s: %s %s %s\n", host, desc, class, isFlappingStr)
				} else if objectType == "hosts" {
					item = fmt.Sprintf("* %s: %s %s\n", host, class, isFlappingStr)
				}

				if config.GetDuration {
					item += fmt.Sprintf("since %s\n", jiffy.DescribeDuration(lastHardStateChangeDuration, jiffy.NewWant()))
				}

				// notifications
				notificationsChannel <- &helpers.Alert{Host: host, Desc: desc, Class: class, NotesURL: notesURL}

				items += item

				itemsMenu = append(itemsMenu, &helpers.Alert{Host: host, Desc: desc, Class: class, NotesURL: notesURL})

			}

			// trim services
			items = strings.TrimRight(items, "\n")

			// make the alert
			alert = AlertStruct{Count: count, Items: items, Class: class}

			serverChannel <- itemsMenu
		}

		// feed the alerts channel
		log.Debug().Msgf("sending %d %s alerts", alert.Count, objectType)

		alertChannel <- alert

		// main sleep
		time.Sleep(time.Duration(localRefresh) * time.Second)
	}
}
