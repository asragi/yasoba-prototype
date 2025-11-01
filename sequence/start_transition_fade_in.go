package sequence

import seqtransition "github.com/asragi/yasoba-prototype/sequence/transition"

type createStartTransitionFadeInEvent func(EventID) *eventUnit

func produceCreateStartTransitionFadeInEventToUnit(
	optionPort seqtransition.OptionPort,
	controller seqtransition.Controller,
) createStartTransitionFadeInEvent {
	return func(id EventID) *eventUnit {
		option := optionPort(string(id))
		waitForComplete := option != nil && option.WaitForComplete()

		return &eventUnit{
			start: func() {
				controller.Start()
			},
			checkIsEnd: func() IsEnd {
				if !waitForComplete {
					return true
				}
				return IsEnd(!controller.IsRunning())
			},
			reset: func() {},
		}
	}
}
