package invoke

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/sequence"
)

type InvokeTiming int

const (
	InvokeTimingStartBattle InvokeTiming = iota
	InvokeTimingStartTurn
	InvokeTimingEveryAction
)

// 起動するべきSequenceのIdを返す
type checkEventInvoke func(
	invokeTiming InvokeTiming,
	turn core.TurnCount,
) []sequence.SequenceId

type produceCheckEventInvokeFunc func(
	core.BattleId,
	checkActorCondition,
) checkEventInvoke

func checkConditions(
	conditions []*eventConditionData,
	checkActorCondition checkActorCondition,
) bool {
	for _, condition := range conditions {
		if condition.conditionType == conditionTypeActorHp {
			if !checkActorCondition(condition.id) {
				return false
			}
		}
	}
	return true
}

func initializeProduceCheckEventInvoke(
	getConditions getConditionsPort,
) produceCheckEventInvokeFunc {
	return func(
		battleId core.BattleId,
		checkActorCondition checkActorCondition,
	) checkEventInvoke {
		sequenceAndConditions := getConditions(battleId)
		return func(
			invokeTiming InvokeTiming,
			turn core.TurnCount,
		) []sequence.SequenceId {
			invokingSequenceIds := []sequence.SequenceId{}
			for sequenceId, conditions := range sequenceAndConditions {
				if !checkConditions(conditions, checkActorCondition) {
					continue
				}
				invokingSequenceIds = append(invokingSequenceIds, sequenceId)
			}
			return invokingSequenceIds
		}
	}
}
