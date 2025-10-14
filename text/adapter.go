package text

import (
	"os"

	"gopkg.in/yaml.v3"
)

type textDataYAML struct {
	Texts []textEntry `yaml:"texts"`
}

type textEntry struct {
	ID   string `yaml:"id"`
	Text string `yaml:"text"`
}

// LoadFromYAML loads text data from the provided YAML file path.
func LoadFromYAML(filePath string) (ServeTextDataFunc, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var yamlData textDataYAML
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, err
	}

	dict := make(map[TextId]*Data)
	for _, entry := range yamlData.Texts {
		dict[TextId(entry.ID)] = &Data{
			Id:   TextId(entry.ID),
			Text: String(entry.Text),
		}
	}

	return func(id TextId) *Data {
		if textData, ok := dict[id]; ok {
			return textData
		}
		panic("text not found: " + string(id))
	}, nil
}
