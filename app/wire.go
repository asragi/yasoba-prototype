// go:build wireinject
//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire ./app

package app

import (
	"math/rand"

	drawingadapter "github.com/asragi/yasoba-prototype/adapter/ebiten/drawing/adapter"
	sequenceadapter "github.com/asragi/yasoba-prototype/adapter/ebiten/sequence/adapter"
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/combination"
	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/battle/decision"
	"github.com/asragi/yasoba-prototype/battle/enemy"
	"github.com/asragi/yasoba-prototype/battle/invoke"
	"github.com/asragi/yasoba-prototype/battle/partner"
	"github.com/asragi/yasoba-prototype/battle/setup"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/common/character"
	"github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/util"
	actorview "github.com/asragi/yasoba-prototype/view/battle/actor"
	damageview "github.com/asragi/yasoba-prototype/view/battle/damage"
	enemyview "github.com/asragi/yasoba-prototype/view/battle/enemy"
	eventview "github.com/asragi/yasoba-prototype/view/battle/event"
	hpview "github.com/asragi/yasoba-prototype/view/battle/hp"
	windowview "github.com/asragi/yasoba-prototype/view/battle/window"
	messageview "github.com/asragi/yasoba-prototype/view/common/message"
	"github.com/asragi/yasoba-prototype/view/common/selection"
	transitionview "github.com/asragi/yasoba-prototype/view/common/transition"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/google/wire"
)

func initializeApp(cfg Config) (*App, error) {
	wire.Build(
		drawing.NewDrawing,
		debug.CreateDrawParameters,
		buildApp,

		frontend.CreateResourceManager,
		drawingadapter.NewDrawTextFunc,
		widget.CreateNewText,
		messageview.StandByNewMessageWindow,
		makeSelectCursor,
		selection.StandByNewSelectWindow,
		windowview.StandByNewBattleSelectWindow,
		makeFaceWindowFactory,
		makeTransitionView,
		hpview.CreateNewBattleHPDisplay,
		actorview.CreateNewBattleParameterDisplay,
		damageview.CreateNewDisplayDamage,
		actorview.CreateNewBattleActorDisplay,
		actorview.CreateNewBattleSubActorDisplay,
		makeBattleEnemyGraphics,
		enemyview.CreateNewBattleEnemyDisplay,
		widget.CreateServeEffectData,
		widget.NewEffectManager,
		enemyview.NewServeEnemyViewData,
		enemyview.CreateGetEnemyGraphics,

		sequenceadapter.InitializeProduceCreateSequence,
		loadTextServer,
		makeProduceCreateSequence,

		actor.NewInMemoryActorServer,
		makeActorSupplier,
		makeActorUpdater,
		battle.CreateProcessPlayerCommand,
		makeServeBattleState,
		battle.CreateDecideActionOrder,

		character.CreateCharacterServer,
		enemy.CreateEnemyServer,
		enemy.CreateNameServer,
		makePrepareService,

		skill.NewSkillServer,
		skill.CreateSkillApply,
		decision.CreateChoiceSkillTarget,
		decision.StandByCreateRandomAction,
		decision.CreateNewChoiceAction,
		combination.CreateCheckCombination,
		wire.Value(util.EmitRandomFunc(rand.Float64)),

		partner.StandByNewPartnerActionServer,
		makePartnerActionServer,
		makeInitializeBattle,
		makeProcessBattle,

		invoke.InitializeProduceCheckInvokeSequence,
		eventview.CreateServeBattleEventSequence,
		eventview.CreateExecBattleEventSequence,
		wire.Value(eventview.SkillToSequenceFunc(eventview.ToEventSequenceId)),
		config.NewServer,
		scene.InitializeCreateBattleScene,
		scene.InitializeCreateDebugScene,

		wire.Value(frontend.MaruMinya),
		wire.Bind(new(frontend.ResourceManagerInterface), new(*frontend.ResourceManager)),
		wire.Bind(new(battle.AllActorServer), new(*actor.InMemoryActorServer)),
		wire.Bind(new(widget.FontProvider), new(*frontend.ResourceManager)),
		makeWindowFunc,
	)
	return nil, nil
}

