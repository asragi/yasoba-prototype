package core

type EnemyId string

const (
	EnemyPunchingBagId EnemyId = "punching-bag"
)

type ServeEnemyData func(id EnemyId) *EnemyData

type EnemyData struct {
	Id     EnemyId
	MaxHP  MaxHP
	Atk    ATK
	Mag    MAG
	Def    DEF
	Spd    SPD
	Skills []SkillId
}

func CreateEnemyServer() ServeEnemyData {
	dict := make(map[EnemyId]*EnemyData)
	dict[EnemyPunchingBagId] = &EnemyData{
		Id:     EnemyPunchingBagId,
		MaxHP:  1000,
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

type EnemyNameServer func(EnemyId) TextId

func CreateEnemyNameServer() EnemyNameServer {
	dict := make(map[EnemyId]TextId)
	dict[EnemyPunchingBagId] = TextIdPunchingBagName
	return func(id EnemyId) TextId {
		return dict[id]
	}
}
