package scene

import (
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/battle/enemy"
	"github.com/asragi/yasoba-prototype/battle/invoke"
	"github.com/asragi/yasoba-prototype/battle/setup"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/common/character"
	"github.com/asragi/yasoba-prototype/common/character/hero"
	battleemotion "github.com/asragi/yasoba-prototype/common/emotion"
	"github.com/asragi/yasoba-prototype/sequence"
	seqtransition "github.com/asragi/yasoba-prototype/sequence/transition"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/toolkit/input"
	battleactor "github.com/asragi/yasoba-prototype/view/battle/actor"
	battledialogue "github.com/asragi/yasoba-prototype/view/battle/dialogue"
	battleenemy "github.com/asragi/yasoba-prototype/view/battle/enemy"
	battleevent "github.com/asragi/yasoba-prototype/view/battle/event"
	"github.com/asragi/yasoba-prototype/view/battle/window"
	"github.com/asragi/yasoba-prototype/view/common/constant"
	"github.com/asragi/yasoba-prototype/view/common/message"
	"github.com/asragi/yasoba-prototype/view/common/selection"
	"github.com/asragi/yasoba-prototype/view/common/selection/option"
	"github.com/asragi/yasoba-prototype/view/common/shake"
	transitionview "github.com/asragi/yasoba-prototype/view/common/transition"
	"github.com/asragi/yasoba-prototype/widget"
)

const battleTransitionFrameCount = 15

type mpPort interface {
	InitialMP() hero.InitialMP
	CurrentMP() hero.MP
	Recover()
	Consume(cost command.Cost)
	MaxMP() hero.MaxMP
}

