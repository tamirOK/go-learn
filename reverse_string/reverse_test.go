package reverse_string

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseStringWithOddLength(t *testing.T) {
	assert.Equal(t, ReverseString("abcde"), "edcba")
}

func TestReverseStringWithEvenLength(t *testing.T) {
	assert.Equal(t, ReverseString("qwer"), "rewq")
}

func TestReverseEmptyString(t *testing.T) {
	assert.Equal(t, ReverseString(""), "")
}

func TestReverseStringWithRunes(t *testing.T) {
	want := "gnaloG ьтавозьлопси ьсучу Я"

	assert.Equal(t, ReverseString("Я учусь использовать Golang"), want)
}
