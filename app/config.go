package app

import (
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/view/common/constant"
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
		GameWidth:       constant.GameWidth,
		GameHeight:      constant.GameHeight,
		TextDataPath:    "assets/data/text_data.yaml",
		BattleSettingID: config.IdTest,
		BattleID:        battle.BattleIdTest001,
	}
}
