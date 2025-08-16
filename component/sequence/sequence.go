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
	if s.index >= len(s.events) {
		return true
	}

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

	return s.index >= len(s.events)
}

type isEnd bool
type render func()
type start func()
type update func() isEnd

type createSequence func(id sequenceId) *sequence

func initializeProduceCreateSequence(
	sequencesDataPort sequenceDataPort,
) func(createPartnerDialogueEvent) createSequence {
	sequenceDataArray := sequencesDataPort()
	return func(createPartnerDialogueEvent createPartnerDialogueEvent) createSequence {
		sequences := make(map[sequenceId]*sequence)
		for _, seq := range sequenceDataArray {
			events := []*eventUnit{}
			for _, event := range seq.events {
				events = append(events, createPartnerDialogueEvent(event.id))
			}
			sequences[seq.id] = &sequence{
				id:     seq.id,
				events: events,
			}
		}

		return func(id sequenceId) *sequence {
			if seq, ok := sequences[id]; ok {
				return seq
			}

			fmt.Println("sequence not found")
			return nil
		}
	}
}

// sequenceはn個のeventUnitを持つ
type eventUnit struct {
	update update
	render render
	start  start
}
