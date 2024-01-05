package render

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/wohrdle/states"
)

type Renderer struct {
	xSpacing int
	ySpacing int
	defStyle tcell.Style
}

func NewRenderer() *Renderer {
	return &Renderer{
		xSpacing: 4,
		ySpacing: 2,
	}
}

func (r *Renderer) DrawMenu(s tcell.Screen, p *states.Parameters) {
	s.Clear()
	defer s.Show()

	style := tcell.StyleDefault
	style_faded := style.Foreground(tcell.ColorGrey)
	width, _ := s.Size()

	welcome := "Welcome to WOHRDLE!"
	instructions := "Please select word length and max guesses."
	welX := startingX(width, welcome)
	insX := startingX(width, instructions)
	drawTextWrapping(s, welX, r.ySpacing, welX+len(welcome), style, welcome)
	drawTextWrapping(s, insX, 2*r.ySpacing, insX+len(instructions), style, instructions)

	starting_dynamic_offset := 4

	// dynamic portion ------
	for i := range p.Fields {
		title := p.Fields[i].Name
		title = title + ": " + strconv.Itoa(p.Fields[i].Value)
		x_start := startingX(width, title)
		x_end := x_start + len(title)
		drawTextWrapping(s, x_start, (starting_dynamic_offset+i)*r.ySpacing, x_end, determineMenuStyle(i, p), title)
	}
	// -------

	help_text_offset := starting_dynamic_offset + len(p.Fields) + 1
	bindsMenu := "Navigate with arrow keys, 'wasd', or 'hjkl'. <return> to start."
	bindsGame := "Type words. <esc> clears whole word. <ctrl-c> to go back."
	menX := startingX(width, bindsMenu)
	gamX := startingX(width, bindsGame)
	drawTextWrapping(s, menX, help_text_offset*r.ySpacing, menX+len(bindsMenu), style_faded, bindsMenu)
	drawTextWrapping(s, gamX, (help_text_offset+1)*r.ySpacing, gamX+len(bindsGame), style_faded, bindsGame)
}

func determineMenuStyle(curDisplayingIdx int, p *states.Parameters) tcell.Style {
	if curDisplayingIdx == p.CurEditingIdx {
		return tcell.StyleDefault.Reverse(true)
	}
	return tcell.StyleDefault
}

func startingX(width int, str string) int {
	return (width - len(str)) / 2
}

func (r *Renderer) DrawGameSession(s tcell.Screen, gs *states.GameSession) {
	s.Clear()
	style := tcell.StyleDefault
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
		letterStyle = tcell.StyleDefault.Bold(true)
	case states.CORRECT:
		letterStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	case states.PARTIAL:
		letterStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorYellow).Bold(true)
	default:
		letterStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
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
	for _, cell := range gs.SeenChars {
		char := cell.Char
		switch cell.GetState() {
		case states.CORRECT:
			style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
		case states.DEFAULT:
			style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		case states.PARTIAL:
			style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorYellow)
		case states.USED:
			char = ' '
		}
		s.SetContent(col, row, char, nil, style)
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
