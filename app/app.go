package app

import (
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/scene"
)

type App struct {
	Drawing     *frontend.Drawing
	BattleScene *scene.BattleScene
	Debug       *debug.Debug
}

func InitializeApp(cfg Config) (*App, error) {
	return initializeApp(cfg)
}

func buildApp(drawing *frontend.Drawing, battleScene *scene.BattleScene, debug *debug.Debug) *App {
	return &App{
		Drawing:     drawing,
		BattleScene: battleScene,
		Debug:       debug,
	}
}
