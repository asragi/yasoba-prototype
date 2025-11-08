package command

import "github.com/asragi/yasoba-prototype/text"

// Id uniquely identifies a battle command.
type Id string

const (
	IdAttack  Id = "attack"
	IdFire    Id = "fire"
	IdThunder Id = "thunder"
	IdBarrier Id = "barrier"
	IdWind    Id = "wind"
	IdFocus   Id = "focus"
	IdDefend  Id = "defend"
)

// ToTextId converts a command id into its localized text identifier.
func (id Id) ToTextId() text.TextId {
	switch id {
	case IdAttack:
		return text.TextIdBattleCommandAttack
	case IdFire:
		return text.TextIdBattleCommandFire
	case IdThunder:
		return text.TextIdBattleCommandThunder
	case IdBarrier:
		return text.TextIdBattleCommandBarrier
	case IdWind:
		return text.TextIdBattleCommandWind
	case IdFocus:
		return text.TextIdBattleCommandFocus
	case IdDefend:
		return text.TextIdBattleCommandDefend
	default:
		return ""
	}
}
