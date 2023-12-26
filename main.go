package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/zahra_bday/render"
	"gitlab.com/daneofmanythings/zahra_bday/states"
	"gitlab.com/daneofmanythings/zahra_bday/utils"
)

func main() {
	s, err := render.CreateScreen()
	if err != nil {
		panic(err)
	}
	renderer := render.NewRenderer()
	parameters := states.Parameters{5, 6}
	renderAlert := make(chan bool)
	gs := states.NewGameSession(renderAlert, &parameters)

	// gs.PushRune('A')
	// renderer.DrawGameSession(s, gs)

	go func() {
		for {
			// Process event
			switch ev := s.PollEvent().(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					s.Fini()
					os.Exit(0)
				} else if utils.RuneIsAlpha(ev.Rune()) {
					gs.PushRune(ev.Rune())
				} else if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
					gs.PopRune()
				} else if ev.Key() == tcell.KeyEnter {
					if !gs.IsValidWord() {
						continue
					}
					// if gs.CurWordGuess() == gs.Word {
					// 	// win condition. make menu to play again or go back
					// 	gs.nextGuessPrep() // placeholder action
					// } else {
					// 	gs.nextGuessPrep()
					// }
				}
			}
		}
	}()

	for {
		renderer.DrawGameSession(s, gs)
		<-renderAlert
	}
}
