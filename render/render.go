package render

import (
	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/worhdle/states"
)

type Renderer struct {
	xSpacing int
	ySpacing int
}

func NewRenderer() *Renderer {
	return &Renderer{
		xSpacing: 4,
		ySpacing: 2,
	}
}

func (r *Renderer) DrawGameSession(s tcell.Screen, gs *states.GameSession) {
	s.Clear()
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	// colorStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
	defer s.Show()

	width, height := s.Size()

	x1 := (width - (gs.WordLen+2)*r.xSpacing) / 2
	y1 := (height - (gs.NumGuesses+1)*r.ySpacing) / 2
	x2 := x1 + r.xSpacing*gs.WordLen
	y2 := y1 + r.ySpacing*gs.NumGuesses

	// draw horizontal ticks
	for y := 0; y <= gs.NumGuesses; y++ {
		for x := 0; x < gs.WordLen; x++ {
			for n := 1; n < r.xSpacing; n++ {
				s.SetContent(x*r.xSpacing+x1+n, y*r.ySpacing+y1, tcell.RuneHLine, nil, style)
			}
		}
	}

	// draw vertical ticks
	for x := 0; x <= gs.WordLen; x++ {
		for y := 0; y < gs.NumGuesses; y++ {
			s.SetContent(x*r.xSpacing+x1, y*r.ySpacing+y1+1, tcell.RuneVLine, nil, style)
		}
	}

	// draw corners
	s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
	s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
	s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
	s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)

	// draw tees
	// top
	for i := 1; i < gs.WordLen; i++ {
		s.SetContent(i*r.xSpacing+x1, y1, tcell.RuneTTee, nil, style)
	}
	// bottom
	for i := 1; i < gs.WordLen; i++ {
		s.SetContent(i*r.xSpacing+x1, y2, tcell.RuneBTee, nil, style)
	}
	// left
	for i := 1; i < gs.NumGuesses; i++ {
		s.SetContent(x1, i*r.ySpacing+y1, tcell.RuneLTee, nil, style)
	}
	// Right
	for i := 1; i < gs.NumGuesses; i++ {
		s.SetContent(x2, i*r.ySpacing+y1, tcell.RuneRTee, nil, style)
	}

	// fill middle with pluses
	for j := 1; j < gs.NumGuesses; j++ {
		for i := 1; i < gs.WordLen; i++ {
			s.SetContent(i*r.xSpacing+x1, j*r.ySpacing+y1, tcell.RunePlus, nil, style)
		}
	}

	// draw cell characters for session grid
	for j, row := range gs.Grid {
		for i, cell := range row {
			drawCellChar(&cell, x1+i*r.xSpacing+r.xSpacing/2, y1+j*r.ySpacing+r.ySpacing/2, s)
		}
	}

	helpMessageX := (width - len(gs.HelpText)) / 2 // centering text
	drawHelpMessage(helpMessageX, y2+r.ySpacing, s, gs)

	drawSeenChars(x2+r.xSpacing, y1+r.ySpacing, x2+r.xSpacing+3, s, gs)
}

func drawCellChar(cell *states.Cell, x, y int, s tcell.Screen) {
	var letterStyle tcell.Style
	switch cell.GetState() {
	case states.DEFAULT:
		letterStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	case states.CORRECT:
		letterStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
	case states.PARTIAL:
		letterStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorYellow)
	}

	s.SetContent(x, y, cell.Char, nil, letterStyle)
}

func drawHelpMessage(x, y int, s tcell.Screen, gs *states.GameSession) {
	var style tcell.Style
	switch gs.GetState() {
	case states.ACTIVE:
		style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorYellow)
	case states.VICTORY:
		style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
	case states.LOSS:
		style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed)
	}
	drawTextWrapping(s, x, y, x+len(gs.HelpText), style, gs.HelpText)
}

func drawSeenChars(x, y, x2 int, s tcell.Screen, gs *states.GameSession) {
	row := y
	col := x
	var style tcell.Style
	for _, r := range gs.SeenChars {
		switch r.GetState() {
		case states.CORRECT:
			style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
		case states.DEFAULT:
			style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		case states.PARTIAL:
			style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorYellow)
		case states.USED:
			// TODO: fix this to properly reverse only the foreground color. hard coded atm
			style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorBlack)
		}
		s.SetContent(col, row, r.Char, nil, style)
		col++
		if col >= x2 {
			row++
			col = x
		}
	}
}

func drawTextWrapping(s tcell.Screen, x1, y1, x2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
	}
}

func CreateScreen() (tcell.Screen, error) {
	screen, creationErr := tcell.NewScreen()
	if creationErr != nil {
		return nil, creationErr
	}
	if initErr := screen.Init(); initErr != nil {
		return nil, initErr
	}
	screen.DisableMouse()
	screen.DisablePaste()

	return screen, nil
}
