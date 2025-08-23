package sequence

type SequenceManager struct {
	sequences []*sequence
}

func (m *SequenceManager) Update() IsEnd {
	for _, sequence := range m.sequences {
		sequence.Update()
	}
	m.sequences = func() []*sequence {
		sequences := []*sequence{}
		for _, sequence := range m.sequences {
			if sequence.Update() {
				continue
			}
			sequences = append(sequences, sequence)
		}
		return sequences
	}()

	return len(m.sequences) == 0
}

func (m *SequenceManager) AddSequence(sequence *sequence) {
	m.sequences = append(m.sequences, sequence)
}

func CreateSequenceManager() *SequenceManager {
	return &SequenceManager{
		sequences: []*sequence{},
	}
}
