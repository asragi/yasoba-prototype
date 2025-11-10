package game

import (
	"github.com/asragi/yasoba-prototype/adapter/ebiten/drawing/adapter"
	"github.com/asragi/yasoba-prototype/app"
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/common/character/hero"
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/hajimehoshi/ebiten/v2"
)

type mpPort interface {
	InitialMP() hero.InitialMP
	RecoverMP() hero.RecoverMP
	MaxMP() hero.MaxMP
}

type commonManager interface {
	GetMPManager() mpPort
}

type commonManagerImpl struct {
	mpManager mpPort
}

func (c *commonManagerImpl) GetMPManager() mpPort {
	return c.mpManager
}

type Game struct {
	app           *app.App
	cfg           app.Config
	activeScene   scene.Scene
	debugOverlay  *debug.Debug
	drawing       *drawing.Drawing
	width         int
	height        int
	commonManager commonManager
}

func New(appInstance *app.App, cfg app.Config) *Game {
	commonMPManager := hero.NewCommonMPManager(
		hero.InitialMP(1),
		hero.RecoverMP(1),
		hero.MaxMP(4),
	)
	game := &Game{
		app:           appInstance,
		cfg:           cfg,
		debugOverlay:  appInstance.DebugOverlay,
		drawing:       appInstance.Drawing,
		width:         cfg.GameWidth,
		height:        cfg.GameHeight,
		commonManager: &commonManagerImpl{mpManager: commonMPManager},
	}
	debugScene := game.newDebugScene()
	game.activeScene = debugScene
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
	g.drawing.DrawEnd(adapter.NewEbitenImage(screen))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func (g *Game) ActivateBattleScene() {
	g.SwitchToBattleScene(g.cfg.BattleID)
}

func (g *Game) SwitchToBattleScene(battleID battle.BattleId) {
	g.activeScene = g.newBattleScene(battleID)
}

func (g *Game) SwitchToDebugScene() {
	g.activeScene = g.newDebugScene()
}

func (g *Game) newBattleScene(battleID battle.BattleId) *scene.BattleScene {
	if g.app == nil {
		panic("game: app is nil")
	}
	battleMpManager := command.NewMpManager(g.commonManager.GetMPManager())
	option := scene.NewBattleOption(
		battleMpManager,
		nil,
		g.cfg.BattleSettingID,
		battleID,
		g.SwitchToBattleScene,
		g.SwitchToDebugScene,
	)
	battleScene := g.app.CreateBattleScene(option)
	if battleScene == nil {
		panic("game: failed to create battle scene")
	}
	return battleScene
}

func (g *Game) newDebugScene() *scene.DebugScene {
	if g.app == nil {
		panic("game: app is nil")
	}
	create := g.app.CreateDebugScene
	if create == nil {
		panic("game: CreateDebugScene is nil")
	}
	debugScene := create()
	if debugScene == nil {
		panic("game: failed to create debug scene")
	}
	debugScene.SetOnSelectBattle(func() {
		g.SwitchToBattleScene(g.cfg.BattleID)
	})
	return debugScene
}
