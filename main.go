package main

import (
	"log"

	"github.com/asragi/yasoba-prototype/app"
	"github.com/asragi/yasoba-prototype/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const DrawRate = 1

func main() {
	cfg := app.DefaultConfig()
	application, err := app.InitializeApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(cfg.GameWidth*DrawRate, cfg.GameHeight*DrawRate)
	ebiten.SetWindowTitle("yasoba-prototype")

	g := game.New(application, cfg)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
