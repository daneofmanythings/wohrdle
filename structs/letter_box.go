package structs

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/zahra_bday/helpers/drawing"
)

type Box interface {
	Draw()
}

type LetterBox struct {
	letter         string
	x1, y1, x2, y2 int
	screen         tcell.Screen
	style          tcell.Style
}

func NewLetterBox(x, y int, s tcell.Screen, style tcell.Style) LetterBox {
	return LetterBox{
		letter: "",
		x1:     x,
		y1:     y,
		x2:     x + 2,
		y2:     y + 2,
		screen: s,
		style:  style,
	}
}

func (lb *LetterBox) Draw() {
	drawing.DrawBox(lb.screen, lb.x1, lb.y1, lb.x2, lb.y2, lb.style, lb.letter)
}

func (lb *LetterBox) updateLetter(s string) error {
	// if len(s) > 1 || !strings.Contains(validLetters, s) {
	// 	// TODO: write this error for safety
	// 	return nil
	// }
	lb.letter = strings.ToUpper(s)
	return nil
}

type WordBox struct {
	boxes  []LetterBox
	curBox int
	x1, y1 int
	screen tcell.Screen
	style  tcell.Style
}

func NewWordBox(mL, x, y int, s tcell.Screen, style tcell.Style) WordBox {
	lr := WordBox{
		x1:     x,
		y1:     y,
		screen: s,
		style:  style,
		boxes:  make([]LetterBox, mL),
		curBox: 0,
	}
	for i := range lr.boxes {
		lr.boxes[i] = NewLetterBox(x+(3*i+2), y, s, style)
	}
	return lr
}

func (lr *WordBox) Draw() {
	for _, lb := range lr.boxes {
		lb.Draw()
	}
}

func (lr *WordBox) pushLetter(s string) error {
	if lr.curBox == lr.Len() { // bounds checking the word
		return nil
	}
	err := lr.boxes[lr.curBox].updateLetter(s)
	lr.curBox += 1
	return err
}

func (lr *WordBox) popLetter() error {
	if lr.curBox == 0 { // bounds checking the word
		return nil
	}
	lr.curBox -= 1
	return lr.boxes[lr.curBox].updateLetter("")
}

func (lr *WordBox) filled() bool {
	return lr.curBox == lr.Len()
}

func (lr *WordBox) empty() bool {
	return lr.curBox == 0
}

func (lr *WordBox) Len() int {
	return len(lr.boxes)
}

type LetterRows struct {
	word   string
	screen tcell.Screen
	style  tcell.Style
	grid   []WordBox
	// maxWords, maxLetters int
	curWord        int
	x1, y1, x2, y2 int
}

func NewLetterBoxGrid(word string, y, x, mL, mW int, s tcell.Screen, style tcell.Style) LetterRows {
	lbg := LetterRows{
		word:    word,
		x1:      x,
		y1:      y,
		x2:      x + (3*mL + 3),
		y2:      y + (3*mW + 1),
		curWord: 0,
		screen:  s,
		style:   style,
		grid:    make([]WordBox, mW),
	}

	for y := range lbg.grid {
		lbg.grid[y] = NewWordBox(mL, lbg.x1, lbg.y1+3*y+1, s, style)
	}

	return lbg
}

func (lbg *LetterRows) Draw() {
	drawing.DrawBox(lbg.screen, lbg.x1, lbg.y1, lbg.x2, lbg.y2, lbg.style, "")
	for _, lr := range lbg.grid {
		lr.Draw()
	}
}

func (g *LetterRows) PushLetter(s string) {
	g.grid[g.curWord].pushLetter(s)
}

func (g *LetterRows) PopLetter() {
	g.grid[g.curWord].popLetter()
}

// This is going to do a lot of work. needs to end game on loss too probably
func (g *LetterRows) ValidateWord() {
	if g.grid[g.curWord].filled() {
		// TODO: add word validation. only english words are allowed.

		// check if word is in the dictionary

		// if yes: send the row off to be processed
		//
		//
		g.curWord += 1
	}
}
