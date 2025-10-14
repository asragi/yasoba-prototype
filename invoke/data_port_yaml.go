package invoke

import (
	"fmt"
	"os"

	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/sequence"
	"gopkg.in/yaml.v3"
)

type conditionActorHpYAML struct {
	ConditionActorHps []struct {
		OwnerId        string  `yaml:"ownerId"`
		ActorId        string  `yaml:"actorId"`
		ThresholdRatio float64 `yaml:"thresholdRatio"`
	} `yaml:"condition_actor_hps"`
}

type eventConditionYAML struct {
	EventConditions []struct {
		Id            string `yaml:"id"`
		OwnerId       string `yaml:"ownerId"`
		ConditionType string `yaml:"conditionType"`
	} `yaml:"event_conditions"`
}

type battleSequenceRelationYAML struct {
	Relations []struct {
		BattleId   string `yaml:"battleId"`
		SequenceId string `yaml:"sequenceId"`
	} `yaml:"battle_sequence_relations"`
}

func createConditionActorHpDataPortFromYAML(path string) conditionActorHpDataPort {
	return func() []*conditionActorHp {
		raw, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		var yamlData conditionActorHpYAML
		if err := yaml.Unmarshal(raw, &yamlData); err != nil {
			panic(err)
		}
		result := make([]*conditionActorHp, 0, len(yamlData.ConditionActorHps))
		for _, row := range yamlData.ConditionActorHps {
			result = append(result, &conditionActorHp{
				ownerId:        conditionId(row.OwnerId),
				actorId:        actor.ActorId(row.ActorId),
				thresholdRatio: row.ThresholdRatio,
			})
		}
		return result
	}
}

func createEventConditionDataPortFromYAML(path string) conditionDataPort {
	return func() []*eventConditionData {
		raw, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		var yamlData eventConditionYAML
		if err := yaml.Unmarshal(raw, &yamlData); err != nil {
			panic(err)
		}
		result := make([]*eventConditionData, 0, len(yamlData.EventConditions))
		for _, row := range yamlData.EventConditions {
			var condType conditionType
			switch row.ConditionType {
			case string(conditionTypeActorHp):
				condType = conditionTypeActorHp
			default:
				panic(fmt.Sprintf("unsupported condition type: %s", row.ConditionType))
			}
			result = append(result, &eventConditionData{
				id:            conditionId(row.Id),
				ownerId:       sequence.SequenceId(row.OwnerId),
				conditionType: condType,
			})
		}
		return result
	}
}

func createBattleSequenceRelationDataPortFromYAML(path string) battleSequenceRelationDataPort {
	return func() []*battleSequenceRelation {
		raw, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		var yamlData battleSequenceRelationYAML
		if err := yaml.Unmarshal(raw, &yamlData); err != nil {
			panic(err)
		}
		result := make([]*battleSequenceRelation, 0, len(yamlData.Relations))
		for _, row := range yamlData.Relations {
			result = append(result, &battleSequenceRelation{
				ownerId:    core.BattleId(row.BattleId),
				sequenceId: sequence.SequenceId(row.SequenceId),
			})
		}
		return result
	}
}
