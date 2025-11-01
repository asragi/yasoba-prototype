package enemy

import (
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/component"
	battleemotion "github.com/asragi/yasoba-prototype/component/battle/emotion"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/enemy"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleEnemyGraphics struct {
	emotion          battleemotion.Queued
	animation        map[battleemotion.BattleEmotionType]*widget.Animation
	displayDamage    *component.DisplayDamage
	shake            *frontend.EmitShake
	disappearShader  *drawing.Shader
	parentPosition   *drawing.Vector
	relativePosition *drawing.Vector
}

type BattleEnemyGraphicsInterface interface {
	SetDamage(battleSkill.Damage)
	DoShake()
	widget.PositionUpdater
	widget.Drawer
	SetEmotion(battleemotion.BattleEmotionType)
	SetDisappear()
	GetDefinitivePosition() *drawing.Vector
}

func (g *BattleEnemyGraphics) GetDefinitivePosition() *drawing.Vector {
	return g.parentPosition.Add(g.relativePosition)
}

func (g *BattleEnemyGraphics) getCurrentAnimation() *widget.Animation {
	animation, ok := g.animation[g.emotion.Current()]
	if !ok {
		return g.animation[battleemotion.BattleEmotionNormal]
	}
	return animation
}

func (g *BattleEnemyGraphics) SetDamage(damage battleSkill.Damage) {
	g.displayDamage.DisplayDamage(damage)
}

func (g *BattleEnemyGraphics) DoShake() {
	g.shake.Shake(frontend.ShakeDefaultAmplitude, frontend.ShakeDefaultPeriod)
}

func (g *BattleEnemyGraphics) Update(parentCenterPosition *drawing.Vector) {
	g.shake.Update()
	g.displayDamage.Update(parentCenterPosition.Add(g.relativePosition))
	g.parentPosition = parentCenterPosition
	animation := g.emotion.Apply(func(emotion battleemotion.BattleEmotionType) *widget.Animation {
		return g.animation[emotion]
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

func (g *BattleEnemyGraphics) SetEmotion(emotion battleemotion.BattleEmotionType) {
	g.emotion.Enqueue(emotion)
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
	newDisplayDamage component.NewDisplayDamageFunc,
) NewBattleEnemyGraphicsFunc {
	return func(
		relativePosition *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		enemyId enemy.EnemyId,
	) BattleEnemyGraphicsInterface {
		enemyGraphicsData := getEnemyGraphics(enemyId)
		animations := func() map[battleemotion.BattleEmotionType]*widget.Animation {
			result := map[battleemotion.BattleEmotionType]*widget.Animation{}
			for _, data := range enemyGraphicsData {
				texture := resource.GetTexture(data.texture)
				animation := resource.GetAnimationData(data.animation)
				anim := widget.NewAnimation(
					relativePosition,
					pivot,
					depth,
					texture,
					animation,
				)
				anim.SetRenderTargetFactory(func() drawing.Image {
					// TODO: Test
					return resource.NewEmptyImage(384, 288)
				})
				result[data.emotion] = anim
			}
			return result
		}()
		return &BattleEnemyGraphics{
			emotion:          battleemotion.NewQueued(battleemotion.BattleEmotionNormal),
			animation:        animations,
			shake:            frontend.NewShake(),
			parentPosition:   drawing.VectorZero,
			relativePosition: relativePosition,
			displayDamage:    newDisplayDamage(),
			disappearShader:  resource.GetShader(drawing.ShaderDisappear),
		}
	}
}

type BattleActorAnimationSet struct {
	emotion   battleemotion.BattleEmotionType
	texture   frontend.TextureId
	animation frontend.AnimationId
}

type GetEnemyGraphicsFunc func(enemy.EnemyId) []*BattleActorAnimationSet

func CreateGetEnemyGraphics() GetEnemyGraphicsFunc {
	dict := map[enemy.EnemyId][]*BattleActorAnimationSet{
		enemy.EnemyPunchingBagId: {
			{
				emotion:   battleemotion.BattleEmotionNormal,
				texture:   frontend.TextureMarshmallowNormal,
				animation: frontend.AnimationMarshmallowNormal,
			},
			{
				emotion:   battleemotion.BattleEmotionDamage,
				texture:   frontend.TextureMarshmallowDamage,
				animation: frontend.AnimationMarshmallowDamage,
			},
		},
	}
	return func(enemyId enemy.EnemyId) []*BattleActorAnimationSet {
		return dict[enemyId]
	}
}
