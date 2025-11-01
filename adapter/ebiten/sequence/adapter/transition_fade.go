package adapter

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TransitionFadeRecord struct {
	EventID         string
	WaitForComplete bool
}

type transitionFadeModelYAML struct {
	Events []struct {
		EventID         string `yaml:"id"`
		WaitForComplete bool   `yaml:"waitForComplete"`
	} `yaml:"events"`
}

func TransitionFadeRecordsFromYAML(filePath string) map[string]*TransitionFadeRecord {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var yamlData transitionFadeModelYAML
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		panic(err)
	}

	result := make(map[string]*TransitionFadeRecord, len(yamlData.Events))
	for _, record := range yamlData.Events {
		result[record.EventID] = &TransitionFadeRecord{
			EventID:         record.EventID,
			WaitForComplete: record.WaitForComplete,
		}
	}
	return result
}
