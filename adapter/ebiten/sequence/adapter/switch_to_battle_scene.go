package adapter

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SwitchToBattleSceneRecord struct {
	EventID  string
	BattleID string
}

type switchToBattleSceneModelYAML struct {
	Events []struct {
		EventID  string `yaml:"id"`
		BattleID string `yaml:"battleId"`
	} `yaml:"events"`
}

func SwitchToBattleSceneRecordsFromYAML(filePath string) map[string]*SwitchToBattleSceneRecord {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var yamlData switchToBattleSceneModelYAML
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		panic(err)
	}

	result := make(map[string]*SwitchToBattleSceneRecord, len(yamlData.Events))
	for _, record := range yamlData.Events {
		result[record.EventID] = &SwitchToBattleSceneRecord{
			EventID:  record.EventID,
			BattleID: record.BattleID,
		}
	}
	return result
}
