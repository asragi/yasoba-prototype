package sequence

import "github.com/asragi/yasoba-prototype/text"

type PartnerDialogueModel struct {
	eventId EventID
	textId  text.TextId
}

type PartnerDialogueDataPort func(EventID) *PartnerDialogueModel

type createPartnerDialogueEvent func(EventID) *eventUnit

type SetPartnerDialogue func(text.String) *SetPartnerDialogueResponse

type SetPartnerDialogueResponse struct {
	CheckIsEnd checkIsEnd
}

func produceCreatePartnerDialogueEventToUnit(
	serveTextData text.ServeTextDataFunc,
	partnerDialogueDataPort PartnerDialogueDataPort,
	setPartnerDialogue SetPartnerDialogue,
) createPartnerDialogueEvent {
	return func(id EventID) *eventUnit {
		model := partnerDialogueDataPort(id)
		textData := serveTextData(model.textId)
		var textUpdate checkIsEnd
		start := func() {
			res := setPartnerDialogue(textData.Text)
			textUpdate = res.CheckIsEnd
		}
		checkIsEnd := func() IsEnd {
			if textUpdate == nil {
				return false
			}
			return textUpdate()
		}
		return &eventUnit{
			start:      start,
			checkIsEnd: checkIsEnd,
			reset:      func() {},
		}
	}
}

func NewPartnerDialogueModel(eventID EventID, textID text.TextId) *PartnerDialogueModel {
	return &PartnerDialogueModel{
		eventId: eventID,
		textId:  textID,
	}
}
