//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire ./app

package app

import (
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
	battleactor "github.com/asragi/yasoba-prototype/component/battle/actor"
	battleenemy "github.com/asragi/yasoba-prototype/component/battle/enemy"
	battleevent "github.com/asragi/yasoba-prototype/component/battle/event"
	battlehp "github.com/asragi/yasoba-prototype/component/battle/hp"
	battleselect "github.com/asragi/yasoba-prototype/component/battle/window"
	"github.com/asragi/yasoba-prototype/debug"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/game/enemy"
	gameSkill "github.com/asragi/yasoba-prototype/game/skill"
	"github.com/asragi/yasoba-prototype/invoke"
	"github.com/asragi/yasoba-prototype/scene"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/util"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/google/wire"
)

func initializeApp(cfg Config) (*App, error) {
	wire.Build(
		frontend.NewDrawing,
		debug.CreateDrawParameters,
		makeBattleScene,
		makeDebugScene,
		buildApp,

		frontend.CreateResourceManager,
		widget.CreateNewText,
		component.StandByNewMessageWindow,
		component.StandByNewSelectWindow,
		battleselect.StandByNewBattleSelectWindow,
		battleactor.StandByNewFaceWindow,
		battlehp.CreateNewBattleHPDisplay,
		battleactor.CreateNewBattleParameterDisplay,
		component.CreateNewDisplayDamage,
		battleactor.CreateNewBattleActorDisplay,
		battleactor.CreateNewBattleSubActorDisplay,
		component.StandByNewVariableMessageWindow,
		battleenemy.NewBattleActorGraphics,
		battleenemy.CreateNewBattleEnemyDisplay,
		widget.CreateServeEffectData,
		widget.NewEffectManager,
		battleenemy.NewServeEnemyViewData,
		battleenemy.CreateGetEnemyGraphics,

		sequence.InitializeProduceCreateSequence,
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

		gameSkill.NewSkillServer,
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
		battleevent.CreateServeBattleEventSequence,
		battleevent.CreateExecBattleEventSequence,
		wire.Value(battleevent.SkillToSequenceFunc(battleevent.ToEventSequenceId)),
		config.NewServer,
		scene.InitializeCreateBattleScene,
		scene.InitializeCreateDebugScene,

		wire.Value(frontend.MaruMinya),
		wire.Bind(new(frontend.ResourceManagerInterface), new(*frontend.ResourceManager)),
		wire.Bind(new(battle.AllActorServer), new(*actor.InMemoryActorServer)),
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
	return widget.CreateNewWindow(resource, cfg.GameWidth, cfg.GameHeight)
}

func makeDebugScene(create scene.CreateDebugScene) *scene.DebugScene {
	return create()
}

func makeBattleScene(cfg Config, create scene.CreateBattleScene) *scene.BattleScene {
	return create(&scene.BattleOption{
		OnEnd:           nil,
		BattleSettingId: cfg.BattleSettingID,
		BattleId:        cfg.BattleID,
	})
}
