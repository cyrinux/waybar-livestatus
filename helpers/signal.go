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

func togglePause(sig os.Signal) {
	Pause = !Pause
	log.Infof("Signal %v, pause %v", sig, Pause)
}

func stopSignal() {
	log.Debugf("Cleaning unix sock file and exit")
	os.Remove(os.Getenv("XDG_RUNTIME_DIR") + "/waybar-livestatus.sock")
	os.Exit(0)
}

// SignalHandler handle the SIGUSR1 to pause the app
func SignalHandler() {
	// channel to trap signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGTERM, os.Interrupt)

	for {
		sig := <-sigs
		if sig == syscall.SIGUSR1 {
			togglePause(sig)
		} else {
			stopSignal()
		}
		time.Sleep(500 * time.Millisecond)
	}
}
