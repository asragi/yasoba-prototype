package scene

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/battle/setup"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	battleenemy "github.com/asragi/yasoba-prototype/component/battle/enemy"
	battleevent "github.com/asragi/yasoba-prototype/component/battle/event"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/enemy"
	"github.com/asragi/yasoba-prototype/invoke"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/text"
)

type BattleScene struct {
	ui                  battleUI
	battleSequence      *battleevent.BattleEventSequencer
	enemyData           []*setup.EnemyIdPair
	actorNames          map[actor.ActorId]text.TextId
	endState            battle.BattleEndType
	createSequence      sequence.CreateSequence
	sequences           *sequence.SequenceManager
	checkInvokeSequence invoke.CheckInvokeSequence
	invokedSequences    map[sequence.SequenceId]bool
	turnCount           battle.TurnCount
}

func (s *BattleScene) onTurnEnd() {
	s.turnCount++
	if s.endState == battle.BattleEndTypeWin {
		s.sequences.AddSequence(s.createSequence("test_sequence_0100"))
		return
	}
	if s.endState == battle.BattleEndTypeLose {
		s.ui.messageWindow.SetText("やられてしまった……", false)
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

func (s *BattleScene) Draw(drawFunc frontend.DrawFunc) {
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

func advanceBattleSequence(sequence *battleevent.BattleEventSequencer) bool {
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
}

type CreateBattleScene func(*BattleOption) *BattleScene

type playBattleSequenceFunc func([]*battleSkill.SkillApplyResult)

// SkillApplyResultに基づいて戦闘の演出を行う
func createPlayBattleSequence(
	skillToSequence battleevent.SkillToSequenceFunc,
	newBattleSequence battleevent.NewBattleSequenceFunc,
	addBattleSequence func(battleevent.BattleSequenceFunc),
	actorIdToEnemy map[actor.ActorId]enemy.EnemyId,
	serveEnemyView battleenemy.ServeEnemyViewData,
) playBattleSequenceFunc {
	return func(skillApplyResultSet []*battleSkill.SkillApplyResult) {
		for _, skillApplyResult := range skillApplyResultSet {
			skillId := skillApplyResult.SkillId
			sequenceId := skillToSequence(skillId)
			damageInformation := func() []*battleevent.DamageInformation {
				result := make([]*battleevent.DamageInformation, 0)
				for _, row := range skillApplyResult.Rows {
					result = append(
						result, &battleevent.DamageInformation{
							Target:  row.TargetId,
							Damage:  row.Damage,
							AfterHP: row.AfterHp,
						},
					)
				}
				return result
			}()
			sequence := newBattleSequence(
				&battleevent.EventSequenceArgs{
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
				if targetSide == actor.ActorSideEnemy {
					enemyId := actorIdToEnemy[actualTarget]
					viewData := serveEnemyView(enemyId)
					beatenSequence := newBattleSequence(
						&battleevent.EventSequenceArgs{
							SequenceId: viewData.BeatenSequenceId,
							Actor:      actualTarget,
							Target: []*battleevent.DamageInformation{
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
	playSequence func([]*battleSkill.SkillApplyResult),
	processBattle battle.ProcessBattleFunc,
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

		resetBattleSequence()
		playSequence(response.SkillApplyResults)
	}
}
