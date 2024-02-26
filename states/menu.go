package states

import (
	"os"
	"slices"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

const (
	MAX_GUESSES int = 20
	MAX_FAILS   int = 20

	TRUE  int = 1
	FALSE int = 0
)

type Field struct {
	Name  string
	Value int
}

var defaultFields []Field = []Field{
	{"word length", 5},
	{"num guesses", 6},
	{"num failed words", 5},
	{"hard-mode", 0},
}

type Parameters struct {
	// Field[0] >> word length
	// Field[1] >> number of guesses
	// Field[2] >> failed word attempts
	Fields        []Field
	CurEditingIdx int

	WordRepo   map[string][]string
	MinWordLen int
	MaxWordLen int
}

func NewDefaultParameters(wordRepo map[string][]string) *Parameters {
	// finding the bounds of wordLen for menu wrapping
	word_lengths := []int{}
	for str_len := range wordRepo {
		word_len, err := strconv.Atoi(str_len)
		if err != nil {
			panic("WHY ARE WE PANICKING HERE. SOMETHING HAS GONE TERRIBLY WRONG")
		}
		word_lengths = append(word_lengths, word_len)
	}

	return &Parameters{
		Fields:        defaultFields,
		CurEditingIdx: 0,
		WordRepo:      wordRepo,
		MinWordLen:    slices.Min(word_lengths),
		MaxWordLen:    slices.Max(word_lengths),
	}
}

func (p *Parameters) ValidWords() []string {
	return p.WordRepo[strconv.Itoa(p.Fields[0].Value)]
}

func (p *Parameters) IncCurField() {
	p.CurEditingIdx -= 1
	// modulus in go doesnt wrap negatives correctly
	if p.CurEditingIdx < 0 {
		p.CurEditingIdx = len(p.Fields) - 1
	}
}

// NOTE: This must be updated when menu items are added
func (p *Parameters) IncValAtCurField() {
	switch p.CurEditingIdx {
	case 0: // word length
		val := &p.Fields[0].Value
		if *val == p.MaxWordLen {
			*val = 1
		} else {
			*val += 1
		}
	case 1: // number of guesses
		val := &p.Fields[1].Value
		if *val == MAX_GUESSES {
			*val = 1
		} else {
			*val += 1
		}
	case 2: // number of failed words
		val := &p.Fields[2].Value
		if *val == MAX_FAILS {
			*val = 1
		} else {
			*val += 1
		}
	case 3: // hard-mode flag
		val := &p.Fields[3].Value
		if *val == FALSE {
			*val = TRUE
		} else {
			*val = FALSE
		}
	}
}

func (p *Parameters) DecCurField() {
	p.CurEditingIdx += 1
	p.CurEditingIdx %= len(p.Fields)
}

// NOTE: This must be updated when a menu item is added
func (p *Parameters) DecValAtCorField() {
	switch p.CurEditingIdx {
	case 0: // word length
		val := &p.Fields[0].Value
		if *val == p.MinWordLen {
			*val = p.MaxWordLen
		} else {
			*val -= 1
		}
	case 1: // number of guesses
		val := &p.Fields[1].Value
		if *val == 1 {
			*val = MAX_GUESSES
		} else {
			*val -= 1
		}
	case 2: // number of failed words
		val := &p.Fields[2].Value
		if *val == 1 {
			*val = MAX_FAILS
		} else {
			*val -= 1
		}
	case 3: // hard-mode flag
		val := &p.Fields[3].Value
		if *val == FALSE {
			*val = TRUE
		} else {
			*val = FALSE
		}
	}
}

var (
	upBinds    []rune = []rune{'k', 'K', 'w', 'W'}
	downBinds  []rune = []rune{'j', 'J', 's', 'S'}
	leftBinds  []rune = []rune{'h', 'H', 'a', 'A'}
	rightBinds []rune = []rune{'l', 'L', 'd', 'D'}
)

func (p *Parameters) HandleEventKey(ev *tcell.EventKey, s tcell.Screen) bool {
	if ev.Key() == tcell.KeyUp || slices.Contains(upBinds, ev.Rune()) {
		p.IncCurField()
	} else if ev.Key() == tcell.KeyDown || slices.Contains(downBinds, ev.Rune()) {
		p.DecCurField()
	} else if ev.Key() == tcell.KeyLeft || slices.Contains(leftBinds, ev.Rune()) {
		p.DecValAtCorField()
	} else if ev.Key() == tcell.KeyRight || slices.Contains(rightBinds, ev.Rune()) {
		p.IncValAtCurField()
	} else if ev.Key() == tcell.KeyEnter {
		return true
	} else if ev.Key() == tcell.KeyCtrlC {
		s.Fini()
		os.Exit(0)
	}
	return false
}
