package enemy

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/battle/enemy"
	"github.com/asragi/yasoba-prototype/view/battle/event"
)

type EnemyViewData struct {
	EnemyId          enemy.EnemyId
	BeatenSequenceId event.EventSequenceId
}

type ServeEnemyViewData func(id enemy.EnemyId) *EnemyViewData

func NewServeEnemyViewData() ServeEnemyViewData {
	dict := map[enemy.EnemyId]*EnemyViewData{}
	register := func(id enemy.EnemyId, beatenSequenceId event.EventSequenceId) {
		dict[id] = &EnemyViewData{
			EnemyId:          id,
			BeatenSequenceId: beatenSequenceId,
		}
	}
	register(enemy.EnemyPunchingBagId, event.EventSequenceIdPunchingBagBeaten)
	return func(id enemy.EnemyId) *EnemyViewData {
		data, ok := dict[id]
		if !ok {
			panic(fmt.Sprintf("enemy view data not found: %v", id))
		}
		return data
	}
}
