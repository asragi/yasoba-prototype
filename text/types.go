package text

// TextId identifies a localized text asset.
type TextId string

const (
	TextIdBattleCommandAttack  TextId = "battle_command_attack"
	TextIdBattleCommandFire    TextId = "battle_command_fire"
	TextIdBattleCommandThunder TextId = "battle_command_thunder"
	TextIdBattleCommandBarrier TextId = "battle_command_barrier"
	TextIdBattleCommandWind    TextId = "battle_command_wind"
	TextIdBattleCommandFocus   TextId = "battle_command_focus"
	TextIdBattleCommandDefend  TextId = "battle_command_defend"
	TextIdLuneName             TextId = "lune_name"
	TextIdSunnyName            TextId = "sunny_name"
	TextIdPunchingBagName      TextId = "enemy_punching_bag_name"
	TextIdLuneAttackDesc       TextId = "lune_attack_desc"
	TextIdLuneFireDesc         TextId = "lune_fire_desc"
	TextIdCombinationThunder   TextId = "combination_thunder"
	TextIdEnemyBeaten          TextId = "enemy_beaten_desc"
	TextIdBattleWin            TextId = "battle_win"
	TextIdBattleLose           TextId = "battle_lose"
	TextIdBattleDialogText     TextId = "battle_dialog_text"
)

// ServeTextDataFunc loads text data by ID.
type ServeTextDataFunc func(id TextId) *Data

// String is a wrapper around text string content.
type String string

// String returns the underlying string value.
func (t String) String() string {
	return string(t)
}

// Data stores text metadata for lookup.
type Data struct {
	Id   TextId
	Text String
}
