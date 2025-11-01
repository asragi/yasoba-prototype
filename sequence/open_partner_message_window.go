package sequence

type createOpenPartnerMessageWindowEvent func(EventID) *eventUnit

type OpenPartnerMessageWindow func()

func produceCreateOpenPartnerMessageWindowEventToUnit(
	setOpenPartnerMessageWindow OpenPartnerMessageWindow,
) createOpenPartnerMessageWindowEvent {
	return func(id EventID) *eventUnit {
		start := func() {
			setOpenPartnerMessageWindow()
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
