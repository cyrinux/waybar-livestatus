package helpers

import (
	"fmt"
	notify "github.com/TheCreeper/go-notify"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Alert define a event
type Alert struct {
	Host     string
	Desc     string
	Class    string
	NotesURL string
}

func send(alert *Alert, icon string) (notification *notify.Notification, err error) {
	if icon == "" {
		icon = "/usr/share/icons/Adwaita/32x32/legacy/dialog-warning.png"
	}
	ntf := notify.NewNotification(alert.Host, alert.Desc)
	ntf.AppIcon = icon
	ntf.Hints = make(map[string]interface{})
	ntf.Hints[notify.HintSoundFile] = "/usr/share/sounds/freedesktop/stereo/dialog-warning.oga"
	if strings.Contains(alert.Class, "critical") {
		ntf.Hints[notify.HintUrgency] = notify.UrgencyCritical
	} else if strings.Contains(alert.Class, "warning") {
		ntf.Hints[notify.HintUrgency] = notify.UrgencyNormal
	} else {
		ntf.Hints[notify.HintUrgency] = notify.UrgencyLow
	}

	if _, err = ntf.Show(); err != nil {
		return
	}

	return
}

// SendNotification send a notification
func SendNotification(notifications chan *Alert, config *CONFIG) {

	alertsWithCounter := make(map[Alert]int)

	if config.Debug {
		startAlert := Alert{Host: "Livestatus", Desc: fmt.Sprintf("starting version %v", Version)}
		if notification, err := send(&startAlert, ""); err != nil {
			log.Errorf("Error sending notification: %v", notification)
		}
	}

	for {
		notification := <-notifications

		// check if notification not null
		if notification == nil {
			continue
		}

		// check if notification was already sent, and polled 10 times
		if alertsWithCounter[*notification] > config.NotificationSnoozeCycle {
			delete(alertsWithCounter, *notification)
		}

		if notification.Class == "ok" || notification.Class == "error" {
			continue
		}

		if alertsWithCounter[*notification] == 0 {
			if notification, err := send(notification, ""); err != nil {
				log.Errorf("Error sending notification: %v", notification)
			}
		}

		alertsWithCounter[*notification]++

		log.Debugf("%v sent %d times", notification, alertsWithCounter[*notification])
	}
}
