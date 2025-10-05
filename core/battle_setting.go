package core

import (
	"github.com/asragi/yasoba-prototype/frontend"
)

type EnemySetting struct {
	EnemyId  EnemyId
	Position *frontend.Vector
}

type BattleSetting struct {
	Enemies []*EnemySetting
}

type ServeBattleSetting func(BattleId) *BattleSetting

func CreateServeBattleSetting() ServeBattleSetting {
	dict := make(map[BattleId]*BattleSetting)
	dict[BattleIdTest001] = &BattleSetting{
		Enemies: []*EnemySetting{
			{
				EnemyId:  EnemyPunchingBagId,
				Position: frontend.VectorZero,
			},
		},
	}
	dict[BattleIdTripleTest] = &BattleSetting{
		Enemies: []*EnemySetting{
			{
				EnemyId:  EnemyPunchingBagId,
				Position: &frontend.Vector{X: -100},
			},
			{
				EnemyId:  EnemyPunchingBagId,
				Position: frontend.VectorZero,
			},
			{
				EnemyId:  EnemyPunchingBagId,
				Position: &frontend.Vector{X: 100},
			},
		},
	}

	return func(id BattleId) *BattleSetting {
		return dict[id]
	}
}
