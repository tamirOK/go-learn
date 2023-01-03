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
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	if m < 0 {
		m = 0
	}

	go func() {
		defer close(tasksCh)
		for _, task := range tasks {
			mu.Lock()
			if m < 0 {
				mu.Unlock()
				break
			}
			mu.Unlock()
			tasksCh <- task
		}
	}()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				result := task()

				mu.Lock()
				if result != nil {
					m--
				}

				if m < 0 {
					mu.Unlock()
					return
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if m < 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
