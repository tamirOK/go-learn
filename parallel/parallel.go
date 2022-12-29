package parallel

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// worker gets task from tasks channel, runs it and puts result into results channel
// When it receives message from quit channel, it stops.
func Run(tasks []Task, n int, m int) error {
	tasksCh := make(chan func() error)
	resultsCh := make(chan error)
	quitCh := make(chan struct{})
	wg := sync.WaitGroup{}

	if m < 0 {
		m = 0
	}

	go func() {
		defer close(tasksCh)
		for _, task := range tasks {
			select {
			case tasksCh <- task:
			case <-quitCh:
				return
			}
		}
	}()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				select {
				case resultsCh <- task():
				case <-quitCh:
					return
				}
			}
		}()
	}

	go func() {
		for result := range resultsCh {
			if result != nil {
				m--

				if m < 0 {
					close(quitCh)
					return
				}
			}
		}
	}()

	wg.Wait()

	if m < 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
