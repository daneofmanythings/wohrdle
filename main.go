package main

import (
	_ "embed"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/daneofmanythings/wohrdle/render"
	"gitlab.com/daneofmanythings/wohrdle/states"
	"gitlab.com/daneofmanythings/wohrdle/static"
	"gitlab.com/daneofmanythings/wohrdle/utils"
)

const wordRepoPath string = "./static/words.json"

func main() {
	screen, err := render.CreateScreen()
	if err != nil {
		panic(err)
	}
	defer screen.Fini()

	renderer := render.NewRenderer()

	// wordRepo, err := utils.LoadWordRepoFromJSON(wordRepoPath)
	wordRepo, err := utils.LoadEmbeddedWordRepo(static.WordRepoBytes)
	if err != nil {
		panic(err)
	}
	parameters := states.NewDefaultParameters(wordRepo.Words)

	// the application loop
	for {
		runMainMenu(parameters, renderer, screen)
		gs := states.NewGameSession(parameters)
		for {
			if shouldRunMenu := runGameSession(gs, renderer, screen); shouldRunMenu {
				break
			}
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

func runMainMenu(p *states.Parameters, r *render.Renderer, s tcell.Screen) {
	for {
		// the menu loop
		r.DrawMenu(s, p)

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if shouldReturn := p.HandleEventKey(ev, s); shouldReturn {
				return
			}
		default:
			// nothing
		}
	}
}
