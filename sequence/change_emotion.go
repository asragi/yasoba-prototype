package sequence

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/common/emotion"
)

type ChangeEmotion struct {
	actorId actor.ActorId
	emotion emotion.EmotionType
}

type ChangeEmotionDataPort func(EventID) *ChangeEmotion

type createChangeEmotionEvent func(EventID) *eventUnit

type SetEmotion func(actor.ActorId, emotion.EmotionType)

func produceCreateChangeEmotionEventToUnit(
	changeEmotionDataPort ChangeEmotionDataPort,
	setEmotion SetEmotion,
) createChangeEmotionEvent {
	return func(eventId EventID) *eventUnit {
		model := changeEmotionDataPort(eventId)
		return &eventUnit{
			start: func() {
				setEmotion(model.actorId, model.emotion)
			},
			checkIsEnd: func() IsEnd { return true },
			reset:      func() {},
		}
	}
}

func NewChangeEmotion(actorId actor.ActorId, emotion emotion.EmotionType) *ChangeEmotion {
	return &ChangeEmotion{
		actorId: actorId,
		emotion: emotion,
	}
}
