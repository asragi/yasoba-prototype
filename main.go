package main

import (
	"log"
	"math/rand"

	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/combination"
	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/battle/decision"
	"github.com/asragi/yasoba-prototype/battle/partner"
	"github.com/asragi/yasoba-prototype/battle/setup"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/game/enemy"
	gameSkill "github.com/asragi/yasoba-prototype/game/skill"
	"github.com/asragi/yasoba-prototype/invoke"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameWidth  = 384
	GameHeight = 288
	DrawRate   = 1
)

var (
	drawing     *frontend.Drawing
	battleScene *scene.BattleScene
	debugParams *debug.Debug
)

func init() {
	drawing = frontend.NewDrawing()
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		log.Fatal(err)
	}
	prepareProduceCreateSequence := sequence.InitializeProduceCreateSequence()
	actorServer := actor.NewInMemoryActorServer()
	textServer, err := text.LoadFromYAML("data/text_data.yaml")
	if err != nil {
		log.Fatal(err)
	}
	produceCreateSequence := prepareProduceCreateSequence(textServer)
	characterServer := character.CreateCharacterServer()
	enemyServer := enemy.CreateEnemyServer()
	prepareActor := setup.NewPrepareService(characterServer, enemyServer, actorServer)
	processCommand := battle.CreateProcessPlayerCommand(actorServer.Get)
	newWindow := widget.CreateNewWindow(resource, GameWidth, GameHeight)
	newText := widget.CreateNewText(resource)
	newMessageWindow := component.StandByNewMessageWindow(newText, newWindow)
	newSelectWindow := component.StandByNewSelectWindow(resource, newText, textServer)
	newBattleSelectWindow := component.StandByNewBattleSelectWindow(newSelectWindow)
	newFaceWindow := component.StandByNewFaceWindow(resource, newWindow)
	battleSettingServer := config.NewServer()
	skillServer := gameSkill.NewSkillServer()
	random := rand.Float64
	applySkill := skill.CreateSkillApply(skillServer, actorServer.Get, actorServer.Upsert, random)
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
	choiceTarget := decision.CreateChoiceSkillTarget(random)
	newChoiceRandomAction := decision.StandByCreateRandomAction(
		random,
		skillServer,
		choiceTarget,
	)
	choiceAction := decision.CreateNewChoiceAction(newChoiceRandomAction)
	decideActionOrder := battle.CreateDecideActionOrder(actorServer)
	serveBattleState := decision.CreateServeBattleState(actorServer)
	checkCombination := combination.CreateCheckCombination()
	newPartnerActionServer := partner.StandByNewPartnerActionServer(random, serveBattleState)
	partnerActionServer := newPartnerActionServer()
	initializeBattle := battle.CreateInitializeBattle(prepareActor, partnerActionServer.DecidePlan)
	newProcessBattle := battle.StandByCreateProcessBattle(
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
	newBattleScene := scene.InitializeCreateBattleScene(
		newMessageWindow,
		newSelectWindow,
		newBattleSelectWindow,
		newBattleActorDisplay,
		newBattleSubActorDisplay,
		enemy.CreateNameServer(),
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
	battleScene = newBattleScene(
		&scene.BattleOption{
			OnEnd:           nil,
			BattleSettingId: config.IdTest,
			BattleId:        battle.BattleIdTest001,
			//BattleSettingId: config.IdTripleTest,
		},
	)
	debugParams = debug.CreateDrawParameters(newText)
}

type Game struct{}

func (g *Game) Update() error {
	battleScene.Update()
	debugParams.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	debugParams.Draw(drawing.Draw)
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
