package core

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
)

type ServeTextDataFunc func(id TextId) *TextData

type TextData struct {
	Id   TextId
	Text string
}

func CreateServeTextData() ServeTextDataFunc {
	dict := map[TextId]*TextData{}

	register := func(id TextId, text string) {
		dict[id] = &TextData{
			Id:   id,
			Text: text,
		}
	}

	register(TextIdBattleCommandAttack, "こうげき")
	register(TextIdBattleCommandFire, "ファイア")
	register(TextIdBattleCommandThunder, "サンダー")
	register(TextIdBattleCommandBarrier, "バリア")
	register(TextIdBattleCommandWind, "ウィンド")
	register(TextIdBattleCommandFocus, "おちつく")
	register(TextIdBattleCommandDefend, "まもる")
	register(TextIdLuneName, "ルーネ")
	register(TextIdSunnyName, "サニー")
	register(TextIdPunchingBagName, "マシュマロス")

	return func(id TextId) *TextData {
		return dict[id]
	}
}
