package adapter

import (
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/common/emotion"
	"github.com/asragi/yasoba-prototype/sequence"
	seqtransition "github.com/asragi/yasoba-prototype/sequence/transition"
	"github.com/asragi/yasoba-prototype/text"
)

type ProduceCreateSequence = sequence.ProduceCreateSequence
type PrepareProduceCreateSequence = sequence.PrepareProduceCreateSequence

// DIにあたる処理を行いpackage外にexportする
func InitializeProduceCreateSequence() PrepareProduceCreateSequence {
	sequenceModelRecordsPort := SequenceModelPortFromYAML("assets/data/sequence.yaml")
	eventDataRecordsPort := EventDataModelPortFromYAML("assets/data/event_data.yaml")
	sequenceModelPort := adaptSequenceModelPort(sequenceModelRecordsPort)
	eventDataModelPort := adaptEventDataModelPort(eventDataRecordsPort)
	partnerDialogueRecords := PartnerDialogueRecordsFromYAML("assets/data/partner_dialogue.yaml")
	partnerDialogueDataPort := adaptPartnerDialogueDataPort(partnerDialogueRecords)
	changeEmotionRecords := ChangeEmotionRecordsFromYAML("assets/data/change_emotion.yaml")
	changeEmotionDataPort := adaptChangeEmotionDataPort(changeEmotionRecords)
	transitionFadeRecords := TransitionFadeRecordsFromYAML("assets/data/transition_fade.yaml")
	transitionFadeOptionPort := adaptTransitionFadeOptionPort(transitionFadeRecords)
	switchToBattleSceneRecords := SwitchToBattleSceneRecordsFromYAML("assets/data/switch_to_battle_scene.yaml")
	switchToBattleSceneDataPort := adaptSwitchToBattleSceneDataPort(switchToBattleSceneRecords)

	return sequence.InitializeProduceCreateSequence(
		sequenceModelPort,
		eventDataModelPort,
		partnerDialogueDataPort,
		changeEmotionDataPort,
		switchToBattleSceneDataPort,
		transitionFadeOptionPort,
	)
}

func adaptSequenceModelPort(port GetSequenceModelFunc) sequence.SequenceModelPort {
	return func() []*sequence.SequenceModel {
		records := port()
		result := make([]*sequence.SequenceModel, 0, len(records))
		for _, record := range records {
			result = append(result, sequence.NewSequenceModel(
				sequence.SequenceId(record.ID),
				sequence.EventID(record.HeadEventID),
			))
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
				sequence.EventID(record.NextEventID),
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
			toEmotionType(record.Emotion),
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

func adaptTransitionFadeOptionPort(records map[string]*TransitionFadeRecord) seqtransition.OptionPort {
	options := make(map[string]*seqtransition.Option, len(records))
	for _, record := range records {
		options[record.EventID] = seqtransition.NewOption(record.WaitForComplete)
	}
	return func(id string) *seqtransition.Option {
		return options[id]
	}
}

func adaptSwitchToBattleSceneDataPort(records map[string]*SwitchToBattleSceneRecord) sequence.SwitchToBattleSceneDataPort {
	models := make(map[sequence.EventID]*sequence.SwitchToBattleSceneModel, len(records))
	for _, record := range records {
		eventID := sequence.EventID(record.EventID)
		models[eventID] = sequence.NewSwitchToBattleSceneModel(battle.BattleId(record.BattleID))
	}
	return func(id sequence.EventID) *sequence.SwitchToBattleSceneModel {
		return models[id]
	}
}

func toEmotionType(value string) emotion.EmotionType {
	switch value {
	case "normal":
		return emotion.EmotionNormal
	case "damage":
		return emotion.EmotionDamage
	case "smile":
		return emotion.EmotionSmile
	case "angry":
		return emotion.EmotionAngry
	case "annoyed":
		return emotion.EmotionAnnoyed
	default:
		return emotion.EmotionNormal
	}
}
