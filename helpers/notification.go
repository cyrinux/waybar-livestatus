package helpers

import (
	notification "github.com/blakek/go-notifier"
	log "github.com/sirupsen/logrus"
)

// SendNotification send a notification
func SendNotification(title, message, icon string) (popup *notification.Notification, err error) {
	if icon == "" {
		icon = "/usr/share/icons/Adwaita/32x32/legacy/dialog-warning.png"
	}
	popup = &notification.Notification{Title: title, Message: message, ImagePath: icon}
	notifier, err := notification.NewNotifier()
	if err != nil {
		log.Error(err)
		return
	}
	if err := notifier.DeliverNotification(*popup); err != nil {
		log.Error(err)
	}
	return
}
