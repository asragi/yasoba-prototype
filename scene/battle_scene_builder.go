package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/invoke"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/widget"
)

func InitializeCreateBattleScene(
	newMessageWindow component.NewMessageWindowFunc,
	newSelectWindow component.NewSelectWindowFunc,
	newBattleSelectWindow component.NewBattleSelectWindowFunc,
	newBattleActorDisplay component.NewBattleActorDisplayFunc,
	newBattleSubActorDisplay component.NewBattleSubActorDisplayFunc,
	serveEnemyName core.EnemyNameServer,
	initializeBattle core.InitializeBattleFunc,
	getBattleSetting core.ServeBattleSetting,
	createNewBattleSequence component.PrepareBattleEventSequenceFunc,
	skillToSequence component.SkillToSequenceFunc,
	newBattleEnemyDisplay component.NewBattleEnemyDisplayFunc,
	effectManager *widget.EffectManager,
	serveEnemyView component.ServeEnemyViewData,
	serveActor core.ActorSupplier,
	newVariableMessageWindow component.NewVariableMessageWindowFunc,
	newProcessBattle core.NewProcessBattleFunc,
	produceCreateSequence sequence.ProduceCreateSequence,
	produceCheckInvokeSequence invoke.ProduceCheckInvokeSequence,
) CreateBattleScene {
	return func(option *BattleOption) *BattleScene {
		// 戦闘設定の取得と初期化
		battleSetting := getBattleSetting(option.BattleSettingId)
		enemyIds := extractEnemyIds(battleSetting)

		// 戦闘の初期化
		initializeRequest := &core.InitializeBattleRequest{
			// TODO: variables must be provided by args
			MainActorCharacterId: core.CharacterLuneId,
			SubActorCharacterId:  core.CharacterSunnyId,
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

		// UIコンポーネントの初期化
		messageWindow := createMessageWindow(newMessageWindow)
		battleEnemyDisplay := createBattleEnemyDisplay(newBattleEnemyDisplay, battleResponse.EnemyIds, battleSetting.Enemies)

		// バトル選択ウィンドウの設定
		input := &frontend.KeyBoardInput{}
		var selectedCommand core.PlayerCommand
		var targetSelectWindow *component.SelectWindow
		onSubmit := func(command core.PlayerCommand) {
			selectedCommand = command
			targetSelectWindow.Open()
			input.Set(targetSelectWindow)
		}

		actorDisplay := newBattleActorDisplay(mainActor)
		subActorDisplay := newBattleSubActorDisplay(subActor)
		subActorDialog := component.CreateNewBattlePartnerDialogue(newMessageWindow)()

		// エフェクト関数の定義
		playEffect := createPlayEffectFunction(serveActor, battleEnemyDisplay, subActorDisplay, actorDisplay, effectManager)
		shake := frontend.NewShake()
		doShake := createDoShakeFunction(serveActor, battleEnemyDisplay, subActorDisplay, shake)
		setDamage := createSetDamageFunction(serveActor, battleEnemyDisplay, subActorDisplay, actorDisplay)
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

		battleSelectWindow := createBattleSelectWindow(newBattleSelectWindow, input, onSubmit)

		// パートナーダイアログの設定
		setPartnerDialogue := createSetPartnerDialogueFunction(subActorDialog)
		createSequence := produceCreateSequence(
			setPartnerDialogue,
			setEmotion,
			subActorDialog.Open,
			subActorDialog.Close,
		)
		seq := createSequence("test_sequence_0000")

		// バトルシーンの作成
		battleScene := &BattleScene{
			ui: battleUI{
				messageWindow:      messageWindow,
				battleSelectWindow: battleSelectWindow,
				actorDisplay:       actorDisplay,
				subActorDisplay:    subActorDisplay,
				subActorDialog:     subActorDialog,
				input:              input,
				battleEnemyDisplay: battleEnemyDisplay,
				effectManager:      effectManager,
				shake:              shake,
			},
			battleSequence: component.NewBattleEventSequencer(),
			enemyData:      battleResponse.EnemyIds,
			actorNames:     actorNames,
			createSequence: createSequence,
			sequences:      sequence.CreateSequenceManager(),
		}
		battleScene.sequences.AddSequence(seq)

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
		closeWindowOnTargetSelect := createOnSubmitTargetSelect(battleSelectWindow, input)
		onTargetSelect := createOnTargetSelect(
			closeWindowOnTargetSelect,
			func(index int) core.ActorId { return allActorId[index] },
			func() core.PlayerCommand { return selectedCommand },
			battleScene.battleSequence.Reset,
			playSequence,
			processBattle,
		)
		targetSelectWindow = newSelectWindow(
			&frontend.Vector{X: 80, Y: 0},
			frontend.PivotBottomLeft,
			frontend.DepthWindow,
			allTextId,
			onTargetSelect,
			true,
		)
		battleScene.ui.targetSelectWindow = targetSelectWindow

		return battleScene
	}
}

func extractEnemyIds(battleSetting *core.BattleSetting) []core.EnemyId {
	ids := make([]core.EnemyId, len(battleSetting.Enemies))
	for i, set := range battleSetting.Enemies {
		ids[i] = set.EnemyId
	}
	return ids
}

func createActorIdToEnemyMapping(enemyIds []*core.EnemyIdPair) map[core.ActorId]core.EnemyId {
	result := make(map[core.ActorId]core.EnemyId)
	for _, pair := range enemyIds {
		result[pair.ActorId] = pair.EnemyId
	}
	return result
}

func createActorNamesMapping(enemyIds []*core.EnemyIdPair, serveEnemyName core.EnemyNameServer) map[core.ActorId]core.TextId {
	names := make(map[core.ActorId]core.TextId)
	names[core.ActorLuneId] = core.TextIdLuneName
	names[core.ActorSunnyId] = core.TextIdSunnyName
	for _, pair := range enemyIds {
		names[pair.ActorId] = serveEnemyName(pair.EnemyId)
	}
	return names
}

func createAllActorIdList(battleResponse *core.InitializeBattleResponse) []core.ActorId {
	ids := []core.ActorId{battleResponse.MainActorId}
	if battleResponse.SubActorId != core.ActorEmptyId {
		ids = append(ids, battleResponse.SubActorId)
	}
	enemyActorIds := make([]core.ActorId, len(battleResponse.EnemyIds))
	for i, pair := range battleResponse.EnemyIds {
		enemyActorIds[i] = pair.ActorId
	}
	return append(ids, enemyActorIds...)
}

func createAllTextIdList(allActorId []core.ActorId, actorNames map[core.ActorId]core.TextId) []core.TextId {
	texts := make([]core.TextId, 0)
	for _, id := range allActorId {
		texts = append(texts, actorNames[id])
	}
	return texts
}

func createMessageWindow(newMessageWindow component.NewMessageWindowFunc) *component.MessageWindow {
	messageWindow := newMessageWindow(
		&frontend.Vector{X: 192, Y: 0},
		&frontend.Vector{X: 292, Y: 62},
		frontend.DepthWindow,
		frontend.PivotTopCenter,
	)
	testString := "あのイーハトーヴォのすきとおった風\n夏でも底に冷たさをもつ青いそら\nうつくしい森で飾られたモリーオ市"
	messageWindow.SetText(testString, false)
	messageWindow.Open()
	return messageWindow
}

func createBattleEnemyDisplay(newBattleEnemyDisplay component.NewBattleEnemyDisplayFunc, enemyIds []*core.EnemyIdPair, enemySettings []*core.EnemySetting) *component.BattleEnemyDisplay {
	displayArgs := component.ToDisplayArgs(enemyIds, enemySettings)
	return newBattleEnemyDisplay(
		displayArgs,
		frontend.DepthEnemy,
	)
}

func createBattleSelectWindow(
	newBattleSelectWindow component.NewBattleSelectWindowFunc,
	input frontend.InputManager,
	onSubmit func(core.PlayerCommand),
) *component.BattleSelectWindow {
	battleSelectWindow := newBattleSelectWindow(
		&frontend.Vector{X: 0, Y: 0},
		frontend.PivotBottomLeft,
		frontend.DepthWindow,
		[]core.PlayerCommand{
			core.PlayerCommandAttack,
			core.PlayerCommandFire,
			core.PlayerCommandBarrier,
			core.PlayerCommandThunder,
			core.PlayerCommandWind,
			core.PlayerCommandFocus,
			core.PlayerCommandDefend,
		},
		onSubmit,
	)
	input.Set(battleSelectWindow)
	battleSelectWindow.Open()
	return battleSelectWindow
}

func createPlayEffectFunction(serveActor core.ActorSupplier, battleEnemyDisplay *component.BattleEnemyDisplay, subActorDisplay *component.BattleSubActorDisplay, actorDisplay *component.BattleActorDisplay, effectManager *widget.EffectManager) func(widget.EffectId, core.ActorId) {
	return func(effectId widget.EffectId, target core.ActorId) {
		actor := serveActor(target)
		position := func() *frontend.Vector {
			if actor.Side == core.ActorSideEnemy {
				return battleEnemyDisplay.GetPosition(target)
			}
			if actor.Id == core.ActorSunnyId {
				return subActorDisplay.GetCenterPosition()
			}
			return actorDisplay.GetMainCharacterPosition()
		}()
		effectManager.CallEffect(effectId, position)
	}
}

func createDoShakeFunction(serveActor core.ActorSupplier, battleEnemyDisplay *component.BattleEnemyDisplay, subActorDisplay *component.BattleSubActorDisplay, shake *frontend.EmitShake) func(core.ActorId) {
	return func(actorId core.ActorId) {
		actor := serveActor(actorId)
		if actor.IsEnemy() {
			battleEnemyDisplay.DoShake(actorId)
			return
		}
		if actor.IsSubActor() {
			subActorDisplay.Shake()
			return
		}
		// メインキャラクターのシェイク処理
		shake.Shake(frontend.ShakeDefaultAmplitude, frontend.ShakeDefaultPeriod)
	}
}

func createSetDamageFunction(serveActor core.ActorSupplier, battleEnemyDisplay *component.BattleEnemyDisplay, subActorDisplay *component.BattleSubActorDisplay, actorDisplay *component.BattleActorDisplay) func(core.ActorId, core.Damage, core.HP) {
	return func(actorId core.ActorId, damage core.Damage, afterHp core.HP) {
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

func createSetEmotionFunction(serveActor core.ActorSupplier, battleEnemyDisplay *component.BattleEnemyDisplay, subActorDisplay *component.BattleSubActorDisplay, actorDisplay *component.BattleActorDisplay) func(core.ActorId, component.BattleEmotionType) {
	return func(actorId core.ActorId, emotion component.BattleEmotionType) {
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

func createSetPartnerDialogueFunction(subActorDialog *component.BattlePartnerDialogue) func(core.TextString) *sequence.SetPartnerDialogueResponse {
	return func(text core.TextString) *sequence.SetPartnerDialogueResponse {
		subActorDialog.SetText(text.String(), false)
		return &sequence.SetPartnerDialogueResponse{
			CheckIsEnd: func() sequence.IsEnd {
				result := subActorDialog.IsTextEnd()
				return sequence.IsEnd(result)
			},
		}
	}
}

func createOnSubmitTargetSelect(
	battleSelectWindow *component.BattleSelectWindow,
	input frontend.InputManager,
) func() {
	return func() {
		battleSelectWindow.Close()
		input.Set(frontend.InputReceiverEmptyInstance)
	}
}
