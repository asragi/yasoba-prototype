package scene

import (
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/battle/enemy"
	"github.com/asragi/yasoba-prototype/battle/invoke"
	"github.com/asragi/yasoba-prototype/battle/setup"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/common/character"
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
	"github.com/asragi/yasoba-prototype/view/common/message"
	"github.com/asragi/yasoba-prototype/view/common/selection"
	"github.com/asragi/yasoba-prototype/view/common/shake"
	transitionview "github.com/asragi/yasoba-prototype/view/common/transition"
	"github.com/asragi/yasoba-prototype/widget"
)

const battleTransitionFrameCount = 15

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
	inputManager input.InputManager,
	effectManager *widget.EffectManager,
	serveEnemyView battleenemy.ServeEnemyViewData,
	serveActor actor.ActorSupplier,
	newProcessBattle battle.NewProcessBattleFunc,
	produceCreateSequence sequence.ProduceCreateSequence,
	produceCheckInvokeSequence invoke.ProduceCheckInvokeSequence,
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
		allTextId := createAllTextIdList(allActorId, actorNames)

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
		battleEnemyDisplay := createBattleEnemyDisplay(newBattleEnemyDisplay, battleResponse.EnemyIds, battleSetting.Enemies)

		// バトル選択ウィンドウの設定
		var selectedCommand battle.PlayerCommand
		var targetSelectWindow *selection.SelectWindow
		onSubmit := func(command battle.PlayerCommand) {
			selectedCommand = command
			targetSelectWindow.Open()
			inputManager.Set(targetSelectWindow)
		}

		actorDisplay := newBattleActorDisplay(mainActor)
		subActorDisplay := newBattleSubActorDisplay(subActor)
		subActorDialog := battledialogue.CreateNewBattlePartnerDialogue(newMessageWindow)()
		transitionView := newTransitionView()
		battleTransition := transitionview.New(
			battleTransitionFrameCount,
			transitionView.Draw,
			transitionview.InitialStateOpaque,
		)

		// エフェクト関数の定義
		playEffect := createPlayEffectFunction(serveActor, battleEnemyDisplay, subActorDisplay, actorDisplay, effectManager)
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

		battleSelectWindow := createBattleSelectWindow(newBattleSelectWindow, inputManager, onSubmit)

		showPlayerCommandWindow := func() {
			battleSelectWindow.Open()
			inputManager.Set(battleSelectWindow)
		}

		hidePlayerCommandWindow := func() {
			battleSelectWindow.Close()
			inputManager.Set(input.InputReceiverEmptyInstance)
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
		}
		battleScene.sequences.AddSequence(seq)
		battleScene.sequences.AddSequence(commandWindowSequence)

		// バトル処理の設定
		processBattle := newProcessBattle(battleResponse, battleScene.onBattleEnd)
		playSequence := createPlayBattleSequence(
			skillToSequence,
			newBattleSequence,
			battleScene.battleSequence.Add,
			actorIdToEnemy,
			serveEnemyView,
		)

		// ターゲット選択の設定
		closeWindowOnTargetSelect := createOnSubmitTargetSelect(battleSelectWindow, inputManager)
		onTargetSelect := createOnTargetSelect(
			closeWindowOnTargetSelect,
			func(index int) actor.ActorId { return allActorId[index] },
			func() battle.PlayerCommand { return selectedCommand },
			battleScene.battleSequence.Reset,
			playSequence,
			processBattle,
		)
		targetSelectWindow = newSelectWindow(
			&drawing.Vector{X: 80, Y: 0},
			drawing.PivotBottomLeft,
			drawing.DepthWindow,
			allTextId,
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

func createActorNamesMapping(enemyIds []*setup.EnemyIdPair, serveEnemyName enemy.NameServer) map[actor.ActorId]text.TextId {
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
		&drawing.Vector{X: 192, Y: 0},
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

func createBattleEnemyDisplay(newBattleEnemyDisplay battleenemy.NewBattleEnemyDisplayFunc, enemyIds []*setup.EnemyIdPair, enemySettings []*config.EnemySetting) *battleenemy.BattleEnemyDisplay {
	displayArgs := battleenemy.ToDisplayArgs(enemyIds, enemySettings)
	return newBattleEnemyDisplay(
		displayArgs,
		drawing.DepthEnemy,
	)
}

func createBattleSelectWindow(
	newBattleSelectWindow window.NewBattleSelectWindowFunc,
	inputManager input.InputManager,
	onSubmit func(battle.PlayerCommand),
) *window.BattleSelectWindow {
	battleSelectWindow := newBattleSelectWindow(
		&drawing.Vector{X: 0, Y: 0},
		drawing.PivotBottomLeft,
		drawing.DepthWindow,
		[]battle.PlayerCommand{
			battle.PlayerCommandAttack,
			battle.PlayerCommandFire,
			battle.PlayerCommandBarrier,
			battle.PlayerCommandThunder,
			battle.PlayerCommandWind,
			battle.PlayerCommandFocus,
			battle.PlayerCommandDefend,
		},
		onSubmit,
	)
	battleSelectWindow.Close()
	inputManager.Set(input.InputReceiverEmptyInstance)
	return battleSelectWindow
}

func createPlayEffectFunction(serveActor actor.ActorSupplier, battleEnemyDisplay *battleenemy.BattleEnemyDisplay, subActorDisplay *battleactor.BattleSubActorDisplay, actorDisplay *battleactor.BattleActorDisplay, effectManager *widget.EffectManager) func(widget.EffectId, actor.ActorId) {
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

func createDoShakeFunction(serveActor actor.ActorSupplier, battleEnemyDisplay *battleenemy.BattleEnemyDisplay, subActorDisplay *battleactor.BattleSubActorDisplay, shakeEmitter *shake.EmitShake) func(actor.ActorId) {
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

func createSetEmotionFunction(serveActor actor.ActorSupplier, battleEnemyDisplay *battleenemy.BattleEnemyDisplay, subActorDisplay *battleactor.BattleSubActorDisplay, actorDisplay *battleactor.BattleActorDisplay) func(actor.ActorId, battleemotion.EmotionType) {
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
	inputManager input.InputManager,
) func() {
	return func() {
		battleSelectWindow.Close()
		inputManager.Set(input.InputReceiverEmptyInstance)
	}
}

// TODO: View非依存のLogic部分だけ抽出してCoreに移動したい
