package sequence

import (
	"testing"
)

func TestSequence_Update(t *testing.T) {
	tests := []struct {
		name     string
		events   []*eventUnit
		expected bool
	}{
		{
			name:     "空のイベントリストの場合、即座に終了",
			events:   []*eventUnit{},
			expected: true,
		},
		{
			name: "1つのイベントが終了する場合",
			events: []*eventUnit{
				{
					start:      func() {},
					checkIsEnd: func() IsEnd { return true },
				},
			},
			expected: true,
		},
		{
			name: "複数のイベントが順次実行される場合",
			events: []*eventUnit{
				{
					start:      func() {},
					checkIsEnd: func() IsEnd { return true },
				},
				{
					start:      func() {},
					checkIsEnd: func() IsEnd { return true },
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := &sequence{
				id:      "test-sequence",
				index:   0,
				isStart: true,
				events:  tt.events,
			}

			result := seq.Update()
			if result != IsEnd(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestProvideCreateSequence(t *testing.T) {
	// テスト用のモックデータ
	mockId := SequenceId("test-sequence-1")
	mockSequencesData := []*sequenceData{
		{
			id: mockId,
			events: []*EventDataModel{
				{
					id:        "event-1",
					eventType: partnerDialogueEvent,
					ownerId:   mockId,
					order:     1,
				},
				{
					id:        "event-2",
					eventType: openPartnerMessageWindowEvent,
					ownerId:   mockId,
					order:     2,
				},
				{
					id:        "event-3",
					eventType: closePartnerMessageWindowEvent,
					ownerId:   mockId,
					order:     3,
				},
			},
		},
	}

	mockSequenceDataPort := func() []*sequenceData {
		return mockSequencesData
	}

	mockCreatePartnerDialogueEvent := func(id EventID) *eventUnit {
		return &eventUnit{
			start:      func() {},
			checkIsEnd: func() IsEnd { return false },
			reset:      func() {},
		}
	}

	mockCreateChangeEmotionEvent := func(id EventID) *eventUnit {
		return &eventUnit{
			start:      func() {},
			checkIsEnd: func() IsEnd { return false },
			reset:      func() {},
		}
	}

	mockCreateOpenPartnerMessageWindowEvent := func(id EventID) *eventUnit {
		return &eventUnit{
			start:      func() {},
			checkIsEnd: func() IsEnd { return false },
			reset:      func() {},
		}
	}

	mockCreateClosePartnerMessageWindowEvent := func(id EventID) *eventUnit {
		return &eventUnit{
			start:      func() {},
			checkIsEnd: func() IsEnd { return false },
			reset:      func() {},
		}
	}

	// initializeProduceCreateSequenceの正しい呼び出し方法に修正
	produceCreateSequence := initializeProduceCreateSequence(mockSequenceDataPort)
	createSeq := produceCreateSequence(mockCreatePartnerDialogueEvent, mockCreateChangeEmotionEvent, mockCreateOpenPartnerMessageWindowEvent, mockCreateClosePartnerMessageWindowEvent)

	// 存在するシーケンスIDでテスト
	seq := createSeq(SequenceId("test-sequence-1"))
	if seq == nil {
		t.Fatal("expected sequence to be created, got nil")
	}
	if seq.id != mockId {
		t.Errorf("expected sequence id '%s', got '%s'", mockId, seq.id)
	}
	if len(seq.events) != 3 {
		t.Errorf("expected 3 events, got %d", len(seq.events))
	}

	// 存在しないシーケンスIDでテスト
	nonExistentSeq := createSeq(SequenceId("non-existent"))
	if nonExistentSeq != nil {
		t.Error("expected nil for non-existent sequence, got sequence")
	}
}

func TestOpenPartnerMessageWindowEvent(t *testing.T) {
	openCalled := false

	setOpenPartnerMessageWindow := func() {
		openCalled = true
	}

	createEvent := produceCreateOpenPartnerMessageWindowEventToUnit(setOpenPartnerMessageWindow)
	event := createEvent(EventID("test-event"))

	// startのテスト
	event.start()
	if !openCalled {
		t.Error("setOpenPartnerMessageWindow function was not called")
	}

	// checkIsEndのテスト
	result := event.checkIsEnd()
	if !result {
		t.Error("expected checkIsEnd to return true")
	}
}

func TestClosePartnerMessageWindowEvent(t *testing.T) {
	closeCalled := false

	setClosePartnerMessageWindow := func() {
		closeCalled = true
	}

	createEvent := produceCreateClosePartnerMessageWindowEventToUnit(setClosePartnerMessageWindow)
	event := createEvent(EventID("test-event"))

	// startのテスト
	event.start()
	if !closeCalled {
		t.Error("setClosePartnerMessageWindow function was not called")
	}

	// checkIsEndのテスト
	result := event.checkIsEnd()
	if !result {
		t.Error("expected checkIsEnd to return true")
	}
}

func TestEventUnit(t *testing.T) {
	startCalled := false
	updateCalled := false
	resetCalled := false

	event := &eventUnit{
		start:      func() { startCalled = true },
		checkIsEnd: func() IsEnd { updateCalled = true; return true },
		reset:      func() { resetCalled = true },
	}

	// startのテスト
	event.start()
	if !startCalled {
		t.Error("start function was not called")
	}

	// updateのテスト
	result := event.checkIsEnd()
	if !updateCalled {
		t.Error("update function was not called")
	}
	if !result {
		t.Error("expected update to return true")
	}

	// resetのテスト
	event.reset()
	if !resetCalled {
		t.Error("reset function was not called")
	}
}
