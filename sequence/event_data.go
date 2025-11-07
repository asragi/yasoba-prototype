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
	setMessageWindowTextEvent      EventType = "set_message_window_text"
	waitFrameEvent                 EventType = "wait_frame"
)

type EventDataModel struct {
	id          EventID
	eventType   EventType
	nextEventID EventID
}

type EventDataModelPort func() []*EventDataModel

type SequenceModel struct {
	id          SequenceId
	headEventID EventID
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
	buildEventMap := func(events []*EventDataModel) map[EventID]*EventDataModel {
		eventMap := make(map[EventID]*EventDataModel, len(events))
		for _, event := range events {
			if _, exists := eventMap[event.id]; exists {
				panic("duplicate event id: " + string(event.id))
			}
			eventMap[event.id] = event
		}
		return eventMap
	}

	orderEventsByNextLink := func(
		sequenceID SequenceId,
		head EventID,
		eventMap map[EventID]*EventDataModel,
		globalAssignment map[EventID]SequenceId,
	) []*EventDataModel {
		if head == "" {
			return nil
		}

		ordered := []*EventDataModel{}
		visited := make(map[EventID]struct{})
		currentID := head
		for currentID != "" {
			event, ok := eventMap[currentID]
			if !ok {
				panic("sequence " + string(sequenceID) + " references unknown event: " + string(currentID))
			}

			if _, seen := visited[event.id]; seen {
				panic("sequence " + string(sequenceID) + " contains a cycle")
			}
			visited[event.id] = struct{}{}

			if owner, assigned := globalAssignment[event.id]; assigned && owner != sequenceID {
				panic("event " + string(event.id) + " already assigned to sequence " + string(owner))
			}
			globalAssignment[event.id] = sequenceID

			ordered = append(ordered, event)
			currentID = event.nextEventID
		}

		return ordered
	}

	return func() []*sequenceData {
		sequences := sequenceModelPort()
		events := eventDataPort()

		eventMap := buildEventMap(events)
		globalAssignment := make(map[EventID]SequenceId, len(events))

		var result []*sequenceData
		for _, sequence := range sequences {
			orderedEvents := orderEventsByNextLink(sequence.id, sequence.headEventID, eventMap, globalAssignment)

			result = append(result, &sequenceData{
				id:     sequence.id,
				events: orderedEvents,
			})
		}

		return result
	}
}

func NewSequenceModel(id SequenceId, headEventID EventID) *SequenceModel {
	return &SequenceModel{id: id, headEventID: headEventID}
}

func NewEventDataModel(id EventID, eventType EventType, nextEventID EventID) *EventDataModel {
	return &EventDataModel{
		id:          id,
		eventType:   eventType,
		nextEventID: nextEventID,
	}
}
