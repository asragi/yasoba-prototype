package widget

import "github.com/asragi/yasoba-prototype/frontend"

type EffectId string

const (
	EffectIdLuneAttack EffectId = "lune_attack"
)

type Effect struct {
	animation *Animation
}

func (e *Effect) Update() {
	e.animation.Update(frontend.VectorZero)
}

func (e *Effect) Draw(drawFunc frontend.DrawFunc) {
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

type ServeParentPosition func() *frontend.Vector

func (m *EffectManager) CallEffect(
	effectId EffectId,
	position *frontend.Vector,
) {
	effectData := m.serveEffect(effectId)
	animationData := m.resource.GetAnimationData(effectData.AnimationId)
	animation := NewAnimation(
		position,
		frontend.PivotCenter,
		frontend.DepthEffect,
		m.resource.GetTexture(animationData.TextureId),
		animationData,
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

func (m *EffectManager) Draw(drawFunc frontend.DrawFunc) {
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
	}
	return func(effectId EffectId) *EffectData {
		return dict[effectId]
	}
}
