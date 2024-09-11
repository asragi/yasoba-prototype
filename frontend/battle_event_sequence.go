package frontend

import "github.com/asragi/yasoba-prototype/core"

type BattleTextDisplay interface {
	SetText(text string)
}

func CreateServeBattleEventSequence() ServeBattleEventSequenceFunc {
	dict := map[EventSequenceId]*BattleEventSequence{}
	testData := []BattleEventRow{
		&DisplayMessageEvent{
			Frame: 0,
			Text:  "test-text",
		},
	}
	dict["test"] = &BattleEventSequence{
		Id:   "test",
		Rows: testData,
	}
	return func(id EventSequenceId) *BattleEventSequence {
		return dict[id]
	}
}

type EventSequenceArgs struct {
	sequenceId EventSequenceId
}
type EventSequenceResult struct{}
type ExecBattleEventSequenceFunc func(*EventSequenceArgs, int) *EventSequenceResult

func CreateExecBattleEventSequence(
	display BattleTextDisplay,
	textServer core.ServeTextDataFunc,
	serveEvent ServeBattleEventSequenceFunc,
) ExecBattleEventSequenceFunc {
	return func(args *EventSequenceArgs, frame int) *EventSequenceResult {
		sequence := serveEvent(args.sequenceId)
		for _, row := range sequence.Rows {
			if !row.IsActive(frame) {
				continue
			}
			switch r := row.(type) {
			case *DisplayMessageEvent:
				text := textServer(r.Text)
				display.SetText(text.Text)
			}
		}
		return &EventSequenceResult{}
	}
}

type BattleEventRow interface {
	IsActive(frame int) bool
}

type BattleEventSequence struct {
	Id   EventSequenceId
	Rows []BattleEventRow
}

type ServeBattleEventSequenceFunc func(id EventSequenceId) *BattleEventSequence

type DisplayMessageEvent struct {
	Frame int
	Text  core.TextId
}

func (e *DisplayMessageEvent) IsActive(frame int) bool {
	return e.Frame == frame
}

type EventSequenceId string