func loadTextServer(cfg Config) (text.ServeTextDataFunc, error) {
	return text.LoadFromYAML(cfg.TextDataPath)
}

func makeProduceCreateSequence(prepare sequence.PrepareProduceCreateSequence, textServer text.ServeTextDataFunc) sequence.ProduceCreateSequence {
	return prepare(textServer)
}

func makeActorSupplier(server *actor.InMemoryActorServer) actor.ActorSupplier {
	return server.Get
}

func makeActorUpdater(server *actor.InMemoryActorServer) actor.UpdateActorFunc {
	return server.Upsert
}

func makeServeBattleState(server *actor.InMemoryActorServer) decision.ServeBattleState {
	return decision.CreateServeBattleState(server)
}

func makePrepareService(
	serveCharacter character.ServeCharacterFunc,
	serveEnemy enemy.ServeEnemyData,
	server *actor.InMemoryActorServer,
) setup.Service {
	return setup.NewPrepareService(serveCharacter, serveEnemy, server)
}

func makePartnerActionServer(factory partner.NewPartnerActionServer) *partner.PartnerActionServer {
	return factory()
}

func makeInitializeBattle(prepare setup.Service, partnerServer *partner.PartnerActionServer) battle.InitializeBattleFunc {
	return battle.CreateInitializeBattle(prepare, partnerServer.DecidePlan)
}

func makeProcessBattle(
	getActor actor.ActorSupplier,
	serveBattleState decision.ServeBattleState,
	processCommand battle.ProcessPlayerCommandFunc,
	partnerServer *partner.PartnerActionServer,
	checkCombination combination.CheckFunc,
	skillApply skill.SkillApplyFunc,
	decideActionOrder battle.DecideActionOrderFunc,
	choiceAction decision.NewChoiceActionFunc,
) battle.NewProcessBattleFunc {
	return battle.StandByCreateProcessBattle(
		getActor,
		serveBattleState,
		processCommand,
		partnerServer.GetPlan,
		checkCombination,
		skillApply,
		decideActionOrder,
		choiceAction,
	)
}

func makeWindowFunc(resource *frontend.ResourceManager, cfg Config) widget.NewWindowFunc {
	return widget.CreateNewWindow(resource.GetTexture, cfg.GameWidth, cfg.GameHeight)
}

func makeBattleEnemyGraphics(
	resource frontend.ResourceManagerInterface,
	getEnemyGraphics enemyview.GetEnemyGraphicsFunc,
	newDisplayDamage damageview.NewDisplayDamageFunc,
	cfg Config,
) enemyview.NewBattleEnemyGraphicsFunc {
	return enemyview.NewBattleActorGraphics(
		resource,
		getEnemyGraphics,
		newDisplayDamage,
		cfg.GameWidth,
		cfg.GameHeight,
	)
}

func makeTransitionView(resource frontend.ResourceManagerInterface, cfg Config) transitionview.NewTransitionViewFunc {
	createImage := func(width, height int) transitionview.OverlayImage {
		img := resource.NewEmptyImage(width, height)
		overlay, ok := img.(transitionview.OverlayImage)
		if !ok {
			panic("transition: created image does not satisfy overlay interface")
		}
		return overlay
	}
	return transitionview.CreateNewView(
		cfg.GameWidth,
		cfg.GameHeight,
		drawing.DepthTransition,
		createImage,
	)
}

func makeSelectCursor(resource *frontend.ResourceManager) selection.NewCursor {
	return func(
		relativePosition *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
	) selection.Cursor {
		return widget.NewImage(
			relativePosition,
			pivot,
			depth,
		resource.GetTexture(texture.Cursor),
		)
	}
}

func makeFaceWindowFactory(resource *frontend.ResourceManager, newWindow widget.NewWindowFunc) actorview.NewFaceWindowFunc {
	return actorview.StandByNewFaceWindow(resource, newWindow)
}
