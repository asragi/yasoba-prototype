package invoke

import (
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/sequence"
)

type battleSequenceRelation struct {
	ownerId    battle.BattleId
	sequenceId sequence.SequenceId
}

type battleSequenceRelationDataPort func() []*battleSequenceRelation

// BattleIdに対するSequenceIdの関係性をportから取り出す
type battleSequenceConditionPort func(battle.BattleId) []sequence.SequenceId

func initializeBattleSequenceCondition(
	dataPort battleSequenceRelationDataPort,
) battleSequenceConditionPort {
	dict := make(map[battle.BattleId][]sequence.SequenceId)
	allData := dataPort()
	for _, d := range allData {
		dict[d.ownerId] = append(dict[d.ownerId], d.sequenceId)
	}
	return func(battleId battle.BattleId) []sequence.SequenceId {
		return dict[battleId]
	}
}
