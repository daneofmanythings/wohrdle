package states

import "unicode"

type Parameters struct {
	WordLen    int
	NumGuesses int
}

type CellState int

const (
	DEFAULT CellState = iota
	CORRECT
	PARTIAL
	USED
)

type Cell struct {
	Char  rune
	State CellState
}

type GameState int

const (
	ACTIVE = iota
	VICTORY
	LOSS
)

type GameSession struct {
	AlertChan chan bool

	WordLen    int
	NumGuesses int
	Word       []rune

	Grid      [][]Cell
	curIdx    int
	seenChars map[rune]CellState

	state GameState
}

func NewGameSession(alertchan chan bool, params *Parameters) *GameSession {
	gs := &GameSession{
		AlertChan:  alertchan,
		WordLen:    params.WordLen,
		NumGuesses: params.NumGuesses,
		curIdx:     0,
		state:      ACTIVE,
	}
	// TODO: generate word from parameters
	gs.Word = []rune("hello") // placeholder

	gs.Grid = make([][]Cell, gs.NumGuesses)
	for i := range gs.Grid {
		gs.Grid[i] = []Cell{}
	}

	return gs
}

func (gs *GameSession) PushRune(r rune) {
	if len(gs.Grid[gs.curIdx]) == gs.WordLen { // bounds checking
		return
	}
	cell := Cell{
		Char:  unicode.ToUpper(r),
		State: DEFAULT,
	}
	gs.Grid[gs.curIdx] = append(gs.Grid[gs.curIdx], cell)

	go func() {
		gs.AlertChan <- true
	}()
}

func (gs *GameSession) PopRune() {
	if len(gs.Grid[gs.curIdx]) == 0 { // bounds checking
		return
	}
	gs.Grid[gs.curIdx] = gs.Grid[gs.curIdx][:len(gs.Grid[gs.curIdx])-1]

	go func() {
		gs.AlertChan <- true
	}()
}

func (gs *GameSession) IsValidWord() bool {
	return false
}

// TODO: figure out this comparison
func (gs *GameSession) CurWordGuess() []rune {
	curWordGuess := make([]rune, gs.WordLen)
	for _, cell := range gs.Grid[gs.curIdx] {
		curWordGuess = append(curWordGuess, cell.Char)
	}
	return curWordGuess
}
