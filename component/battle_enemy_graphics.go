package component

import (
	"github.com/asragi/yasoba-prototype/battle_skill"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

// BattleEnemyGraphics is a component that displays a battle enemy.
type BattleEnemyGraphics struct {
	emotion          queuedEmotion
	animation        map[BattleEmotionType]*widget.Animation
	displayDamage    *DisplayDamage
	shake            *frontend.EmitShake
	disappearShader  *frontend.Shader
	parentPosition   *frontend.Vector
	relativePosition *frontend.Vector
}

type BattleEnemyGraphicsInterface interface {
	SetDamage(battle_skill.Damage)
	DoShake()
	widget.PositionUpdater
	widget.Drawer
	SetEmotion(BattleEmotionType)
	SetDisappear()
	GetDefinitivePosition() *frontend.Vector
}

func (g *BattleEnemyGraphics) GetDefinitivePosition() *frontend.Vector {
	return g.parentPosition.Add(g.relativePosition)
}

func (g *BattleEnemyGraphics) getCurrentAnimation() *widget.Animation {
	animation, ok := g.animation[g.emotion.current()]
	if !ok {
		return g.animation[BattleEmotionNormal]
	}
	return animation
}

func (g *BattleEnemyGraphics) SetDamage(damage battle_skill.Damage) {
	g.displayDamage.DisplayDamage(damage)
}

func (g *BattleEnemyGraphics) DoShake() {
	g.shake.Shake(frontend.ShakeDefaultAmplitude, frontend.ShakeDefaultPeriod)
}

func (g *BattleEnemyGraphics) Update(parentCenterPosition *frontend.Vector) {
	g.shake.Update()
	g.displayDamage.Update(parentCenterPosition.Add(g.relativePosition))
	g.parentPosition = parentCenterPosition
	animation := g.emotion.apply(func(emotion BattleEmotionType) *widget.Animation {
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

func (g *BattleEnemyGraphics) SetEmotion(emotion BattleEmotionType) {
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
	core.EnemyId,
) BattleEnemyGraphicsInterface

func NewBattleActorGraphics(
	resource frontend.ResourceManagerInterface,
	getEnemyGraphics GetEnemyGraphicsFunc,
	newDisplayDamage NewDisplayDamageFunc,
) NewBattleEnemyGraphicsFunc {
	return func(
		relativePosition *frontend.Vector,
		pivot *frontend.Pivot,
		depth frontend.Depth,
		enemyId core.EnemyId,
	) BattleEnemyGraphicsInterface {
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
		return &BattleEnemyGraphics{
			emotion:          newQueuedEmotion(BattleEmotionNormal),
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
