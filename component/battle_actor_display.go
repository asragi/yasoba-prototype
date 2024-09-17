package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleActorDisplay struct {
	actorIds      []core.ActorId
	actorGraphics map[core.ActorId]BattleActorGraphics
}

func (d *BattleActorDisplay) DoShake(actorId core.ActorId) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.DoShake()
}

func (d *BattleActorDisplay) SetEmotion(actorId core.ActorId, emotion BattleEmotionType) {
	graphics, ok := d.actorGraphics[actorId]
	if !ok {
		return
	}
	graphics.SetEmotion(emotion)
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
	shake          *frontend.EmitShake
}

func (g *BattleActorGraphics) getCurrentAnimation() *widget.Animation {
	animation, ok := g.animation[g.currentEmotion]
	if !ok {
		return g.animation[BattleEmotionNormal]
	}
	return animation
}

func (g *BattleActorGraphics) DoShake() {
	const amplitude = 3
	const period = 12
	g.shake.Shake(amplitude, period)
}

func (g *BattleActorGraphics) Update(parentCenterPosition *frontend.Vector) {
	g.shake.Update()
	position := parentCenterPosition.Add(g.shake.Delta())
	g.getCurrentAnimation().Update(position)
}

func (g *BattleActorGraphics) Draw(drawFunc frontend.DrawFunc) {
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
			shake:     frontend.NewShake(),
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
