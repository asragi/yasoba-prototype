package battle

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/text"
)

// PlayerCommand is a command that the player can select in the battle.
type PlayerCommand int

const (
	PlayerCommandAttack PlayerCommand = iota
	PlayerCommandFire
	PlayerCommandThunder
	PlayerCommandBarrier
	PlayerCommandWind
	PlayerCommandFocus
	PlayerCommandDefend
)

// ToTextId converts a player command into the associated text ID.
func (b *PlayerCommand) ToTextId() text.TextId {
	switch *b {
	case PlayerCommandAttack:
		return text.TextIdBattleCommandAttack
	case PlayerCommandFire:
		return text.TextIdBattleCommandFire
	case PlayerCommandThunder:
		return text.TextIdBattleCommandThunder
	case PlayerCommandBarrier:
		return text.TextIdBattleCommandBarrier
	case PlayerCommandWind:
		return text.TextIdBattleCommandWind
	case PlayerCommandFocus:
		return text.TextIdBattleCommandFocus
	case PlayerCommandDefend:
		return text.TextIdBattleCommandDefend
	}
	return ""
}

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
	return func(command *PostCommandRequest) *BattlePlayerCommandResult {
		decidedSkillId := func() skill.SkillId {
			if isToEnemy(command.TargetId) {
				switch command.Command {
				case PlayerCommandAttack:
					return skill.SkillIdLuneAttack
				case PlayerCommandFire:
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
				Actor:    command.ActorId,
				SubActor: actor.ActorEmptyId,
				Target:   command.TargetId,
			},
		}
	}
}
