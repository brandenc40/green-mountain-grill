package scheduler

import (
	"github.com/jasonlvhit/gocron"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

// New -
func New() *gocron.Scheduler {
	return gocron.NewScheduler()
}
