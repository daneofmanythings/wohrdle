package utils

import (
	"testing"
)

func TestRuneToAlpha(t *testing.T) {
	r := 'a'
	s := "a"

	if RuneToAlpha(r) != s {
		t.Fatalf("rune=%c does not equal string %s", r, s)
	}
}

func TestRuneSliceToUpper(t *testing.T) {
	input := []rune("Hello")
	answer := []rune("HELLO")
	output := RuneSliceToUpper(input)

	if len(input) != len(answer) {
		t.Fatalf("input len=%d did not match output len=%d", len(input), len(output))
	}
	if len(output) != len(answer) {
		t.Fatalf("output len=%d did not match answer len=%d", len(output), len(answer))
	}
	for i := range input {
		if output[i] != answer[i] {
			t.Fatalf("output=%s did not match answer=%s", string([]rune{input[i]}), string([]rune{answer[i]}))
		}
	}
}

func TestLoadWordRepoFromJSON(t *testing.T) {
	// this is error prone. make it better
	input_path := "./test_json.json"
	expected_words := []string{"test", "work"}

	wr, e := LoadWordRepoFromJSON(input_path)
	if e != nil {
		t.Fatal(e)
	}

	wr_slice := wr.Words["4"]
	for i := range expected_words {
		if wr_slice[i] != expected_words[i] {
			t.Fatalf("expected word not recieved. expected=%s, got=%s", expected_words[i], wr_slice[i])
		}
	}
}
