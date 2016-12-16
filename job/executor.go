package job

import (
	"sync"
)

type JobExecutor interface {
	// Submit the job for execution.
	Submit(Job) bool
	// Unsubmit a running job.
	Unsubmit(string) (Job, bool)
	// Get a running job instance throught its name.
	GetJob(string) (Job, bool)
	// Stop all running jobs and shutdown this executor.
	Shutdown()
	// Block and wait for the finalization of all running jobs
	// since the shutdown may occur in a different goroutine.
	AwaitTermination()
}

// Default implementation of the JobExecutor.
// Should be sufficient in most cases.
type simpleExecutor struct {
	jobs     map[string]Job
	jobsLock sync.RWMutex
	wg       sync.WaitGroup
}

func (executor *simpleExecutor) Submit(job Job) bool {
	executor.jobsLock.Lock()
	defer executor.jobsLock.Unlock()
	if _, ok := executor.jobs[job.Name()]; ok {
		return false
	} else {
		go func() {
			executor.wg.Add(1)
			job.Run()
			executor.wg.Done()
		}()
		executor.jobs[job.Name()] = job
		return true
	}
}

func (executor *simpleExecutor) Unsubmit(name string) (Job, bool) {
	executor.jobsLock.Lock()
	defer executor.jobsLock.Unlock()
	if job, ok := executor.jobs[name]; ok {
		job.Stop()
		delete(executor.jobs, name)
		return job, ok
	} else {
		return nil, false
	}
}

func (executor *simpleExecutor) GetJob(name string) (Job, bool) {
	executor.jobsLock.RLock()
	defer executor.jobsLock.RUnlock()
	job, ok := executor.jobs[name]
	return job, ok
}

func (executor *simpleExecutor) Shutdown() {
	executor.jobsLock.Lock()
	defer executor.jobsLock.Unlock()
	for _, job := range executor.jobs {
		job.Stop()
	}
}

func (executor *simpleExecutor) AwaitTermination() {
	executor.wg.Wait()
}

var defaultExecutor = &simpleExecutor{jobs: make(map[string]Job)}

func DefaultExecutor() JobExecutor {
	return defaultExecutor
}
