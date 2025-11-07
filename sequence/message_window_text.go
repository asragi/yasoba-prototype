package sequence

import "github.com/asragi/yasoba-prototype/text"

type MessageWindowTextModel struct {
	eventId         EventID
	textId          text.TextId
	waitForComplete bool
}

type MessageWindowTextDataPort func(EventID) *MessageWindowTextModel

type createSetMessageWindowTextEvent func(EventID) *eventUnit

type SetMessageWindowText func(text.String) *SetMessageWindowTextResponse

type SetMessageWindowTextResponse struct {
	CheckIsEnd checkIsEnd
}

func NewMessageWindowTextModel(eventID EventID, textID text.TextId, waitForComplete bool) *MessageWindowTextModel {
	return &MessageWindowTextModel{
		eventId:         eventID,
		textId:          textID,
		waitForComplete: waitForComplete,
	}
}

func produceCreateSetMessageWindowTextEventToUnit(
	serveTextData text.ServeTextDataFunc,
	dataPort MessageWindowTextDataPort,
	setMessageWindowText SetMessageWindowText,
) createSetMessageWindowTextEvent {
	return func(id EventID) *eventUnit {
		model := dataPort(id)
		textData := serveTextData(model.textId)
		var isComplete checkIsEnd

		start := func() {
			response := setMessageWindowText(textData.Text)
			if response == nil {
				isComplete = nil
				return
			}
			isComplete = response.CheckIsEnd
		}

		checkIsEnd := func() IsEnd {
			if !model.waitForComplete {
				return true
			}
			if isComplete == nil {
				return false
			}
			return isComplete()
		}

		return &eventUnit{
			start:      start,
			checkIsEnd: checkIsEnd,
			reset: func() {
				isComplete = nil
			},
		}
	}
}
