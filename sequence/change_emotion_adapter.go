package sequence

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
)

// YAMLデータ構造体
type changeEmotionModelYAML struct {
	ChangeEmotions []struct {
		EventID string `yaml:"eventId"`
		ActorID string `yaml:"actorId"`
		Emotion string `yaml:"emotion"`
	} `yaml:"changeEmotions"`
}

// createChangeEmotionDataPortFromYAML はYAMLファイルからchangeEmotionDataPortを作成します
func createChangeEmotionDataPortFromYAML(filePath string) changeEmotionDataPort {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var yamlData changeEmotionModelYAML
	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		panic(err)
	}

	emotionMap := make(map[eventId]*ChangeEmotion)
	for _, emotion := range yamlData.ChangeEmotions {
		eventId := eventId(emotion.EventID)

		// emotion文字列をBattleEmotionTypeに変換
		var emotionType component.BattleEmotionType
		switch emotion.Emotion {
		case "normal":
			emotionType = component.BattleEmotionNormal
		case "damage":
			emotionType = component.BattleEmotionDamage
		default:
			emotionType = component.BattleEmotionNormal // デフォルト値
		}

		emotionMap[eventId] = &ChangeEmotion{
			actorId: core.ActorId(emotion.ActorID),
			emotion: emotionType,
		}
	}

	return func(id eventId) *ChangeEmotion {
		return emotionMap[id]
	}
}
