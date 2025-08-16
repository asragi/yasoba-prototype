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

type setPartnerDialogueResponse struct {
	update update
	render render
}

type setPartnerDialogue func(core.TextString) *setPartnerDialogueResponse

func produceCreatePartnerDialogueEventToUnit(
	serveTextData core.ServeTextDataFunc,
	partnerDialogueDataPort partnerDialogueDataPort,
	setPartnerDialogue setPartnerDialogue,
) createPartnerDialogueEvent {
	return func(id eventId) *eventUnit {
		model := partnerDialogueDataPort(id)
		textData := serveTextData(model.textId)
		var textUpdate update
		var textRender render
		start := func() {
			res := setPartnerDialogue(textData.Text)
			textUpdate = res.update
			textRender = res.render
		}
		update := func() isEnd {
			if textUpdate == nil {
				return false
			}
			return textUpdate()
		}
		render := func() {
			if textRender == nil {
				return
			}
			textRender()
		}
		return &eventUnit{
			start:  start,
			render: render,
			update: update,
		}
	}
}
