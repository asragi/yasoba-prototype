package enemy

import (
	"github.com/asragi/yasoba-prototype/adapter/ebiten/frontend"
	"github.com/asragi/yasoba-prototype/battle/enemy"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	commonanimation "github.com/asragi/yasoba-prototype/common/animation"
	"github.com/asragi/yasoba-prototype/common/emotion"
	commontexture "github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/damage"
	viewemotion "github.com/asragi/yasoba-prototype/view/battle/emotion"
	anim "github.com/asragi/yasoba-prototype/view/common/animation"
	"github.com/asragi/yasoba-prototype/view/common/shake"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleEnemyGraphics struct {
	emotion          viewemotion.Queued
	animation        map[emotion.EmotionType]*anim.Animation
	displayDamage    *damage.DisplayDamage
	shake            *shake.EmitShake
	disappearShader  *drawing.Shader
	parentPosition   *drawing.Vector
	relativePosition *drawing.Vector
}

type BattleEnemyGraphicsInterface interface {
	SetDamage(battleSkill.Damage)
	DoShake()
	widget.PositionUpdater
	widget.Drawer
	SetEmotion(emotion.EmotionType)
	SetDisappear()
	GetDefinitivePosition() *drawing.Vector
}

func (g *BattleEnemyGraphics) GetDefinitivePosition() *drawing.Vector {
	return g.parentPosition.Add(g.relativePosition)
}

func (g *BattleEnemyGraphics) getCurrentAnimation() *anim.Animation {
	animation, ok := g.animation[g.emotion.Current()]
	if !ok {
		return g.animation[emotion.EmotionNormal]
	}
	return animation
}

func (g *BattleEnemyGraphics) SetDamage(damage battleSkill.Damage) {
	g.displayDamage.DisplayDamage(damage)
}

func (g *BattleEnemyGraphics) DoShake() {
	g.shake.Shake(shake.ShakeDefaultAmplitude, shake.ShakeDefaultPeriod)
}

func (g *BattleEnemyGraphics) Update(parentCenterPosition *drawing.Vector) {
	g.shake.Update()
	g.displayDamage.Update(parentCenterPosition.Add(g.relativePosition))
	g.parentPosition = parentCenterPosition
	animation := g.emotion.Apply(func(value emotion.EmotionType) *anim.Animation {
		return g.animation[value]
	})
	if animation == nil {
		return
	}
	position := parentCenterPosition.Add(g.shake.Delta())
	animation.Update(position)
}

func (g *BattleEnemyGraphics) Draw(drawFunc drawing.DrawFunc) {
	g.displayDamage.Draw(drawFunc)
	g.getCurrentAnimation().Draw(drawFunc)
}

func (g *BattleEnemyGraphics) SetEmotion(value emotion.EmotionType) {
	g.emotion.Enqueue(value)
}

func (g *BattleEnemyGraphics) SetDisappear() {
	for _, anim := range g.animation {
		anim.SetShader(g.disappearShader)
	}
}

type NewBattleEnemyGraphicsFunc func(
	*drawing.Vector,
	*drawing.Pivot,
	drawing.Depth,
	enemy.EnemyId,
) BattleEnemyGraphicsInterface

func NewBattleActorGraphics(
	resource frontend.ResourceManagerInterface,
	getEnemyGraphics GetEnemyGraphicsFunc,
	newDisplayDamage damage.NewDisplayDamageFunc,
	screenWidth int,
	screenHeight int,
) NewBattleEnemyGraphicsFunc {
	return func(
		relativePosition *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		enemyId enemy.EnemyId,
	) BattleEnemyGraphicsInterface {
		enemyGraphicsData := getEnemyGraphics(enemyId)
		animations := func() map[emotion.EmotionType]*anim.Animation {
			result := map[emotion.EmotionType]*anim.Animation{}
			for _, data := range enemyGraphicsData {
				texture := resource.GetTexture(data.texture)
				animation := resource.GetAnimationData(data.animation)
				animSprite := anim.New(
					func(
						relativePosition *drawing.Vector,
						pivot *drawing.Pivot,
						depth drawing.Depth,
						image drawing.Image,
					) anim.Sprite {
						return widget.NewImage(relativePosition, pivot, depth, image)
					},
					relativePosition,
					pivot,
					depth,
					texture,
					animation,
					func() drawing.Image {
						return resource.NewEmptyImage(screenWidth, screenHeight)
					},
				)
				result[data.emotion] = animSprite
			}
			return result
		}()
		return &BattleEnemyGraphics{
			emotion:          viewemotion.NewQueued(emotion.EmotionNormal),
			animation:        animations,
			shake:            shake.NewShake(),
			parentPosition:   drawing.VectorZero,
			relativePosition: relativePosition,
			displayDamage:    newDisplayDamage(),
			disappearShader:  resource.GetShader(drawing.ShaderDisappear),
		}
	}
}

type BattleActorAnimationSet struct {
	emotion   emotion.EmotionType
	texture   commontexture.ID
	animation commonanimation.ID
}

type GetEnemyGraphicsFunc func(enemy.EnemyId) []*BattleActorAnimationSet

func CreateGetEnemyGraphics() GetEnemyGraphicsFunc {
	dict := map[enemy.EnemyId][]*BattleActorAnimationSet{
		enemy.EnemyPunchingBagId: {
			{
				emotion:   emotion.EmotionNormal,
				texture:   commontexture.MarshmallowNormal,
				animation: commonanimation.MarshmallowNormal,
			},
			{
				emotion:   emotion.EmotionDamage,
				texture:   commontexture.MarshmallowDamage,
				animation: commonanimation.MarshmallowDamage,
			},
		},
	}
	return func(enemyId enemy.EnemyId) []*BattleActorAnimationSet {
		return dict[enemyId]
	}
}
