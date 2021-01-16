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
	q.Columns("host_name", "state", "description", "notes_url", "is_flapping")
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
			log.Debugf("filter host %s", hostName)
			hostFilter := fmt.Sprintf("host_name ~ %s", hostName)
			q.Filter(hostFilter)
		}
		q.Or(countHosts)
	}

	return
}

// GetResponse parse the response from the livestatus client
// objectType: hosts or services
func GetResponse(objectType string, c *livestatus.Client, q *livestatus.Query, config *helpers.CONFIG) (resp *livestatus.Response, localRefresh int, err error) {

	sanitizeObjectType(objectType)

	t0 := time.Now()
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

	t1 := time.Now()
	duration := t1.Sub(t0)

	// if the query is two slow, use backoff
	if int(duration.Seconds()) >= config.Refresh {
		localRefresh = config.LongRefresh
		log.Debugf("end of LQL %s query, took %v seconds, too slow, next refresh in %d seconds", objectType, duration.Seconds(), localRefresh)
	} else {
		log.Debugf("end of LQL %s query, took %v seconds", objectType, duration.Seconds())
		localRefresh = config.Refresh
	}

	return
}

// GetItems open livestatus connection then query it
func GetItems(objectType string, config *helpers.CONFIG, channel chan AlertStruct) {

	sanitizeObjectType(objectType)

	// services alerts
	client := livestatus.NewClient("tcp", config.Server)
	defer client.Close()

	// LQL query
	query := getQuery(objectType, *config)

	for {

		var class string

		if helpers.Pause {
			time.Sleep(5 * time.Second)
			continue
		}

		// make the LQL requests
		resp, localRefresh, err := GetResponse(objectType, client, query, config)
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

				notes_url, err := r.GetString("notes_url")
				if err != nil {
					log.Warn(err)
				} else {
					if len(notes_url) > 0 {
						log.Infof("%s: %s: %s", host, desc, notes_url)
					}
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

				if objectType == "services" {
					item = fmt.Sprintf("* %s: %s %s %s\n", host, desc, class, isFlappingStr)
				} else if objectType == "hosts" {
					item = fmt.Sprintf("* %s: %s %s\n", host, class, isFlappingStr)
				}
				items += item
			}
			// trim services
			items = strings.TrimRight(items, "\n")

			// make the alert
			alert = AlertStruct{Count: count, Items: items, Class: class}
		}

		// feed the alerts channel
		log.Debugf("sending %d %s alerts", alert.Count, objectType)
		channel <- alert

		// main sleep
		time.Sleep(time.Duration(localRefresh) * time.Second)
	}
}
