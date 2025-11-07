package adapter

import (
	"os"

	"gopkg.in/yaml.v3"
)

type WaitFrameRecord struct {
	EventID string
	Frame   int
}

type waitFrameYAML struct {
	WaitFrames []struct {
		EventID string `yaml:"eventId"`
		Frame   int    `yaml:"frame"`
	} `yaml:"waitFrames"`
}

func WaitFrameRecordsFromYAML(filePath string) map[string]*WaitFrameRecord {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var yamlData waitFrameYAML
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		panic(err)
	}

	result := make(map[string]*WaitFrameRecord, len(yamlData.WaitFrames))
	for _, record := range yamlData.WaitFrames {
		result[record.EventID] = &WaitFrameRecord{
			EventID: record.EventID,
			Frame:   record.Frame,
		}
	}
	return result
}
