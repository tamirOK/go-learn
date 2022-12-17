package parallel

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTask(durationSeconds int) Task {
	return func() error {
		time.Sleep(time.Duration(durationSeconds) * time.Second)
		return nil
	}
}

func createTaskWithError(durationSeconds int) Task {
	return func() error {
		time.Sleep(time.Duration(durationSeconds) * time.Second)
		return errors.New("Error during calculation")
	}
}

func TestRunWithMultipleTasks(t *testing.T) {
	var tasks []Task

	for i := 0; i < 1000; i++ {
		tasks = append(tasks, createTask(rand.Intn(5)))
	}

	assert.Nil(t, Run(tasks, 100, 0))
}

func TestRunWhenTasksLessThanGoroutines(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(5)))
	}

	assert.Nil(t, Run(tasks, 1000, 0))
}

func TestRunMultipleTasksWithErrors(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(5)))
	}

	for i := 0; i < 50; i++ {
		tasks = append(tasks, createTaskWithError(rand.Intn(5)))
	}

	assert.Nil(t, Run(tasks, 100, 100))
}

func TestRunMultipleTasksWithTooManyErrors(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(5)))
	}

	for i := 0; i < 50; i++ {
		tasks = append(tasks, createTaskWithError(rand.Intn(5)))
	}

	assert.Equal(t, Run(tasks, 100, 10), ErrErrorsLimitExceeded)
}

func TestRunWhenErrorsCountZero(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(5)))
	}

	tasks = append(tasks, createTaskWithError(rand.Intn(5)))

	assert.Equal(t, Run(tasks, 100, 0), ErrErrorsLimitExceeded)
}

func TestRunWhenErrorsCountNegative(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(5)))
	}

	tasks = append(tasks, createTaskWithError(rand.Intn(5)))

	assert.Equal(t, Run(tasks, 100, -12), ErrErrorsLimitExceeded)
}

func TestRunEnsureTasksRunningConcurrently(t *testing.T) {
	var tasks []Task
	var totalTasksDuration time.Duration = 0

	for i := 0; i < 100; i++ {
		duration := rand.Intn(5)
		tasks = append(tasks, createTask(duration))
		totalTasksDuration += (time.Duration(duration) * time.Second)
	}

	start := time.Now()

	Run(tasks, 100, 0)

	elapsed := time.Since(start)

	assert.LessOrEqual(t, elapsed, totalTasksDuration)
}
