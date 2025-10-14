package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/core"
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
func (b *PlayerCommand) ToTextId() core.TextId {
	switch *b {
	case PlayerCommandAttack:
		return "battle_command_attack"
	case PlayerCommandFire:
		return "battle_command_fire"
	case PlayerCommandThunder:
		return "battle_command_thunder"
	case PlayerCommandBarrier:
		return "battle_command_barrier"
	case PlayerCommandWind:
		return "battle_command_wind"
	case PlayerCommandFocus:
		return "battle_command_focus"
	case PlayerCommandDefend:
		return "battle_command_defend"
	}
	return ""
}

// BattlePlayerCommandResult holds the SelectedAction triggered by a command.
type BattlePlayerCommandResult struct {
	SkillApplyArgs *core.SelectedAction
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
		decidedSkillId := func() core.SkillId {
			if isToEnemy(command.TargetId) {
				switch command.Command {
				case PlayerCommandAttack:
					return core.SkillIdLuneAttack
				case PlayerCommandFire:
					return core.SkillIdLuneFireEnemy
				default:
					panic("not implemented")
				}
			}
			// TODO: implement friendly-target commands
			return core.SkillIdLuneAttack
		}()
		return &BattlePlayerCommandResult{
			SkillApplyArgs: &core.SelectedAction{
				Id:       decidedSkillId,
				Actor:    command.ActorId,
				SubActor: actor.ActorEmptyId,
				Target:   command.TargetId,
			},
		}
	}
}
