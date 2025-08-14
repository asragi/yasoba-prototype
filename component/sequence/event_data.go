package sequence

type eventId string
type eventType string

const (
	partnerDialogueEvent eventType = "partnerDialogue"
)

type eventDataModel struct {
	id        eventId
	eventType eventType
	ownerId   sequenceId
	order     int
}

type eventDataPort func() []*eventDataModel

type sequenceModel struct {
	id sequenceId
}

type sequenceModelPort func() []*sequenceModel

type sequenceData struct {
	id     sequenceId
	events []*eventDataModel
}

type sequenceDataPort func() []*sequenceData

func createSequenceDataAdapter(
	sequenceModelPort sequenceModelPort,
	eventDataPort eventDataPort,
) sequenceDataPort {
	groupEventsBySequenceId := func(events []*eventDataModel) map[sequenceId][]*eventDataModel {
		eventMap := make(map[sequenceId][]*eventDataModel)
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
