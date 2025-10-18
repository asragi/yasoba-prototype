//go:build wireinject
// +build wireinject

package app

import "github.com/google/wire"

func initializeApp(cfg Config) (*App, error) {
	wire.Build(
		provideDrawing,
		provideBattleSceneResult,
		provideDebug,
		buildApp,
	)
	return nil, nil
}
