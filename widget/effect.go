package widget

import (
	"fmt"

	commonanimation "github.com/asragi/yasoba-prototype/common/animation"
	commontexture "github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	viewanimation "github.com/asragi/yasoba-prototype/view/common/animation"
)

type effectResource interface {
	GetAnimationData(commonanimation.ID) *commonanimation.AnimationData
	GetTexture(commontexture.ID) drawing.Image
}

type EffectId string

const (
	EffectIdLuneAttack EffectId = "lune_attack"
	EffectIdLuneFire   EffectId = "lune_fire"
	EffectIdExplode    EffectId = "explode"
)

type Effect struct {
	animation *viewanimation.Animation
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
	resource    effectResource
	effects     map[EffectId]*Effect
}

func NewEffectManager(
	serveEffect ServeEffectDataFunc,
	resource effectResource,
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
	animation := viewanimation.New(
		func(
			relativePosition *drawing.Vector,
			pivot *drawing.Pivot,
			depth drawing.Depth,
			image drawing.Image,
		) viewanimation.Sprite {
			return NewImage(relativePosition, pivot, depth, image)
		},
		position,
		drawing.PivotCenter,
		drawing.DepthEffect,
		m.resource.GetTexture(animationData.TextureID),
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
	AnimationId commonanimation.ID
}

type ServeEffectDataFunc func(EffectId) *EffectData

func CreateServeEffectData() ServeEffectDataFunc {
	dict := map[EffectId]*EffectData{
		EffectIdLuneAttack: {
			EffectId:    EffectIdLuneAttack,
			AnimationId: commonanimation.BattleEffectImpact,
		},
		EffectIdLuneFire: {
			EffectId:    EffectIdLuneFire,
			AnimationId: commonanimation.BattleEffectFire,
		},
		EffectIdExplode: {
			EffectId:    EffectIdExplode,
			AnimationId: commonanimation.BattleEffectExplode,
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
