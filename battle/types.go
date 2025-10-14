package battle

import "github.com/asragi/yasoba-prototype/battle/decision"

// BattleId identifies a battle scenario.
type BattleId string

const (
	BattleIdTest001 BattleId = "test_battle_001"
)

// TurnCount is zero-based from the beginning of a battle.
type TurnCount int

// BattleEndType mirrors the decision package end state type.
type BattleEndType = decision.BattleEndType

const (
	BattleEndTypeNone = decision.BattleEndTypeNone
	BattleEndTypeWin  = decision.BattleEndTypeWin
	BattleEndTypeLose = decision.BattleEndTypeLose
)
