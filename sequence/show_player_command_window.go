package sequence

type createShowPlayerCommandWindowEvent func(EventID) *eventUnit

type ShowPlayerCommandWindow func()

func produceCreateShowPlayerCommandWindowEventToUnit(
	showPlayerCommandWindow ShowPlayerCommandWindow,
) createShowPlayerCommandWindowEvent {
	return func(id EventID) *eventUnit {
		start := func() {
			showPlayerCommandWindow()
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
