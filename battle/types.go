package battle

import "github.com/asragi/yasoba-prototype/battle_decision"

// BattleId identifies a battle scenario.
type BattleId string

const (
	BattleIdTest001 BattleId = "test_battle_001"
)

// TurnCount is zero-based from the beginning of a battle.
type TurnCount int

// BattleEndType mirrors the decision package end state type.
type BattleEndType = battle_decision.BattleEndType

const (
	BattleEndTypeNone = battle_decision.BattleEndTypeNone
	BattleEndTypeWin  = battle_decision.BattleEndTypeWin
	BattleEndTypeLose = battle_decision.BattleEndTypeLose
)
