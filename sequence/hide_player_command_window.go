package sequence

type createHidePlayerCommandWindowEvent func(EventID) *eventUnit

type HidePlayerCommandWindow func()

func produceCreateHidePlayerCommandWindowEventToUnit(
	hidePlayerCommandWindow HidePlayerCommandWindow,
) createHidePlayerCommandWindowEvent {
	return func(id EventID) *eventUnit {
		start := func() {
			hidePlayerCommandWindow()
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
