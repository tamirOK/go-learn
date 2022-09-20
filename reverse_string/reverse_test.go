package reverse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseStringWithOddLength(t *testing.T) {
	assert.Equal(t, String("abcde"), "edcba")
}

func TestReverseStringWithEvenLength(t *testing.T) {
	assert.Equal(t, String("qwer"), "rewq")
}

func TestReverseEmptyString(t *testing.T) {
	assert.Equal(t, String(""), "")
}

func TestReverseStringWithRunes(t *testing.T) {
	want := "gnaloG ьтавозьлопси ьсучу Я"

	assert.Equal(t, String("Я учусь использовать Golang"), want)
}
