package invoke

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/sequence"
)

type battleSequenceRelation struct {
	ownerId    core.BattleId
	sequenceId sequence.SequenceId
}

type battleSequenceRelationDataPort func() []*battleSequenceRelation

// BattleIdに対するSequenceIdの関係性をportから取り出す
type battleSequenceConditionPort func(core.BattleId) []sequence.SequenceId

func initializeBattleSequenceCondition(
	dataPort battleSequenceRelationDataPort,
) battleSequenceConditionPort {
	dict := make(map[core.BattleId][]sequence.SequenceId)
	allData := dataPort()
	for _, d := range allData {
		dict[d.ownerId] = append(dict[d.ownerId], d.sequenceId)
	}
	return func(battleId core.BattleId) []sequence.SequenceId {
		return dict[battleId]
	}
}
