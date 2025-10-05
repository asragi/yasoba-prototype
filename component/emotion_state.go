package component

import "github.com/asragi/yasoba-prototype/widget"

type queuedEmotion struct {
	currentType BattleEmotionType
	pending     BattleEmotionType
	hasPending  bool
}

func newQueuedEmotion(initial BattleEmotionType) queuedEmotion {
	return queuedEmotion{currentType: initial, pending: initial}
}

func (q *queuedEmotion) Enqueue(emotion BattleEmotionType) {
	q.pending = emotion
	q.hasPending = true
}

func (q *queuedEmotion) consume() (bool, BattleEmotionType) {
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

func (q *queuedEmotion) current() BattleEmotionType {
	return q.currentType
}

func (q *queuedEmotion) apply(fetch func(BattleEmotionType) *widget.Animation) *widget.Animation {
	changed, emotion := q.consume()
	animation := fetch(emotion)
	if changed {
		animation.Reset()
	}
	return animation
}
