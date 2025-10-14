package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle_skill"
	"github.com/asragi/yasoba-prototype/skilldata"
	textpkg "github.com/asragi/yasoba-prototype/text"
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
func (b *PlayerCommand) ToTextId() textpkg.TextId {
	switch *b {
	case PlayerCommandAttack:
		return textpkg.TextIdBattleCommandAttack
	case PlayerCommandFire:
		return textpkg.TextIdBattleCommandFire
	case PlayerCommandThunder:
		return textpkg.TextIdBattleCommandThunder
	case PlayerCommandBarrier:
		return textpkg.TextIdBattleCommandBarrier
	case PlayerCommandWind:
		return textpkg.TextIdBattleCommandWind
	case PlayerCommandFocus:
		return textpkg.TextIdBattleCommandFocus
	case PlayerCommandDefend:
		return textpkg.TextIdBattleCommandDefend
	}
	return ""
}

// BattlePlayerCommandResult holds the SelectedAction triggered by a command.
type BattlePlayerCommandResult struct {
	SkillApplyArgs *battle_skill.SelectedAction
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
		decidedSkillId := func() skilldata.SkillId {
			if isToEnemy(command.TargetId) {
				switch command.Command {
				case PlayerCommandAttack:
					return skilldata.SkillIdLuneAttack
				case PlayerCommandFire:
					return skilldata.SkillIdLuneFireEnemy
				default:
					panic("not implemented")
				}
			}
			// TODO: implement friendly-target commands
			return skilldata.SkillIdLuneAttack
		}()
		return &BattlePlayerCommandResult{
			SkillApplyArgs: &battle_skill.SelectedAction{
				Id:       decidedSkillId,
				Actor:    command.ActorId,
				SubActor: actor.ActorEmptyId,
				Target:   command.TargetId,
			},
		}
	}
}
