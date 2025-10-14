package enemydata

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/skilldata"
	"github.com/asragi/yasoba-prototype/text"
)

type EnemyId string

const (
	EnemyPunchingBagId EnemyId = "punching-bag"
)

type ServeEnemyData func(id EnemyId) *EnemyData

type EnemyData struct {
	Id     EnemyId
	MaxHP  actor.MaxHP
	Atk    actor.ATK
	Mag    actor.MAG
	Def    actor.DEF
	Spd    actor.SPD
	Skills []skilldata.SkillId
}

func CreateEnemyServer() ServeEnemyData {
	dict := make(map[EnemyId]*EnemyData)
	dict[EnemyPunchingBagId] = &EnemyData{
		Id:     EnemyPunchingBagId,
		MaxHP:  1000,
		Atk:    25,
		Mag:    10,
		Def:    10,
		Spd:    10,
		Skills: []skilldata.SkillId{},
	}
	return func(id EnemyId) *EnemyData {
		return dict[id]
	}
}

type NameServer func(EnemyId) text.TextId

func CreateNameServer() NameServer {
	dict := make(map[EnemyId]text.TextId)
	dict[EnemyPunchingBagId] = text.TextIdPunchingBagName
	return func(id EnemyId) text.TextId {
		return dict[id]
	}
}