func InitializeCreateBattleScene(
	newMessageWindow message.NewMessageWindowFunc,
	newSelectWindow selection.NewSelectWindowFunc,
	newBattleSelectWindow window.NewBattleSelectWindowFunc,
	newBattleActorDisplay battleactor.NewBattleActorDisplayFunc,
	newBattleSubActorDisplay battleactor.NewBattleSubActorDisplayFunc,
	newTransitionView transitionview.NewTransitionViewFunc,
	serveEnemyName enemy.NameServer,
	initializeBattle battle.InitializeBattleFunc,
	getBattleSetting config.ServeFunc,
	createNewBattleSequence battleevent.PrepareBattleEventSequenceFunc,
	skillToSequence battleevent.SkillToSequenceFunc,
	newBattleEnemyDisplay battleenemy.NewBattleEnemyDisplayFunc,
	inputManager input.Manager,
	effectManager *widget.EffectManager,
	serveEnemyView battleenemy.ServeEnemyViewData,
	serveActor actor.ActorSupplier,
	newProcessBattle battle.NewProcessBattleFunc,
	produceCreateSequence sequence.ProduceCreateSequence,
	produceCheckInvokeSequence invoke.ProduceCheckInvokeSequence,
	newTextItem option.NewTextItemFunc,
) CreateBattleScene {
	return func(option *BattleOption) *BattleScene {
		if option == nil {
			panic("battle scene option is nil")
		}
		if option.SwitchToBattle == nil {
			panic("battle scene option SwitchToBattle is nil")
		}
		if option.SwitchToDebug == nil {
			panic("battle scene option SwitchToDebug is nil")
		}
		if inputManager == nil {
			panic("battle scene input manager is nil")
		}
		mpManager := option.MpManager
		if mpManager == nil {
			panic("battle scene mp manager is nil")
		}
		// 戦闘設定の取得と初期化
		battleSetting := getBattleSetting(option.BattleSettingId)
		enemyIds := extractEnemyIds(battleSetting)

		// 戦闘の初期化
		initializeRequest := &battle.InitializeBattleRequest{
			// TODO: variables must be provided by args
			MainActorCharacterId: character.CharacterLuneId,
			SubActorCharacterId:  character.CharacterSunnyId,
			EnemyIds:             enemyIds,
		}
		battleResponse := initializeBattle(initializeRequest)

		// アクターデータの取得
		mainActor := serveActor(battleResponse.MainActorId)
		subActor := serveActor(battleResponse.SubActorId)

		// データマッピングの作成
		actorIdToEnemy := createActorIdToEnemyMapping(battleResponse.EnemyIds)
		actorNames := createActorNamesMapping(battleResponse.EnemyIds, serveEnemyName)
		allActorId := createAllActorIdList(battleResponse)

		displayedHp := func() map[actor.ActorId]character.HP {
			hp := make(map[actor.ActorId]character.HP, len(allActorId))
			for _, id := range allActorId {
				actor := serveActor(id)
				if actor == nil {
					continue
				}
				hp[id] = actor.HP
			}
			return hp
		}()

		// UIコンポーネントの初期化
		messageWindow := createMessageWindow(newMessageWindow)
		setMessageWindowText := createSetMessageWindowTextFunction(messageWindow)
		battleEnemyDisplay := createBattleEnemyDisplay(
			newBattleEnemyDisplay,
			battleResponse.EnemyIds,
			battleSetting.Enemies,
		)

		// バトル選択ウィンドウの設定
		var selectedCommand command.Id
		var targetSelectWindow *selection.SelectWindow
		onSubmit := func(cmd command.Id) {
			selectedCommand = cmd
			targetSelectWindow.Open()
			inputManager.Set(targetSelectWindow)
		}

		actorDisplay := newBattleActorDisplay(mainActor, mpManager.InitialMP(), mpManager.MaxMP())
		subActorDisplay := newBattleSubActorDisplay(subActor)
		subActorDialog := battledialogue.CreateNewBattlePartnerDialogue(newMessageWindow)()
		transitionView := newTransitionView()
		battleTransition := transitionview.New(
			battleTransitionFrameCount,
			transitionView.Draw,
			transitionview.InitialStateOpaque,
		)

		// エフェクト関数の定義
		playEffect := createPlayEffectFunction(
			serveActor,
			battleEnemyDisplay,
			subActorDisplay,
			actorDisplay,
			effectManager,
		)
		uiShake := shake.NewShake()
		doShake := createDoShakeFunction(serveActor, battleEnemyDisplay, subActorDisplay, uiShake)
		setDamage := createSetDamageFunction(serveActor, battleEnemyDisplay, subActorDisplay, actorDisplay, displayedHp)
		setEmotion := createSetEmotionFunction(serveActor, battleEnemyDisplay, subActorDisplay, actorDisplay)

		// バトルシーケンスの設定
		newBattleSequence := createNewBattleSequence(
			messageWindow,
			doShake,
			setEmotion,
			setDamage,
			playEffect,
			battleEnemyDisplay.SetDisappear,
		)

		battleSelectWindow := createBattleSelectWindow(newBattleSelectWindow, inputManager, onSubmit, mpManager)
		showPlayerCommandWindow := func() {
			battleSelectWindow.Open(mpManager.CurrentMP())
			inputManager.Set(battleSelectWindow)
		}

		hidePlayerCommandWindow := func() {
			battleSelectWindow.Close()
			inputManager.Set(input.EmptyInstance)
		}

		startTransitionFadeOut := seqtransition.NewController(
			battleTransition.FadeOut,
			battleTransition.IsFadingOut,
		)

		startTransitionFadeIn := seqtransition.NewController(
			battleTransition.FadeIn,
			battleTransition.IsFadingIn,
		)

		// パートナーダイアログの設定
		setPartnerDialogue := createSetPartnerDialogueFunction(subActorDialog)
		createSequence := produceCreateSequence(
			setPartnerDialogue,
			setMessageWindowText,
			setEmotion,
			subActorDialog.Open,
			subActorDialog.Close,
			showPlayerCommandWindow,
			hidePlayerCommandWindow,
			option.SwitchToBattle,
			option.SwitchToDebug,
			startTransitionFadeOut,
			startTransitionFadeIn,
		)
		seq := createSequence("test_sequence_0000")
		commandWindowSequence := createSequence("test_sequence_transition_command")

		// バトルシーンの作成
		initialMainActor := serveActor(battleResponse.MainActorId)
		if initialMainActor == nil {
			panic("battle main actor not found")
		}

		battleScene := &BattleScene{
			ui: battleUI{
				messageWindow:      messageWindow,
				battleSelectWindow: battleSelectWindow,
				actorDisplay:       actorDisplay,
				subActorDisplay:    subActorDisplay,
				subActorDialog:     subActorDialog,
				input:              inputManager,
				battleEnemyDisplay: battleEnemyDisplay,
				effectManager:      effectManager,
				shake:              uiShake,
				transition:         battleTransition,
			},
			battleSequence: battleevent.NewBattleEventSequencer(),
			enemyData:      battleResponse.EnemyIds,
			actorNames:     actorNames,
			createSequence: createSequence,
			sequences:      sequence.CreateSequenceManager(),
			playerBeaten:   initialMainActor.IsBeaten(),
			mpManager:      mpManager,
		}
		battleScene.sequences.AddSequence(seq)
		battleScene.sequences.AddSequence(commandWindowSequence)

		// バトル処理の設定
		processBattle := newProcessBattle(battleResponse, battleScene.onBattleEnd, mpManager)
		playSequence := createPlayBattleSequence(
			skillToSequence,
			newBattleSequence,
			battleScene.battleSequence.Add,
			actorIdToEnemy,
			serveEnemyView,
		)
		battleScene.autoPlayNonPlayer = func() {
			response := processBattle(nil)
			battleScene.playerBeaten = response.IsMainActorBeaten
			battleScene.battleSequence.Reset()
			playSequence(response.SkillApplyResults)
		}

		// ターゲット選択の設定
		closeWindowOnTargetSelect := createOnSubmitTargetSelect(battleSelectWindow, inputManager)
		onTargetSelect := createOnTargetSelect(
			closeWindowOnTargetSelect,
			func(index int) actor.ActorId { return allActorId[index] },
			func() command.Id { return selectedCommand },
			battleScene.battleSequence.Reset,
			playSequence,
			processBattle,
			func(beaten bool) {
				battleScene.playerBeaten = beaten
			},
			func(mp hero.MP) {
				battleScene.ui.actorDisplay.RefreshMp(mp)
			},
		)
		textItems := func() []selection.Item {
			items := make([]selection.Item, len(allActorId))
			for i, id := range allActorId {
				items[i] = newTextItem(actorNames[id])
			}
			return items
		}()
		targetSelectWindow = newSelectWindow(
			&drawing.Vector{X: 80, Y: 0},
			drawing.PivotBottomLeft,
			drawing.DepthWindow,
			0,
			textItems,
			onTargetSelect,
			true,
		)
		battleScene.ui.targetSelectWindow = targetSelectWindow

		getActorHpRatio := func(actorId actor.ActorId) character.HPRatio {
			actor := serveActor(actorId)
			if actor == nil {
				panic("actor not found: " + string(actorId))
			}
			displayHp, ok := displayedHp[actorId]
			if !ok {
				displayHp = actor.HP
			}
			return displayHp.Ratio(actor.MaxHP)
		}
		battleScene.invokedSequences = make(map[sequence.SequenceId]bool)
		battleScene.checkInvokeSequence = produceCheckInvokeSequence(option.BattleId, getActorHpRatio)
		battleScene.checkAndStartSequences(invoke.InvokeTimingStartBattle)

		return battleScene
	}
}

