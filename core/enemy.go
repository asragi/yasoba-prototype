package core

type EnemyId string

type ServeEnemyData func(id EnemyId) *EnemyData

type EnemyData struct {
	Id                EnemyId
	Name              TextId
	AppearanceMessage TextId
	HP                HP
	Atk               ATK
	Mag               MAG
	Def               DEF
	Spd               SPD
	Skills            []SkillId
}
