package states

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/wohrdle/utils"
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

func (c *Cell) isEqualTo(other Cell) bool {
	return c.Char == other.Char && c.state == other.state
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

var gameStates = []GameState{ACTIVE, VICTORY, LOSS}

type GameSession struct {
	Parameters Parameters

	WordLen     int
	NumGuesses  int
	MaxNumFails int
	HardMode    int

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
		Parameters:  *params,
		WordLen:     params.Fields[0].Value,
		NumGuesses:  params.Fields[1].Value,
		MaxNumFails: params.Fields[2].Value,
		HardMode:    params.Fields[3].Value,
		curIdx:      0,
		validWords:  params.ValidWords(),
		state:       ACTIVE,
		SeenChars:   NewSeenCharRecord(),
	}
	word := gs.validWords[rand.Intn(len(gs.validWords))]
	// word := "volts"
	gs.targetWordAsString = strings.ToUpper(word)
	gs.targetWordAsRunes = utils.RuneSliceToUpper([]rune(word))

	gs.Grid = make([][]Cell, gs.NumGuesses)
	for i := range gs.Grid {
		gs.Grid[i] = []Cell{}
	}

	return gs
}

func (gs *GameSession) setState(state GameState) {
	if !slices.Contains(gameStates, state) {
		return
	}
	gs.state = state
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

// helper function for debugging
func (gs *GameSession) getCurrentRow() []Cell {
	return gs.Grid[gs.curIdx]
}

func (gs *GameSession) ClearCurrentGuess() {
	gs.Grid[gs.curIdx] = nil
	gs.HelpText = ""
}

func (gs *GameSession) UpdateGamestate() {
	gs.HelpText = ""
	failed_entry_loss := "Out of failed entries. %s was the word! [c]ontinue | go b[a]ck"
	failed_entry := "%s not in word list. %d failed entries left"
	victory := "%s is correct! [c]ontinue | go b[a]ck"
	guess_loss := "%s was the word! [c]ontinue | go b[a]ck"
	gave_up_loss := "Aborted. [c]ontinue | go b[a]ck"
	hardmode_violated := "Hard-mode violated. %d failed entries left"

	if gs.GetState() == LOSS {
		gs.HelpText = fmt.Sprint(gave_up_loss)
	}

	if !gs.isValidWord() {
		if len(gs.curGuessAsLowerString()) == gs.WordLen {
			gs.MaxNumFails -= 1
			if gs.MaxNumFails == 0 {
				gs.setState(LOSS)
				gs.HelpText = fmt.Sprintf(failed_entry_loss, gs.targetWordAsString)
			} else {
				gs.HelpText = fmt.Sprintf(failed_entry, gs.curGuessAsUpperString(), gs.MaxNumFails)
			}
		}
		return
	}

	if gs.HardMode == 1 {
		if !gs.isHardModeSatisfied() {
			gs.MaxNumFails -= 1
			if gs.MaxNumFails == 0 {
				gs.setState(LOSS)
				gs.HelpText = fmt.Sprintf(failed_entry_loss, gs.targetWordAsString)
			} else {
				gs.HelpText = fmt.Sprintf(hardmode_violated, gs.MaxNumFails)
			}
			return
		}
	}

	if gs.IsWinner() {
		gs.setState(VICTORY)
		gs.HelpText = fmt.Sprintf(victory, gs.curGuessAsUpperString())
	}

	gs.finalizeCurRow()

	if gs.curIdx == gs.NumGuesses {
		gs.setState(LOSS)
		gs.HelpText = fmt.Sprintf(guess_loss, gs.targetWordAsString)
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

func (gs *GameSession) isHardModeSatisfied() bool {
	// Can't fail on the first guess
	if gs.curIdx == 0 {
		return true
	}

	// Need this to detect missed PARTIALS
	countByRune := gs.countMapForCurrRow()

	// Making things easier to reason about in the code
	prevRow := &gs.Grid[gs.curIdx-1]
	currRow := &gs.Grid[gs.curIdx]
	// First pass to see if any previously correct are missing and to update the countMap
	// for the second pass
	for i := range *prevRow {
		// Making things easier to reason about in the code
		prevRowCell := (*prevRow)[i]
		currRowCell := (*currRow)[i]
		if prevRowCell.GetState() != CORRECT {
			continue
		}
		// Since the cell is correct, the chars should match
		if prevRowCell.Char != currRowCell.Char {
			return false
		}
		// they matched, so decrement the countMap
		countByRune[currRowCell.Char] -= 1
	}

	// Second pass to catch any missing PARTIALS. looking at the cells of the previous row
	// in relation to how many are left in the countMap of the current row
	for _, cell := range *prevRow {
		// dont care if it isnt a PARTIAL
		if cell.GetState() != PARTIAL {
			continue
		}
		if countByRune[cell.Char] < 1 {
			// We found a partial that isnt represented in the current row.
			// IT HAS TO BE REPRESENTED
			return false
		}
		// it is represented, so we decrement the count for that PARTIAL
		countByRune[cell.Char] -= 1
	}

	return true
}

func (gs *GameSession) isValidWord() bool {
	return slices.Contains(gs.validWords, gs.curGuessAsLowerString())
}

func (gs *GameSession) finalizeCurRow() {
	// This populates the cells in the current row with thier correct stylings for the renderer

	countByRune := gs.countMapForTargetWord() // This is to track repeat letters from ISSUE#1
	// First pass
	for i := range gs.Grid[gs.curIdx] {
		cell := &gs.Grid[gs.curIdx][i]
		idx := utils.Find[rune](AllRunes, cell.Char) // finding the location in the seen char tracker
		if cell.Char == gs.targetWordAsRunes[i] {
			countByRune[cell.Char] -= 1
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
	// Second pass to remove potential false positives of PARTIALS when they have
	// all been correctly guessed by comparing the remaining unfound CORRECTS
	for i := range gs.Grid[gs.curIdx] {
		cell := &gs.Grid[gs.curIdx][i]
		idx := utils.Find[rune](AllRunes, cell.Char) // finding the location in the seen char tracker
		if cell.GetState() != PARTIAL {
			continue
		}
		// We know it is a PARTIAL
		if countByRune[cell.Char] < 1 {
			cell.SetState(USED)
			gs.SeenChars[idx].SetState(CORRECT)
		}
		countByRune[cell.Char] -= 1
	}
	gs.curIdx += 1
}

func (gs *GameSession) countMapForTargetWord() map[rune]int {
	countByRune := map[rune]int{}
	for i := range gs.targetWordAsRunes {
		countByRune[gs.targetWordAsRunes[i]] += 1
	}
	return countByRune
}

func (gs *GameSession) countMapForCurrRow() map[rune]int {
	countByRune := map[rune]int{}
	for i := range gs.curGuessAsUpperString() {
		countByRune[rune(gs.curGuessAsUpperString()[i])] += 1
	}
	return countByRune
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

func (gs *GameSession) Reset() {
	gs.curIdx = 0
	gs.setState(ACTIVE)
	gs.MaxNumFails = gs.Parameters.Fields[2].Value // NOTE: update this if the menu position of max fails changes
	for i := range gs.Grid {
		gs.Grid[i] = nil
	}
	for i := range gs.SeenChars {
		gs.SeenChars[i].SetState(DEFAULT)
	}

	word := gs.validWords[rand.Intn(len(gs.validWords))]
	gs.targetWordAsString = strings.ToUpper(word)
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
	} else if ev.Rune() == 'a' || ev.Rune() == 'A' || ev.Key() == tcell.KeyCtrlC {
		return true
	}
	// fallthrough. nothing happens
	return false
}

func (gs *GameSession) activeEventKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyCtrlC {
		gs.state = LOSS
		gs.UpdateGamestate()
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
