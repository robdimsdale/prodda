package schedule

import (
	"os"

	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

type Runner struct {
	logger lager.Logger
	c      *cron.Cron
}

func NewRunner(c *cron.Cron, logger lager.Logger) Runner {
	return Runner{
		logger: logger,
		c:      c,
	}
}

func (a Runner) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	a.c.Start()
	a.logger.Info("Scheduler started")

	close(ready)

	<-signals
	a.c.Stop()
	return nil
}
