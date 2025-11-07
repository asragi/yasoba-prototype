package app

import (
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/config"
)

type Config struct {
	GameWidth       int
	GameHeight      int
	TextDataPath    string
	BattleSettingID config.Id
	BattleID        battle.BattleId
}

var sizeMap = map[int]map[int]int{
	0: {
		0: 384,
		1: 288,
	},
	1: {
		0: 416,
		1: 312,
	},
}

func DefaultConfig() Config {
	sizeType := 1
	return Config{
		GameWidth:       sizeMap[sizeType][0],
		GameHeight:      sizeMap[sizeType][1],
		TextDataPath:    "assets/data/text_data.yaml",
		BattleSettingID: config.IdTest,
		BattleID:        battle.BattleIdTest001,
	}
}
