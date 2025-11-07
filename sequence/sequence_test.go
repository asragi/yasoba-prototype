package sequence

import (
	"testing"

	"github.com/asragi/yasoba-prototype/battle"
	seqtransition "github.com/asragi/yasoba-prototype/sequence/transition"
	"github.com/asragi/yasoba-prototype/text"
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
					id:          "event-1",
					eventType:   partnerDialogueEvent,
					nextEventID: "event-2",
				},
				{
					id:          "event-2",
					eventType:   openPartnerMessageWindowEvent,
					nextEventID: "event-3",
				},
				{
					id:          "event-3",
					eventType:   closePartnerMessageWindowEvent,
					nextEventID: "event-4",
				},
				{
					id:          "event-4",
					eventType:   showPlayerCommandWindowEvent,
					nextEventID: "event-5",
				},
				{
					id:          "event-5",
					eventType:   hidePlayerCommandWindowEvent,
					nextEventID: "event-6",
				},
				{
					id:          "event-6",
					eventType:   startTransitionFadeOutEvent,
					nextEventID: "event-7",
				},
				{
					id:          "event-7",
					eventType:   startTransitionFadeInEvent,
					nextEventID: "event-8",
				},
				{
					id:          "event-8",
					eventType:   switchToBattleSceneEvent,
					nextEventID: "event-9",
				},
				{
					id:          "event-9",
					eventType:   switchToDebugSceneEvent,
					nextEventID: "",
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

	mockCreateSetMessageWindowTextEvent := func(id EventID) *eventUnit {
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
	createEventUnit := func() *eventUnit {
		return &eventUnit{
			start:      func() {},
			checkIsEnd: func() IsEnd { return false },
			reset:      func() {},
		}
	}

	createSeq := produceCreateSequence(
		mockCreatePartnerDialogueEvent,
		mockCreateSetMessageWindowTextEvent,
		mockCreateChangeEmotionEvent,
		mockCreateOpenPartnerMessageWindowEvent,
		mockCreateClosePartnerMessageWindowEvent,
		func(EventID) *eventUnit { return createEventUnit() },
		func(EventID) *eventUnit { return createEventUnit() },
		func(EventID) *eventUnit { return createEventUnit() },
		func(EventID) *eventUnit { return createEventUnit() },
		func(EventID) *eventUnit { return createEventUnit() },
		func(EventID) *eventUnit { return createEventUnit() },
		func(EventID) *eventUnit { return createEventUnit() },
	)

	// 存在するシーケンスIDでテスト
	seq := createSeq(SequenceId("test-sequence-1"))
	if seq == nil {
		t.Fatal("expected sequence to be created, got nil")
	}
	if seq.id != mockId {
		t.Errorf("expected sequence id '%s', got '%s'", mockId, seq.id)
	}
	if len(seq.events) != 9 {
		t.Errorf("expected 9 events, got %d", len(seq.events))
	}

	// 存在しないシーケンスIDでテスト
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-existent sequence")
		}
	}()
	createSeq(SequenceId("non-existent"))
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

func TestSetMessageWindowTextEvent(t *testing.T) {
	textId := text.TextIdBattleLose
	dataPort := func(id EventID) *MessageWindowTextModel {
		return NewMessageWindowTextModel(id, textId, true)
	}
	serveTextData := func(id text.TextId) *text.Data {
		return &text.Data{
			Id:   id,
			Text: "lose",
		}
	}

	waitCount := 0
	setter := func(value text.String) *SetMessageWindowTextResponse {
		if value.String() != "lose" {
			t.Fatalf("unexpected text: %s", value)
		}
		return &SetMessageWindowTextResponse{
			CheckIsEnd: func() IsEnd {
				waitCount++
				return waitCount >= 2
			},
		}
	}

	createEvent := produceCreateSetMessageWindowTextEventToUnit(serveTextData, dataPort, setter)
	event := createEvent(EventID("battle_lose_message"))

	event.start()
	if event.checkIsEnd() {
		t.Fatal("expected event to wait for completion")
	}
	if !event.checkIsEnd() {
		t.Fatal("expected event to finish after completion")
	}
}

