package adapter

import (
	"github.com/asragi/yasoba-prototype/actor"
	battleemotion "github.com/asragi/yasoba-prototype/component/battle/emotion"
	"github.com/asragi/yasoba-prototype/sequence"
	"github.com/asragi/yasoba-prototype/text"
)

type ProduceCreateSequence = sequence.ProduceCreateSequence
type PrepareProduceCreateSequence = sequence.PrepareProduceCreateSequence

// DIにあたる処理を行いpackage外にexportする
func InitializeProduceCreateSequence() PrepareProduceCreateSequence {
	sequenceModelRecordsPort := SequenceModelPortFromYAML("data/sequence.yaml")
	eventDataRecordsPort := EventDataModelPortFromYAML("data/event_data.yaml")
	sequenceModelPort := adaptSequenceModelPort(sequenceModelRecordsPort)
	eventDataModelPort := adaptEventDataModelPort(eventDataRecordsPort)
	partnerDialogueRecords := PartnerDialogueRecordsFromYAML("data/partner_dialogue.yaml")
	partnerDialogueDataPort := adaptPartnerDialogueDataPort(partnerDialogueRecords)
	changeEmotionRecords := ChangeEmotionRecordsFromYAML("data/change_emotion.yaml")
	changeEmotionDataPort := adaptChangeEmotionDataPort(changeEmotionRecords)

	return sequence.InitializeProduceCreateSequence(
		sequenceModelPort,
		eventDataModelPort,
		partnerDialogueDataPort,
		changeEmotionDataPort,
	)
}

func adaptSequenceModelPort(port GetSequenceModelFunc) sequence.SequenceModelPort {
	return func() []*sequence.SequenceModel {
		records := port()
		result := make([]*sequence.SequenceModel, 0, len(records))
		for _, record := range records {
			result = append(result, sequence.NewSequenceModel(sequence.SequenceId(record.ID)))
		}
		return result
	}
}

func adaptEventDataModelPort(port GetEventDataModelFunc) sequence.EventDataModelPort {
	return func() []*sequence.EventDataModel {
		records := port()
		result := make([]*sequence.EventDataModel, 0, len(records))
		for _, record := range records {
			result = append(result, sequence.NewEventDataModel(
				sequence.EventID(record.ID),
				sequence.EventType(record.EventType),
				sequence.SequenceId(record.OwnerID),
				record.Order,
			))
		}
		return result
	}
}

func adaptPartnerDialogueDataPort(records map[string]*PartnerDialogueRecord) sequence.PartnerDialogueDataPort {
	dialogueMap := make(map[sequence.EventID]*sequence.PartnerDialogueModel, len(records))
	for _, record := range records {
		eventID := sequence.EventID(record.EventID)
		dialogueMap[eventID] = sequence.NewPartnerDialogueModel(eventID, text.TextId(record.TextID))
	}
	return func(id sequence.EventID) *sequence.PartnerDialogueModel {
		return dialogueMap[id]
	}
}

func adaptChangeEmotionDataPort(records map[string]*ChangeEmotionRecord) sequence.ChangeEmotionDataPort {
	emotionMap := make(map[sequence.EventID]*sequence.ChangeEmotion, len(records))
	for _, record := range records {
		eventID := sequence.EventID(record.EventID)
		emotionMap[eventID] = sequence.NewChangeEmotion(
			actor.ActorId(record.ActorID),
			toBattleEmotionType(record.Emotion),
		)
	}
	return func(id sequence.EventID) *sequence.ChangeEmotion {
		emotion := emotionMap[id]
		if emotion == nil {
			panic("event not found: " + string(id))
		}
		return emotion
	}
}

func toBattleEmotionType(value string) battleemotion.BattleEmotionType {
	switch value {
	case "normal":
		return battleemotion.BattleEmotionNormal
	case "damage":
		return battleemotion.BattleEmotionDamage
	case "smile":
		return battleemotion.BattleEmotionSmile
	case "angry":
		return battleemotion.BattleEmotionAngry
	case "annoyed":
		return battleemotion.BattleEmotionAnnoyed
	default:
		return battleemotion.BattleEmotionNormal
	}
}
