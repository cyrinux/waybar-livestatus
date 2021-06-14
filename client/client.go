package client

import (
	"fmt"
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
		fmt.Fprintf(os.Stderr, "could not connect: %s", err)
		return
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
			fmt.Fprintf(os.Stderr, "Error when calling GetAlertsList: %s", err)
		}
		fmt.Println(response.List)
		return
	}

	var response *alert.ResponseAlert
	response, err = c.GetNotesURL(context.Background(), &message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling GetNotesUrl: %s", err)
		return
	}
	if response.NotesUrl != "" {
		fmt.Print(response.NotesUrl)
	}

	return
}
