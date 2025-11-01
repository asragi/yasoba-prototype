package sequence

import "github.com/asragi/yasoba-prototype/text"

type ProduceCreateSequence func(
	SetPartnerDialogue,
	SetEmotion,
	OpenPartnerMessageWindow,
	ClosePartnerMessageWindow,
) CreateSequence

type PrepareProduceCreateSequence func(text.ServeTextDataFunc) ProduceCreateSequence

func InitializeProduceCreateSequence(
	sequenceModelPort SequenceModelPort,
	eventDataModelPort EventDataModelPort,
	partnerDialogueDataPort PartnerDialogueDataPort,
	changeEmotionDataPort ChangeEmotionDataPort,
) PrepareProduceCreateSequence {
	sequencesDataAdapter := initializeSequenceDataAdapter(
		sequenceModelPort,
		eventDataModelPort,
	)
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
