package reverse_string

import "testing"

func TestReverseStringWithOddLength(t *testing.T) {
	want := "edcba"

	if got := ReverseString("abcde"); got != want {
		t.Errorf("reversed string is not %q, got=%q", want, got)
	}
}

func TestReverseStringWithEvenLength(t *testing.T) {
	want := "rewq"

	if got := ReverseString("qwer"); got != want {
		t.Errorf("reversed string is not %q, got=%q", want, got)
	}
}

func TestReverseEmptyString(t *testing.T) {
	want := ""

	if got := ReverseString(""); got != want {
		t.Errorf("reversed string is not %q, got=%q", want, got)
	}
}

func TestReverseStringWithRunes(t *testing.T) {
	want := "gnaloG ьтавозьлопси ьсучу Я"

	if got := ReverseString("Я учусь использовать Golang"); got != want {
		t.Errorf("reversed string is not %q, got=%q", want, got)
	}
}
