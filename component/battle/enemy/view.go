package enemy

import (
	"fmt"

	battleevent "github.com/asragi/yasoba-prototype/component/battle/event"
	"github.com/asragi/yasoba-prototype/game/enemy"
)

type EnemyViewData struct {
	EnemyId          enemy.EnemyId
	BeatenSequenceId battleevent.EventSequenceId
}

type ServeEnemyViewData func(id enemy.EnemyId) *EnemyViewData

func NewServeEnemyViewData() ServeEnemyViewData {
	dict := map[enemy.EnemyId]*EnemyViewData{}
	register := func(id enemy.EnemyId, beatenSequenceId battleevent.EventSequenceId) {
		dict[id] = &EnemyViewData{
			EnemyId:          id,
			BeatenSequenceId: beatenSequenceId,
		}
	}
	register(enemy.EnemyPunchingBagId, battleevent.EventSequenceIdPunchingBagBeaten)
	return func(id enemy.EnemyId) *EnemyViewData {
		data, ok := dict[id]
		if !ok {
			panic(fmt.Sprintf("enemy view data not found: %v", id))
		}
		return data
	}
}
