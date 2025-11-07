package emotion

import (
	commonemotion "github.com/asragi/yasoba-prototype/common/emotion"
	"github.com/asragi/yasoba-prototype/view/common/animation"
)

type Queued struct {
	currentType commonemotion.EmotionType
	pending     commonemotion.EmotionType
	hasPending  bool
}

func NewQueued(initial commonemotion.EmotionType) Queued {
	return Queued{
		currentType: initial,
		pending:     initial,
	}
}

func (q *Queued) Enqueue(value commonemotion.EmotionType) {
	q.pending = value
	q.hasPending = true
}

func (q *Queued) Current() commonemotion.EmotionType {
	return q.currentType
}

func (q *Queued) Apply(fetch func(commonemotion.EmotionType) *animation.Animation) *animation.Animation {
	changed, current := q.consume()
	anim := fetch(current)
	if changed && anim != nil {
		anim.Reset()
	}
	return anim
}

func (q *Queued) consume() (bool, commonemotion.EmotionType) {
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
