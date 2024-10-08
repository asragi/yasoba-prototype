package core

import "strconv"

type CharacterId string

const (
	CharacterEmptyId CharacterId = "empty"
	CharacterLuneId  CharacterId = "lune"
	CharacterSunnyId CharacterId = "sunny"
)

type MaxHP int

func (h MaxHP) ToHP() HP {
	return HP(h)
}

type HP int

func (h HP) String() string {
	return strconv.Itoa(int(h))
}

type ATK int
type MAG int
type DEF int
type SPD int

type ServeCharacterFunc func(CharacterId) *CharacterData

func CreateCharacterServer() ServeCharacterFunc {
	dict := make(map[CharacterId]*CharacterData)
	dict[CharacterLuneId] = &CharacterData{
		Id:    CharacterLuneId,
		Name:  TextIdLuneName,
		MaxHP: 110,
		HP:    110,
		ATK:   6,
		MAG:   30,
		DEF:   6,
		SPD:   6,
	}
	dict[CharacterSunnyId] = &CharacterData{
		Id:    CharacterSunnyId,
		Name:  TextIdSunnyName,
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
	Name  TextId
	MaxHP MaxHP
	HP    HP
	ATK   ATK
	MAG   MAG
	DEF   DEF
	SPD   SPD
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
