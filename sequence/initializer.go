package sequence

import (
	seqtransition "github.com/asragi/yasoba-prototype/sequence/transition"
	"github.com/asragi/yasoba-prototype/text"
)

type ProduceCreateSequence func(
	SetPartnerDialogue,
	SetEmotion,
	OpenPartnerMessageWindow,
	ClosePartnerMessageWindow,
	ShowPlayerCommandWindow,
	HidePlayerCommandWindow,
	seqtransition.Controller,
	seqtransition.Controller,
) CreateSequence

type PrepareProduceCreateSequence func(text.ServeTextDataFunc) ProduceCreateSequence

func InitializeProduceCreateSequence(
	sequenceModelPort SequenceModelPort,
	eventDataModelPort EventDataModelPort,
	partnerDialogueDataPort PartnerDialogueDataPort,
	changeEmotionDataPort ChangeEmotionDataPort,
	transitionFadeOptionPort seqtransition.OptionPort,
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
			showPlayerCommandWindow ShowPlayerCommandWindow,
			hidePlayerCommandWindow HidePlayerCommandWindow,
			startTransitionFadeOut seqtransition.Controller,
			startTransitionFadeIn seqtransition.Controller,
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
				produceCreateShowPlayerCommandWindowEventToUnit(
					showPlayerCommandWindow,
				),
				produceCreateHidePlayerCommandWindowEventToUnit(
					hidePlayerCommandWindow,
				),
				produceCreateStartTransitionFadeOutEventToUnit(
					transitionFadeOptionPort,
					startTransitionFadeOut,
				),
				produceCreateStartTransitionFadeInEventToUnit(
					transitionFadeOptionPort,
					startTransitionFadeIn,
				),
			)
		}
	}
}
