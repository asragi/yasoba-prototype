package main

import (
	"log"
	"math/rand"

	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/invoke"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameWidth  = 384
	GameHeight = 288
	DrawRate   = 1
)

type sceneRunner interface {
	Update()
	Draw(frontend.DrawFunc)
}

var (
	drawing      *frontend.Drawing
	currentScene sceneRunner
	debugParams  *debug.Debug
)

func init() {
	drawing = frontend.NewDrawing()
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		log.Fatal(err)
	}
	prepareProduceCreateSequence := sequence.InitializeProduceCreateSequence()
	actorServer := core.NewInMemoryActorServer()
	textServer, err := core.LoadTextDataFromYaml("data/text_data.yaml")
	if err != nil {
		log.Fatal(err)
	}
	produceCreateSequence := prepareProduceCreateSequence(textServer)
	characterServer := core.CreateCharacterServer()
	enemyServer := core.CreateEnemyServer()
	prepareActor := core.CreatePrepareActorService(characterServer, enemyServer, actorServer)
	processCommand := core.CreateProcessPlayerCommand(actorServer.Get)
	newWindow := widget.CreateNewWindow(resource, GameWidth, GameHeight)
	newText := widget.CreateNewText(resource)
	newMessageWindow := component.StandByNewMessageWindow(newText, newWindow)
	newSelectWindow := component.StandByNewSelectWindow(resource, newText, textServer)
	newBattleSelectWindow := component.StandByNewBattleSelectWindow(newSelectWindow)
	newFaceWindow := component.StandByNewFaceWindow(resource, newWindow)
	battleSettingServer := core.CreateServeBattleSetting()
	skillServer := core.NewSkillServer()
	random := rand.Float64
	applySkill := core.CreateSkillApply(skillServer, actorServer.Get, actorServer.Upsert, random)
	battleSequenceServer := component.CreateServeBattleEventSequence()
	prepareBattleSequence := component.CreateExecBattleEventSequence(
		textServer,
		battleSequenceServer,
	)
	enemyGraphicServer := component.CreateGetEnemyGraphics()
	newDisplayDamage := component.CreateNewDisplayDamage(newText)
	newBattleActorGraphics := component.NewBattleActorGraphics(
		resource,
		enemyGraphicServer,
		newDisplayDamage,
	)
	newHPDisplay := component.CreateNewBattleHPDisplay(frontend.MaruMinya, newText)
	newParameterDisplay := component.CreateNewBattleParameterDisplay(newWindow, newHPDisplay)
	newBattleActorDisplay := component.CreateNewBattleActorDisplay(newFaceWindow, newDisplayDamage, newParameterDisplay)
	newBattleSubActorDisplay := component.CreateNewBattleSubActorDisplay(
		newFaceWindow,
		newDisplayDamage,
		newParameterDisplay,
	)
	newBattleEnemyDisplay := component.CreateNewBattleEnemyDisplay(newBattleActorGraphics)
	effectData := widget.CreateServeEffectData()
	effectManager := widget.NewEffectManager(effectData, resource)
	serveEnemyView := component.NewServeEnemyViewData()
	choiceTarget := core.CreateChoiceSkillTarget(random)
	newChoiceRandomAction := core.StandByCreateRandomAction(
		random,
		skillServer,
		choiceTarget,
	)
	choiceAction := core.CreateNewChoiceAction(newChoiceRandomAction)
	decideActionOrder := core.CreateDecideActionOrder(actorServer)
	serveBattleState := core.CreateServeBattleState(actorServer)
	checkCombination := core.CreateCheckCombination()
	newPartnerActionServer := core.StandByNewPartnerActionServer(random, serveBattleState)
	partnerActionServer := newPartnerActionServer()
	initializeBattle := core.CreateInitializeBattle(prepareActor, partnerActionServer.DecidePlan)
	newProcessBattle := core.StandByCreateProcessBattle(
		actorServer.Get,
		serveBattleState,
		processCommand,
		partnerActionServer.GetPlan,
		checkCombination,
		applySkill,
		decideActionOrder,
		choiceAction,
	)
	newVariableMessageWindow := component.StandByNewVariableMessageWindow(newWindow, newText, textServer)
	produceCheckInvokeSequence := invoke.InitializeProduceCheckInvokeSequence()
	createBattleScene := scene.InitializeCreateBattleScene(
		newMessageWindow,
		newSelectWindow,
		newBattleSelectWindow,
		newBattleActorDisplay,
		newBattleSubActorDisplay,
		core.CreateEnemyNameServer(),
		initializeBattle,
		battleSettingServer,
		prepareBattleSequence,
		component.ToEventSequenceId,
		newBattleEnemyDisplay,
		effectManager,
		serveEnemyView,
		actorServer.Get,
		newVariableMessageWindow,
		newProcessBattle,
		produceCreateSequence,
		produceCheckInvokeSequence,
	)
	createDebugScene := scene.InitializeCreateDebugScene(newSelectWindow)
	debugScene := createDebugScene(
		&scene.DebugOption{
			OnSelectBattle: func() {
				currentScene = createBattleScene(
					&scene.BattleOption{
						OnEnd:    nil,
						BattleId: core.BattleIdTest001,
						//BattleId: core.BattleIdTripleTest,
					},
				)
			},
		},
	)
	currentScene = debugScene
	debugParams = debug.CreateDrawParameters(newText)
}

type Game struct{}

func (g *Game) Update() error {
	currentScene.Update()
	debugParams.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	currentScene.Draw(drawing.Draw)
	debugParams.Draw(drawing.Draw)
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
