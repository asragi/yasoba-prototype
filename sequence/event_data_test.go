package sequence

import (
	"reflect"
	"testing"
)

func Test_createSequenceDataAdapter(t *testing.T) {
	sequenceModelPort := func() []*SequenceModel {
		return []*SequenceModel{
			{id: "sequence1"},
			{id: "sequence2"},
		}
	}

	eventDataPort := func() []*EventDataModel {
		return []*EventDataModel{
			{id: "event1", ownerId: "sequence1", order: 2},
			{id: "event2", ownerId: "sequence1", order: 1},
			{id: "event3", ownerId: "sequence2", order: 1},
		}
	}

	expected := []*sequenceData{
		{id: "sequence1", events: []*EventDataModel{
			{id: "event2", ownerId: "sequence1", order: 1},
			{id: "event1", ownerId: "sequence1", order: 2},
		}},
		{id: "sequence2", events: []*EventDataModel{
			{id: "event3", ownerId: "sequence2", order: 1},
		}},
	}

	sequenceDataPort := initializeSequenceDataAdapter(sequenceModelPort, eventDataPort)
	sequenceData := sequenceDataPort()

	if !reflect.DeepEqual(sequenceData, expected) {
		t.Errorf("expected %+v, got %+v", expected, sequenceData)
	}
}
