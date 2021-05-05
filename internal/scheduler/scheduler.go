package scheduler

type Scheduler interface {
	Schedule(spec interface{}, job func()) (id interface{}, err error)
	Unschedule(id interface{}) error
	Start()
	Stop()
}
