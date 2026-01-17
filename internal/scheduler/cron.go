package scheduler

import (
	"errors"

	"github.com/robfig/cron/v3"
)

// CronScheduler implements scheduler using
// a cron-like schedule spec syntax.
type CronScheduler struct {
	c *cron.Cron
}

// NewCronScheduler returns a new CronScheduler instance.
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		c: cron.New(),
	}
}

func (t *CronScheduler) Schedule(spec any, job func()) (id any, err error) {
	specStr, ok := spec.(string)
	if !ok {
		return nil, errors.New("invalid spec type: must be a string")
	}
	return t.c.AddFunc(specStr, job)
}

func (t *CronScheduler) UnSchedule(id any) error {
	cid, ok := id.(cron.EntryID)
	if !ok {
		return errors.New("invalid id type")
	}
	t.c.Remove(cid)
	return nil
}

func (t *CronScheduler) Start() {
	t.c.Start()
}

func (t *CronScheduler) Stop() {
	t.c.Stop()
}
