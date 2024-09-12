package core

type EnemyId string

const (
	EnemyPunchingBagId EnemyId = "punching-bag"
)

type ServeEnemyData func(id EnemyId) *EnemyData

type EnemyData struct {
	Id                EnemyId
	Name              TextId
	AppearanceMessage TextId
	MaxHP             MaxHP
	Atk               ATK
	Mag               MAG
	Def               DEF
	Spd               SPD
	Skills            []SkillId
}

func CreateEnemyServer() ServeEnemyData {
	dict := make(map[EnemyId]*EnemyData)
	dict[EnemyPunchingBagId] = &EnemyData{
		Id:     EnemyPunchingBagId,
		Name:   TextIdPunchingBagName,
		MaxHP:  10000,
		Atk:    10,
		Mag:    10,
		Def:    10,
		Spd:    10,
		Skills: []SkillId{},
	}
	return func(id EnemyId) *EnemyData {
		return dict[id]
	}
}
