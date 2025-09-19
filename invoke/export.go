package invoke

import "github.com/asragi/yasoba-prototype/core"

type CheckInvokeSequence = checkEventInvoke

type ProduceCheckInvokeSequence func(
	core.BattleId,
	checkActorCondition,
) CheckInvokeSequence

func InitializeProduceCheckInvokeSequence() ProduceCheckInvokeSequence {
	tmpConditionDataPort := func() []*eventConditionData { return nil }
	tmpBattleSequenceRelationDataPort := func() []*battleSequenceRelation { return nil }
	battleSequenceRelationAdapter := initializeBattleSequenceCondition(
		tmpBattleSequenceRelationDataPort,
	)
	getConditionAdapter := initializeGetConditionsPort(
		tmpConditionDataPort,
		battleSequenceRelationAdapter,
	)
	produceCheckEventInvokeAdapter := initializeProduceCheckEventInvoke(
		getConditionAdapter,
	)
	return func(
		battleId core.BattleId,
		checkActorCondition checkActorCondition,
	) CheckInvokeSequence {
		checkEventInvoke := produceCheckEventInvokeAdapter(
			battleId,
			checkActorCondition,
		)
		return checkEventInvoke
	}
}
