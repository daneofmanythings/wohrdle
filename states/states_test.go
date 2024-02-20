package states

import (
	"fmt"
	"testing"

	"gitlab.com/daneofmanythings/wohrdle/utils"
)

const (
	wordTests string = "tests"
	wordWrong string = "wrong"
	wordVolts string = "volts"
)

func mockNewGameSession(word string) *GameSession {
	wordRepo := map[string][]string{
		fmt.Sprintf("%d", len(word)): {
			word,
		},
	}
	params := NewDefaultParameters(wordRepo)
	return NewGameSession(params)
}

func TestIsValidWord(t *testing.T) {
	gs := mockNewGameSession(wordTests)
	var convenience_word_string string
	for _, r := range wordTests {
		gs.PushRune(r)
		convenience_word_string += utils.RuneToAlpha(r)
	}

	if !gs.isValidWord() {
		t.Fatalf("word=%s was not found in validWords", convenience_word_string)
	}
}

func TestIsWinner(t *testing.T) {
	gs := mockNewGameSession(wordTests)
	for _, r := range wordTests {
		gs.PushRune(r)
	}

	if !gs.IsWinner() {
		t.Fatal("winner was not detected")
	}

	for i := 0; i < len(wordTests); i++ {
		gs.PopRune()
	}

	for _, r := range wordWrong {
		gs.PushRune(r)
	}
	if gs.IsWinner() {
		t.Fatal("winner was incorrectly detected")
	}
}

func TestFinalizeCurRow(t *testing.T) {
	// WARN: these tests do not use the word constants because the output is also
	// dependant the guess.
	testCases := []struct {
		name  string
		word  string
		guess string
		row   []Cell
	}{
		{
			"sanityCheck",
			"volts",
			"volts",
			[]Cell{{'V', CORRECT}, {'O', CORRECT}, {'L', CORRECT}, {'T', CORRECT}, {'S', CORRECT}},
		},
		{
			"twoInputOneOutputPreceeding",
			"volts",
			"lusts",
			[]Cell{{'L', PARTIAL}, {'U', USED}, {'S', USED}, {'T', CORRECT}, {'S', CORRECT}},
		},
		{
			"twoInputOneOutputFollowing",
			"stims",
			"sassy",
			[]Cell{{'S', CORRECT}, {'A', USED}, {'S', PARTIAL}, {'S', USED}, {'Y', USED}},
		},
		{
			"twoInputOneOutputFollowingWithOnePartial",
			"stims",
			"sissy",
			[]Cell{{'S', CORRECT}, {'I', PARTIAL}, {'S', PARTIAL}, {'S', USED}, {'Y', USED}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gs := mockNewGameSession(tc.word)
			for _, r := range tc.guess {
				gs.PushRune(r)
			}
			gs.finalizeCurRow()
			gs.curIdx -= 1
			curRow := gs.Grid[gs.curIdx]
			for i := range curRow {
				if !curRow[i].isEqualTo(tc.row[i]) {
					t.Fatalf(
						"unexpected cell value at %d, got=%v. expected=%v. word=%s, guess=%s",
						i,
						curRow[i],
						tc.row[i],
						tc.word,
						tc.guess,
					)
				}
			}
		})
	}
}

func TestCountByRune(t *testing.T) {
	gs := mockNewGameSession(wordTests)
	for _, r := range gs.targetWordAsRunes {
		gs.PushRune(r)
	}
	// WARN: this is done manually. if wordTests changes, this will also need to change
	targetMap := map[rune]int{
		'T': 2,
		'E': 1,
		'S': 2,
	}
	recievedMap := gs.countByRuneForCurRow()
	for k := range targetMap {
		if targetMap[k] != recievedMap[k] {
			t.Fatalf("unexpected count map for %s:\n%v\nexpected:\n%v", gs.targetWordAsString, recievedMap, targetMap)
		}
	}
	for k := range recievedMap {
		if targetMap[k] != recievedMap[k] {
			t.Fatalf("unexpected count map for %s:\n%v\nexpected:\n%v", gs.targetWordAsString, recievedMap, targetMap)
		}
	}
}
