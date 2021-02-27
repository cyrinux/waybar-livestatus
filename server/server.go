package server

import (
	"github.com/cyrinux/waybar-livestatus/alert"
	"github.com/cyrinux/waybar-livestatus/helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

// Listen start the gRPC server
func GRPCListen(alertsChan chan []*helpers.Alert, config *helpers.CONFIG) {
	go func() {
		for {
			alert.MenuEntries = <-alertsChan
		}
	}()
	sock := os.Getenv("XDG_RUNTIME_DIR") + "/waybar-livestatus.sock"
	lis, err := net.Listen("unix", sock)
	if err != nil {
		log.Fatalf("Failed to listen on unix socket: %v", err)
	}
	defer lis.Close()
	as := alert.Server{Config: config}

	gs := grpc.NewServer()

	// Allow method discovery
	reflection.Register(gs)

	alert.RegisterAlertServer(gs, &as)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC over unix socket: %v", err)
	}
}
