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

func DefaultConfig() Config {
	return Config{
		GameWidth:       384,
		GameHeight:      288,
		TextDataPath:    "assets/data/text_data.yaml",
		BattleSettingID: config.IdTest,
		BattleID:        battle.BattleIdTest001,
	}
}
