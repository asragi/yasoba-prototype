package adapter

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PartnerDialogueRecord struct {
	EventID string
	TextID  string
}

type partnerDialogueModelYAML struct {
	PartnerDialogues []struct {
		EventID string `yaml:"eventId"`
		TextID  string `yaml:"textId"`
	} `yaml:"partnerDialogues"`
}

func PartnerDialogueRecordsFromYAML(filePath string) map[string]*PartnerDialogueRecord {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var yamlData partnerDialogueModelYAML
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		panic(err)
	}

	result := make(map[string]*PartnerDialogueRecord, len(yamlData.PartnerDialogues))
	for _, dialogue := range yamlData.PartnerDialogues {
		result[dialogue.EventID] = &PartnerDialogueRecord{
			EventID: dialogue.EventID,
			TextID:  dialogue.TextID,
		}
	}
	return result
}
