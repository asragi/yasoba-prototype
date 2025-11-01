package sequence

type createClosePartnerMessageWindowEvent func(EventID) *eventUnit

type ClosePartnerMessageWindow func()

func produceCreateClosePartnerMessageWindowEventToUnit(
	setClosePartnerMessageWindow ClosePartnerMessageWindow,
) createClosePartnerMessageWindowEvent {
	return func(id EventID) *eventUnit {
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
