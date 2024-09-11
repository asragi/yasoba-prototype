package core

type CharacterId string

const (
	CharacterLuneId  CharacterId = "lune"
	CharacterSunnyId CharacterId = "sunny"
)

type MaxHP int
type HP int
type ATK int
type MAG int
type DEF int
type SPD int

type ServeCharacterFunc func(CharacterId) *CharacterData

type CharacterData struct {
	Id     CharacterId
	Name   TextId
	MaxHP  MaxHP
	HP     HP
	ATK    ATK
	MAG    MAG
	DEF    DEF
	SPD    SPD
	Skills []SkillId
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
