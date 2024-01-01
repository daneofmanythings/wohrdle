package states

import (
	"math"
	"os"
	"slices"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

type MenuSession struct{}

type Parameters struct {
	WordLen    int
	NumGuesses int
	WordRepo   map[string][]string

	CurField   int
	MinWordLen int
	MaxWordLen int
}

func NewParameters(wordLen int, numGuesses int, wordRepo map[string][]string) *Parameters {
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
		WordLen:    wordLen,
		NumGuesses: numGuesses,
		WordRepo:   wordRepo,
		CurField:   0,
		MinWordLen: slices.Min(word_lengths),
		MaxWordLen: slices.Max(word_lengths),
	}
}

func (p *Parameters) ValidWords() []string {
	return p.WordRepo[strconv.Itoa(p.WordLen)]
}

// TODO: fix this abstraction. It will leak everywhere
func (p *Parameters) IncCurField() {
	p.CurField += 1
	p.CurField %= 2 // change this to the number of fields to edit
}

func (p *Parameters) IncValAtCurField() {
	if p.CurField == 0 {
		if p.WordLen == p.MaxWordLen {
			p.WordLen = 1
		} else {
			p.WordLen += 1
		}
	} else {
		if p.NumGuesses == 20 {
			return
		}
		p.NumGuesses += 1
	}
}

func (p *Parameters) DecCurField() {
	p.CurField -= 1
	p.CurField = int(math.Abs(float64(p.CurField)))
	p.CurField %= 2 // change this to the number of fields to edit
}

func (p *Parameters) DecValAtCorField() {
	if p.CurField == 0 {
		if p.WordLen == p.MinWordLen {
			p.WordLen = p.MaxWordLen
		} else {
			p.WordLen -= 1
		}
	} else {
		if p.NumGuesses == 1 {
			return
		}
		p.NumGuesses -= 1
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
