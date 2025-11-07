package sequence

type WaitFrameModel struct {
	eventId    EventID
	frameCount int
}

type WaitFrameDataPort func(EventID) *WaitFrameModel

type createWaitFrameEvent func(EventID) *eventUnit

func NewWaitFrameModel(eventID EventID, frameCount int) *WaitFrameModel {
	return &WaitFrameModel{
		eventId:    eventID,
		frameCount: frameCount,
	}
}

func produceCreateWaitFrameEventToUnit(
	dataPort WaitFrameDataPort,
) createWaitFrameEvent {
	return func(id EventID) *eventUnit {
		model := dataPort(id)
		remaining := 0

		start := func() {
			remaining = model.frameCount
		}

		checkIsEnd := func() IsEnd {
			if remaining <= 0 {
				return true
			}
			remaining--
			return remaining <= 0
		}

		reset := func() {
			remaining = 0
		}

		return &eventUnit{
			start:      start,
			checkIsEnd: checkIsEnd,
			reset:      reset,
		}
	}
}
