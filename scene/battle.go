package scene

import (
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/battle/enemy"
	"github.com/asragi/yasoba-prototype/battle/invoke"
	"github.com/asragi/yasoba-prototype/battle/setup"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	battleenemy "github.com/asragi/yasoba-prototype/view/battle/enemy"
	"github.com/asragi/yasoba-prototype/view/battle/event"
)

type BattleScene struct {
	ui                  battleUI
	battleSequence      *event.BattleEventSequencer
	enemyData           []*setup.EnemyIdPair
	actorNames          map[actor.ActorId]text.TextId
	endState            battle.BattleEndType
	createSequence      sequence.CreateSequence
	sequences           *sequence.SequenceManager
	checkInvokeSequence invoke.CheckInvokeSequence
	invokedSequences    map[sequence.SequenceId]bool
	turnCount           battle.TurnCount
	playerBeaten        bool
	autoPlayNonPlayer   func()
}

func (s *BattleScene) onTurnEnd() {
	s.turnCount++
	if s.endState == battle.BattleEndTypeWin {
		s.sequences.AddSequence(s.createSequence("test_sequence_0100"))
		return
	}
	if s.endState == battle.BattleEndTypeLose {
		s.sequences.AddSequence(s.createSequence("battle_sequence_lose"))
		return
	}
	if s.playerBeaten {
		if s.autoPlayNonPlayer == nil {
			panic("auto play func is required when player is beaten")
		}
		s.autoPlayNonPlayer()
		s.checkAndStartSequences(invoke.InvokeTimingStartTurn)
		return
	}
	s.ui.input.Set(s.ui.battleSelectWindow)
	s.ui.battleSelectWindow.Open()
	s.checkAndStartSequences(invoke.InvokeTimingStartTurn)
}

func (s *BattleScene) onBattleEnd(endType battle.BattleEndType) {
	s.endState = endType
}

func (s *BattleScene) Update() {
	s.ui.Update()
	if advanceBattleSequence(s.battleSequence) {
		s.onTurnEnd()
	}
	s.sequences.Update()
	s.checkAndStartSequences(invoke.InvokeTimingEveryAction)
}

func (s *BattleScene) Draw(drawFunc drawing.DrawFunc) {
	s.ui.Draw(drawFunc)
}

func (s *BattleScene) checkAndStartSequences(timing invoke.InvokeTiming) {
	ids := s.checkInvokeSequence(timing, s.turnCount)
	for _, id := range ids {
		if s.invokedSequences[id] {
			continue
		}
		sequence := s.createSequence(id)
		s.sequences.AddSequence(sequence)
		s.invokedSequences[id] = true
	}
}

func advanceBattleSequence(sequence *event.BattleEventSequencer) bool {
	if !sequence.IsRun() {
		return false
	}
	sequence.Update()
	return sequence.IsEnd()
}

type BattleResult struct{}
type OnEndBattle func(BattleResult)

type BattleOption struct {
	OnEnd           OnEndBattle
	BattleSettingId config.Id
	BattleId        battle.BattleId
	SwitchToBattle  func(battle.BattleId)
	SwitchToDebug   func()
}

type CreateBattleScene func(*BattleOption) *BattleScene

type playBattleSequenceFunc func([]*skill.SkillApplyResult)

// SkillApplyResultに基づいて戦闘の演出を行う
func createPlayBattleSequence(
	skillToSequence event.SkillToSequenceFunc,
	newBattleSequence event.NewBattleSequenceFunc,
	addBattleSequence func(event.BattleSequenceFunc),
	actorIdToEnemy map[actor.ActorId]enemy.EnemyId,
	serveEnemyView battleenemy.ServeEnemyViewData,
) playBattleSequenceFunc {
	return func(skillApplyResultSet []*skill.SkillApplyResult) {
		for _, skillApplyResult := range skillApplyResultSet {
			skillId := skillApplyResult.SkillId
			sequenceId := skillToSequence(skillId)
			damageInformation := func() []*event.DamageInformation {
				result := make([]*event.DamageInformation, 0)
				for _, row := range skillApplyResult.Rows {
					result = append(
						result, &event.DamageInformation{
							Target:  row.TargetId,
							Damage:  row.Damage,
							AfterHP: row.AfterHp,
						},
					)
				}
				return result
			}()
			battleSequence := newBattleSequence(
				&event.EventSequenceArgs{
					SequenceId: sequenceId,
					Actor:      skillApplyResult.Actor,
					Target:     damageInformation,
				},
			)
			addBattleSequence(battleSequence)
			for _, row := range skillApplyResult.Rows {
				if !row.IsTargetBeaten {
					continue
				}
				actualTarget := row.TargetId
				targetSide := row.TargetSide
				if targetSide == actor.ActorSideEnemy {
					enemyId := actorIdToEnemy[actualTarget]
					viewData := serveEnemyView(enemyId)
					beatenSequence := newBattleSequence(
						&event.EventSequenceArgs{
							SequenceId: viewData.BeatenSequenceId,
							Actor:      actualTarget,
							Target: []*event.DamageInformation{
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
	indexToActor func(int) actor.ActorId,
	serveSelectedCommand func() battle.PlayerCommand,
	resetBattleSequence func(),
	playSequence func([]*skill.SkillApplyResult),
	processBattle battle.ProcessBattleFunc,
	setPlayerBeaten func(bool),
) func(int) {
	return func(index int) {
		closeWindow()
		target := indexToActor(index)
		command := serveSelectedCommand()
		response := processBattle(
			&battle.ProcessBattleRequest{
				TargetId: []actor.ActorId{target},
				Command:  command,
			},
		)

		setPlayerBeaten(response.IsMainActorBeaten)
		resetBattleSequence()
		playSequence(response.SkillApplyResults)
	}
}
