package sequence

import (
	"github.com/asragi/yasoba-prototype/core"
)

type partnerDialogueModel struct {
	eventId eventId
	textId  core.TextId
}

type partnerDialogueDataPort func(eventId) *partnerDialogueModel

type createPartnerDialogueEvent func(eventId) *eventUnit

type SetPartnerDialogueResponse struct {
	CheckIsEnd checkIsEnd
}

func produceCreatePartnerDialogueEventToUnit(
	serveTextData core.ServeTextDataFunc,
	partnerDialogueDataPort partnerDialogueDataPort,
	setPartnerDialogue SetPartnerDialogue,
) createPartnerDialogueEvent {
	return func(id eventId) *eventUnit {
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
