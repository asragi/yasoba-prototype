package component

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/enemydata"
)

type EnemyViewData struct {
	EnemyId          enemydata.EnemyId
	BeatenSequenceId EventSequenceId
}

type ServeEnemyViewData func(id enemydata.EnemyId) *EnemyViewData

func NewServeEnemyViewData() ServeEnemyViewData {
	dict := map[enemydata.EnemyId]*EnemyViewData{}
	register := func(id enemydata.EnemyId, beatenSequenceId EventSequenceId) {
		dict[id] = &EnemyViewData{
			EnemyId:          id,
			BeatenSequenceId: beatenSequenceId,
		}
	}
	register(enemydata.EnemyPunchingBagId, EventSequenceIdPunchingBagBeaten)
	return func(id enemydata.EnemyId) *EnemyViewData {
		data, ok := dict[id]
		if !ok {
			panic(fmt.Sprintf("enemy view data not found: %v", id))
		}
		return data
	}
}
