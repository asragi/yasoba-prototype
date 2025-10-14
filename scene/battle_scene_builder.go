package scene

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle_skill"
	"github.com/asragi/yasoba-prototype/battleconfig"
	"github.com/asragi/yasoba-prototype/battlesetup"
	"github.com/asragi/yasoba-prototype/characterdata"
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/enemydata"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/invoke"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/widget"
)

func InitializeCreateBattleScene(
	newMessageWindow component.NewMessageWindowFunc,
	newSelectWindow component.NewSelectWindowFunc,
	newBattleSelectWindow component.NewBattleSelectWindowFunc,
	newBattleActorDisplay component.NewBattleActorDisplayFunc,
	newBattleSubActorDisplay component.NewBattleSubActorDisplayFunc,
	serveEnemyName enemydata.NameServer,
	initializeBattle battle.InitializeBattleFunc,
	getBattleSetting battleconfig.ServeFunc,
	createNewBattleSequence component.PrepareBattleEventSequenceFunc,
	skillToSequence component.SkillToSequenceFunc,
	newBattleEnemyDisplay component.NewBattleEnemyDisplayFunc,
	effectManager *widget.EffectManager,
	serveEnemyView component.ServeEnemyViewData,
	serveActor actor.ActorSupplier,
	newVariableMessageWindow component.NewVariableMessageWindowFunc,
	newProcessBattle battle.NewProcessBattleFunc,
	produceCreateSequence sequence.ProduceCreateSequence,
	produceCheckInvokeSequence invoke.ProduceCheckInvokeSequence,
) CreateBattleScene {
	return func(option *BattleOption) *BattleScene {
		// 戦闘設定の取得と初期化
		battleSetting := getBattleSetting(option.BattleSettingId)
		enemyIds := extractEnemyIds(battleSetting)

		// 戦闘の初期化
		initializeRequest := &battle.InitializeBattleRequest{
			// TODO: variables must be provided by args
			MainActorCharacterId: characterdata.CharacterLuneId,
			SubActorCharacterId:  characterdata.CharacterSunnyId,
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

		displayedHp := func() map[actor.ActorId]actor.HP {
			hp := make(map[actor.ActorId]actor.HP, len(allActorId))
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
		battleEnemyDisplay := createBattleEnemyDisplay(newBattleEnemyDisplay, battleResponse.EnemyIds, battleSetting.Enemies)

		// バトル選択ウィンドウの設定
		input := &frontend.KeyBoardInput{}
		var selectedCommand battle.PlayerCommand
		var targetSelectWindow *component.SelectWindow
		onSubmit := func(command battle.PlayerCommand) {
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
			func(index int) actor.ActorId { return allActorId[index] },
			func() battle.PlayerCommand { return selectedCommand },
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

		getActorHpRatio := func(actorId actor.ActorId) actor.HPRatio {
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

func extractEnemyIds(battleSetting *battleconfig.Setting) []enemydata.EnemyId {
	ids := make([]enemydata.EnemyId, len(battleSetting.Enemies))
	for i, set := range battleSetting.Enemies {
		ids[i] = set.EnemyId
	}
	return ids
}

func createActorIdToEnemyMapping(enemyIds []*battlesetup.EnemyIdPair) map[actor.ActorId]enemydata.EnemyId {
	result := make(map[actor.ActorId]enemydata.EnemyId)
	for _, pair := range enemyIds {
		result[pair.ActorId] = pair.EnemyId
	}
	return result
}

func createActorNamesMapping(enemyIds []*battlesetup.EnemyIdPair, serveEnemyName enemydata.NameServer) map[actor.ActorId]text.TextId {
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

func createBattleEnemyDisplay(newBattleEnemyDisplay component.NewBattleEnemyDisplayFunc, enemyIds []*battlesetup.EnemyIdPair, enemySettings []*battleconfig.EnemySetting) *component.BattleEnemyDisplay {
	displayArgs := component.ToDisplayArgs(enemyIds, enemySettings)
	return newBattleEnemyDisplay(
		displayArgs,
		frontend.DepthEnemy,
	)
}

func createBattleSelectWindow(
	newBattleSelectWindow component.NewBattleSelectWindowFunc,
	input frontend.InputManager,
	onSubmit func(battle.PlayerCommand),
) *component.BattleSelectWindow {
	battleSelectWindow := newBattleSelectWindow(
		&frontend.Vector{X: 0, Y: 0},
		frontend.PivotBottomLeft,
		frontend.DepthWindow,
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
	input.Set(battleSelectWindow)
	battleSelectWindow.Open()
	return battleSelectWindow
}

func createPlayEffectFunction(serveActor actor.ActorSupplier, battleEnemyDisplay *component.BattleEnemyDisplay, subActorDisplay *component.BattleSubActorDisplay, actorDisplay *component.BattleActorDisplay, effectManager *widget.EffectManager) func(widget.EffectId, actor.ActorId) {
	return func(effectId widget.EffectId, target actor.ActorId) {
		targetActor := serveActor(target)
		position := func() *frontend.Vector {
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

func createDoShakeFunction(serveActor actor.ActorSupplier, battleEnemyDisplay *component.BattleEnemyDisplay, subActorDisplay *component.BattleSubActorDisplay, shake *frontend.EmitShake) func(actor.ActorId) {
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
		// メインキャラクターのシェイク処理
		shake.Shake(frontend.ShakeDefaultAmplitude, frontend.ShakeDefaultPeriod)
	}
}

func createSetDamageFunction(
	serveActor actor.ActorSupplier,
	battleEnemyDisplay *component.BattleEnemyDisplay,
	subActorDisplay *component.BattleSubActorDisplay,
	actorDisplay *component.BattleActorDisplay,
	displayedHp map[actor.ActorId]actor.HP,
) func(actor.ActorId, battle_skill.Damage, actor.HP) {
	return func(actorId actor.ActorId, damage battle_skill.Damage, afterHp actor.HP) {
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

func createSetEmotionFunction(serveActor actor.ActorSupplier, battleEnemyDisplay *component.BattleEnemyDisplay, subActorDisplay *component.BattleSubActorDisplay, actorDisplay *component.BattleActorDisplay) func(actor.ActorId, component.BattleEmotionType) {
	return func(actorId actor.ActorId, emotion component.BattleEmotionType) {
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

func createSetPartnerDialogueFunction(subActorDialog *component.BattlePartnerDialogue) func(text.String) *sequence.SetPartnerDialogueResponse {
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
	battleSelectWindow *component.BattleSelectWindow,
	input frontend.InputManager,
) func() {
	return func() {
		battleSelectWindow.Close()
		input.Set(frontend.InputReceiverEmptyInstance)
	}
}
