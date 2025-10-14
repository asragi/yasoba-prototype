package sequence

import "github.com/asragi/yasoba-prototype/text"

type SetPartnerDialogue func(text.String) *SetPartnerDialogueResponse

type ProduceCreateSequence func(SetPartnerDialogue, SetEmotion, OpenPartnerMessageWindow, ClosePartnerMessageWindow) CreateSequence

type PrepareProduceCreateSequence func(text.ServeTextDataFunc) ProduceCreateSequence

// DIにあたる処理を行いpackage外にexportする
func InitializeProduceCreateSequence() PrepareProduceCreateSequence {
	sequenceModelPort := createSequenceModelPortFromYAML("data/sequence.yaml")
	eventDataModelPort := createEventDataModelPortFromYAML("data/event_data.yaml")
	// Modelから1対多の構造のデータに変換する
	sequencesDataAdapter := initializeSequenceDataAdapter(
		sequenceModelPort,
		eventDataModelPort,
	)
	partnerDialogueDataPort := createPartnerDialogueDataPortFromYAML("data/partner_dialogue.yaml")
	changeEmotionDataPort := createChangeEmotionDataPortFromYAML("data/change_emotion.yaml")

	produceCreateSequence := initializeProduceCreateSequence(
		sequencesDataAdapter,
	)
	return func(
		serveTextData text.ServeTextDataFunc,
	) ProduceCreateSequence {
		return func(
			setPartnerDialogue SetPartnerDialogue,
			setEmotion SetEmotion,
			setOpenPartnerMessageWindow OpenPartnerMessageWindow,
			setClosePartnerMessageWindow ClosePartnerMessageWindow,
		) CreateSequence {
			return produceCreateSequence(
				produceCreatePartnerDialogueEventToUnit(
					serveTextData,
					partnerDialogueDataPort,
					setPartnerDialogue,
				),
				produceCreateChangeEmotionEventToUnit(
					changeEmotionDataPort,
					setEmotion,
				),
				produceCreateOpenPartnerMessageWindowEventToUnit(
					setOpenPartnerMessageWindow,
				),
				produceCreateClosePartnerMessageWindowEventToUnit(
					setClosePartnerMessageWindow,
				),
			)
		}
	}
}
