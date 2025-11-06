package sequence

type EventID string
type EventType string

const (
	partnerDialogueEvent           EventType = "partner_dialogue"
	changeEmotionEvent             EventType = "change_emotion"
	openPartnerMessageWindowEvent  EventType = "open_partner_message_window"
	closePartnerMessageWindowEvent EventType = "close_partner_message_window"
	showPlayerCommandWindowEvent   EventType = "show_player_command_window"
	hidePlayerCommandWindowEvent   EventType = "hide_player_command_window"
	startTransitionFadeOutEvent    EventType = "start_transition_fade_out"
	startTransitionFadeInEvent     EventType = "start_transition_fade_in"
	switchToBattleSceneEvent       EventType = "switch_to_battle_scene"
	switchToDebugSceneEvent        EventType = "switch_to_debug_scene"
)

type EventDataModel struct {
	id        EventID
	eventType EventType
	ownerId   SequenceId
	order     int
}

type EventDataModelPort func() []*EventDataModel

type SequenceModel struct {
	id SequenceId
}

type SequenceModelPort func() []*SequenceModel

type sequenceData struct {
	id     SequenceId
	events []*EventDataModel
}

// sequenceModelとeventDataModelを組み合わせてsequenceDataを作成します
type sequenceDataPort func() []*sequenceData

func initializeSequenceDataAdapter(
	sequenceModelPort SequenceModelPort,
	eventDataPort EventDataModelPort,
) sequenceDataPort {
	groupEventsBySequenceId := func(events []*EventDataModel) map[SequenceId][]*EventDataModel {
		eventMap := make(map[SequenceId][]*EventDataModel)
		for _, event := range events {
			eventMap[event.ownerId] = append(eventMap[event.ownerId], event)
		}
		return eventMap
	}

	sortEventsByOrder := func(events []*EventDataModel) []*EventDataModel {
		for i := 0; i < len(events)-1; i++ {
			for j := i + 1; j < len(events); j++ {
				if events[i].order > events[j].order {
					events[i], events[j] = events[j], events[i]
				}
			}
		}
		return events
	}

	return func() []*sequenceData {
		sequences := sequenceModelPort()
		events := eventDataPort()

		// イベントをsequenceIdでグループ化
		eventMap := groupEventsBySequenceId(events)

		// sequenceDataを作成
		var result []*sequenceData
		for _, sequence := range sequences {
			sequenceEvents := eventMap[sequence.id]

			// orderに基づいてソート
			sortedEvents := sortEventsByOrder(sequenceEvents)

			result = append(result, &sequenceData{
				id:     sequence.id,
				events: sortedEvents,
			})
		}

		return result
	}
}

func NewSequenceModel(id SequenceId) *SequenceModel {
	return &SequenceModel{id: id}
}

func NewEventDataModel(id EventID, eventType EventType, ownerId SequenceId, order int) *EventDataModel {
	return &EventDataModel{
		id:        id,
		eventType: eventType,
		ownerId:   ownerId,
		order:     order,
	}
}
