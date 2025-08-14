package sequence

import "fmt"

type sequenceId string

type sequence struct {
	id      sequenceId
	index   int
	isStart bool
	events  []*eventUnit
}

func (s *sequence) render() {
	for _, event := range s.events {
		event.render()
	}
}

func (s *sequence) update() isEnd {
	if s.isStart {
		s.isStart = false
		s.events[s.index].start()
	}

	for i, event := range s.events {
		if i != s.index {
			event.update()
			continue
		}

		isEnd := event.update()
		if isEnd {
			s.index++
			s.isStart = true
		}
	}

	if s.index >= len(s.events) {
		return true
	}

	return false
}

type isEnd bool
type render func()
type start func()
type update func() isEnd

type createSequence func(id sequenceId) *sequence

func provideCreateSequence(
	sequencesDataPort sequenceDataPort,
	createPartnerDialogueEvent createPartnerDialogueEvent,
) createSequence {
	sequences := sequencesDataPort()
	sequencesMap := make(map[sequenceId]*sequence)
	for _, seq := range sequences {
		events := []*eventUnit{}
		for _, event := range seq.events {
			if event.eventType == partnerDialogueEvent {
				events = append(events, createPartnerDialogueEvent(event.id))
			}
		}
		sequencesMap[seq.id] = &sequence{
			id:     seq.id,
			events: events,
		}
	}

	return func(id sequenceId) *sequence {
		if seq, ok := sequencesMap[id]; ok {
			return seq
		}

		fmt.Println("sequence not found")
		return nil
	}
}

// sequenceはn個のeventUnitを持つ
type eventUnit struct {
	update update
	render render
	start  start
}

type provideCreatePartnerDialogueEvent func() createPartnerDialogueEvent