func extractEnemyIds(battleSetting *config.Setting) []enemy.EnemyId {
	ids := make([]enemy.EnemyId, len(battleSetting.Enemies))
	for i, set := range battleSetting.Enemies {
		ids[i] = set.EnemyId
	}
	return ids
}

func createActorIdToEnemyMapping(enemyIds []*setup.EnemyIdPair) map[actor.ActorId]enemy.EnemyId {
	result := make(map[actor.ActorId]enemy.EnemyId)
	for _, pair := range enemyIds {
		result[pair.ActorId] = pair.EnemyId
	}
	return result
}

func createActorNamesMapping(
	enemyIds []*setup.EnemyIdPair,
	serveEnemyName enemy.NameServer,
) map[actor.ActorId]text.TextId {
	names := make(map[actor.ActorId]text.TextId)
	names[actor.ActorLuneId] = text.TextIdLuneName
	names[actor.ActorSunnyId] = text.TextIdSunnyName
	for _, pair := range enemyIds {
		names[pair.ActorId] = serveEnemyName(pair.EnemyId)
	}
	return names
}

func createAllActorIdList(battleResponse *battle.InitializeBattleResponse) []actor.ActorId {
	ids := []actor.ActorId{battleResponse.MainActorId}
	if battleResponse.SubActorId != actor.ActorEmptyId {
		ids = append(ids, battleResponse.SubActorId)
	}
	enemyActorIds := make([]actor.ActorId, len(battleResponse.EnemyIds))
	for i, pair := range battleResponse.EnemyIds {
		enemyActorIds[i] = pair.ActorId
	}
	return append(ids, enemyActorIds...)
}

