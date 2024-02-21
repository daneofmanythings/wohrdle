package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"unicode"
)

type WordRepository struct {
	Words map[string][]string // the keys are the len of the contained words
}

func Find[T comparable](s []T, t T) int {
	for i := range s {
		if s[i] == t {
			return i
		}
	}
	return len(s)
}

func RuneIsAlpha(r rune) bool {
	return unicode.In(r, unicode.Latin)
}

func RuneToAlpha(r rune) string {
	return fmt.Sprintf("%c", r)
}

func RuneSliceToUpper(rs []rune) []rune {
	result := make([]rune, 0)
	for _, r := range rs {
		result = append(result, unicode.ToUpper(r))
	}
	return result
}

func LoadWordRepoFromJSON(path string) (WordRepository, error) {
	file, err := os.Open(path)
	if err != nil {
		return WordRepository{}, err
	}

	wr := WordRepository{}

	byteVal, _ := io.ReadAll(file)

	json.Unmarshal(byteVal, &wr)

	return wr, nil
}

func LoadEmbeddedWordRepo(bytes []byte) (WordRepository, error) {
	wr := WordRepository{}

	err := json.Unmarshal(bytes, &wr)
	if err != nil {
		return WordRepository{}, err
	}

	return wr, nil
}
