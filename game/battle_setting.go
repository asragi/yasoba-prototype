package game

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
)

type EnemySetting struct {
	EnemyId  core.EnemyId
	Position *frontend.Vector
}

type BattleSettingId string

const (
	BattleSettingTest BattleSettingId = "test"
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
				EnemyId:  core.EnemyPunchingBagId,
				Position: frontend.VectorZero,
			},
		},
	}

	return func(id BattleSettingId) *BattleSetting {
		return dict[id]
	}
}
