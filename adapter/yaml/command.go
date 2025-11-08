package yaml

import (
	"os"

	"github.com/asragi/yasoba-prototype/battle/command"
	"gopkg.in/yaml.v3"
)

type commandYAML struct {
	Commands []commandEntry `yaml:"commands"`
}

type commandEntry struct {
	ID   string `yaml:"id"`
	Cost int    `yaml:"cost"`
}

// CommandModelsFromYAML loads battle commands from the provided YAML file path.
func CommandModelsFromYAML(path string) []command.Model {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var raw commandYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		panic(err)
	}

	models := make([]command.Model, 0, len(raw.Commands))
	for _, entry := range raw.Commands {
		models = append(models, command.NewModel(entry.ID, command.Cost(entry.Cost)))
	}
	return models
}
