package core

type SkillId string

type ServeSkillData func(id SkillId) *SkillData

type SkillData struct {
	Id         SkillId
	Name       TextId
	SequenceId EventSequenceId
}

type EventSequenceId string

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
	Text  TextId
}

func (e *DisplayMessageEvent) IsActive(frame int) bool {
	return e.Frame == frame
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
