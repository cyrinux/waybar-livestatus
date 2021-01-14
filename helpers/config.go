package helpers

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	toml "github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// TOML define the configuration content
type TOML struct {
	Server             string
	Refresh            int `default:30`
	HostsPattern       []string
	HostsPatternString string
	Debug              bool `default:false`
	Popup              bool `default:true`
	Warnings           bool `default:true`
	Version            bool `default:false`
}

// GetConfig merge config from file and flag
// and return `config`
func GetConfig() (config TOML) {
	// set log formatter
	log.SetFormatter(&log.JSONFormatter{})

	// read toml config file
	user, _ := user.Current()
	homeDir := user.HomeDir

	// default values
	config.Debug = false
	config.Popup = true
	config.Refresh = 30
	config.Warnings = true

	// try to read config file
	configFile := homeDir + "/.config/waybar/livestatus.toml"
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Error(err)
	}
	fmt.Fprintf(os.Stderr, "Config file loaded: %s\n", configFile)
	config.HostsPatternString = strings.Join(config.HostsPattern, ",")

	// override config values with values from cli
	flag.StringVar(&config.Server, "s", config.Server, "Livestatus 'server:port'")
	flag.BoolVar(&config.Warnings, "w", config.Debug, "Get also state warnings.")
	flag.BoolVar(&config.Popup, "n", config.Popup, "Send popup alert if too many alerts unhandled.")
	flag.BoolVar(&config.Debug, "d", config.Debug, "Get debug log.")
	flag.IntVar(&config.Refresh, "r", config.Refresh, "Refresh rate in seconds. Min 15.")
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

	return
}
