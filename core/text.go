package core

type TextId string

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

	register("battle_command_attack", "Attack")
	register("battle_command_fire", "Fire")
	register("battle_command_thunder", "Thunder")
	register("battle_command_barrier", "Barrier")
	register("battle_command_wind", "Wind")
	register("battle_command_focus", "Focus")
	register("battle_command_defend", "Defend")

	return func(id TextId) *TextData {
		return dict[id]
	}
}
