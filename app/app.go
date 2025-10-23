package app

import (
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/scene"
)

type App struct {
	Drawing      *frontend.Drawing
	BattleScene  *scene.BattleScene
	DebugScene   *scene.DebugScene
	DebugOverlay *debug.Debug
}

func InitializeApp(cfg Config) (*App, error) {
	return initializeApp(cfg)
}

func buildApp(
	drawing *frontend.Drawing,
	battleScene *scene.BattleScene,
	debugScene *scene.DebugScene,
	debugOverlay *debug.Debug,
) *App {
	return &App{
		Drawing:      drawing,
		BattleScene:  battleScene,
		DebugScene:   debugScene,
		DebugOverlay: debugOverlay,
	}
}
