package sequence

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/asragi/yasoba-prototype/core"
)

// YAMLデータ構造体
type partnerDialogueModelYAML struct {
	PartnerDialogues []struct {
		EventID string `yaml:"eventId"`
		TextID  string `yaml:"textId"`
	} `yaml:"partnerDialogues"`
}

// createPartnerDialogueDataPortFromYAML はYAMLファイルからpartnerDialogueDataPortを作成します
func createPartnerDialogueDataPortFromYAML(filePath string) partnerDialogueDataPort {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var yamlData partnerDialogueModelYAML
	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		panic(err)
	}

	dialogueMap := make(map[eventId]*partnerDialogueModel)
	for _, dialogue := range yamlData.PartnerDialogues {
		eventId := eventId(dialogue.EventID)
		dialogueMap[eventId] = &partnerDialogueModel{
			eventId: eventId,
			textId:  core.TextId(dialogue.TextID),
		}
	}

	return func(id eventId) *partnerDialogueModel {
		return dialogueMap[id]
	}
}
