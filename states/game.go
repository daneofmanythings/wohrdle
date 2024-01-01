package states

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/worhdle/utils"
)

type CellState int

const (
	DEFAULT CellState = iota
	CORRECT
	PARTIAL
	USED
)

var cellStates []CellState = []CellState{DEFAULT, CORRECT, PARTIAL, USED}

type Cell struct {
	Char  rune
	state CellState
}

func (c *Cell) SetState(state CellState) {
	if !slices.Contains(cellStates, state) {
		return
	}
	c.state = state
}

func (c *Cell) GetState() CellState {
	return c.state
}

var AllRunes []rune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewSeenCharRecord() []Cell {
	seenCharRecord := []Cell{}
	for _, r := range AllRunes {
		seenCharRecord = append(seenCharRecord, Cell{r, DEFAULT})
	}
	return seenCharRecord
}

type GameState int

const (
	ACTIVE = iota
	VICTORY
	LOSS
)

type GameSession struct {
	Parameters Parameters

	WordLen            int
	NumGuesses         int
	targetWordAsRunes  []rune
	targetWordAsString string

	Grid       [][]Cell
	curIdx     int
	SeenChars  []Cell
	validWords []string
	HelpText   string

	state GameState
}

func NewGameSession(params *Parameters) *GameSession {
	gs := &GameSession{
		Parameters: *params,
		WordLen:    params.WordLen,
		NumGuesses: params.NumGuesses,
		curIdx:     0,
		validWords: params.ValidWords(),
		state:      ACTIVE,
		SeenChars:  NewSeenCharRecord(),
	}
	word := gs.validWords[rand.Intn(len(gs.validWords))]
	gs.targetWordAsString = strings.ToUpper(word)
	gs.targetWordAsRunes = utils.RuneSliceToUpper([]rune(word))

	gs.Grid = make([][]Cell, gs.NumGuesses)
	for i := range gs.Grid {
		gs.Grid[i] = []Cell{}
	}

	return gs
}

func (gs *GameSession) GetState() GameState {
	return gs.state
}

func (gs *GameSession) PushRune(r rune) {
	if len(gs.Grid[gs.curIdx]) == gs.WordLen { // bounds checking
		return
	}
	cell := Cell{
		Char:  unicode.ToUpper(r),
		state: DEFAULT,
	}
	gs.Grid[gs.curIdx] = append(gs.Grid[gs.curIdx], cell)

	gs.HelpText = ""
}

func (gs *GameSession) PopRune() {
	if len(gs.Grid[gs.curIdx]) == 0 { // bounds checking
		return
	}
	gs.Grid[gs.curIdx] = gs.Grid[gs.curIdx][:len(gs.Grid[gs.curIdx])-1]

	gs.HelpText = ""
}

func (gs *GameSession) ClearCurrentGuess() {
	gs.Grid[gs.curIdx] = nil
	gs.HelpText = ""
}

func (gs *GameSession) UpdateGamestate() {
	gs.HelpText = ""

	if !gs.isValidWord() {
		if len(gs.curGuessAsLowerString()) == gs.WordLen {
			gs.HelpText = fmt.Sprintf("%s not in word list", gs.curGuessAsUpperString())
		}
		return
	}

	if gs.IsWinner() {
		gs.state = VICTORY
		gs.HelpText = fmt.Sprintf("%s is correct! [c]ontinue | go b[a]ck", gs.curGuessAsUpperString())
	}

	gs.finalizeCurRow()

	if gs.curIdx == gs.NumGuesses {
		gs.state = LOSS
		gs.HelpText = fmt.Sprintf("%s was the word! [c]ontinue | go b[a]ck", gs.targetWordAsString)
	}
}

func (gs *GameSession) curGuessAsLowerString() string {
	var word string
	for _, cell := range gs.Grid[gs.curIdx] {
		word += utils.RuneToAlpha(unicode.ToLower(cell.Char))
	}
	return word
}

func (gs *GameSession) curGuessAsUpperString() string {
	var word string
	for _, cell := range gs.Grid[gs.curIdx] {
		word += utils.RuneToAlpha(unicode.ToUpper(cell.Char))
	}
	return word
}

func (gs *GameSession) isValidWord() bool {
	return slices.Contains(gs.validWords, gs.curGuessAsLowerString())
}

func (gs *GameSession) finalizeCurRow() {
	for i := range gs.Grid[gs.curIdx] {
		cell := &gs.Grid[gs.curIdx][i]
		idx := utils.Find[rune](AllRunes, cell.Char) // finding the location in the char tracker
		if cell.Char == gs.targetWordAsRunes[i] {
			cell.SetState(CORRECT)
			gs.SeenChars[idx].SetState(CORRECT)
		} else if slices.Contains(gs.targetWordAsRunes, cell.Char) {
			cell.SetState(PARTIAL)
			gs.SeenChars[idx].SetState(PARTIAL)
		} else {
			cell.SetState(USED)
			gs.SeenChars[idx].SetState(USED)
		}
	}
	gs.curIdx += 1
}

func (gs *GameSession) IsWinner() bool {
	if len(gs.targetWordAsRunes) != len(gs.Grid[gs.curIdx]) {
		panic("len of word and guess do not match")
	}
	for i := range gs.targetWordAsRunes {
		if gs.targetWordAsRunes[i] != gs.Grid[gs.curIdx][i].Char {
			return false
		}
	}
	return true
}

func (gs *GameSession) IsGameOver() bool {
	return gs.state != ACTIVE
}

func (gs *GameSession) Reset() {
	gs.curIdx = 0
	gs.state = ACTIVE
	for i := range gs.Grid {
		gs.Grid[i] = nil
	}
	for i := range gs.SeenChars {
		gs.SeenChars[i].SetState(DEFAULT)
	}

	word := gs.validWords[rand.Intn(len(gs.validWords))]
	gs.targetWordAsRunes = utils.RuneSliceToUpper([]rune(word))
	gs.HelpText = ""
}

func (gs *GameSession) HandleEventKey(ev *tcell.EventKey) bool {
	if gs.state == ACTIVE {
		if shouldExit := gs.activeEventKey(ev); shouldExit {
			return true
		}
	} else {
		if shouldExit := gs.gameOverEventKey(ev); shouldExit {
			return true
		}
	}
	return false
}

func (gs *GameSession) gameOverEventKey(ev *tcell.EventKey) bool {
	if ev.Rune() == 'c' || ev.Rune() == 'C' {
		gs.Reset()
		return false
	} else if ev.Rune() == 'a' || ev.Rune() == 'A' {
		return true
	}
	// fallthrough. nothing happens
	return false
}

func (gs *GameSession) activeEventKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyCtrlC {
		return true
	} else if ev.Key() == tcell.KeyEscape {
		gs.ClearCurrentGuess()
	} else if utils.RuneIsAlpha(ev.Rune()) {
		gs.PushRune(ev.Rune())
	} else if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
		gs.PopRune()
	} else if ev.Key() == tcell.KeyEnter {
		gs.UpdateGamestate()
	}
	// fallthrough. nothing happens
	return false
}
