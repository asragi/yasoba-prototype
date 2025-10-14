package sequence

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/component"
)

type ChangeEmotion struct {
	actorId actor.ActorId
	emotion component.BattleEmotionType
}

type changeEmotionDataPort func(eventId) *ChangeEmotion

type createChangeEmotionEvent func(eventId) *eventUnit

type SetEmotion func(actor.ActorId, component.BattleEmotionType)

func produceCreateChangeEmotionEventToUnit(
	changeEmotionDataPort changeEmotionDataPort,
	setEmotion SetEmotion,
) createChangeEmotionEvent {
	return func(eventId eventId) *eventUnit {
		model := changeEmotionDataPort(eventId)
		return &eventUnit{
			start: func() {
				fmt.Println("change emotion", model.actorId, model.emotion)
				fmt.Println("setEmotion", setEmotion)
				setEmotion(model.actorId, model.emotion)
			},
			checkIsEnd: func() IsEnd { return true },
			reset:      func() {},
		}
	}
}
