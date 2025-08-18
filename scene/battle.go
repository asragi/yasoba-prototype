package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleScene struct {
	messageWindow      *component.MessageWindow
	battleSelectWindow *component.BattleSelectWindow
	actorDisplay       *component.BattleActorDisplay
	subActorDisplay    *component.BattleSubActorDisplay
	subActorDialog     *component.BattlePartnerDialogue
	enemyData          []*core.EnemyIdPair
	actorNames         map[core.ActorId]core.TextId
	targetSelectWindow *component.SelectWindow
	input              frontend.InputManager
	battleSequence     *component.BattleEventSequencer
	battleEnemyDisplay *component.BattleEnemyDisplay
	effectManager      *widget.EffectManager
	shake              *frontend.EmitShake
	endState           core.BattleEndType
}

func (s *BattleScene) OnSequenceEnd() {
	if s.endState == core.BattleEndTypeWin {
		s.messageWindow.SetText("しょうりした！", false)
		return
	}
	if s.endState == core.BattleEndTypeLose {
		s.messageWindow.SetText("やられてしまった……", false)
		return
	}
	s.input.Set(s.battleSelectWindow)
	s.battleSelectWindow.Open()
}

func (s *BattleScene) Update() {
	s.shake.Update()
	delta := s.shake.Delta()
	zeroVector := frontend.VectorZero
	zeroVector = zeroVector.Add(delta)
	bottomLeft := &frontend.Vector{X: 0, Y: 288}
	bottomLeft = bottomLeft.Add(delta)
	bottomRight := &frontend.Vector{X: 384, Y: 288}
	bottomRight = bottomRight.Add(delta)
	center := &frontend.Vector{X: 192, Y: 144}
	center = center.Add(delta)
	mainCharacterTopLeftPosition := s.actorDisplay.GetMainCharacterTopLeftPosition()
	mainCharacterTopLeftPosition = mainCharacterTopLeftPosition.Add(delta)
	s.messageWindow.Update(zeroVector)
	s.actorDisplay.Update(bottomLeft)
	s.subActorDisplay.Update(bottomRight)
	s.subActorDialog.Update(s.subActorDisplay.GetTopCenterPosition())
	s.battleEnemyDisplay.Update(center)
	s.battleSelectWindow.Update(mainCharacterTopLeftPosition)
	s.targetSelectWindow.Update(mainCharacterTopLeftPosition)
	s.input.Update()
	if s.battleSequence.IsRun() {
		s.battleSequence.Update()
		if s.battleSequence.IsEnd() {
			s.OnSequenceEnd()
		}
	}
	s.effectManager.Update()
}

func (s *BattleScene) Draw(drawFunc frontend.DrawFunc) {
	s.messageWindow.Draw(drawFunc)
	s.battleSelectWindow.Draw(drawFunc)
	s.targetSelectWindow.Draw(drawFunc)
	s.battleEnemyDisplay.Draw(drawFunc)
	s.subActorDisplay.Draw(drawFunc)
	s.subActorDialog.Draw(drawFunc)
	s.actorDisplay.Draw(drawFunc)
	s.effectManager.Draw(drawFunc)
}

type BattleResult struct{}
type OnEndBattle func(BattleResult)

type BattleOption struct {
	OnEnd           OnEndBattle
	BattleSettingId core.BattleSettingId
}

type NewBattleScene func(*BattleOption) *BattleScene

