package sequence

import seqtransition "github.com/asragi/yasoba-prototype/sequence/transition"

type createStartTransitionFadeOutEvent func(EventID) *eventUnit

func produceCreateStartTransitionFadeOutEventToUnit(
	optionPort seqtransition.OptionPort,
	controller seqtransition.Controller,
) createStartTransitionFadeOutEvent {
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
