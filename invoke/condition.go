package invoke

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/sequence"
)

type conditionType string

const (
	conditionTypeActorHp conditionType = "actor_hp"
)

type conditionId string

type eventConditionData struct {
	id            conditionId
	ownerId       sequence.SequenceId
	conditionType conditionType
}

// BattleIdに対してSequenceとConditionのmapを返す
type getConditionsPort func(core.BattleId) map[sequence.SequenceId][]*eventConditionData

func initializeGetConditionsPort(
	conditionDataPort conditionDataPort,
	getBattleRelation battleSequenceConditionPort,
) getConditionsPort {
	dict := make(map[sequence.SequenceId][]*eventConditionData)
	for _, condition := range conditionDataPort() {
		dict[condition.ownerId] = append(dict[condition.ownerId], condition)
	}

	getCondition := func(sequenceIds []sequence.SequenceId) map[sequence.SequenceId][]*eventConditionData {
		result := make(map[sequence.SequenceId][]*eventConditionData)
		for _, sequenceId := range sequenceIds {
			result[sequenceId] = dict[sequenceId]
		}
		return result
	}

	return func(battleId core.BattleId) map[sequence.SequenceId][]*eventConditionData {
		sequenceIds := getBattleRelation(battleId)
		return getCondition(sequenceIds)
	}
}

type conditionDataPort func() []*eventConditionData