func StandByNewBattleScene(
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
) NewBattleScene {
	return func(option *BattleOption) *BattleScene {
		onEnd := func(endType core.BattleEndType) {
			// TODO: Implement end of battle
			option.OnEnd(BattleResult{})
		}
		battleSetting := getBattleSetting(option.BattleSettingId)
		enemyIds := func() []core.EnemyId {
			ids := make([]core.EnemyId, len(battleSetting.Enemies))
			for i, set := range battleSetting.Enemies {
				ids[i] = set.EnemyId
			}
			return ids
		}()
		subCharacterId := core.CharacterSunnyId
		initializeRequest := &core.InitializeBattleRequest{
			// TODO: variables must be provided by args
			MainActorCharacterId: core.CharacterLuneId,
			SubActorCharacterId:  subCharacterId,
			EnemyIds:             enemyIds,
		}
		battleResponse := initializeBattle(initializeRequest)
		processBattle := newProcessBattle(battleResponse, onEnd)
		mainActorId := battleResponse.MainActorId
		mainActor := serveActor(mainActorId)
		subActorId := battleResponse.SubActorId
		subActor := serveActor(subActorId)
		actorIdToEnemy := func() map[core.ActorId]core.EnemyId {
			result := make(map[core.ActorId]core.EnemyId)
			for _, pair := range battleResponse.EnemyIds {
				result[pair.ActorId] = pair.EnemyId
			}
			return result
		}()
		actorNames := func() map[core.ActorId]core.TextId {
			names := make(map[core.ActorId]core.TextId)
			names[core.ActorLuneId] = core.TextIdLuneName
			names[core.ActorSunnyId] = core.TextIdSunnyName
			for _, pair := range battleResponse.EnemyIds {
				names[pair.ActorId] = serveEnemyName(pair.EnemyId)
			}
			return names
		}()
		allActorId := func() []core.ActorId {
			ids := []core.ActorId{
				battleResponse.MainActorId,
			}
			if battleResponse.SubActorId != core.ActorEmptyId {
				ids = append(ids, battleResponse.SubActorId)
			}
			enemyActorIds := func() []core.ActorId {
				result := make([]core.ActorId, len(battleResponse.EnemyIds))
				for i, pair := range battleResponse.EnemyIds {
					result[i] = pair.ActorId
				}
				return result
			}()
			return append(ids, enemyActorIds...)
		}()
		allTextId := func() []core.TextId {
			texts := make([]core.TextId, 0)
			for _, id := range allActorId {
				texts = append(texts, actorNames[id])
			}
			return texts
		}()

		messageWindow := newMessageWindow(
			&frontend.Vector{X: 192, Y: 0},
			&frontend.Vector{X: 292, Y: 62},
			frontend.DepthWindow,
			frontend.PivotTopCenter,
		)
		testString := "あのイーハトーヴォのすきとおった風\n夏でも底に冷たさをもつ青いそら\nうつくしい森で飾られたモリーオ市"
		messageWindow.SetText(testString, false)

		displayArgs := component.ToDisplayArgs(battleResponse.EnemyIds, battleSetting.Enemies)
		battleEnemyDisplay := newBattleEnemyDisplay(
			displayArgs,
			frontend.DepthEnemy,
		)

		var battleSelectWindow *component.BattleSelectWindow
		var selectWindow *component.SelectWindow
		var selectedCommand core.PlayerCommand

		input := &frontend.KeyBoardInput{}

		onSubmit := func(command core.PlayerCommand) {
			selectedCommand = command
			selectWindow.Open()
			input.Set(selectWindow)
		}

		battleSelectWindow = newBattleSelectWindow(
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

		actorDisplay := newBattleActorDisplay(mainActor)
		subActorDisplay := newBattleSubActorDisplay(subActor)
		subActorDialog := component.CreateNewBattlePartnerDialogue(newMessageWindow)()
		subActorDialog.Open()
		subActorDialog.SetText("こんにちは！！\nこれは改行テストだよ！\n３行くらいまでの表示を想定しているよ", false)

		playEffect := func(effectId widget.EffectId, target core.ActorId) {
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

		var battleScene *BattleScene
		doShake := func(actorId core.ActorId) {
			actor := serveActor(actorId)
			if actor.IsEnemy() {
				battleEnemyDisplay.DoShake(actorId)
				return
			}
			if actor.IsSubActor() {
				subActorDisplay.Shake()
				return
			}
			battleScene.shake.Shake(frontend.ShakeDefaultAmplitude, frontend.ShakeDefaultPeriod)
		}

		setDamage := func(actorId core.ActorId, damage core.Damage, afterHp core.HP) {
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

		setEmotion := func(actorId core.ActorId, emotion component.BattleEmotionType) {
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

		// TODO: Expand these functions to handle player actor
		newBattleSequence := createNewBattleSequence(
			messageWindow,
			doShake,
			setEmotion,
			setDamage,
			playEffect,
			battleEnemyDisplay.SetDisappear,
		)

		battleScene = &BattleScene{
			messageWindow:      messageWindow,
			battleSelectWindow: battleSelectWindow,
			actorDisplay:       actorDisplay,
			enemyData:          battleResponse.EnemyIds,
			actorNames:         actorNames,
			input:              input,
			battleSequence:     component.NewBattleEventSequencer(),
			battleEnemyDisplay: battleEnemyDisplay,
			subActorDisplay:    subActorDisplay,
			effectManager:      effectManager,
			shake:              frontend.NewShake(),
			subActorDialog:     subActorDialog,
		}

		playSequence := createPlayBattleSequence(
			skillToSequence,
			newBattleSequence,
			battleScene.battleSequence.Add,
			actorIdToEnemy,
			serveEnemyView,
		)
		closeWindowOnTargetSelect := func() {
			battleSelectWindow.Close()
			selectWindow.Close()
			input.Set(frontend.InputReceiverEmptyInstance)
		}
		onTargetSelect := createOnTargetSelect(
			closeWindowOnTargetSelect,
			func(index int) core.ActorId { return allActorId[index] },
			func() core.PlayerCommand { return selectedCommand },
			battleScene.battleSequence.Reset,
			playSequence,
			processBattle,
		)
		selectWindow = newSelectWindow(
			&frontend.Vector{X: 80, Y: 0},
			frontend.PivotBottomLeft,
			frontend.DepthWindow,
			allTextId,
			onTargetSelect,
		)
		battleScene.targetSelectWindow = selectWindow
		return battleScene
	}
}

func createPlayBattleSequence(
	skillToSequence component.SkillToSequenceFunc,
	newBattleSequence component.NewBattleSequenceFunc,
	addBattleSequence func(component.BattleSequenceFunc),
	actorIdToEnemy map[core.ActorId]core.EnemyId,
	serveEnemyView component.ServeEnemyViewData,
) func([]*core.SkillApplyResult) {
	return func(skillApplyResultSet []*core.SkillApplyResult) {
		for _, skillApplyResult := range skillApplyResultSet {
			skillId := skillApplyResult.SkillId
			sequenceId := skillToSequence(skillId)
			damageInformation := func() []*component.DamageInformation {
				result := make([]*component.DamageInformation, 0)
				for _, row := range skillApplyResult.Rows {
					result = append(
						result, &component.DamageInformation{
							Target:  row.TargetId,
							Damage:  row.Damage,
							AfterHP: row.AfterHp,
						},
					)
				}
				return result
			}()
			sequence := newBattleSequence(
				&component.EventSequenceArgs{
					SequenceId: sequenceId,
					Actor:      skillApplyResult.Actor,
					Target:     damageInformation,
				},
			)
			addBattleSequence(sequence)
			for _, row := range skillApplyResult.Rows {
				if !row.IsTargetBeaten {
					continue
				}
				actualTarget := row.TargetId
				targetSide := row.TargetSide
				if targetSide == core.ActorSideEnemy {
					enemyId := actorIdToEnemy[actualTarget]
					viewData := serveEnemyView(enemyId)
					beatenSequence := newBattleSequence(
						&component.EventSequenceArgs{
							SequenceId: viewData.BeatenSequenceId,
							Actor:      actualTarget,
							Target: []*component.DamageInformation{
								{
									Target: actualTarget,
									Damage: 0,
								},
							},
						},
					)
					addBattleSequence(beatenSequence)
					continue
				}
				// TODO: Implement player beaten sequence
			}
		}
	}
}

// TODO: View非依存のLogic部分だけ抽出してCoreに移動したい
func createOnTargetSelect(
	closeWindow func(),
	indexToActor func(int) core.ActorId,
	serveSelectedCommand func() core.PlayerCommand,
	resetBattleSequence func(),
	playSequence func([]*core.SkillApplyResult),
	processBattle core.ProcessBattleFunc,
) func(int) {
	return func(index int) {
		closeWindow()
		target := indexToActor(index)
		command := serveSelectedCommand()
		response := processBattle(
			&core.ProcessBattleRequest{
				TargetId: []core.ActorId{target},
				Command:  command,
			},
		)

		resetBattleSequence()
		playSequence(response.SkillApplyResults)
	}
}
