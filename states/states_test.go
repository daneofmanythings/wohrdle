package states

import (
	"testing"

	"gitlab.com/daneofmanythings/wohrdle/utils"
)

const (
	testWordCorrect string = "tests"
	testWordWrong   string = "wrong"
)

func mockNewGameSession() *GameSession {
	wordRepo := map[string][]string{
		"5": {
			testWordCorrect,
		},
	}
	params := NewDefaultParameters(wordRepo)
	return NewGameSession(params)
}

func TestIsValidWord(t *testing.T) {
	gs := mockNewGameSession()
	var convenience_word_string string
	for _, r := range testWordCorrect {
		gs.PushRune(r)
		convenience_word_string += utils.RuneToAlpha(r)
	}

	if !gs.isValidWord() {
		t.Fatalf("word=%s was not found in validWords", convenience_word_string)
	}
}

func TestIsWinner(t *testing.T) {
	gs := mockNewGameSession()
	for _, r := range testWordCorrect {
		gs.PushRune(r)
	}

	if !gs.IsWinner() {
		t.Fatal("winner was not detected")
	}

	for i := 0; i < len(testWordCorrect); i++ {
		gs.PopRune()
	}

	for _, r := range testWordWrong {
		gs.PushRune(r)
	}
	if gs.IsWinner() {
		t.Fatal("winner was incorrectly detected")
	}
}

func TestFinalizeCurRow(t *testing.T) {
	gs := mockNewGameSession()
	gs.finalizeCurRow()
	gs.curIdx -= 1
	for i, cell := range gs.Grid[gs.curIdx] {
		if cell.GetState() != CORRECT {
			t.Fatalf("incorrect state detected after FinalizeCurRow. cell.Char=%c, gs.Word[i]=%c", cell.Char, gs.targetWordAsRunes[i])
		}
	}
}

func TestCountByRune(t *testing.T) {
	gs := mockNewGameSession()
	for _, r := range gs.targetWordAsRunes {
		gs.PushRune(r)
	}
	// WARN: this is done manually. if testWordCorrect changes, this will also need to change
	targetMap := map[rune]int{
		'T': 2,
		'E': 1,
		'S': 2,
	}
	recievedMap := gs.countByRuneForCurRow()
	for k := range targetMap {
		if targetMap[k] != recievedMap[k] {
			t.Fatalf("incorrect count map generated:\n%v\nexpected:\n%v", recievedMap, targetMap)
		}
	}
}
