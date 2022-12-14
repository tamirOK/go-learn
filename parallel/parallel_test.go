package parallel

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const MaxSleepTime = 3

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

	for i := 0; i < 10; i++ {
		tasks = append(tasks, createTask(rand.Intn(MaxSleepTime)))
	}

	assert.Nil(t, Run(tasks, 10, 0))
}

func TestRunWhenTasksLessThanGoroutines(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(MaxSleepTime)))
	}

	assert.Nil(t, Run(tasks, 1000, 0))
}

func TestRunMultipleTasksWithErrors(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(MaxSleepTime)))
	}

	for i := 0; i < 50; i++ {
		tasks = append(tasks, createTaskWithError(rand.Intn(MaxSleepTime)))
	}

	assert.Nil(t, Run(tasks, 100, 100))
}

func TestRunMultipleTasksWithTooManyErrors(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(MaxSleepTime)))
	}

	for i := 0; i < 50; i++ {
		tasks = append(tasks, createTaskWithError(rand.Intn(MaxSleepTime)))
	}

	assert.Equal(t, Run(tasks, 100, 10), ErrErrorsLimitExceeded)
}

func TestRunWhenErrorsCountZero(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(MaxSleepTime)))
	}

	tasks = append(tasks, createTaskWithError(rand.Intn(MaxSleepTime)))

	assert.Equal(t, Run(tasks, 100, 0), ErrErrorsLimitExceeded)
}

func TestRunWhenErrorsCountNegative(t *testing.T) {
	var tasks []Task

	for i := 0; i < 100; i++ {
		tasks = append(tasks, createTask(rand.Intn(MaxSleepTime)))
	}

	tasks = append(tasks, createTaskWithError(rand.Intn(MaxSleepTime)))

	assert.Equal(t, Run(tasks, 100, -12), ErrErrorsLimitExceeded)
}

func TestRunEnsureTasksRunningConcurrently(t *testing.T) {
	var tasks []Task
	var totalTasksDuration time.Duration

	for i := 0; i < 100; i++ {
		duration := rand.Intn(MaxSleepTime)
		tasks = append(tasks, createTask(duration))
		totalTasksDuration += (time.Duration(duration) * time.Second)
	}

	start := time.Now()

	Run(tasks, 100, 0)

	elapsed := time.Since(start)

	assert.LessOrEqual(t, elapsed, totalTasksDuration)
}