func createAllTextIdList(allActorId []actor.ActorId, actorNames map[actor.ActorId]text.TextId) []text.TextId {
	texts := make([]text.TextId, 0)
	for _, id := range allActorId {
		texts = append(texts, actorNames[id])
	}
	return texts
}

func createMessageWindow(newMessageWindow message.NewMessageWindowFunc) *message.MessageWindow {
	messageWindow := newMessageWindow(
		&drawing.Vector{X: float64(constant.GameWidthHalf), Y: 0},
		&drawing.Vector{X: 292, Y: 62},
		drawing.DepthWindow,
		drawing.PivotTopCenter,
	)
	testString := "あのイーハトーヴォのすきとおった風\n夏でも底に冷たさをもつ青いそら\nうつくしい森で飾られたモリーオ市"
	messageWindow.SetText(testString, false)
	messageWindow.Open()
	return messageWindow
}

func createSetMessageWindowTextFunction(messageWindow *message.MessageWindow) sequence.SetMessageWindowText {
	return func(value text.String) *sequence.SetMessageWindowTextResponse {
		messageWindow.SetText(value.String(), false)
		messageWindow.Open()
		return &sequence.SetMessageWindowTextResponse{
			CheckIsEnd: func() sequence.IsEnd {
				return sequence.IsEnd(messageWindow.IsTextEnd())
			},
		}
	}
}

func createBattleEnemyDisplay(
	newBattleEnemyDisplay battleenemy.NewBattleEnemyDisplayFunc,
	enemyIds []*setup.EnemyIdPair,
	enemySettings []*config.EnemySetting,
) *battleenemy.BattleEnemyDisplay {
	displayArgs := battleenemy.ToDisplayArgs(enemyIds, enemySettings)
	return newBattleEnemyDisplay(
		displayArgs,
		drawing.DepthEnemy,
	)
}

func createBattleSelectWindow(
	newBattleSelectWindow window.NewBattleSelectWindowFunc,
	inputManager input.Manager,
	onSubmit func(command.Id),
	mpManager mpPort,
) *window.BattleSelectWindow {
	battleSelectWindow := newBattleSelectWindow(
		mpManager.CurrentMP(),
		&drawing.Vector{X: 0, Y: 0},
		drawing.PivotBottomLeft,
		drawing.DepthWindow,
		[]command.Id{
			command.IdAttack,
			command.IdFire,
			command.IdBarrier,
			command.IdThunder,
			command.IdWind,
			command.IdFocus,
			command.IdDefend,
		},
		onSubmit,
	)
	battleSelectWindow.Close()
	inputManager.Set(input.EmptyInstance)
	return battleSelectWindow
}

