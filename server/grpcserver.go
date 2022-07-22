package server

import (
	"net"
	"os"

	"github.com/cyrinux/waybar-livestatus/alert"
	"github.com/cyrinux/waybar-livestatus/helpers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// GRPCListen start the gRPC server
func GRPCListen(alertsChan chan []*helpers.Alert, config *helpers.CONFIG) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	go func() {
		for {
			alert.MenuEntries = <-alertsChan
		}
	}()
	sock := os.Getenv("XDG_RUNTIME_DIR") + "/waybar-livestatus.sock"
	lis, err := net.Listen("unix", sock)

	defer func() {
		lis.Close()
		os.Remove(sock)
	}()

	if err != nil {
		log.Fatal().Msgf("Failed to listen on unix socket: %v", err)
	}
	as := alert.Server{Config: config}

	gs := grpc.NewServer()

	alert.RegisterAlertServer(gs, &as)

	if err := gs.Serve(lis); err != nil {
		log.Fatal().Msgf("Failed to serve gRPC over unix socket: %v", err)
	}
}
