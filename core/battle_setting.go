package core

import (
	"github.com/asragi/yasoba-prototype/frontend"
)

type EnemySetting struct {
	EnemyId  EnemyId
	Position *frontend.Vector
}

type BattleSettingId string

const (
	BattleSettingTest       BattleSettingId = "test"
	BattleSettingTripleTest BattleSettingId = "triple_test"
)

type BattleSetting struct {
	Enemies []*EnemySetting
}

type ServeBattleSetting func(BattleSettingId) *BattleSetting

func CreateServeBattleSetting() ServeBattleSetting {
	dict := make(map[BattleSettingId]*BattleSetting)
	dict[BattleSettingTest] = &BattleSetting{
		Enemies: []*EnemySetting{
			{
				EnemyId:  EnemyPunchingBagId,
				Position: frontend.VectorZero,
			},
		},
	}
	dict[BattleSettingTripleTest] = &BattleSetting{
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

	return func(id BattleSettingId) *BattleSetting {
		return dict[id]
	}
}
