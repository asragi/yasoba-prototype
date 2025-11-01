package adapter

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SequenceModelRecord struct {
	ID string
}

type EventDataRecord struct {
	ID        string
	EventType string
	OwnerID   string
	Order     int
}

type sequenceModelYAML struct {
	Sequences []struct {
		ID string `yaml:"id"`
	} `yaml:"sequences"`
}

type eventDataModelYAML struct {
	Events []struct {
		ID        string `yaml:"id"`
		EventType string `yaml:"eventType"`
		OwnerID   string `yaml:"ownerId"`
		Order     int    `yaml:"order"`
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
			result = append(result, &SequenceModelRecord{ID: seq.ID})
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
				ID:        event.ID,
				EventType: event.EventType,
				OwnerID:   event.OwnerID,
				Order:     event.Order,
			})
		}
		return result
	}
}
