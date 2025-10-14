package battle

// BattleId identifies a battle scenario.
type BattleId string

const (
	BattleIdTest001 BattleId = "test_battle_001"
)

// TurnCount is zero-based from the beginning of a battle.
type TurnCount int
