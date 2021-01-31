package alert

import (
	"fmt"
	"github.com/cyrinux/waybar-livestatus/helpers"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// MenuEntries contains the alerts list
var MenuEntries []*helpers.Alert

var menuEntriesString string

// Server define a gRPC server
type Server struct {
	Config *helpers.CONFIG
}

// GetAlertsList return the alerts list
func (s *Server) GetAlertsList(ctx context.Context, message *RequestAlert) (*ResponseAlertsList, error) {
	var tempMenusEntryString string
	for _, alert := range MenuEntries {
		if s.Config.NotesURL {
			if alert.NotesURL != "" {
				tempMenusEntryString += fmt.Sprintf("%q %q: %s\n", alert.Host, alert.Desc, alert.NotesURL)
			} else {
				tempMenusEntryString += fmt.Sprintf("%q %q\n", alert.Host, alert.Desc)
			}
		} else {
			tempMenusEntryString += fmt.Sprintf("%q %q\n", alert.Host, alert.Desc)
		}
	}

	if tempMenusEntryString != "" {
		menuEntriesString = tempMenusEntryString
	}

	return &ResponseAlertsList{List: menuEntriesString}, nil
}

// GetNotesURL return the notes_url from a host / service
func (s *Server) GetNotesURL(ctx context.Context, message *RequestAlert) (*ResponseAlert, error) {
	log.Debugf("Received message body from client: %s / %s", message.Host, message.Service)

	for _, alert := range MenuEntries {
		if alert.Host == message.Host && alert.Desc == message.Service {
			return &ResponseAlert{NotesUrl: alert.NotesURL}, nil
		}
	}

	return &ResponseAlert{NotesUrl: ""}, nil
}
