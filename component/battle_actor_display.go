package component

import (
	"fmt"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleActorDisplay struct {
	actorIds      []core.ActorId
	actorGraphics map[core.ActorId]BattleActorGraphics
}

func (d *BattleActorDisplay) Update(parentCenterPosition *frontend.Vector) {
	for _, id := range d.actorIds {
		graphics := d.actorGraphics[id]
		graphics.Update(parentCenterPosition)
	}
}

func (d *BattleActorDisplay) Draw(drawFunc frontend.DrawFunc) {
	for _, id := range d.actorIds {
		graphics := d.actorGraphics[id]
		graphics.Draw(drawFunc)
	}
}

type BattleDisplayArgs struct {
	ActorId  core.ActorId
	EnemyId  core.EnemyId
	Position *frontend.Vector
}

type NewBattleActorDisplayFunc func([]*BattleDisplayArgs, frontend.Depth) *BattleActorDisplay

func CreateNewBattleActorDisplay(
	newBattleActorGraphics NewBattleActorGraphicsFunc,
) NewBattleActorDisplayFunc {
	return func(
		enemies []*BattleDisplayArgs,
		depth frontend.Depth,
	) *BattleActorDisplay {
		actorIds := make([]core.ActorId, len(enemies))
		actorGraphics := map[core.ActorId]BattleActorGraphics{}
		for i, enemy := range enemies {
			graphics := newBattleActorGraphics(
				enemy.Position,
				frontend.PivotCenter,
				depth,
				enemy.EnemyId,
			)
			actorGraphics[enemy.ActorId] = *graphics
			actorIds[i] = enemy.ActorId
		}
		return &BattleActorDisplay{
			actorGraphics: actorGraphics,
			actorIds:      actorIds,
		}
	}
}

type BattleActorGraphics struct {
	currentEmotion BattleEmotionType
	animation      map[BattleEmotionType]*widget.Animation
}

func (g *BattleActorGraphics) getCurrentAnimation() *widget.Animation {
	animation, ok := g.animation[g.currentEmotion]
	if !ok {
		return g.animation[BattleEmotionNormal]
	}
	return animation
}

func (g *BattleActorGraphics) Update(parentCenterPosition *frontend.Vector) {
	g.getCurrentAnimation().Update(parentCenterPosition)
}

func (g *BattleActorGraphics) Draw(drawFunc frontend.DrawFunc) {
	fmt.Printf("draw actor graphics\n")
	g.getCurrentAnimation().Draw(drawFunc)
}

func (g *BattleActorGraphics) SetEmotion(emotion BattleEmotionType) {
	g.currentEmotion = emotion
	g.getCurrentAnimation().Reset()
}

type NewBattleActorGraphicsFunc func(
	*frontend.Vector,
	*frontend.Pivot,
	frontend.Depth,
	core.EnemyId,
) *BattleActorGraphics

func NewBattleActorGraphics(
	resource *frontend.ResourceManager,
	getEnemyGraphics GetEnemyGraphicsFunc,
) NewBattleActorGraphicsFunc {
	return func(
		relativePosition *frontend.Vector,
		pivot *frontend.Pivot,
		depth frontend.Depth,
		enemyId core.EnemyId,
	) *BattleActorGraphics {
		enemyGraphicsData := getEnemyGraphics(enemyId)
		animations := func() map[BattleEmotionType]*widget.Animation {
			result := map[BattleEmotionType]*widget.Animation{}
			for _, data := range enemyGraphicsData {
				texture := resource.GetTexture(data.texture)
				animation := resource.GetAnimationData(data.animation)
				result[data.emotion] = widget.NewAnimation(
					relativePosition,
					pivot,
					depth,
					texture,
					animation,
				)
			}
			return result
		}()
		return &BattleActorGraphics{
			animation: animations,
		}
	}
}

type BattleActorAnimationSet struct {
	emotion   BattleEmotionType
	texture   frontend.TextureId
	animation frontend.AnimationId
}

type GetEnemyGraphicsFunc func(core.EnemyId) []*BattleActorAnimationSet

func CreateGetEnemyGraphics() GetEnemyGraphicsFunc {
	dict := map[core.EnemyId][]*BattleActorAnimationSet{
		core.EnemyPunchingBagId: {
			{
				emotion:   BattleEmotionNormal,
				texture:   frontend.TextureMarshmallowNormal,
				animation: frontend.AnimationMarshmallowNormal,
			},
			{
				emotion:   BattleEmotionDamage,
				texture:   frontend.TextureMarshmallowDamage,
				animation: frontend.AnimationMarshmallowDamage,
			},
		},
	}
	return func(enemyId core.EnemyId) []*BattleActorAnimationSet {
		return dict[enemyId]
	}
}
