package lql

import (
	"fmt"
	"github.com/cyrinux/waybar-livestatus/helpers"
	log "github.com/sirupsen/logrus"
	livestatus "github.com/vbatoufflet/go-livestatus"
	"strconv"
	"strings"
	"time"
)

const (
	statusCritical = 2
	statusWarning  = 1
	statusOK       = 0
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
		log.Fatal("Bad objectType, valid are 'services' or 'hosts'")
	}
}

func getQuery(objectType string, config helpers.CONFIG) (q *livestatus.Query) {

	sanitizeObjectType(objectType)

	q = livestatus.NewQuery(objectType)
	columns := []string{"host_name", "state", "description", "is_flapping"}
	if config.NotesURL {
		columns = append(columns, "notes_url")
	}
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
	if countHosts > 1 {
		for _, hostName := range config.HostsPattern {
			log.Debugf("filter host_name ~ '%s'", hostName)
			hostFilter := fmt.Sprintf("host_name ~ %s", hostName)
			q.Filter(hostFilter)
		}
		q.Or(countHosts)
	}
	q.Limit(config.Limit)

	if config.Debug {
		log.Debug(q)
	}

	return
}

// getResponse parse the response from the livestatus client
// objectType: hosts or services
func getResponse(objectType string, c *livestatus.Client, q *livestatus.Query, config *helpers.CONFIG) (resp *livestatus.Response, localRefresh int, err error) {

	sanitizeObjectType(objectType)

	startTime := time.Now()
	log.Debugf("start of LQL %s query", objectType)
	resp, err = c.Exec(q)
	if err != nil {
		log.Error(err)
		return
	}
	if resp.Status != 200 {
		log.Errorf("bad LQL query status response code %v", resp.Status)
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	// if the query is two slow, use backoff
	if int(duration.Seconds()) >= config.Refresh {
		localRefresh = config.LongRefresh
		log.Debugf("end of LQL %s query, took %v seconds, too slow, next refresh in %d seconds",
			objectType, duration.Seconds(), localRefresh)
	} else {
		localRefresh = config.Refresh
		log.Debugf("end of LQL %s query, took %v seconds",
			objectType, duration.Seconds())
	}

	return
}

// GetItems open livestatus connection then query it
func GetItems(objectType string, config *helpers.CONFIG, alertChannel chan AlertStruct, notificationsChannel chan *helpers.Alert) {
	sanitizeObjectType(objectType)

	// services alerts
	client := livestatus.NewClient("tcp", config.Server)
	defer client.Close()

	// LQL query
	query := getQuery(objectType, *config)

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

		items := ""

		if err != nil {
			class = "error"
			count := 0
			alert = AlertStruct{Count: count, Items: items, Class: class, Error: err}
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
						log.Warn(err)
					}
					lastHardStateChangeDuration = now.Sub(lastHardStateChange)
				}

				host, err := r.GetString("host_name")
				if err != nil {
					log.Warn(err)
				}

				state, err := r.GetInt("state")
				if err != nil {
					log.Warn(err)
				}

				desc, err := r.GetString("description")
				if err != nil {
					log.Warn(err)
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
					log.Warn(err)
				}

				isFlappingStr := strconv.FormatInt(isFlapping, 10)
				if isFlappingStr == "1" {
					isFlappingStr = config.FlappingIcon
				} else {
					isFlappingStr = ""
				}

				if config.NotesURL {
					notesURL, err := r.GetString("notes_url")
					if err != nil {
						log.Warn(err)
					} else {
						if len(notesURL) > 0 {
							log.Infof("%s: %s: %s", host, desc, notesURL)
						}
					}
				}

				if objectType == "services" {
					item = fmt.Sprintf("* %s: %s %s %s\n", host, desc, class, isFlappingStr)
				} else if objectType == "hosts" {
					item = fmt.Sprintf("* %s: %s %s\n", host, class, isFlappingStr)
				}

				if config.GetDuration {
					item += fmt.Sprintf("  since %s\n", lastHardStateChangeDuration.Round(1*time.Second))
				}

				// notifications
				notification := helpers.Alert{Host: host, Desc: desc, Class: class}
				notificationsChannel <- &notification

				items += item

			}

			// trim services
			items = strings.TrimRight(items, "\n")

			// make the alert
			alert = AlertStruct{Count: count, Items: items, Class: class}
		}

		// feed the alerts channel
		log.Debugf("sending %d %s alerts", alert.Count, objectType)

		alertChannel <- alert

		// main sleep
		time.Sleep(time.Duration(localRefresh) * time.Second)
	}
}
