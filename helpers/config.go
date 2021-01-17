package helpers

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	toml "github.com/pelletier/go-toml"

	log "github.com/sirupsen/logrus"
)

// CONFIG define the configuration content
type CONFIG struct {
	Server                        string
	Refresh                       int      `toml:"refresh" default:"30"`
	LongRefresh                   int      `toml:"long_refresh" default:"60"`
	HostsPattern                  []string `toml:"hosts_pattern"`
	HostsPatternString            string
	Debug                         bool   `toml:"debug" default:"false"`
	Popup                         bool   `toml:"popup" default:"true"`
	Warnings                      bool   `toml:"warnings" default:"true"`
	Version                       bool   `default:"false"`
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
}

// GetConfig merge config from file and flag
// and return `config`
func GetConfig() *CONFIG {

	config := &CONFIG{}
	// set log formatter
	log.SetFormatter(&log.JSONFormatter{})

	// read toml config file
	user, _ := user.Current()
	homeDir := user.HomeDir

	// try to read config file
	configFile := homeDir + "/.config/waybar/livestatus.toml"
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Error(err)
	}
	err = toml.Unmarshal(file, config)
	if err != nil {
		log.Error(err)
	}
	fmt.Fprintf(os.Stderr, "Config file loaded: %s\n", configFile)
	config.HostsPatternString = strings.Join(config.HostsPattern, ",")

	// override config values with values from cli
	flag.StringVar(&config.Server, "s", config.Server, "Livestatus 'server:port'.")
	flag.BoolVar(&config.Warnings, "w", config.Debug, "Get also state warnings. Default show critical only.")
	flag.BoolVar(&config.Popup, "n", config.Popup, "Disable popup alert.")
	flag.BoolVar(&config.Debug, "d", config.Debug, "Get debug log.")
	flag.IntVar(&config.Refresh, "r", config.Refresh, "Refresh rate in seconds. Min 15.")
	flag.IntVar(&config.LongRefresh, "R", config.Refresh, "Long refresh rate in seconds.")
	flag.StringVar(&config.HostsPatternString, "H", config.HostsPatternString, "Hostname pattern comma separated.")
	flag.BoolVar(&config.Version, "V", false, "Print version and exit")
	flag.Parse() // Parse flags

	if config.Server == "" {
		fmt.Fprintf(os.Stderr, "The server can't be empty!")
		os.Exit(1)
	}
	// set log level
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// sanitize refresh
	if config.Refresh < 15 {
		log.Info("Refresh rate can't be under 15 seconds ! Fallback to default: 30 seconds")
		config.Refresh = 30
	}
	if config.LongRefresh < 30 {
		log.Info("Long refresh rate can't be under 30 seconds ! Fallback to default: 60 seconds")
		config.LongRefresh = 60
	}
	if config.ServicesOnly && config.HostsOnly {
		log.Error("services_only and hosts_only can't be set together")
		config.HostsOnly = false
		config.ServicesOnly = false
	}

	return config
}
