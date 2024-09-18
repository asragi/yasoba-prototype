package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/widget"
)

type EventSequenceId string

const (
	EventSequenceIdLuneAttack EventSequenceId = "lune-attack"
)

type BattleTextDisplay interface {
	SetText(text string, displayAll bool)
}

type ShakeActor func(core.ActorId)
type ChangeEmotion func(core.ActorId, BattleEmotionType)
type ShakeScreen func()
type DisplayDamageFunc func(core.ActorId, core.Damage)
type PlayEffect func(widget.EffectId, core.ActorId)

type SkillToSequenceFunc func(core.SkillId) EventSequenceId

func CreateSkillToSequenceId() SkillToSequenceFunc {
	dict := map[core.SkillId]EventSequenceId{
		core.SkillIdLuneAttack: EventSequenceIdLuneAttack,
	}
	return func(id core.SkillId) EventSequenceId {
		return dict[id]
	}
}

func CreateServeBattleEventSequence() ServeBattleEventSequenceFunc {
	dict := map[EventSequenceId]*BattleEventSequence{}
	testData := []BattleEventRow{
		&DisplayMessageEvent{
			Frame: 1,
			Text:  core.TextIdLuneAttackDesc,
		},
		&PlayEffectEvent{
			Frame:    1,
			EffectId: widget.EffectIdLuneAttack,
		},
		&ShakeActorAnimationEvent{
			Frame: 30,
		},
		&DisplayDamageEvent{
			Frame: 30,
		},
		&ChangeEmotionEvent{
			Frame:       30,
			EmotionType: BattleEmotionDamage,
		},
		&ChangeEmotionEvent{
			Frame:       60,
			EmotionType: BattleEmotionNormal,
		},
	}
	dict[EventSequenceIdLuneAttack] = &BattleEventSequence{
		Id:   EventSequenceIdLuneAttack,
		Rows: testData,
	}
	return func(id EventSequenceId) *BattleEventSequence {
		return dict[id]
	}
}

type DamageInformation struct {
	Target core.ActorId
	Damage core.Damage
}

type EventSequenceArgs struct {
	SequenceId EventSequenceId
	Actor      core.ActorId
	Target     []*DamageInformation
}

type EventSequenceResult struct {
	IsEnd bool
}
type BattleSequenceFunc func() *EventSequenceResult
type NewBattleSequenceFunc func(*EventSequenceArgs) BattleSequenceFunc
type PrepareBattleEventSequenceFunc func(
	BattleTextDisplay,
	ShakeActor,
	ChangeEmotion,
	DisplayDamageFunc,
	PlayEffect,
) NewBattleSequenceFunc

func CreateExecBattleEventSequence(
	textServer core.ServeTextDataFunc,
	serveEvent ServeBattleEventSequenceFunc,
) PrepareBattleEventSequenceFunc {
	return func(
		display BattleTextDisplay,
		shakeActor ShakeActor,
		changeEmotion ChangeEmotion,
		displayDamage DisplayDamageFunc,
		playEffect PlayEffect,
	) NewBattleSequenceFunc {
		return func(args *EventSequenceArgs) BattleSequenceFunc {
			sequence := serveEvent(args.SequenceId)
			frame := 0
			return func() *EventSequenceResult {
				frame++
				isEnd := true
				for _, row := range sequence.Rows {
					if !row.IsEnd(frame) {
						isEnd = false
					}
					if !row.IsActive(frame) {
						continue
					}
					switch r := row.(type) {
					case *DisplayMessageEvent:
						text := textServer(r.Text)
						display.SetText(text.Text, false)
					case *ShakeActorAnimationEvent:
						target := args.Target[0]
						shakeActor(target.Target)
					case *DisplayDamageEvent:
						target := args.Target[0]
						displayDamage(target.Target, target.Damage)
					case *ChangeEmotionEvent:
						target := args.Target[0]
						changeEmotion(target.Target, r.EmotionType)
					case *PlayEffectEvent:
						target := args.Target[0]
						playEffect(r.EffectId, target.Target)
					}
				}
				return &EventSequenceResult{
					IsEnd: isEnd,
				}
			}
		}
	}
}

type BattleEventRow interface {
	IsActive(frame int) bool
	IsEnd(frame int) bool
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

func (e *DisplayMessageEvent) IsEnd(frame int) bool {
	return e.Frame < frame
}

type ShakeActorAnimationEvent struct {
	Frame int
}

func (e *ShakeActorAnimationEvent) IsActive(frame int) bool {
	return e.Frame == frame
}

func (e *ShakeActorAnimationEvent) IsEnd(frame int) bool {
	return e.Frame < frame
}

type DisplayDamageEvent struct {
	Frame int
}

func (e *DisplayDamageEvent) IsActive(frame int) bool {
	return e.Frame == frame
}

func (e *DisplayDamageEvent) IsEnd(frame int) bool {
	return e.Frame < frame
}

type PlayEffectEvent struct {
	Frame    int
	EffectId widget.EffectId
}

func (e *PlayEffectEvent) IsActive(frame int) bool {
	return e.Frame == frame
}

func (e *PlayEffectEvent) IsEnd(frame int) bool {
	return e.Frame < frame
}

type ChangeEmotionEvent struct {
	Frame       int
	EmotionType BattleEmotionType
}

func (e *ChangeEmotionEvent) IsActive(frame int) bool {
	return e.Frame == frame
}

func (e *ChangeEmotionEvent) IsEnd(frame int) bool {
	return e.Frame < frame
}
