package core

import "fmt"

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

type ServeTextDataFunc func(id TextId) *TextData

type TextString string

func (t TextString) String() string {
	return string(t)
}

type TextData struct {
	Id   TextId
	Text TextString
}

func CreateServeTextData() ServeTextDataFunc {
	dict := map[TextId]*TextData{}

	register := func(id TextId, text TextString) {
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
	register(TextIdLuneAttackDesc, "たいあたりした！")
	register(TextIdLuneFireDesc, "ファイアをとなえた！")
	register(TextIdCombinationThunder, "サニーはルーネのまほうにあわせた！\nおおきなばくはつがおこった！")
	register(TextIdEnemyBeaten, "てきをやっつけた！")
	register(TextIdBattleWin, "しょうりした！")
	register(TextIdBattleLose, "やられてしまった……")
	register(TextIdBattleDialogText, "かかってこい！")

	return func(id TextId) *TextData {
		if _, ok := dict[id]; !ok {
			panic(fmt.Sprintf("text not found: %s", id))
		}
		return dict[id]
	}
}
