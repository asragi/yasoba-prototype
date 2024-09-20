package component

import "github.com/asragi/yasoba-prototype/core"

type EnemyViewData struct {
	EnemyId          core.EnemyId
	BeatenSequenceId EventSequenceId
}

type ServeEnemyViewData func(id core.EnemyId) *EnemyViewData

func NewServeEnemyViewData() ServeEnemyViewData {
	dict := map[core.EnemyId]*EnemyViewData{}
	register := func(id core.EnemyId, beatenSequenceId EventSequenceId) {
		dict[id] = &EnemyViewData{
			EnemyId:          id,
			BeatenSequenceId: beatenSequenceId,
		}
	}
	register(core.EnemyPunchingBagId, EventSequenceIdPunchingBagBeaten)
	return func(id core.EnemyId) *EnemyViewData {
		return dict[id]
	}
}
