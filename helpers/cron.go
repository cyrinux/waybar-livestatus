package helpers

import (
	cron "github.com/robfig/cron"
)

// AutoPause cron
func AutoPause(config *CONFIG) {
	c := cron.New()
	c.AddFunc(config.AutoStart, func() { SetPause() })
	c.AddFunc(config.AutoStop, func() { SetResume() })
	c.Start()
}
