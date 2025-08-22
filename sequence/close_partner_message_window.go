package sequence

type createClosePartnerMessageWindowEvent func(eventId) *eventUnit

type ClosePartnerMessageWindow func()

func produceCreateClosePartnerMessageWindowEventToUnit(
	setClosePartnerMessageWindow ClosePartnerMessageWindow,
) createClosePartnerMessageWindowEvent {
	return func(id eventId) *eventUnit {
		start := func() {
			setClosePartnerMessageWindow()
		}

		checkIsEnd := func() IsEnd {
			return true
		}

		return &eventUnit{
			start:      start,
			checkIsEnd: checkIsEnd,
			reset:      func() {},
		}
	}
}
