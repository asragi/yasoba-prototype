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

	register("battle_command_attack", "こうげき")
	register("battle_command_fire", "まほう")
	register("battle_command_thunder", "まほう２")
	register("battle_command_barrier", "まほう３")
	register("battle_command_wind", "まほう４")
	register("battle_command_focus", "おちつく")
	register("battle_command_defend", "まもる")

	return func(id TextId) *TextData {
		return dict[id]
	}
}
