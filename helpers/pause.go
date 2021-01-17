package helpers

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// Pause state
var Pause = false

// SetPause pause the polling
func SetPause() {
	log.Info("Pause polling")
	Pause = true
	runtime.GC()
}

// SetResume resume from pause
func SetResume() {
	log.Info("Start polling")
	Pause = false
}

func togglePause() {
	Pause = !Pause
}

// PauseHandler handle the SIGUSR1 to pause the app
func PauseHandler() {
	// channel to trap signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)

	for {
		sig := <-sigs
		togglePause()
		log.Infof("signal %v, pause %v", sig, Pause)
		time.Sleep(1 * time.Second)
	}
}
