package helpers

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"strings"

	toml "github.com/pelletier/go-toml"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Version give the software version
var Version string

// CONFIG define the configuration content
type CONFIG struct {
	Server                        string   `toml:"server" default:""`
	Refresh                       int      `toml:"refresh" default:"30"`
	LongRefresh                   int      `toml:"long_refresh" default:"60"`
	HostsPattern                  []string `toml:"hosts_pattern"`
	HostsPatternString            string
	Client                        bool   `default:"false"`
	Debug                         bool   `toml:"debug" default:"false"`
	Popup                         bool   `toml:"popup" default:"true"`
	Warnings                      bool   `toml:"warnings" default:"true"`
	Version                       bool   `default:"false"`
	NotificationSnoozeCycle       int    `toml:"notification_snooze_cycle" default:"10"`
	Acknowledged                  int    `toml:"acknowledged" default:"0"`
	NotificationsEnabled          int    `toml:"notifications_enabled" default:"1"`
	InNotificationPeriod          int    `toml:"in_notification_period" default:"1"`
	ScheduledDowntimeDepth        int    `toml:"scheduled_downtime_depth" default:"0"`
	ServiceScheduledDowntimeDepth int    `toml:"service_scheduled_downtime_depth" default:"0"`
	HostScheduledDowntimeDepth    int    `toml:"host_scheduled_downtime_depth" default:"0"`
	ServicePrefix                 string `toml:"service_prefix" default:""`
	HostPrefix                    string `toml:"host_prefix" default:""`
	PausePrefix                   string `toml:"pause_prefix" default:""`
	ErrorPrefix                   string `toml:"error_prefix" default:""`
	OkPrefix                      string `toml:"ok_prefix" default:""`
	FlappingIcon                  string `toml:"flapping_icon" default:""`
	HostsOnly                     bool   `toml:"hosts_only" default:"false"`
	ServicesOnly                  bool   `toml:"services_only" default:"false"`
	NotesURL                      bool   `toml:"notes_url" default:"false"`
	Limit                         int    `toml:"limit" default:"0"`
	GetDuration                   bool   `toml:"get_duration" default:"false"`
	// eg: 0 9 * * *
	AutoStart string `toml:"auto_start" default:""`
	// eg: 0 22 * * *
	AutoStop string `toml:"auto_stop" default:""`
}

// GetConfig merge config from file and flag
// and return `config`
func GetConfig() *CONFIG {
	config := &CONFIG{}
	// set log formatter
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// read toml config file
	user, _ := user.Current()
	homeDir := user.HomeDir

	// try to read config file
	configFile := homeDir + "/.config/waybar/livestatus.toml"
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Error().Err(err)
	}
	err = toml.Unmarshal(file, config)
	if err != nil {
		log.Error().Err(err)
	}
	log.Info().Msgf("Config file loaded: %s\n", configFile)
	config.HostsPatternString = strings.Join(config.HostsPattern, ",")
	// check if 'client' not define in the toml config
	if config.Client {
		log.Fatal().Msg("You are not allowed to define 'client' mode in the config file!")
	}

	// override config values with values from cli

	flag.BoolVar(&config.Debug, "d", config.Debug, "Debug mode.")
	flag.BoolVar(&config.Client, "c", false, "Client mode.")
	flag.BoolVar(&config.NotesURL, "u", config.NotesURL, "Display notes_url.")
	flag.BoolVar(&config.Popup, "n", config.Popup, "Disable notification popup alert.")
	flag.BoolVar(&config.Version, "V", false, "Print version and exit.")
	flag.BoolVar(&config.Warnings, "w", config.Warnings, "Get also state warnings. Default show critical only.")

	flag.IntVar(&config.Refresh, "r", config.Refresh, "Refresh rate in seconds. Min 15.")
	flag.IntVar(&config.LongRefresh, "R", config.LongRefresh, "Long refresh rate in seconds.")
	flag.IntVar(&config.NotificationSnoozeCycle, "N", config.NotificationSnoozeCycle, "Notifications snooze cycle.")
	flag.StringVar(&config.Server, "s", config.Server, "Livestatus 'server:port'.")
	flag.StringVar(&config.HostsPatternString, "H", config.HostsPatternString, "Hostname pattern comma separated.")
	flag.Parse() // Parse flags

	// if not server exit
	if config.Server == "" {
		fmt.Fprintf(os.Stderr, "The server can't be empty!")
		os.Exit(1)
	}
	// if string without port add it
	if _, _, err := net.SplitHostPort(config.Server); err != nil {
		config.Server += ":50000"
	}

	// set log level
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// sanitize server
	config.Server = strings.TrimSpace(config.Server)

	// sanitize refresh
	if config.Refresh < 15 {
		log.Info().Msg("Refresh rate can't be under 15 seconds ! Fallback to default: 30 seconds")
		config.Refresh = 30
	}
	if config.LongRefresh < 30 {
		log.Info().Msg("Long refresh rate can't be under 30 seconds ! Fallback to default: 60 seconds")
		config.LongRefresh = 60
	}
	if config.NotificationSnoozeCycle < 0 {
		config.NotificationSnoozeCycle = 10
	}
	if config.ServicesOnly && config.HostsOnly {
		log.Error().Msg("The params services_only and hosts_only can't be set together")
		config.HostsOnly = false
		config.ServicesOnly = false
	}

	return config
}
