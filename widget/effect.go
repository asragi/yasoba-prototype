package widget

import (
	"fmt"

	anim "github.com/asragi/yasoba-prototype/animation"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
)

type EffectId string

const (
	EffectIdLuneAttack EffectId = "lune_attack"
	EffectIdLuneFire   EffectId = "lune_fire"
	EffectIdExplode    EffectId = "explode"
)

type Effect struct {
	animation *anim.Animation
}

func (e *Effect) Update() {
	e.animation.Update(drawing.VectorZero)
}

func (e *Effect) Draw(drawFunc drawing.DrawFunc) {
	if e.animation.IsEnd() {
		return
	}
	e.animation.Draw(drawFunc)
}

type EffectManager struct {
	serveEffect ServeEffectDataFunc
	resource    *frontend.ResourceManager
	effects     map[EffectId]*Effect
}

func NewEffectManager(
	serveEffect ServeEffectDataFunc,
	resource *frontend.ResourceManager,
) *EffectManager {
	return &EffectManager{
		serveEffect: serveEffect,
		resource:    resource,
		effects:     map[EffectId]*Effect{},
	}
}

type ServeParentPosition func() *drawing.Vector

func (m *EffectManager) CallEffect(
	effectId EffectId,
	position *drawing.Vector,
) {
	effectData := m.serveEffect(effectId)
	animationData := m.resource.GetAnimationData(effectData.AnimationId)
	// TODO: インゲーム中にメモリ確保し続けるのは良くないためプールするなどの対策を行う
	animation := anim.New(
		func(
			relativePosition *drawing.Vector,
			pivot *drawing.Pivot,
			depth drawing.Depth,
			image drawing.Image,
		) anim.Sprite {
			return NewImage(relativePosition, pivot, depth, image)
		},
		position,
		drawing.PivotCenter,
		drawing.DepthEffect,
		m.resource.GetTexture(animationData.TextureId),
		animationData,
		nil,
	)
	m.effects[effectId] = &Effect{
		animation: animation,
	}
}

func (m *EffectManager) Update() {
	for _, effect := range m.effects {
		effect.Update()
	}
}

func (m *EffectManager) Draw(drawFunc drawing.DrawFunc) {
	for _, effect := range m.effects {
		effect.Draw(drawFunc)
	}
}

type EffectData struct {
	EffectId    EffectId
	AnimationId frontend.AnimationId
}

type ServeEffectDataFunc func(EffectId) *EffectData

func CreateServeEffectData() ServeEffectDataFunc {
	dict := map[EffectId]*EffectData{
		EffectIdLuneAttack: {
			EffectId:    EffectIdLuneAttack,
			AnimationId: frontend.AnimationBattleEffectImpact,
		},
		EffectIdLuneFire: {
			EffectId:    EffectIdLuneFire,
			AnimationId: frontend.AnimationBattleEffectFire,
		},
		EffectIdExplode: {
			EffectId:    EffectIdExplode,
			AnimationId: frontend.AnimationBattleEffectExplode,
		},
	}
	return func(effectId EffectId) *EffectData {
		data, ok := dict[effectId]
		if !ok {
			panic(fmt.Sprintf("effect data not found: %s", effectId))
		}
		return data
	}
}
