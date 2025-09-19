package sequence

import "fmt"

type sequenceId string

type sequence struct {
	id      sequenceId
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

type CreateSequence func(id sequenceId) *sequence

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
		sequences := make(map[sequenceId]*sequence)
		for _, seq := range sequenceDataArray {
			events := []*eventUnit{}
			for _, event := range seq.events {
				if event.eventType == "partner_dialogue" {
					events = append(events, createPartnerDialogueEvent(event.id))
					continue
				}
				if event.eventType == "change_emotion" {
					events = append(events, createChangeEmotionEvent(event.id))
					continue
				}
				if event.eventType == "open_partner_message_window" {
					events = append(events, createOpenPartnerMessageWindowEvent(event.id))
					continue
				}
				if event.eventType == "close_partner_message_window" {
					events = append(events, createClosePartnerMessageWindowEvent(event.id))
					continue
				}
			}
			sequences[seq.id] = &sequence{
				id:      seq.id,
				events:  events,
				isStart: true,
			}
		}

		return func(id sequenceId) *sequence {
			if seq, ok := sequences[id]; ok {
				seq.Reset()
				return seq
			}

			fmt.Println("sequence not found")
			return nil
		}
	}
}

// sequenceはn個のeventUnitを持つ
type eventUnit struct {
	checkIsEnd checkIsEnd
	start      start
	reset      func()
}
