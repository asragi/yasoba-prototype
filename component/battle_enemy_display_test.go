package component

import (
	"testing"

	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/battle/setup"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/enemy"
)

func TestToDisplayArgs(t *testing.T) {
	type testCase struct {
		enemyIdPair   []*setup.EnemyIdPair
		enemySettings []*config.EnemySetting
		expect        []*BattleDisplayArgs
	}
	enemyId := enemy.EnemyId("enemy_id")
	anotherEnemyId := enemy.EnemyId("another_enemy_id")
	testCases := []testCase{
		{
			enemyIdPair: []*setup.EnemyIdPair{
				{
					EnemyId: enemyId,
					ActorId: "actor1",
				},
				{
					EnemyId: enemyId,
					ActorId: "actor2",
				},
				{
					EnemyId: anotherEnemyId,
					ActorId: "actor3",
				},
			},
			enemySettings: []*config.EnemySetting{
				{
					EnemyId:  enemyId,
					Position: &frontend.Vector{X: 1, Y: 2},
				},
				{
					EnemyId:  enemyId,
					Position: &frontend.Vector{X: 3, Y: 4},
				},
				{
					EnemyId:  anotherEnemyId,
					Position: &frontend.Vector{X: 5, Y: 6},
				},
			},
			expect: []*BattleDisplayArgs{
				{
					EnemyId:  enemyId,
					ActorId:  "actor1",
					Position: &frontend.Vector{X: 1, Y: 2},
				},
				{
					EnemyId:  enemyId,
					ActorId:  "actor2",
					Position: &frontend.Vector{X: 3, Y: 4},
				},
				{
					EnemyId:  anotherEnemyId,
					ActorId:  "actor3",
					Position: &frontend.Vector{X: 5, Y: 6},
				},
			},
		},
	}

	for i, tc := range testCases {
		actual := ToDisplayArgs(tc.enemyIdPair, tc.enemySettings)
		if len(actual) != len(tc.expect) {
			t.Errorf("testCases[%d] expect %v, but actual %v", i, tc.expect, actual)
		}
		for j, a := range actual {
			if a.EnemyId != tc.expect[j].EnemyId {
				t.Errorf("testCases[%d] expect %v, but actual %v", i, tc.expect[j].EnemyId, a.EnemyId)
			}
			if a.ActorId != tc.expect[j].ActorId {
				t.Errorf("testCases[%d] expect %v, but actual %v", i, tc.expect[j].ActorId, a.ActorId)
			}
			if a.Position.X != tc.expect[j].Position.X {
				t.Errorf("testCases[%d] expect %v, but actual %v", i, tc.expect[j].Position.X, a.Position.X)
			}
			if a.Position.Y != tc.expect[j].Position.Y {
				t.Errorf("testCases[%d] expect %v, but actual %v", i, tc.expect[j].Position.Y, a.Position.Y)
			}
		}
	}
}
