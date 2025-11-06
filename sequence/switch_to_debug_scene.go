package sequence

type SwitchToDebugScene func()

type createSwitchToDebugSceneEvent func(EventID) *eventUnit

func produceCreateSwitchToDebugSceneEventToUnit(
	switchToDebug SwitchToDebugScene,
) createSwitchToDebugSceneEvent {
	if switchToDebug == nil {
		panic("sequence: switchToDebug is nil")
	}
	return func(EventID) *eventUnit {
		return &eventUnit{
			start: func() {
				switchToDebug()
			},
			checkIsEnd: func() IsEnd { return true },
			reset:      func() {},
		}
	}
}
