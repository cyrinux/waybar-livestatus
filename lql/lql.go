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
}

func getQuery(objectType string, warnings bool, hostsArray []string) (q *livestatus.Query) {
	q = livestatus.NewQuery(objectType)
	q.Columns("host_name", "state", "description")
	q.Filter("is_problem = 1")
	q.Filter(fmt.Sprintf("state = %v", statusCritical))
	if warnings {
		q.Filter(fmt.Sprintf("state = %v", statusWarning))
		q.Or(2)
	}
	q.Filter("acknowledged = 0")           // not acknowledged
	q.Filter("notifications_enabled = 1")  // notifications are enabled
	q.Filter("in_notification_period = 1") // in notification period
	q.Filter("scheduled_downtime_depth = 0")
	scheduleDowntimeDepthFilter := fmt.Sprintf("%s_scheduled_downtime_depth = 0",
		strings.TrimRight(objectType, "s"))
	q.Filter(scheduleDowntimeDepthFilter)
	for _, hostName := range hostsArray {
		log.Debugf("filter host %s", hostName)
		filter := fmt.Sprintf("host_name ~ %s", hostName)
		q.Filter(filter)
	}
	countHosts := len(hostsArray)
	if countHosts > 1 {
		q.Or(countHosts)
	}

	return
}

// GetResponse parse the response from the livestatus client
// objectType: hosts or services
func GetResponse(objectType string, c livestatus.Client, q *livestatus.Query, refresh int) (resp *livestatus.Response, localRefresh int, err error) {
	longRefresh := 60 // backoff refresh value
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
	if int(duration.Seconds()) >= refresh {
		localRefresh = longRefresh
		log.Debugf("end of LQL %s query, took %v seconds, too slow, next refresh in %d seconds", objectType, duration.Seconds(), localRefresh)
	} else {
		log.Debugf("end of LQL %s query, took %v seconds", objectType, duration.Seconds())
		localRefresh = refresh
	}

	return
}

// GetItems open livestatus connection then query it
func GetItems(objectType string, server string, warnings, debug bool, channel chan AlertStruct, refresh int, hostPattern string) {
	// services alerts
	c := livestatus.NewClient("tcp", server)
	defer c.Close()

	// LQL query
	hostsArray := strings.Split(hostPattern, ",")
	query := getQuery(objectType, warnings, hostsArray)

	for {

		if helpers.Pause {
			time.Sleep(5 * time.Second)
			continue
		}

		// make the LQL requests
		resp, localRefresh, err := GetResponse(objectType, *c, query, refresh)
		if err != nil {
			time.Sleep(10 * time.Second) // backoff
			continue
		}

		items := ""
		class := "ok"

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

			// keep the worse state
			if state == statusCritical {
				class = "critical"
			} else if state == statusWarning && class != "critical" {
				class = "warning"
			}

			// flapping handle
			if objectType == "hosts" {
				isFlapping, err := r.GetInt("is_flapping")
				if err != nil {
					log.Warn(err)
				}

				isFlappingStr := strconv.FormatInt(isFlapping, 10)
				if isFlappingStr == "1" {
					isFlappingStr = " "
				} else {
					isFlappingStr = ""
				}
				item = host + ": " + class + " " + isFlappingStr + "\n"
			} else {
				item = host + ": " + desc + " : " + class + "\n"
			}

			items += item
		}

		// trim services
		items = strings.TrimRight(items, "\n")

		// make the alert
		alert := AlertStruct{Count: count, Items: items, Class: class}

		// feed the alerts channel
		select {
		case channel <- alert:
			log.Debugf("sent alerts %+v", alert.Items)
		default:
		}
		// main sleep
		time.Sleep(time.Duration(localRefresh) * time.Second)

	}
}
