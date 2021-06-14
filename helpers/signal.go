package helpers

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// Pause state
var Pause = false

const SleepTime = 500

// SetPause pause the polling
func SetPause() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Pause polling")
	Pause = true
	runtime.GC()
}

// SetResume resume from pause
func SetResume() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Start polling")
	Pause = false
}

func togglePause(sig os.Signal) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	Pause = !Pause
	log.Info().Msgf("Signal %v, pause %v", sig, Pause)
}

func stopSignal() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msg("Cleaning unix sock file and exit")
	os.Remove(os.Getenv("XDG_RUNTIME_DIR") + "/waybar-livestatus.sock")
	os.Exit(0)
}

// SignalHandler handle the SIGUSR1 to pause the app
func SignalHandler() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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
		time.Sleep(SleepTime * time.Millisecond)
	}
}
