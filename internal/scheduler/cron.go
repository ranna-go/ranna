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

func (lct *CronScheduler) Schedule(spec interface{}, job func()) (id interface{}, err error) {
	specStr, ok := spec.(string)
	if !ok {
		return nil, errors.New("invalid spec type: must be a string")
	}
	return lct.c.AddFunc(specStr, job)
}

func (lct *CronScheduler) Unschedule(id interface{}) error {
	cid, ok := id.(cron.EntryID)
	if !ok {
		return errors.New("invalid id type")
	}
	lct.c.Remove(cid)
	return nil
}

func (lct *CronScheduler) Start() {
	lct.c.Start()
}

func (lct *CronScheduler) Stop() {
	lct.c.Stop()
}
