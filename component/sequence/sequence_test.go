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
					start:  func() {},
					update: func() isEnd { return true },
					render: func() {},
				},
			},
			expected: true,
		},
		{
			name: "複数のイベントが順次実行される場合",
			events: []*eventUnit{
				{
					start:  func() {},
					update: func() isEnd { return true },
					render: func() {},
				},
				{
					start:  func() {},
					update: func() isEnd { return true },
					render: func() {},
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

			result := seq.update()
			if result != isEnd(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSequence_Render(t *testing.T) {
	renderCalled := false
	seq := &sequence{
		id: "test-sequence",
		events: []*eventUnit{
			{
				start:  func() {},
				update: func() isEnd { return false },
				render: func() { renderCalled = true },
			},
		},
	}

	seq.render()

	if !renderCalled {
		t.Error("render function was not called")
	}
}

func TestProvideCreateSequence(t *testing.T) {
	// テスト用のモックデータ
	mockId := sequenceId("test-sequence-1")
	mockSequencesData := []*sequenceData{
		{
			id: mockId,
			events: []*eventDataModel{
				{
					id:        "event-1",
					eventType: partnerDialogueEvent,
					ownerId:   mockId,
					order:     1,
				},
			},
		},
	}

	mockSequenceDataPort := func() []*sequenceData {
		return mockSequencesData
	}

	mockCreatePartnerDialogueEvent := func(id eventId) *eventUnit {
		return &eventUnit{
			start:  func() {},
			update: func() isEnd { return false },
			render: func() {},
		}
	}

	// initializeProduceCreateSequenceの正しい呼び出し方法に修正
	produceCreateSequence := initializeProduceCreateSequence(mockSequenceDataPort)
	createSeq := produceCreateSequence(mockCreatePartnerDialogueEvent)

	// 存在するシーケンスIDでテスト
	seq := createSeq("test-sequence-1")
	if seq == nil {
		t.Fatal("expected sequence to be created, got nil")
	}
	if seq.id != mockId {
		t.Errorf("expected sequence id '%s', got '%s'", mockId, seq.id)
	}
	if len(seq.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(seq.events))
	}

	// 存在しないシーケンスIDでテスト
	nonExistentSeq := createSeq("non-existent")
	if nonExistentSeq != nil {
		t.Error("expected nil for non-existent sequence, got sequence")
	}
}

func TestEventUnit(t *testing.T) {
	startCalled := false
	updateCalled := false
	renderCalled := false

	event := &eventUnit{
		start:  func() { startCalled = true },
		update: func() isEnd { updateCalled = true; return true },
		render: func() { renderCalled = true },
	}

	// startのテスト
	event.start()
	if !startCalled {
		t.Error("start function was not called")
	}

	// updateのテスト
	result := event.update()
	if !updateCalled {
		t.Error("update function was not called")
	}
	if !result {
		t.Error("expected update to return true")
	}

	// renderのテスト
	event.render()
	if !renderCalled {
		t.Error("render function was not called")
	}
}
