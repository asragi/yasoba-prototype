package character

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/text"
)

type CharacterId string

const (
	CharacterEmptyId CharacterId = "empty"
	CharacterLuneId  CharacterId = "lune"
	CharacterSunnyId CharacterId = "sunny"
)

type ServeCharacterFunc func(CharacterId) *CharacterData

func CreateCharacterServer() ServeCharacterFunc {
	dict := make(map[CharacterId]*CharacterData)
	dict[CharacterLuneId] = &CharacterData{
		Id:    CharacterLuneId,
		Name:  text.TextIdLuneName,
		MaxHP: 110,
		HP:    110,
		ATK:   6,
		MAG:   30,
		DEF:   6,
		SPD:   6,
	}
	dict[CharacterSunnyId] = &CharacterData{
		Id:    CharacterSunnyId,
		Name:  text.TextIdSunnyName,
		MaxHP: 220,
		HP:    220,
		ATK:   22,
		MAG:   21,
		DEF:   28,
		SPD:   28,
	}
	return func(id CharacterId) *CharacterData {
		return dict[id]
	}
}

type CharacterData struct {
	Id    CharacterId
	Name  text.TextId
	MaxHP actor.MaxHP
	HP    actor.HP
	ATK   actor.ATK
	MAG   actor.MAG
	DEF   actor.DEF
	SPD   actor.SPD
}

type InitialMP int
type RecoverMP int
type MaxMP int

type MainActorData struct {
	InitialMP InitialMP
	RecoverMP RecoverMP
	MapMP     MaxMP
}

type ServeMainActorFunc func() *MainActorData
