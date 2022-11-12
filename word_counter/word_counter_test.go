package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTop10WordsCounter(t *testing.T) {
	input := "cat and dog, one dog,two cats and one man"
	want := []string{
		"and (2)",
		"one (2)",
		"cat (1)",
		"cats (1)",
		"dog, (1)",
		"dog,two (1)",
		"man (1)",
	}

	assert.Equal(t, GetTop10FrequentWords(input), want)
}

func TestGetTop10FrequentWordsWhenWordInDifferentForms(t *testing.T) {
	input := "cat cats Cat cat,"
	want := []string{"Cat (1)", "cat (1)", "cat, (1)", "cats (1)"}
	assert.Equal(t, GetTop10FrequentWords(input), want)
}

func TestGetTop10FrequentWordsWithLongInput(t *testing.T) {
	input := `
	Plan 9 is a distributed operating system, designed to make a network of heterogeneous and
	geographically separated computers function as a single system.[38] In a typical Plan 9
	installation, users work at terminals running the window system rio, and they access CPU
	servers which handle computation-intensive processes. Permanent data storage is provided
	by additional network hosts acting as file servers and archival storage
	`

	assert.Equal(t, len(GetTop10FrequentWords(input)), 10)
}

func TestGetTop10FrequentWordsWithPunctuationSigns(t *testing.T) {
	input := ", . ,  ! "
	want := []string{", (2)", "! (1)", ". (1)"}

	assert.Equal(t, GetTop10FrequentWords(input), want)
}
