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
	screen, err := render.CreateScreen()
	if err != nil {
		panic(err)
	}
	renderer := render.NewRenderer()
	wordRepo, err := utils.LoadWordRepoFromJSON(wordRepoPath)
	if err != nil {
		panic(err)
	}
	parameters := states.NewParameters(5, 6, wordRepo.Words)

	parameters = *runMainMenu(renderer, screen, &parameters)

	gs := states.NewGameSession(&parameters)

	// the application loop
	for {
		if shouldRunMenu := runGameSession(gs, renderer, screen); shouldRunMenu {
			parameters = *runMainMenu(renderer, screen, &parameters)
			gs = states.NewGameSession(&parameters)
		}
	}
}

func runGameSession(gs *states.GameSession, r *render.Renderer, s tcell.Screen) bool {
	for {
		// the game loop
		r.DrawGameSession(s, gs)
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if shouldExit := gs.HandleEventKey(ev); shouldExit {
				return true
			}
		}
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
		// the menu loop
		r.DrawMenu(s, p)

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			continue
		case *tcell.EventKey:
			// p.HandleKeyPress(ev)
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
