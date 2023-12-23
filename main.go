package main

import (

	// "gitlab.com/daneofmanythings/zahra_bday/example"

	"log"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/zahra_bday/helpers"
	"gitlab.com/daneofmanythings/zahra_bday/structs"
)

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	// boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnablePaste()
	s.Clear()

	grid := structs.NewLetterBoxGrid("a", 2, 15, 7, 6, s, defStyle)
	grid.Draw()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Here's how to get the screen size when you need it.
	// xmax, ymax := s.Size()

	// Here's an example of how to inject a keystroke where it will
	// be picked up by the next PollEvent call.  Note that the
	// queue is LIFO, it has a limited length, and PostEvent() can
	// return an error.
	// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

	// Event loop
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if helpers.RuneIsAlpha(ev.Rune()) {
				grid.PushLetter(helpers.RuneToAlpha(ev.Rune()))
				grid.Draw()
			} else if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
				grid.PopLetter()
				grid.Draw()
			} else if ev.Key() == tcell.KeyEnter {
				grid.ValidateWord()
			}
		}
	}
}