func createPlayEffectFunction(
	serveActor actor.ActorSupplier,
	battleEnemyDisplay *battleenemy.BattleEnemyDisplay,
	subActorDisplay *battleactor.BattleSubActorDisplay,
	actorDisplay *battleactor.BattleActorDisplay,
	effectManager *widget.EffectManager,
) func(widget.EffectId, actor.ActorId) {
	return func(effectId widget.EffectId, target actor.ActorId) {
		targetActor := serveActor(target)
		position := func() *drawing.Vector {
			if targetActor.Side == actor.ActorSideEnemy {
				return battleEnemyDisplay.GetPosition(target)
			}
			if targetActor.Id == actor.ActorSunnyId {
				return subActorDisplay.GetCenterPosition()
			}
			return actorDisplay.GetMainCharacterPosition()
		}()
		effectManager.CallEffect(effectId, position)
	}
}

func createDoShakeFunction(
	serveActor actor.ActorSupplier,
	battleEnemyDisplay *battleenemy.BattleEnemyDisplay,
	subActorDisplay *battleactor.BattleSubActorDisplay,
	shakeEmitter *shake.EmitShake,
) func(actor.ActorId) {
	return func(actorId actor.ActorId) {
		actor := serveActor(actorId)
		if actor.IsEnemy() {
			battleEnemyDisplay.DoShake(actorId)
			return
		}
		if actor.IsSubActor() {
			subActorDisplay.Shake()
			return
		}
		shakeEmitter.Shake(shake.ShakeDefaultAmplitude, shake.ShakeDefaultPeriod)
	}
}

func createSetDamageFunction(
	serveActor actor.ActorSupplier,
	battleEnemyDisplay *battleenemy.BattleEnemyDisplay,
	subActorDisplay *battleactor.BattleSubActorDisplay,
	actorDisplay *battleactor.BattleActorDisplay,
	displayedHp map[actor.ActorId]character.HP,
) func(actor.ActorId, skill.Damage, character.HP) {
	return func(actorId actor.ActorId, damage skill.Damage, afterHp character.HP) {
		displayedHp[actorId] = afterHp
		actor := serveActor(actorId)
		if actor.IsEnemy() {
			battleEnemyDisplay.SetDamage(actorId, damage)
			return
		}
		if actor.IsSubActor() {
			subActorDisplay.SetDamage(damage, afterHp)
			return
		}
		actorDisplay.SetDamage(damage, afterHp)
	}
}

func createSetEmotionFunction(
	serveActor actor.ActorSupplier,
	battleEnemyDisplay *battleenemy.BattleEnemyDisplay,
	subActorDisplay *battleactor.BattleSubActorDisplay,
	actorDisplay *battleactor.BattleActorDisplay,
) func(actor.ActorId, battleemotion.EmotionType) {
	return func(actorId actor.ActorId, emotion battleemotion.EmotionType) {
		actor := serveActor(actorId)
		if actor.IsEnemy() {
			battleEnemyDisplay.SetEmotion(actorId, emotion)
			return
		}
		if actor.IsSubActor() {
			subActorDisplay.SetEmotion(emotion)
			return
		}
		actorDisplay.SetEmotion(emotion)
	}
}

func createSetPartnerDialogueFunction(subActorDialog *battledialogue.BattlePartnerDialogue) func(text.String) *sequence.SetPartnerDialogueResponse {
	return func(textValue text.String) *sequence.SetPartnerDialogueResponse {
		subActorDialog.SetText(textValue.String(), false)
		return &sequence.SetPartnerDialogueResponse{
			CheckIsEnd: func() sequence.IsEnd {
				result := subActorDialog.IsTextEnd()
				return sequence.IsEnd(result)
			},
		}
	}
}

func createOnSubmitTargetSelect(
	battleSelectWindow *window.BattleSelectWindow,
	inputManager input.Manager,
) func() {
	return func() {
		battleSelectWindow.Close()
		inputManager.Set(input.EmptyInstance)
	}
}

// TODO: View非依存のLogic部分だけ抽出してCoreに移動したい
