package core

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
		MaxHP: 40,
		HP:    40,
		ATK:   22,
		MAG:   30,
		DEF:   6,
		SPD:   6,
	}
	dict[CharacterSunnyId] = &CharacterData{
		Id:    CharacterSunnyId,
		Name:  TextIdSunnyName,
		MaxHP: 100,
		HP:    100,
		ATK:   10,
		MAG:   10,
		DEF:   10,
		SPD:   10,
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
