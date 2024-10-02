package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

// BattleEnemyDisplay is a component that displays battle enemies.
type BattleEnemyDisplay struct {
	actorIds      []core.ActorId
	actorGraphics map[core.ActorId]*BattleEnemyGraphics
}

func (d *BattleEnemyDisplay) SetDisappear(actorId core.ActorId) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.SetDisappear()
}

func (d *BattleEnemyDisplay) SetDamage(actorId core.ActorId, damage core.Damage) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.SetDamage(damage)
}

func (d *BattleEnemyDisplay) DoShake(actorId core.ActorId) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.DoShake()
}

func (d *BattleEnemyDisplay) SetEmotion(actorId core.ActorId, emotion BattleEmotionType) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.SetEmotion(emotion)
}

func (d *BattleEnemyDisplay) GetPosition(id core.ActorId) *frontend.Vector {
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
	ActorId  core.ActorId
	EnemyId  core.EnemyId
	Position *frontend.Vector
}

type NewBattleEnemyDisplayFunc func([]*BattleDisplayArgs, frontend.Depth) *BattleEnemyDisplay

func CreateNewBattleEnemyDisplay(
	newBattleActorGraphics NewBattleEnemyGraphicsFunc,
) NewBattleEnemyDisplayFunc {
	return func(
		enemies []*BattleDisplayArgs,
		depth frontend.Depth,
	) *BattleEnemyDisplay {
		actorIds := make([]core.ActorId, len(enemies))
		actorGraphics := map[core.ActorId]*BattleEnemyGraphics{}
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

type BattleEnemyGraphics struct {
	currentEmotion   BattleEmotionType
	animation        map[BattleEmotionType]*widget.Animation
	displayDamage    *DisplayDamage
	shake            *frontend.EmitShake
	disappearShader  *frontend.Shader
	parentPosition   *frontend.Vector
	relativePosition *frontend.Vector
}

func (g *BattleEnemyGraphics) GetDefinitivePosition() *frontend.Vector {
	return g.parentPosition.Add(g.relativePosition)
}

func (g *BattleEnemyGraphics) getCurrentAnimation() *widget.Animation {
	animation, ok := g.animation[g.currentEmotion]
	if !ok {
		return g.animation[BattleEmotionNormal]
	}
	return animation
}

func (g *BattleEnemyGraphics) SetDamage(damage core.Damage) {
	g.displayDamage.DisplayDamage(damage)
}

func (g *BattleEnemyGraphics) DoShake() {
	const amplitude = 3
	const period = 12
	g.shake.Shake(amplitude, period)
}

func (g *BattleEnemyGraphics) Update(parentCenterPosition *frontend.Vector) {
	g.shake.Update()
	g.displayDamage.Update(parentCenterPosition)
	g.parentPosition = parentCenterPosition
	position := parentCenterPosition.Add(g.shake.Delta())
	g.getCurrentAnimation().Update(position)
}

func (g *BattleEnemyGraphics) Draw(drawFunc frontend.DrawFunc) {
	g.displayDamage.Draw(drawFunc)
	g.getCurrentAnimation().Draw(drawFunc)
}

func (g *BattleEnemyGraphics) SetEmotion(emotion BattleEmotionType) {
	g.currentEmotion = emotion
	g.getCurrentAnimation().Reset()
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
) *BattleEnemyGraphics

func NewBattleActorGraphics(
	resource *frontend.ResourceManager,
	getEnemyGraphics GetEnemyGraphicsFunc,
	newDisplayDamage NewDisplayDamageFunc,
) NewBattleEnemyGraphicsFunc {
	return func(
		relativePosition *frontend.Vector,
		pivot *frontend.Pivot,
		depth frontend.Depth,
		enemyId core.EnemyId,
	) *BattleEnemyGraphics {
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
			currentEmotion:   BattleEmotionNormal,
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