func TestSetMessageWindowTextEvent_NoWait(t *testing.T) {
	dataPort := func(id EventID) *MessageWindowTextModel {
		return NewMessageWindowTextModel(id, text.TextIdBattleLose, false)
	}
	serveTextData := func(id text.TextId) *text.Data {
		return &text.Data{
			Id:   id,
			Text: "lose",
		}
	}
	setter := func(value text.String) *SetMessageWindowTextResponse {
		return &SetMessageWindowTextResponse{
			CheckIsEnd: func() IsEnd { return false },
		}
	}
	createEvent := produceCreateSetMessageWindowTextEventToUnit(serveTextData, dataPort, setter)
	event := createEvent(EventID("battle_lose_message"))
	event.start()
	if !event.checkIsEnd() {
		t.Fatal("expected event to finish immediately when waitForComplete is false")
	}
}

func TestWaitFrameEvent(t *testing.T) {
	dataPort := func(id EventID) *WaitFrameModel {
		return NewWaitFrameModel(id, 2)
	}
	createEvent := produceCreateWaitFrameEventToUnit(dataPort)
	event := createEvent(EventID("wait_event"))
	event.start()
	if event.checkIsEnd() {
		t.Fatal("expected wait to continue on first frame")
	}
	if !event.checkIsEnd() {
		t.Fatal("expected wait to finish on second frame")
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

func TestStartTransitionFadeOutEvent_WaitOption(t *testing.T) {
	eventID := EventID("fade-out")
	running := true
	optionPort := func(id string) *seqtransition.Option {
		if id != string(eventID) {
			return nil
		}
		return seqtransition.NewOption(true)
	}
	controller := seqtransition.NewController(
		func() {},
		func() bool { return running },
	)
	createEvent := produceCreateStartTransitionFadeOutEventToUnit(optionPort, controller)
	event := createEvent(eventID)

	event.start()
	if event.checkIsEnd() {
		t.Fatal("expected event to wait while transition is running")
	}

	running = false
	if !event.checkIsEnd() {
		t.Fatal("expected event to finish after transition stops")
	}
}

func TestStartTransitionFadeInEvent_NoWaitByDefault(t *testing.T) {
	eventID := EventID("fade-in")
	optionPort := func(string) *seqtransition.Option {
		return nil
	}
	controller := seqtransition.NewController(
		func() {},
		func() bool { return true },
	)
	createEvent := produceCreateStartTransitionFadeInEventToUnit(optionPort, controller)
	event := createEvent(eventID)

	event.start()
	if !event.checkIsEnd() {
		t.Fatal("expected event to finish immediately when wait option is disabled")
	}
}

func TestProduceCreateSwitchToBattleSceneEventToUnit(t *testing.T) {
	called := false
	dataPort := func(id EventID) *SwitchToBattleSceneModel {
		if id != EventID("switch-battle") {
			return nil
		}
		return NewSwitchToBattleSceneModel(battle.BattleId("battle-extra"))
	}
	createEvent := produceCreateSwitchToBattleSceneEventToUnit(
		dataPort,
		func(received battle.BattleId) {
			if received != battle.BattleId("battle-extra") {
				t.Fatalf("unexpected battle id: %s", received)
			}
			called = true
		},
	)
	event := createEvent(EventID("switch-battle"))
	event.start()
	if !called {
		t.Fatal("expected switchToBattle to be called")
	}
}

func TestProduceCreateSwitchToDebugSceneEventToUnit(t *testing.T) {
	called := false
	createEvent := produceCreateSwitchToDebugSceneEventToUnit(func() {
		called = true
	})
	event := createEvent(EventID("switch-debug"))
	event.start()
	if !called {
		t.Fatal("expected switchToDebug to be called")
	}
}
