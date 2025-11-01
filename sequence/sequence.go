package sequence

type SequenceId string

type sequence struct {
	id      SequenceId
	index   int
	isStart bool
	events  []*eventUnit
}

func (s *sequence) Reset() {
	s.index = 0
	s.isStart = true
	for _, event := range s.events {
		event.reset()
	}
}

func (s *sequence) Update() IsEnd {
	if s.index >= len(s.events) {
		return true
	}

	for i, event := range s.events {
		if s.isStart {
			s.isStart = false
			event.start()
		}

		if i != s.index {
			continue
		}

		isEnd := event.checkIsEnd()
		if isEnd {
			s.index++
			s.isStart = true
		}
	}

	return s.index >= len(s.events)
}

type IsEnd bool
type start func()
type checkIsEnd func() IsEnd

type CreateSequence func(id SequenceId) *sequence

func initializeProduceCreateSequence(
	sequencesDataPort sequenceDataPort,
) func(createPartnerDialogueEvent, createChangeEmotionEvent, createOpenPartnerMessageWindowEvent, createClosePartnerMessageWindowEvent) CreateSequence {
	sequenceDataArray := sequencesDataPort()
	return func(
		createPartnerDialogueEvent createPartnerDialogueEvent,
		createChangeEmotionEvent createChangeEmotionEvent,
		createOpenPartnerMessageWindowEvent createOpenPartnerMessageWindowEvent,
		createClosePartnerMessageWindowEvent createClosePartnerMessageWindowEvent,
	) CreateSequence {
		sequences := make(map[SequenceId]*sequence)
		for _, seqData := range sequenceDataArray {
			events := []*eventUnit{}
			for _, event := range seqData.events {
				switch event.eventType {
				case partnerDialogueEvent:
					events = append(events, createPartnerDialogueEvent(event.id))
				case changeEmotionEvent:
					events = append(events, createChangeEmotionEvent(event.id))
				case openPartnerMessageWindowEvent:
					events = append(events, createOpenPartnerMessageWindowEvent(event.id))
				case closePartnerMessageWindowEvent:
					events = append(events, createClosePartnerMessageWindowEvent(event.id))
				}
			}
			sequences[seqData.id] = &sequence{
				id:      seqData.id,
				events:  events,
				isStart: true,
			}
		}

		return func(id SequenceId) *sequence {
			if seq, ok := sequences[id]; ok {
				seq.Reset()
				return seq
			}

			panic("sequence not found")
		}
	}
}

// sequenceはn個のeventUnitを持つ
type eventUnit struct {
	checkIsEnd checkIsEnd
	start      start
	reset      func()
}
