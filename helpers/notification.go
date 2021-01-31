package helpers

import (
	"fmt"
	notifier "github.com/blakek/go-notifier"
	log "github.com/sirupsen/logrus"
)

// Alert define a event
type Alert struct {
	Host     string
	Desc     string
	Class    string
	NotesURL string
}

func notify(title, message, icon string) (notification *notifier.Notification, err error) {
	if icon == "" {
		icon = "/usr/share/icons/Adwaita/32x32/legacy/dialog-warning.png"
	}
	notification = &notifier.Notification{Title: title, Message: message, ImagePath: icon}
	notifier, err := notifier.NewNotifier()
	if err != nil {
		log.Error(err)
		return
	}
	if err := notifier.DeliverNotification(*notification); err != nil {
		log.Error(err)
	}
	return
}

// SendNotification send a notification
func SendNotification(notifications chan *Alert, config *CONFIG) {

	alertsWithCounter := make(map[Alert]int)

	if config.Debug {
		if notification, err := notify("Livestatus", fmt.Sprintf("starting version %v", Version), ""); err != nil {
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
			if notification, err := notify(notification.Host, notification.Desc, ""); err != nil {
				log.Errorf("Error sending notification: %v", notification)
			}
		}

		alertsWithCounter[*notification]++

		log.Debugf("%v sent %d times", notification, alertsWithCounter[*notification])
	}
}
