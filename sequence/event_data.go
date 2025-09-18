package sequence

type eventId string
type eventType string

const (
	partnerDialogueEvent           eventType = "partner_dialogue"
	openPartnerMessageWindowEvent  eventType = "open_partner_message_window"
	closePartnerMessageWindowEvent eventType = "close_partner_message_window"
)

type eventDataModel struct {
	id        eventId
	eventType eventType
	ownerId   SequenceId
	order     int
}

type eventDataModelPort func() []*eventDataModel

type sequenceModel struct {
	id SequenceId
}

type sequenceModelPort func() []*sequenceModel

type sequenceData struct {
	id     SequenceId
	events []*eventDataModel
}

// sequenceModelとeventDataModelを組み合わせてsequenceDataを作成します
type sequenceDataPort func() []*sequenceData

func initializeSequenceDataAdapter(
	sequenceModelPort sequenceModelPort,
	eventDataPort eventDataModelPort,
) sequenceDataPort {
	groupEventsBySequenceId := func(events []*eventDataModel) map[SequenceId][]*eventDataModel {
		eventMap := make(map[SequenceId][]*eventDataModel)
		for _, event := range events {
			eventMap[event.ownerId] = append(eventMap[event.ownerId], event)
		}
		return eventMap
	}

	sortEventsByOrder := func(events []*eventDataModel) []*eventDataModel {
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
