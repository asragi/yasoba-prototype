package enemy

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/config"
	"github.com/asragi/yasoba-prototype/battle/enemy"
	"github.com/asragi/yasoba-prototype/battle/setup"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/emotion"
)

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

func (d *BattleEnemyDisplay) SetDamage(actorId actor.ActorId, damage skill.Damage) {
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

func (d *BattleEnemyDisplay) SetEmotion(actorId actor.ActorId, emotion emotion.BattleEmotionType) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.SetEmotion(emotion)
}

func (d *BattleEnemyDisplay) GetPosition(id actor.ActorId) *drawing.Vector {
	graphics, ok := d.actorGraphics[id]
	if !ok {
		return nil
	}
	return graphics.GetDefinitivePosition()
}

func (d *BattleEnemyDisplay) Update(parentCenterPosition *drawing.Vector) {
	for _, id := range d.actorIds {
		graphics := d.actorGraphics[id]
		graphics.Update(parentCenterPosition)
	}
}

func (d *BattleEnemyDisplay) Draw(drawFunc drawing.DrawFunc) {
	for _, id := range d.actorIds {
		graphics := d.actorGraphics[id]
		graphics.Draw(drawFunc)
	}
}

type BattleDisplayArgs struct {
	ActorId  actor.ActorId
	EnemyId  enemy.EnemyId
	Position *drawing.Vector
}

func ToDisplayArgs(
	enemyIdPair []*setup.EnemyIdPair,
	enemySettings []*config.EnemySetting,
) []*BattleDisplayArgs {
	enemySettingMap := func() map[enemy.EnemyId][]*config.EnemySetting {
		result := make(map[enemy.EnemyId][]*config.EnemySetting)
		for _, setting := range enemySettings {
			result[setting.EnemyId] = append(result[setting.EnemyId], setting)
		}
		return result
	}()
	result := make([]*BattleDisplayArgs, len(enemyIdPair))
	enemyIndex := func() map[enemy.EnemyId]int {
		result := make(map[enemy.EnemyId]int)
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

type NewBattleEnemyDisplayFunc func([]*BattleDisplayArgs, drawing.Depth) *BattleEnemyDisplay

func CreateNewBattleEnemyDisplay(
	newBattleActorGraphics NewBattleEnemyGraphicsFunc,
) NewBattleEnemyDisplayFunc {
	return func(
		enemies []*BattleDisplayArgs,
		depth drawing.Depth,
	) *BattleEnemyDisplay {
		actorIds := make([]actor.ActorId, len(enemies))
		actorGraphics := map[actor.ActorId]BattleEnemyGraphicsInterface{}
		for i, enemy := range enemies {
			graphics := newBattleActorGraphics(
				enemy.Position,
				drawing.PivotCenter,
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
