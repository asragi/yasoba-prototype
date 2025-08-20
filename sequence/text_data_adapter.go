package sequence

import (
	"os"

	"github.com/asragi/yasoba-prototype/core"
	"gopkg.in/yaml.v3"
)

type TextDataYaml struct {
	Texts []TextEntry `yaml:"texts"`
}

type TextEntry struct {
	ID   string `yaml:"id"`
	Text string `yaml:"text"`
}

func LoadTextDataFromYaml(filePath string) (core.ServeTextDataFunc, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var yamlData TextDataYaml
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, err
	}

	dict := make(map[core.TextId]*core.TextData)
	for _, entry := range yamlData.Texts {
		dict[core.TextId(entry.ID)] = &core.TextData{
			Id:   core.TextId(entry.ID),
			Text: core.TextString(entry.Text),
		}
	}

	return func(id core.TextId) *core.TextData {
		if textData, ok := dict[id]; ok {
			return textData
		}
		panic("text not found: " + string(id))
	}, nil
}

// CreateServeTextDataFromYamlは、デフォルトのtext_data.yamlファイルからテキストデータを読み込みます
func CreateServeTextDataFromYaml() (core.ServeTextDataFunc, error) {
	return LoadTextDataFromYaml("data/text_data.yaml")
}
