package invoke

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/common/character"
)

type conditionActorHp struct {
	ownerId        conditionId
	actorId        actor.ActorId
	thresholdRatio float64
}

type conditionActorHpDataPort func() []*conditionActorHp

type checkActorCondition func(conditionId) bool
type getActorHpRatio func(actor.ActorId) character.HPRatio

type produceCheckActorCondition func(getActorHpRatio) checkActorCondition

func initializeCheckActorCondition(
	conditionActorHpDataPort conditionActorHpDataPort,
) produceCheckActorCondition {
	conditionActorHpData := conditionActorHpDataPort()
	conditionMap := make(map[conditionId]*conditionActorHp)
	for _, condition := range conditionActorHpData {
		conditionMap[condition.ownerId] = condition
	}
	return func(getActorHpRatio getActorHpRatio) checkActorCondition {
		return func(conditionId conditionId) bool {
			if _, ok := conditionMap[conditionId]; !ok {
				panic(fmt.Sprintf("conditionActorHp not found: %s", conditionId))
			}
			actorId := conditionMap[conditionId].actorId
			hpRatio := getActorHpRatio(actorId)
			return hpRatio.Float64() <= conditionMap[conditionId].thresholdRatio
		}
	}
}
