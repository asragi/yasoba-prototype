package enemy

import (
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/component"
	battleemotion "github.com/asragi/yasoba-prototype/component/battle/emotion"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/enemy"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleEnemyGraphics struct {
	emotion          battleemotion.Queued
	animation        map[battleemotion.BattleEmotionType]*widget.Animation
	displayDamage    *component.DisplayDamage
	shake            *frontend.EmitShake
	disappearShader  *frontend.Shader
	parentPosition   *frontend.Vector
	relativePosition *frontend.Vector
}

type BattleEnemyGraphicsInterface interface {
	SetDamage(battleSkill.Damage)
	DoShake()
	widget.PositionUpdater
	widget.Drawer
	SetEmotion(battleemotion.BattleEmotionType)
	SetDisappear()
	GetDefinitivePosition() *frontend.Vector
}

func (g *BattleEnemyGraphics) GetDefinitivePosition() *frontend.Vector {
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

func (g *BattleEnemyGraphics) Update(parentCenterPosition *frontend.Vector) {
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

func (g *BattleEnemyGraphics) Draw(drawFunc frontend.DrawFunc) {
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
	*frontend.Vector,
	*frontend.Pivot,
	frontend.Depth,
	enemy.EnemyId,
) BattleEnemyGraphicsInterface

func NewBattleActorGraphics(
	resource frontend.ResourceManagerInterface,
	getEnemyGraphics GetEnemyGraphicsFunc,
	newDisplayDamage component.NewDisplayDamageFunc,
) NewBattleEnemyGraphicsFunc {
	return func(
		relativePosition *frontend.Vector,
		pivot *frontend.Pivot,
		depth frontend.Depth,
		enemyId enemy.EnemyId,
	) BattleEnemyGraphicsInterface {
		enemyGraphicsData := getEnemyGraphics(enemyId)
		animations := func() map[battleemotion.BattleEmotionType]*widget.Animation {
			result := map[battleemotion.BattleEmotionType]*widget.Animation{}
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
		return &BattleEnemyGraphics{
			emotion:          battleemotion.NewQueued(battleemotion.BattleEmotionNormal),
			animation:        animations,
			shake:            frontend.NewShake(),
			parentPosition:   frontend.VectorZero,
			relativePosition: relativePosition,
			displayDamage:    newDisplayDamage(),
			disappearShader:  resource.GetShader(frontend.ShaderDisappear),
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
