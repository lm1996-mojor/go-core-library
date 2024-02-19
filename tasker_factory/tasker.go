package tasker_factory

import (
	"github.com/robfig/cron/v3"
)

func GetCornTasker(opts ...cron.Option) *cron.Cron {
	if len(opts) > 0 {
		return cron.New(opts...)
	} else {
		return cron.New()
	}
}
