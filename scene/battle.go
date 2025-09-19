package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/sequence"
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
	createSequence     sequence.CreateSequence
	sequences          *sequence.SequenceManager
}

func (s *BattleScene) onTurnEnd() {
	if s.endState == core.BattleEndTypeWin {
		s.sequences.AddSequence(s.createSequence("test_sequence_0100"))
		return
	}
	if s.endState == core.BattleEndTypeLose {
		s.messageWindow.SetText("やられてしまった……", false)
		return
	}
	s.input.Set(s.battleSelectWindow)
	s.battleSelectWindow.Open()
}

func (s *BattleScene) onBattleEnd(endType core.BattleEndType) {
	s.endState = endType
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
			s.onTurnEnd()
		}
	}
	s.effectManager.Update()
	s.sequences.Update()
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
