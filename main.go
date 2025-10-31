package main

import (
	"log"

	"github.com/asragi/yasoba-prototype/app"
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/hajimehoshi/ebiten/v2"
)

const DrawRate = 1

type Game struct {
	activeScene  scene.Scene
	battleScene  *scene.BattleScene
	debugOverlay *debug.Debug
	drawing      *drawing.Drawing
	width        int
	height       int
}

func NewGame(app *app.App, cfg app.Config) *Game {
	game := &Game{
		activeScene:  app.DebugScene,
		battleScene:  app.BattleScene,
		debugOverlay: app.DebugOverlay,
		drawing:      app.Drawing,
		width:        cfg.GameWidth,
		height:       cfg.GameHeight,
	}
	app.DebugScene.SetOnSelectBattle(game.activateBattleScene)
	return game
}

func (g *Game) Update() error {
	g.activeScene.Update()
	g.debugOverlay.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.activeScene.Draw(g.drawing.Draw)
	g.debugOverlay.Draw(g.drawing.Draw)
	g.drawing.DrawEnd(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func (g *Game) activateBattleScene() {
	g.activeScene = g.battleScene
}

func main() {
	cfg := app.DefaultConfig()
	application, err := app.InitializeApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(cfg.GameWidth*DrawRate, cfg.GameHeight*DrawRate)
	ebiten.SetWindowTitle("yasoba-prototype")

	game := NewGame(application, cfg)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
