package emotion

import "github.com/asragi/yasoba-prototype/view/common/animation"

type BattleEmotionType int

const (
	BattleEmotionNormal BattleEmotionType = iota
	BattleEmotionDamage
	BattleEmotionSmile
	BattleEmotionAngry
	BattleEmotionAnnoyed
)

type Queued struct {
	currentType BattleEmotionType
	pending     BattleEmotionType
	hasPending  bool
}

func NewQueued(initial BattleEmotionType) Queued {
	return Queued{
		currentType: initial,
		pending:     initial,
	}
}

func (q *Queued) Enqueue(emotion BattleEmotionType) {
	q.pending = emotion
	q.hasPending = true
}

func (q *Queued) Current() BattleEmotionType {
	return q.currentType
}

func (q *Queued) Apply(fetch func(BattleEmotionType) *animation.Animation) *animation.Animation {
	changed, emotion := q.consume()
	animation := fetch(emotion)
	if changed && animation != nil {
		animation.Reset()
	}
	return animation
}

func (q *Queued) consume() (bool, BattleEmotionType) {
	if !q.hasPending {
		return false, q.currentType
	}
	q.hasPending = false
	if q.pending == q.currentType {
		return false, q.currentType
	}
	q.currentType = q.pending
	return true, q.currentType
}
