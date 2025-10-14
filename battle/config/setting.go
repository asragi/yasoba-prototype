package config

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/enemy"
)

type EnemySetting struct {
	EnemyId  enemy.EnemyId
	Position *frontend.Vector
}

type Id string

const (
	IdTest       Id = "test"
	IdTripleTest Id = "triple_test"
)

type Setting struct {
	Enemies []*EnemySetting
}

type ServeFunc func(Id) *Setting

func NewServer() ServeFunc {
	dict := make(map[Id]*Setting)
	dict[IdTest] = &Setting{
		Enemies: []*EnemySetting{
			{
				EnemyId:  enemy.EnemyPunchingBagId,
				Position: frontend.VectorZero,
			},
		},
	}
	dict[IdTripleTest] = &Setting{
		Enemies: []*EnemySetting{
			{
				EnemyId:  enemy.EnemyPunchingBagId,
				Position: &frontend.Vector{X: -100},
			},
			{
				EnemyId:  enemy.EnemyPunchingBagId,
				Position: frontend.VectorZero,
			},
			{
				EnemyId:  enemy.EnemyPunchingBagId,
				Position: &frontend.Vector{X: 100},
			},
		},
	}
	return func(id Id) *Setting {
		return dict[id]
	}
}
