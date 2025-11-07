package sequence

import (
	"reflect"
	"testing"
)

func Test_createSequenceDataAdapter(t *testing.T) {
	sequenceModelPort := func() []*SequenceModel {
		return []*SequenceModel{
			{id: "sequence1", headEventID: "event1"},
			{id: "sequence2", headEventID: "event3"},
		}
	}

	eventDataPort := func() []*EventDataModel {
		return []*EventDataModel{
			{id: "event2", nextEventID: ""},
			{id: "event1", nextEventID: "event2"},
			{id: "event3", nextEventID: ""},
		}
	}

	expected := []*sequenceData{
		{id: "sequence1", events: []*EventDataModel{
			{id: "event1", nextEventID: "event2"},
			{id: "event2", nextEventID: ""},
		}},
		{id: "sequence2", events: []*EventDataModel{
			{id: "event3", nextEventID: ""},
		}},
	}

	sequenceDataPort := initializeSequenceDataAdapter(sequenceModelPort, eventDataPort)
	sequenceData := sequenceDataPort()

	if !reflect.DeepEqual(sequenceData, expected) {
		t.Errorf("expected %+v, got %+v", expected, sequenceData)
	}
}
