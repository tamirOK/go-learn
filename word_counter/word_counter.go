package counter

import (
	"fmt"
	"sort"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

func (w wordCount) String() string {
	return fmt.Sprintf("%s (%d)", w.word, w.count)
}

func countWordFrequency(words []string) map[string]int {
	counter := make(map[string]int)
	for _, word := range words {
		counter[word]++
	}
	return counter
}

func sortWordsByFrequency(wordCounter map[string]int) []wordCount {
	wordCounts := make([]wordCount, 0, len(wordCounter))

	for word, count := range wordCounter {
		wordCounts = append(wordCounts, wordCount{word: word, count: count})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		left := wordCounts[i]
		right := wordCounts[j]

		if left.count > right.count {
			return true
		}

		// if word frequencies are same, then sort lexicographically
		if left.count == right.count && left.word < right.word {
			return true
		}

		return false
	})

	return wordCounts
}

func takeTop10FrequentWords(sortedWordCounts []wordCount) []string {
	result := make([]string, 0, 10)

	for i, word := range sortedWordCounts {
		if i == 10 {
			break
		}
		result = append(result, word.String())
	}

	return result
}

func GetTop10FrequentWords(input string) []string {
	splitted := strings.Fields(input)
	wordFrequency := countWordFrequency(splitted)
	sortedWords := sortWordsByFrequency(wordFrequency)
	top10FrequentWords := takeTop10FrequentWords(sortedWords)

	return top10FrequentWords
}
