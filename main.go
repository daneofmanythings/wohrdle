package main

import (
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
	defer screen.Fini()

	renderer := render.NewRenderer()

	wordRepo, err := utils.LoadWordRepoFromJSON(wordRepoPath)
	if err != nil {
		panic(err)
	}
	parameters := states.NewParameters(5, 6, wordRepo.Words)

	runMainMenu(parameters, renderer, screen)

	gs := states.NewGameSession(parameters)

	// the application loop
	for {
		if shouldRunMenu := runGameSession(gs, renderer, screen); shouldRunMenu {
			if shouldQuit := runMainMenu(parameters, renderer, screen); shouldQuit {
				break
			}
			gs = states.NewGameSession(parameters)
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
		default:
			// nothing
		}
	}
}

func runMainMenu(p *states.Parameters, r *render.Renderer, s tcell.Screen) bool {
	for {
		// the menu loop
		r.DrawMenu(s, p)

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if shouldReturn := p.HandleEventKey(ev, s); shouldReturn {
				return true
			}
		default:
			// nothing
		}
	}
}
