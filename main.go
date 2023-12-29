package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/worhdle/render"
	"gitlab.com/daneofmanythings/worhdle/states"
	"gitlab.com/daneofmanythings/worhdle/utils"
)

const wordRepoPath string = "./static/words.json"

func main() {
	s, err := render.CreateScreen()
	if err != nil {
		panic(err)
	}
	renderer := render.NewRenderer()
	wordRepo, _ := utils.LoadWordRepoFromJSON(wordRepoPath)
	parameters := states.NewParameters(5, 6, wordRepo.Words)
	renderAlert := make(chan bool)

	gs := states.NewGameSession(renderAlert, &parameters)
	runGameSession(gs, renderer, s, &renderAlert)
}

func runGameSession(gs *states.GameSession, r *render.Renderer, s tcell.Screen, renderAlert *chan bool) bool {
	go func() {
		for {
			// The regular game loop
			switch ev := s.PollEvent().(type) {
			case *tcell.EventResize:
				s.Sync()
				*renderAlert <- true
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC {
					s.Fini()
					os.Exit(0)
				} else if (ev.Key() == 'c' || ev.Key() == 'C') && gs.IsGameOver() {
					gs = states.NewGameSession(*renderAlert, &gs.Parameters)
					continue
				} else if ev.Key() == tcell.KeyEscape {
					gs.ClearCurrentGuess()
				} else if utils.RuneIsAlpha(ev.Rune()) {
					// Logic to continue game or rerun the menu
					if gs.IsGameOver() && (ev.Rune() == 'c' || ev.Rune() == 'C') {
						gs.Reset()
						*renderAlert <- true
						continue
					} else if gs.IsGameOver() && (ev.Rune() == 'a' || ev.Rune() == 'A') {
						// TODO: add logic to rerun main menu and reset game
						s.Fini()
						os.Exit(0)
					}
					gs.PushRune(ev.Rune())
				} else if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
					gs.PopRune()
				} else if ev.Key() == tcell.KeyEnter {
					gs.UpdateGamestate()
					*renderAlert <- true
				}
			}
		}
	}()

	// The draw loop. blocks until a render alert is sent
	for {
		r.DrawGameSession(s, gs)
		<-*renderAlert
	}
}

// func runMainMenu() *states.Parameters {
//
// }
