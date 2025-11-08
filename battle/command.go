package battle

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/battle/skill"
)

// BattlePlayerCommandResult holds the SelectedAction triggered by a command.
type BattlePlayerCommandResult struct {
	SkillApplyArgs *skill.SelectedAction
}

// ProcessPlayerCommandFunc determines the SelectedAction for a player command.
type ProcessPlayerCommandFunc func(*PostCommandRequest) *BattlePlayerCommandResult

// CreateProcessPlayerCommand builds a command processor using the supplied actor lookup.
func CreateProcessPlayerCommand(supplyActor actor.ActorSupplier) ProcessPlayerCommandFunc {
	isToEnemy := func(targets []actor.ActorId) bool {
		if len(targets) == 0 {
			return false
		}
		target := supplyActor(targets[0])
		return target.Side == actor.ActorSideEnemy
	}
	return func(commandReq *PostCommandRequest) *BattlePlayerCommandResult {
		decidedSkillId := func() skill.SkillId {
			if isToEnemy(commandReq.TargetId) {
				switch commandReq.Command {
				case command.IdAttack:
					return skill.SkillIdLuneAttack
				case command.IdFire:
					return skill.SkillIdLuneFireEnemy
				default:
					panic("not implemented")
				}
			}
			// TODO: implement friendly-target commands
			return skill.SkillIdLuneAttack
		}()
		return &BattlePlayerCommandResult{
			SkillApplyArgs: &skill.SelectedAction{
				Id:       decidedSkillId,
				Actor:    commandReq.ActorId,
				SubActor: actor.ActorEmptyId,
				Target:   commandReq.TargetId,
			},
		}
	}
}
