package sequence

import (
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLデータ構造体
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

// createSequenceModelPortFromYAML はYAMLファイルからsequenceModelPortを作成します
func createSequenceModelPortFromYAML(filePath string) sequenceModelPort {
	return func() []*sequenceModel {
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		var yamlData sequenceModelYAML
		err = yaml.Unmarshal(data, &yamlData)
		if err != nil {
			panic(err)
		}

		var result []*sequenceModel
		for _, seq := range yamlData.Sequences {
			result = append(result, &sequenceModel{
				id: SequenceId(seq.ID),
			})
		}

		return result
	}
}

// createEventDataModelPortFromYAML はYAMLファイルからeventDataPortを作成します
func createEventDataModelPortFromYAML(filePath string) eventDataModelPort {
	return func() []*eventDataModel {
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		var yamlData eventDataModelYAML
		err = yaml.Unmarshal(data, &yamlData)
		if err != nil {
			panic(err)
		}

		var result []*eventDataModel
		for _, event := range yamlData.Events {
			result = append(result, &eventDataModel{
				id:        eventId(event.ID),
				eventType: eventType(event.EventType),
				ownerId:   SequenceId(event.OwnerID),
				order:     event.Order,
			})
		}

		return result
	}
}
