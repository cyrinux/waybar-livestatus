package client

import (
	"log"
	"os"

	"github.com/cyrinux/waybar-livestatus/alert"
	"github.com/cyrinux/waybar-livestatus/helpers"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Start start the gRPC client
func Start(config *helpers.CONFIG) (err error) {
	var conn *grpc.ClientConn
	sock := "passthrough:///unix://" + os.Getenv("XDG_RUNTIME_DIR") + "/waybar-livestatus.sock"
	conn, err = grpc.Dial(sock, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()
	c := alert.NewAlertClient(conn)

	var host, service string

	if config.Client && len(os.Args) == 4 {
		host = os.Args[2]
		service = os.Args[3]
	}

	message := alert.RequestAlert{
		Host:    host,
		Service: service,
	}

	if host == "" || service == "" {
		var response *alert.ResponseAlertsList
		response, err = c.GetAlertsList(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling GetAlertsList: %s", err)
		}
		log.Print(response.List)
		return
	}

	var response *alert.ResponseAlert
	response, err = c.GetNotesURL(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling GetNotesUrl: %s", err)
	}
	log.Print(response.NotesUrl)

	return
}
