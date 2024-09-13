package main

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var (
	drawing     *frontend.Drawing
	battleScene *scene.BattleScene
)

func init() {
	drawing = frontend.NewDrawing()
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		log.Fatal(err)
	}
	actorServer := core.NewInMemoryActorServer()
	textServer := core.CreateServeTextData()
	characterServer := core.CreateCharacterServer()
	enemyServer := core.CreateEnemyServer()
	prepareActor := core.CreatePrepareActorService(characterServer, enemyServer, actorServer)
	initializeBattle := core.CreateInitializeBattle(prepareActor)
	processCommand := core.CreateProcessPlayerCommand(actorServer.Get)
	postCommand := core.CreatePostCommand(processCommand)
	newMessageWindow := component.StandByNewMessageWindow(resource)
	newBattleSelectWindow := component.StandByNewBattleSelectWindow(resource, textServer)
	newFaceWindow := component.StandByNewFaceWindow(resource)
	battleSettingServer := game.CreateServeBattleSetting()
	newBattleScene := scene.StandByNewBattleScene(
		newMessageWindow,
		newBattleSelectWindow,
		newFaceWindow,
		initializeBattle,
		postCommand,
		battleSettingServer,
	)
	battleScene = newBattleScene(
		&scene.BattleOption{
			OnEnd:           nil,
			BattleSettingId: game.BattleSettingTest,
		},
	)
}

type Game struct{}

func (g *Game) Update() error {
	battleScene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	battleScene.Draw(drawing.Draw)
	drawing.DrawEnd(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 384, 288
}

func main() {
	ebiten.SetWindowSize(768, 576)
	ebiten.SetWindowTitle("yasoba-prototype")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
