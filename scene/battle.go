package scene

import (
	"fmt"
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type BattleScene struct {
	messageWindow      *component.MessageWindow
	battleSelectWindow *component.BattleSelectWindow
	faceWindow         *component.FaceWindow
	faceSubWindow      *component.FaceWindow
	enemyData          []*core.EnemyIdPair
	actorNames         map[core.ActorId]core.TextId
	targetSelectWindow *component.SelectWindow
	input              frontend.InputManager
	battleSequence     component.BattleSequenceFunc
	battleActorDisplay *component.BattleActorDisplay
	effectManager      *widget.EffectManager
}

func (s *BattleScene) OnSequenceEnd() {
	s.input.Set(s.battleSelectWindow)
	s.battleSelectWindow.Open()
}

func (s *BattleScene) Update() {
	s.messageWindow.Update(frontend.VectorZero)
	s.faceWindow.Update(&frontend.Vector{X: 0, Y: 288})
	s.battleActorDisplay.Update(&frontend.Vector{X: 192, Y: 144})
	s.faceSubWindow.Update(&frontend.Vector{X: 384, Y: 288})
	s.battleSelectWindow.Update(s.faceWindow.GetTopLeftPosition())
	s.targetSelectWindow.Update(s.faceWindow.GetTopLeftPosition())
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.messageWindow.Shake(2, 10)
	}
	s.input.Update()
	if s.battleSequence != nil {
		result := s.battleSequence()
		if result.IsEnd {
			s.battleSequence = nil
			s.OnSequenceEnd()
		}
	}
	s.effectManager.Update()
}

func (s *BattleScene) Draw(drawFunc frontend.DrawFunc) {
	s.messageWindow.Draw(drawFunc)
	s.battleSelectWindow.Draw(drawFunc)
	s.targetSelectWindow.Draw(drawFunc)
	s.battleActorDisplay.Draw(drawFunc)
	s.faceWindow.Draw(drawFunc)
	s.faceSubWindow.Draw(drawFunc)
	s.effectManager.Draw(drawFunc)
}

type BattleResult struct{}
type OnEndBattle func(BattleResult)

type BattleOption struct {
	OnEnd           OnEndBattle
	BattleSettingId game.BattleSettingId
}

type NewBattleScene func(*BattleOption) *BattleScene

func StandByNewBattleScene(
	newMessageWindow component.NewMessageWindowFunc,
	newSelectWindow component.NewSelectWindowFunc,
	newBattleSelectWindow component.NewBattleSelectWindowFunc,
	newFaceWindow component.NewFaceWindowFunc,
	serveEnemyName core.EnemyNameServer,
	initializeBattle core.InitializeBattleFunc,
	postCommand core.PostCommandFunc,
	skillApply core.SkillApplyFunc,
	getBattleSetting game.ServeBattleSetting,
	createNewBattleSequence component.PrepareBattleEventSequenceFunc,
	skillToSequence component.SkillToSequenceFunc,
	newBattleActorDisplay component.NewBattleActorDisplayFunc,
	effectManager *widget.EffectManager,
) NewBattleScene {
	return func(option *BattleOption) *BattleScene {
		battleSetting := getBattleSetting(option.BattleSettingId)
		enemyIds := func() []core.EnemyId {
			ids := make([]core.EnemyId, len(battleSetting.Enemies))
			for i, set := range battleSetting.Enemies {
				ids[i] = set.EnemyId
			}
			return ids
		}()
		initializeRequest := &core.InitializeBattleRequest{
			// TODO: variables must be provided by args
			MainActorCharacterId: core.CharacterLuneId,
			SubActorCharacterId:  core.CharacterSunnyId,
			EnemyIds:             enemyIds,
		}
		battleResponse := initializeBattle(initializeRequest)
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

		displayArgs := func() []*component.BattleDisplayArgs {
			result := make([]*component.BattleDisplayArgs, len(battleResponse.EnemyIds))
			mappedList := make(map[core.EnemyId][]core.ActorId)
			mappedSettingList := make(map[core.EnemyId][]*game.EnemySetting)
			for _, pair := range battleResponse.EnemyIds {
				mappedList[pair.EnemyId] = append(mappedList[pair.EnemyId], pair.ActorId)
			}
			for _, set := range battleSetting.Enemies {
				mappedSettingList[set.EnemyId] = append(mappedSettingList[set.EnemyId], set)
			}
			index := 0
			for _, pair := range battleResponse.EnemyIds {
				enemyId := pair.EnemyId
				actors := mappedList[enemyId]
				setting := mappedSettingList[enemyId]
				for j, actorId := range actors {
					result[index] = &component.BattleDisplayArgs{
						ActorId:  actorId,
						EnemyId:  enemyId,
						Position: setting[j].Position,
					}
					index++
				}
			}
			return result
		}()
		battleActorDisplay := newBattleActorDisplay(
			displayArgs,
			frontend.DepthEnemy,
		)

		playEffect := func(effectId widget.EffectId, target core.ActorId) {
			position := battleActorDisplay.GetPosition(target)
			effectManager.CallEffect(effectId, position)
		}

		newBattleSequence := createNewBattleSequence(
			messageWindow,
			battleActorDisplay.DoShake,
			battleActorDisplay.SetEmotion,
			func(id core.ActorId, damage core.Damage) {
				fmt.Printf("actor: %s, damage: %d\n", id, damage)
			},
			playEffect,
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
				core.PlayerCommandFocus,
				core.PlayerCommandDefend,
			},
			onSubmit,
		)
		input.Set(battleSelectWindow)
		battleSelectWindow.Open()

		faceWindow := newFaceWindow(
			&frontend.Vector{X: 0, Y: 0},
			frontend.DepthWindow,
			frontend.PivotBottomLeft,
			frontend.TextureFaceLuneNormal,
		)
		faceSubWindow := newFaceWindow(
			&frontend.Vector{X: 0, Y: 0},
			frontend.DepthWindow,
			frontend.PivotBottomRight,
			frontend.TextureFaceSunnyNormal,
		)

		battleScene := &BattleScene{
			messageWindow:      messageWindow,
			battleSelectWindow: battleSelectWindow,
			faceWindow:         faceWindow,
			faceSubWindow:      faceSubWindow,
			enemyData:          battleResponse.EnemyIds,
			actorNames:         actorNames,
			input:              input,
			effectManager:      effectManager,
			battleActorDisplay: battleActorDisplay,
		}

		onTargetSelect := func(index int) {
			target := allActorId[index]
			response := postCommand(
				&core.PostCommandRequest{
					ActorId:  core.ActorLuneId,
					TargetId: []core.ActorId{target},
					Command:  selectedCommand,
				},
			)
			skillId := response.Actions.Id
			appliedResponse := skillApply(response.Actions)
			battleSelectWindow.Close()
			selectWindow.Close()
			input.Set(frontend.InputReceiverEmptyInstance)
			sequenceId := skillToSequence(skillId)
			damageInformation := func() []*component.DamageInformation {
				result := make([]*component.DamageInformation, 0)
				for _, row := range appliedResponse.Rows {
					switch r := row.(type) {
					case *core.SkillSingleAttackResult:
						result = append(
							result, &component.DamageInformation{
								Target: r.TargetId,
								Damage: r.Damage,
							},
						)
					}
				}
				return result
			}()
			sequence := newBattleSequence(
				&component.EventSequenceArgs{
					SequenceId: sequenceId,
					Actor:      response.Actions.Actor,
					Target:     damageInformation,
				},
			)
			battleScene.battleSequence = sequence
		}
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
