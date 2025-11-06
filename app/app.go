package app

import (
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
)

type App struct {
	Drawing           *drawing.Drawing
	CreateBattleScene scene.CreateBattleScene
	CreateDebugScene  scene.CreateDebugScene
	DebugOverlay      *debug.Debug
	Config            Config
}

func InitializeApp(cfg Config) (*App, error) {
	return initializeApp(cfg)
}

func buildApp(
	drawing *drawing.Drawing,
	createBattleScene scene.CreateBattleScene,
	createDebugScene scene.CreateDebugScene,
	debugOverlay *debug.Debug,
	cfg Config,
) *App {
	return &App{
		Drawing:           drawing,
		CreateBattleScene: createBattleScene,
		CreateDebugScene:  createDebugScene,
		DebugOverlay:      debugOverlay,
		Config:            cfg,
	}
}
