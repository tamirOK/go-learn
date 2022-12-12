package parallel

import "errors"

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// worker gets task from tasks channel, runs it and puts result into results channel
// When it receives message from quit channel, it stops.
func worker(tasks <-chan Task, quit <-chan bool, results chan<- error) {
	for {
		select {
		case task := <-tasks:
			results <- task()
		case <-quit:
			return
		}
	}
}

// createTasksChannel returns unbuffered channel which contains tasks.
func createTasksChannel(tasks []Task) chan Task {
	tasksChannel := make(chan Task)

	go func() {
		for _, task := range tasks {
			tasksChannel <- task
		}
	}()

	return tasksChannel
}

// runTasks runs n worker goroutines and returns two channels.
// First is a channel for consuming results.
// Second is a channel for stopping worker goroutines.
func runTasks(n int, tasksChannel chan Task) (<-chan error, chan bool) {
	resultsChannel := make(chan error)
	quitChannel := make(chan bool)

	for i := 0; i < n; i++ {
		go worker(tasksChannel, quitChannel, resultsChannel)
	}

	return resultsChannel, quitChannel
}

// checkTaskResults consumes results from resultsChannel taskCount times and
// checks for errors. If number of errors exceedes errorLimit, corresponding error is returned.
func checkTaskResults(
	taskCount int,
	errorLimit int,
	resultsChannel <-chan error,
	quitChannel chan<- bool,
) error {
	// Stop all worker goroutines after returning from this function
	defer func() {
		quitChannel <- true
	}()

	if errorLimit < 0 {
		errorLimit = 0
	}

	for i := 0; i < taskCount; i++ {
		result := <-resultsChannel

		if result != nil {
			errorLimit--
		}

		if errorLimit < 0 {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksChannel := createTasksChannel(tasks)
	resultsChannel, quitChannel := runTasks(n, tasksChannel)

	return checkTaskResults(len(tasks), m, resultsChannel, quitChannel)
}
