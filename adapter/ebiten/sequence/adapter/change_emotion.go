package adapter

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ChangeEmotionRecord struct {
	EventID string
	ActorID string
	Emotion string
}

type changeEmotionModelYAML struct {
	ChangeEmotions []struct {
		EventID string `yaml:"id"`
		ActorID string `yaml:"actorId"`
		Emotion string `yaml:"emotion"`
	} `yaml:"events"`
}

func ChangeEmotionRecordsFromYAML(filePath string) map[string]*ChangeEmotionRecord {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var yamlData changeEmotionModelYAML
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		panic(err)
	}

	result := make(map[string]*ChangeEmotionRecord, len(yamlData.ChangeEmotions))
	for _, emotion := range yamlData.ChangeEmotions {
		result[emotion.EventID] = &ChangeEmotionRecord{
			EventID: emotion.EventID,
			ActorID: emotion.ActorID,
			Emotion: emotion.Emotion,
		}
	}
	return result
}
