package dictionary

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

// Dictionary represents a simple in-memory word store.
type Dictionary struct {
	words map[string]bool
}

// LoadFromFile loads words from a given file path.
func LoadDictionary(path string) (*Dictionary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	words := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words[strings.ToLower(word)] = true
		}
	}

	return &Dictionary{words: words}, scanner.Err()
}

// IsValid checks if a word exists in the dictionary.
func (d *Dictionary) IsValid(word string) bool {
	return d.words[strings.ToLower(word)]
}

func (d *Dictionary) RandomWord() string {
	keys := make([]string, 0, len(d.words))
	for k := range d.words {
		keys = append(keys, k)
	}
	return keys[rand.Intn(len(keys))]
}
