package component

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle_skill"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
)

// BattleEnemyDisplay is a component that displays battle enemies.
type BattleEnemyDisplay struct {
	actorIds      []actor.ActorId
	actorGraphics map[actor.ActorId]BattleEnemyGraphicsInterface
}

func (d *BattleEnemyDisplay) SetDisappear(actorId actor.ActorId) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.SetDisappear()
}

func (d *BattleEnemyDisplay) SetDamage(actorId actor.ActorId, damage battle_skill.Damage) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.SetDamage(damage)
}

func (d *BattleEnemyDisplay) DoShake(actorId actor.ActorId) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.DoShake()
}

func (d *BattleEnemyDisplay) SetEmotion(actorId actor.ActorId, emotion BattleEmotionType) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.SetEmotion(emotion)
}

func (d *BattleEnemyDisplay) GetPosition(id actor.ActorId) *frontend.Vector {
	graphics, ok := d.actorGraphics[id]
	if !ok {
		return nil
	}
	return graphics.GetDefinitivePosition()
}

func (d *BattleEnemyDisplay) Update(parentCenterPosition *frontend.Vector) {
	for _, id := range d.actorIds {
		graphics := d.actorGraphics[id]
		graphics.Update(parentCenterPosition)
	}
}

func (d *BattleEnemyDisplay) Draw(drawFunc frontend.DrawFunc) {
	for _, id := range d.actorIds {
		graphics := d.actorGraphics[id]
		graphics.Draw(drawFunc)
	}
}

type BattleDisplayArgs struct {
	ActorId  actor.ActorId
	EnemyId  core.EnemyId
	Position *frontend.Vector
}

func ToDisplayArgs(
	enemyIdPair []*core.EnemyIdPair,
	enemySettings []*core.EnemySetting,
) []*BattleDisplayArgs {
	enemySettingMap := func() map[core.EnemyId][]*core.EnemySetting {
		result := make(map[core.EnemyId][]*core.EnemySetting)
		for _, setting := range enemySettings {
			result[setting.EnemyId] = append(result[setting.EnemyId], setting)
		}
		return result
	}()
	result := make([]*BattleDisplayArgs, len(enemyIdPair))
	enemyIndex := func() map[core.EnemyId]int {
		result := make(map[core.EnemyId]int)
		for _, pair := range enemyIdPair {
			result[pair.EnemyId] = 0
		}
		return result
	}()
	for i, id := range enemyIdPair {
		actorId := id.ActorId
		enemyId := id.EnemyId
		index := enemyIndex[enemyId]
		setting := enemySettingMap[enemyId][index]
		result[i] = &BattleDisplayArgs{
			ActorId:  actorId,
			EnemyId:  enemyId,
			Position: setting.Position,
		}
		enemyIndex[enemyId]++
	}
	return result
}

type NewBattleEnemyDisplayFunc func([]*BattleDisplayArgs, frontend.Depth) *BattleEnemyDisplay

func CreateNewBattleEnemyDisplay(
	newBattleActorGraphics NewBattleEnemyGraphicsFunc,
) NewBattleEnemyDisplayFunc {
	return func(
		enemies []*BattleDisplayArgs,
		depth frontend.Depth,
	) *BattleEnemyDisplay {
		actorIds := make([]actor.ActorId, len(enemies))
		actorGraphics := map[actor.ActorId]BattleEnemyGraphicsInterface{}
		for i, enemy := range enemies {
			graphics := newBattleActorGraphics(
				enemy.Position,
				frontend.PivotCenter,
				depth,
				enemy.EnemyId,
			)
			actorGraphics[enemy.ActorId] = graphics
			actorIds[i] = enemy.ActorId
		}
		return &BattleEnemyDisplay{
			actorGraphics: actorGraphics,
			actorIds:      actorIds,
		}
	}
}
