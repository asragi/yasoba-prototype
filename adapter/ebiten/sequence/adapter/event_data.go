package adapter

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SequenceModelRecord struct {
	ID          string
	HeadEventID string
}

type EventDataRecord struct {
	ID          string
	EventType   string
	NextEventID string
}

type sequenceModelYAML struct {
	Sequences []struct {
		ID          string `yaml:"id"`
		HeadEventID string `yaml:"headEventId"`
	} `yaml:"sequences"`
}

type eventDataModelYAML struct {
	Events []struct {
		ID          string `yaml:"id"`
		EventType   string `yaml:"eventType"`
		NextEventID string `yaml:"nextEventId"`
	} `yaml:"events"`
}

type GetSequenceModelFunc func() []*SequenceModelRecord
type GetEventDataModelFunc func() []*EventDataRecord

func SequenceModelPortFromYAML(filePath string) GetSequenceModelFunc {
	return func() []*SequenceModelRecord {
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		var yamlData sequenceModelYAML
		if err := yaml.Unmarshal(data, &yamlData); err != nil {
			panic(err)
		}

		result := make([]*SequenceModelRecord, 0, len(yamlData.Sequences))
		for _, seq := range yamlData.Sequences {
			result = append(result, &SequenceModelRecord{
				ID:          seq.ID,
				HeadEventID: seq.HeadEventID,
			})
		}
		return result
	}
}

func EventDataModelPortFromYAML(filePath string) GetEventDataModelFunc {
	return func() []*EventDataRecord {
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		var yamlData eventDataModelYAML
		if err := yaml.Unmarshal(data, &yamlData); err != nil {
			panic(err)
		}

		result := make([]*EventDataRecord, 0, len(yamlData.Events))
		for _, event := range yamlData.Events {
			result = append(result, &EventDataRecord{
				ID:          event.ID,
				EventType:   event.EventType,
				NextEventID: event.NextEventID,
			})
		}
		return result
	}
}
