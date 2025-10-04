package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/invoke"
	"github.com/asragi/yasoba-prototype/sequence"
)

type BattleScene struct {
	ui                  battleUI
	battleSequence      *component.BattleEventSequencer
	enemyData           []*core.EnemyIdPair
	actorNames          map[core.ActorId]core.TextId
	endState            core.BattleEndType
	createSequence      sequence.CreateSequence
	sequences           *sequence.SequenceManager
	checkInvokeSequence invoke.CheckInvokeSequence
	invokedSequences    map[sequence.SequenceId]bool
	turnCount           core.TurnCount
}

func (s *BattleScene) onTurnEnd() {
	s.turnCount++
	if s.endState == core.BattleEndTypeWin {
		s.sequences.AddSequence(s.createSequence("test_sequence_0100"))
		return
	}
	if s.endState == core.BattleEndTypeLose {
		s.ui.messageWindow.SetText("やられてしまった……", false)
		return
	}
	s.ui.input.Set(s.ui.battleSelectWindow)
	s.ui.battleSelectWindow.Open()
	s.checkAndStartSequences(invoke.InvokeTimingStartTurn)
}

func (s *BattleScene) onBattleEnd(endType core.BattleEndType) {
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

func advanceBattleSequence(sequence *component.BattleEventSequencer) bool {
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
	BattleSettingId core.BattleSettingId
	BattleId        core.BattleId
}

type CreateBattleScene func(*BattleOption) *BattleScene

type playBattleSequenceFunc func([]*core.SkillApplyResult)

// SkillApplyResultに基づいて戦闘の演出を行う
func createPlayBattleSequence(
	skillToSequence component.SkillToSequenceFunc,
	newBattleSequence component.NewBattleSequenceFunc,
	addBattleSequence func(component.BattleSequenceFunc),
	actorIdToEnemy map[core.ActorId]core.EnemyId,
	serveEnemyView component.ServeEnemyViewData,
) playBattleSequenceFunc {
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
