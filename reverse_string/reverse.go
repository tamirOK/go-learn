package reverse_string

func ReverseString(input string) string {
	rune := []rune(input)
	size := len(runes)

	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
