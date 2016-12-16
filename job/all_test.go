package job

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type testRunner struct {
	err          error
	stoppingChan chan bool
}

func (r *testRunner) Name() string {
	return "test_runner"
}

func (r *testRunner) Stop() {
	r.LogError(Cancelled(r.Name()))
	close(r.stoppingChan)
}

func (r *testRunner) Start(wg *sync.WaitGroup) error {
	go r.start()
	return nil
}

func (r *testRunner) start() {
	for {
		select {
		case <-r.stoppingChan:
			return
		default:
			fmt.Println("I'm running!")
			time.Sleep(time.Second)
		}
	}
}

func (r *testRunner) LogError(err error) {
	r.err = err
}

func (r *testRunner) Stopping() <-chan bool {
	return r.stoppingChan
}

func (r *testRunner) Error() error {
	return r.err
}

func TestRunner(t *testing.T) {
	r := &testRunner{stoppingChan: make(chan bool)}
	r.Start(nil)
	time.Sleep(time.Duration(5) * time.Second)
	r.Stop()
}
