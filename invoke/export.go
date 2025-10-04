package invoke

import "github.com/asragi/yasoba-prototype/core"

type GetActorHpRatio = getActorHpRatio
type CheckInvokeSequence = checkEventInvoke

type ProduceCheckInvokeSequence func(
	core.BattleId,
	GetActorHpRatio,
) CheckInvokeSequence

func InitializeProduceCheckInvokeSequence() ProduceCheckInvokeSequence {
	conditionDataPort := createEventConditionDataPortFromYAML("data/event_condition.yaml")
	battleSequenceRelationDataPort := createBattleSequenceRelationDataPortFromYAML("data/battle_sequence_relation.yaml")
	conditionActorHpDataPort := createConditionActorHpDataPortFromYAML("data/condition_actor_hp.yaml")
	battleSequenceRelationAdapter := initializeBattleSequenceCondition(
		battleSequenceRelationDataPort,
	)
	getConditionAdapter := initializeGetConditionsPort(
		conditionDataPort,
		battleSequenceRelationAdapter,
	)
	produceCheckEventInvokeAdapter := initializeProduceCheckEventInvoke(
		getConditionAdapter,
	)
	produceCheckActorCondition := initializeCheckActorCondition(conditionActorHpDataPort)
	return func(
		battleId core.BattleId,
		getActorHpRatio GetActorHpRatio,
	) CheckInvokeSequence {
		checkActorCondition := produceCheckActorCondition(getActorHpRatio)
		return produceCheckEventInvokeAdapter(
			battleId,
			checkActorCondition,
		)
	}
}
