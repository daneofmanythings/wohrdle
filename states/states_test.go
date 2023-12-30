package states

import (
	"strconv"
	"testing"

	"gitlab.com/daneofmanythings/worhdle/utils"
)

const (
	testWordCorrect string = "hello"
	testWordWrong   string = "world"
)

var alertChan = make(chan bool)

func testnewGameSession() *GameSession {
	vw := map[string][]string{}
	int_str := strconv.Itoa(len(testWordCorrect))
	vw[int_str] = make([]string, 1)
	vw[int_str] = append(vw[int_str], testWordCorrect)
	p := Parameters{
		5,
		6,
		vw,
		0,
		5,
		5,
	}
	return NewGameSession(alertChan, &p)
}

func TestIsValidWord(t *testing.T) {
	gs := testnewGameSession()
	var convenience_word_string string
	for _, r := range testWordCorrect {
		gs.PushRune(r)
		convenience_word_string += utils.RuneToAlpha(r)
	}

	if !gs.isValidWord() {
		t.Fatalf("word=%s was not found in validWords", convenience_word_string)
	}
}

// TODO: fix this test
func TestIsWinner(t *testing.T) {
	gs := testnewGameSession()
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
	gs := testnewGameSession()
	gs.finalizeCurRow()
	gs.curIdx -= 1
	for i, cell := range gs.Grid[gs.curIdx] {
		if cell.GetState() != CORRECT {
			t.Fatalf("incorrect state detected after FinalizeCurRow. cell.Char=%c, gs.Word[i]=%c", cell.Char, gs.targetWordAsRunes[i])
		}
	}
}
