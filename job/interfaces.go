package job

import (
	"fmt"
	"sync"
)

type Namer interface {
	// Return the human readable name.
	Name() string
}

type Stopper interface {
	// Could be stopped manually.
	Stop()
}

// Runner abstracts the smallest unit of a task of a job.
type Runner interface {
	Namer
	Stopper
	// Start the runner. If the wg is not nil, the wg.Done() method
	// must be called properly whether this method is being executed
	// in a seperate goroutine.
	Start(wg *sync.WaitGroup) error
	// Log error message.
	LogError(err error)
	// Return a channel which will be closed so can be used by goroutines as a
	// notification when the runner ends abnormally, usually because of being stopped
	// manually or getting some error.
	Stopping() <-chan bool
	// Return the error caused the runner to stop or nil if ends successfully.
	Error() error
}

// A single job may included periodically running runners, a series
// of runners or a single runner, etc..
type OnceJob interface {
	Namer
	Stopper
	// Run the job once.
	RunOnce() error
	// Interrupt this job. Unlike the Stop() which actually
	// ends the running of the job, this method may have a coordination
	// behavior as Java interruption.
	Interrupt()
}

type Job interface {
	OnceJob
	// Run the job.
	Run() error
}

// Error type indicates a runner or job has been cancelled.
type Cancelled string

func (e Cancelled) Error() string {
	return fmt.Sprintf("%s has been cancelled.", string(e))
}
