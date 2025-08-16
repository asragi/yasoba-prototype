package sequence

import "github.com/asragi/yasoba-prototype/core"

// DIにあたる処理を行いpackage外にexportする
func InitializeProduceCreateSequence() func(core.ServeTextDataFunc) func(setPartnerDialogue) createSequence {
	sequenceModelPort := createSequenceModelPortFromYAML("data/sequence_model.yaml")
	eventDataModelPort := createEventDataModelPortFromYAML("data/event_data.yaml")
	// Modelから1対多の構造のデータに変換する
	sequencesDataAdapter := initializeSequenceDataAdapter(
		sequenceModelPort,
		eventDataModelPort,
	)
	partnerDialogueDataPort := createPartnerDialogueDataPortFromYAML("data/partner_dialogue_data.yaml")

	produceCreateSequence := initializeProduceCreateSequence(
		sequencesDataAdapter,
	)
	return func(
		serveTextData core.ServeTextDataFunc,
	) func(setPartnerDialogue) createSequence {
		return func(setPartnerDialogue setPartnerDialogue) createSequence {
			return produceCreateSequence(
				produceCreatePartnerDialogueEventToUnit(
					serveTextData,
					partnerDialogueDataPort,
					setPartnerDialogue,
				),
			)
		}
	}
}
