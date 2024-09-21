package main

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
)

const (
	GameWidth  = 384
	GameHeight = 288
	DrawRate   = 1
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
	newSelectWindow := component.StandByNewSelectWindow(resource, textServer)
	newBattleSelectWindow := component.StandByNewBattleSelectWindow(newSelectWindow)
	newFaceWindow := component.StandByNewFaceWindow(resource)
	battleSettingServer := game.CreateServeBattleSetting()
	skillServer := core.NewSkillServer()
	random := rand.Float64
	applySkill := core.CreateSkillApply(skillServer, actorServer.Get, actorServer.Upsert, random)
	battleSequenceServer := component.CreateServeBattleEventSequence()
	prepareBattleSequence := component.CreateExecBattleEventSequence(
		textServer,
		battleSequenceServer,
	)
	skillToSequence := component.CreateSkillToSequenceId()
	enemyGraphicServer := component.CreateGetEnemyGraphics()
	newDisplayDamage := component.CreateNewDisplayDamage(resource)
	newBattleActorGraphics := component.NewBattleActorGraphics(
		resource,
		enemyGraphicServer,
		newDisplayDamage,
	)
	newBattleActorDisplay := component.CreateNewBattleActorDisplay(newBattleActorGraphics)
	effectData := widget.CreateServeEffectData()
	effectManager := widget.NewEffectManager(effectData, resource)
	serveEnemyView := component.NewServeEnemyViewData()
	newBattleScene := scene.StandByNewBattleScene(
		newMessageWindow,
		newSelectWindow,
		newBattleSelectWindow,
		newFaceWindow,
		core.CreateEnemyNameServer(),
		initializeBattle,
		postCommand,
		applySkill,
		battleSettingServer,
		prepareBattleSequence,
		skillToSequence,
		newBattleActorDisplay,
		effectManager,
		serveEnemyView,
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
	return GameWidth, GameHeight
}

func main() {
	ebiten.SetWindowSize(GameWidth*DrawRate, GameHeight*DrawRate)
	ebiten.SetWindowTitle("yasoba-prototype")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
