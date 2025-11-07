package adapter

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MessageWindowTextRecord struct {
	EventID         string
	TextID          string
	WaitForComplete bool
}

type messageWindowTextYAML struct {
	MessageWindowTexts []struct {
		EventID         string `yaml:"eventId"`
		TextID          string `yaml:"textId"`
		WaitForComplete bool   `yaml:"waitForComplete"`
	} `yaml:"messageWindowTexts"`
}

func MessageWindowTextRecordsFromYAML(filePath string) map[string]*MessageWindowTextRecord {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var yamlData messageWindowTextYAML
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		panic(err)
	}

	result := make(map[string]*MessageWindowTextRecord, len(yamlData.MessageWindowTexts))
	for _, record := range yamlData.MessageWindowTexts {
		result[record.EventID] = &MessageWindowTextRecord{
			EventID:         record.EventID,
			TextID:          record.TextID,
			WaitForComplete: record.WaitForComplete,
		}
	}
	return result
}
