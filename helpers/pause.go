package helpers

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// Pause state
var Pause = false

// PauseHandler handle the SIGUSR1 to pause the app
func PauseHandler() {
	// channel to trap signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)

	for {
		sig := <-sigs
		Pause = !Pause // toggle pause
		log.Infof("signal %v, pause %v", sig, Pause)
		time.Sleep(1000 * time.Millisecond)
	}
}
