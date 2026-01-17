package scheduler

// Scheduler defines a simple timed task scheduler to
// execute jobs at given times or intervals.
type Scheduler interface {

	// Schedule a new job with the given scheduling spec.
	// Returns a unique identifier of the scheduled job.
	Schedule(spec any, job func()) (id any, err error)

	// UnSchedule removes a given job from the scheduler.
	UnSchedule(id any) error

	// Start runs the scheduler cycle.
	Start()

	// Stop stops the scheduler cycle.
	Stop()
}
