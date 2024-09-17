package widget

import "github.com/asragi/yasoba-prototype/frontend"

type Effect struct {
	animation *Animation
}

func NewEffect(
	animation *Animation,
) *Effect {
	return &Effect{
		animation: animation,
	}
}

func (e *Effect) Update(parentPosition *frontend.Vector) {
	e.animation.Update(parentPosition)
}

func (e *Effect) Draw(drawFunc frontend.DrawFunc) {
	e.animation.Draw(drawFunc)
}

type EffectManager struct {
	effects []*Effect
}

func NewEffectManager() *EffectManager {
	return &EffectManager{
		effects: []*Effect{},
	}
}

func (m *EffectManager) AddEffect(effect *Effect) {
	m.effects = append(m.effects, effect)
}

func (m *EffectManager) Update(parentPosition *frontend.Vector) {
	for _, effect := range m.effects {
		effect.Update(parentPosition)
	}
}

func (m *EffectManager) Draw(drawFunc frontend.DrawFunc) {
	for _, effect := range m.effects {
		effect.Draw(drawFunc)
	}
}
