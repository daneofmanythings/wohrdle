package main

import (
	"os"
	"slices"

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

	parameters = *runMainMenu(renderer, s, &parameters)

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
					gs = states.NewGameSession(*renderAlert, runMainMenu(r, s, &gs.Parameters))
					*renderAlert <- true
				} else if (ev.Key() == 'c' || ev.Key() == 'C') && gs.IsGameOver() {
					gs = states.NewGameSession(*renderAlert, &gs.Parameters)
					continue
				} else if ev.Key() == tcell.KeyEscape {
					gs.ClearCurrentGuess()
				} else if utils.RuneIsAlpha(ev.Rune()) {
					// Logic to continue game or rerun the menu
					if gs.IsGameOver() {
						if ev.Rune() == 'c' || ev.Rune() == 'C' {
							gs.Reset()
							*renderAlert <- true
							continue
						} else if ev.Rune() == 'a' || ev.Rune() == 'A' {
							gs = states.NewGameSession(*renderAlert, runMainMenu(r, s, &gs.Parameters))
							*renderAlert <- true
						}
					} else {
						gs.PushRune(ev.Rune())
					}
				} else if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
					if gs.IsGameOver() {
						continue
					}
					gs.PopRune()
				} else if ev.Key() == tcell.KeyEnter {
					if gs.IsGameOver() {
						continue
					}
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

var (
	upBinds    []rune = []rune{'k', 'K', 'w', 'W'}
	downBinds  []rune = []rune{'j', 'J', 's', 'S'}
	leftBinds  []rune = []rune{'h', 'H', 'a', 'A'}
	rightBinds []rune = []rune{'l', 'L', 'd', 'D'}
)

func runMainMenu(r *render.Renderer, s tcell.Screen, p *states.Parameters) *states.Parameters {
	for {
		r.DrawMenu(s, p)

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			continue
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyUp || slices.Contains(upBinds, ev.Rune()) {
				p.IncCurField()
			} else if ev.Key() == tcell.KeyDown || slices.Contains(downBinds, ev.Rune()) {
				p.DecCurField()
			} else if ev.Key() == tcell.KeyLeft || slices.Contains(leftBinds, ev.Rune()) {
				p.DecValAtCorField()
			} else if ev.Key() == tcell.KeyRight || slices.Contains(rightBinds, ev.Rune()) {
				p.IncValAtCurField()
			} else if ev.Key() == tcell.KeyEnter {
				return p
			} else if ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				os.Exit(0)
			}
		default:
			// nothing to do here
		}
	}
}
