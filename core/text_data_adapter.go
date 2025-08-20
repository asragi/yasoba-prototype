package core

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type TextDataYaml struct {
	Texts []TextEntry `yaml:"texts"`
}

type TextEntry struct {
	ID   string `yaml:"id"`
	Text string `yaml:"text"`
}

func LoadTextDataFromYaml(filePath string) (ServeTextDataFunc, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var yamlData TextDataYaml
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, err
	}

	dict := make(map[TextId]*TextData)
	for _, entry := range yamlData.Texts {
		dict[TextId(entry.ID)] = &TextData{
			Id:   TextId(entry.ID),
			Text: TextString(entry.Text),
		}
	}

	fmt.Printf("dict: %+v\n", dict)

	return func(id TextId) *TextData {
		if textData, ok := dict[id]; ok {
			return textData
		}
		panic("text not found: " + string(id))
	}, nil
}

// CreateServeTextDataFromYamlは、デフォルトのtext_data.yamlファイルからテキストデータを読み込みます
func CreateServeTextDataFromYaml() (ServeTextDataFunc, error) {
	return LoadTextDataFromYaml("data/text_data.yaml")
}
