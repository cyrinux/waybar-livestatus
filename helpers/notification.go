package helpers

import (
	notifier "github.com/blakek/go-notifier"
	log "github.com/sirupsen/logrus"
)

// Alert define a event
type Alert struct {
	Host  string
	Desc  string
	Class string
}

func notify(title, message, icon string) (popup *notifier.Notification, err error) {
	if icon == "" {
		icon = "/usr/share/icons/Adwaita/32x32/legacy/dialog-warning.png"
	}
	popup = &notifier.Notification{Title: title, Message: message, ImagePath: icon}
	notifier, err := notifier.NewNotifier()
	if err != nil {
		log.Error(err)
		return
	}
	if err := notifier.DeliverNotification(*popup); err != nil {
		log.Error(err)
	}
	return
}

// SendNotification send a notification
func SendNotification(notifications chan *Alert, config *CONFIG) {

	alertsWithCounter := make(map[Alert]int, 20)

	if config.Debug {
		if popup, err := notify("Livestatus", "starting", ""); err != nil {
			log.Errorf("Error sending notification: %v", popup)
		}
	}

	for {
		notification := <-notifications

		// // check if notification not null
		if notification == nil {
			continue
		}

		// check if notification was already sent, and polled 10 times
		if alertsWithCounter[*notification] > 10 {
			// if polled > 10 times, then delete to be able to resent the notification
			delete(alertsWithCounter, *notification)
		}

		if notification.Class == "ok" || notification.Class == "error" {
			continue
		}

		if alertsWithCounter[*notification] == 0 {
			if popup, err := notify(notification.Host, notification.Desc, ""); err != nil {
				log.Errorf("Error sending notification: %v", popup)
			}
		}

		alertsWithCounter[*notification]++

		if config.Debug {
			log.Debugf("%+v sent %d times", notification, alertsWithCounter[*notification])
		}
	}
}
