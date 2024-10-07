package component

import (
	"fmt"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/widget"
)

// EventSequenceId is an identifier for the event sequence.
// One Skill has one sequence, so it is the same as SkillId.
type EventSequenceId string

const (
	EventSequenceIdPunchingBagBeaten EventSequenceId = "punching_bag_beaten"
)

type BattleTextDisplay interface {
	SetText(text string, displayAll bool)
}

type ShakeActor func(core.ActorId)
type ChangeEmotion func(core.ActorId, BattleEmotionType)
type ShakeScreen func()
type DisplayDamageFunc func(core.ActorId, core.Damage, core.HP)
type PlayEffect func(widget.EffectId, core.ActorId)
type SetDisappear func(core.ActorId)

type SkillToSequenceFunc func(core.SkillId) EventSequenceId

func ToEventSequenceId(skillId core.SkillId) EventSequenceId {
	return EventSequenceId(skillId)
}

func CreateServeBattleEventSequence() ServeBattleEventSequenceFunc {
	dict := map[EventSequenceId]*BattleEventSequence{}
	register := func(id core.SkillId, rows []BattleEventRow) {
		eventId := ToEventSequenceId(id)
		dict[eventId] = &BattleEventSequence{
			Id:   eventId,
			Rows: rows,
		}
	}
	normalAttack := []BattleEventRow{
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
	register(core.SkillIdLuneAttack, normalAttack)
	register(core.SkillIdNormalTackle, normalAttack)
	luneFire := []BattleEventRow{
		&DisplayMessageEvent{
			Frame: 1,
			Text:  core.TextIdLuneFireDesc,
		},
		&PlayEffectEvent{
			Frame:    1,
			EffectId: widget.EffectIdLuneFire,
		},
		&ShakeActorAnimationEvent{
			Frame: 66,
		},
		&DisplayDamageEvent{
			Frame: 66,
		},
		&ChangeEmotionEvent{
			Frame:       22,
			EmotionType: BattleEmotionDamage,
		},
		&ChangeEmotionEvent{
			Frame:       96,
			EmotionType: BattleEmotionNormal,
		},
	}
	register(core.SkillIdLuneFireEnemy, luneFire)
	combinationThunder := []BattleEventRow{
		&DisplayMessageEvent{
			Frame: 1,
			Text:  core.TextIdCombinationThunder,
		},
		&PlayEffectEvent{
			Frame:    60,
			EffectId: widget.EffectIdExplode,
		},
		&ChangeEmotionEvent{
			Frame:       126,
			EmotionType: BattleEmotionDamage,
		},
		&ShakeActorAnimationEvent{
			Frame: 126,
		},
		&DisplayDamageEvent{
			Frame: 126,
		},
		&ChangeEmotionEvent{
			Frame:       180,
			EmotionType: BattleEmotionNormal,
		},
	}
	register(core.SkillIdCombinationThunder, combinationThunder)
	punchingBagBeaten := []BattleEventRow{
		&ChangeEmotionEvent{
			Frame:       1,
			EmotionType: BattleEmotionDamage,
		},
		&EnemyDisappearEvent{
			Frame: 1,
		},
		&DisplayMessageEvent{
			Frame: 1,
			Text:  core.TextIdEnemyBeaten,
		},
	}
	dict[EventSequenceIdPunchingBagBeaten] = &BattleEventSequence{
		Id:   EventSequenceIdPunchingBagBeaten,
		Rows: punchingBagBeaten,
	}
	return func(id EventSequenceId) *BattleEventSequence {
		data, ok := dict[id]
		if !ok {
			panic(fmt.Sprintf("event sequence not found: %s", id))
		}
		return data
	}
}

type DamageInformation struct {
	Target  core.ActorId
	Damage  core.Damage
	AfterHP core.HP
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
	SetDisappear,
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
		setDisappear SetDisappear,
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
						damageMap := func() map[core.ActorId][]*DamageInformation {
							result := map[core.ActorId][]*DamageInformation{}
							for _, target := range args.Target {
								result[target.Target] = append(result[target.Target], target)
							}
							return result
						}()
						for target, damages := range damageMap {
							allDamage := func() core.Damage {
								var result core.Damage = 0
								for _, d := range damages {
									result += d.Damage
								}
								return result
							}()
							finalAfterHp := damages[len(damages)-1].AfterHP
							displayDamage(target, allDamage, finalAfterHp)
						}
					case *ChangeEmotionEvent:
						target := args.Target[0]
						changeEmotion(target.Target, r.EmotionType)
					case *PlayEffectEvent:
						target := args.Target[0]
						playEffect(r.EffectId, target.Target)
					case *EnemyDisappearEvent:
						target := args.Target[0]
						setDisappear(target.Target)
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

type ReferencePoint int

const (
	ReferencePointActorCenter ReferencePoint = iota
	ReferencePointScreenCenter
)

type PlayEffectEvent struct {
	Frame          int
	EffectId       widget.EffectId
	ReferencePoint ReferencePoint
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

type EnemyDisappearEvent struct {
	Frame int
}

func (e *EnemyDisappearEvent) IsActive(frame int) bool {
	return e.Frame == frame
}

func (e *EnemyDisappearEvent) IsEnd(frame int) bool {
	const disappearFrame = 60
	return e.Frame+disappearFrame < frame
}

type BattleEventSequencer struct {
	index    int
	sequence []BattleSequenceFunc
}

func NewBattleEventSequencer() *BattleEventSequencer {
	return &BattleEventSequencer{
		index:    0,
		sequence: []BattleSequenceFunc{},
	}
}

func (s *BattleEventSequencer) Add(sequence BattleSequenceFunc) {
	s.sequence = append(s.sequence, sequence)
}

func (s *BattleEventSequencer) Update() {
	if s.index >= len(s.sequence) {
		return
	}
	result := s.sequence[s.index]()
	if result.IsEnd {
		s.index++
	}
}

func (s *BattleEventSequencer) IsEnd() bool {
	return s.index >= len(s.sequence)
}

func (s *BattleEventSequencer) IsRun() bool {
	return s.index < len(s.sequence)
}

func (s *BattleEventSequencer) Reset() {
	s.index = 0
	s.sequence = []BattleSequenceFunc{}
}
