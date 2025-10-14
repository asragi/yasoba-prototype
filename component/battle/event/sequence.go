package event

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/actor"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	battleemotion "github.com/asragi/yasoba-prototype/component/battle/emotion"
	gameSkill "github.com/asragi/yasoba-prototype/game/skill"
	"github.com/asragi/yasoba-prototype/text"
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

type ShakeActor func(actor.ActorId)
type ChangeEmotion func(actor.ActorId, battleemotion.BattleEmotionType)
type ShakeScreen func()
type DisplayDamageFunc func(actor.ActorId, battleSkill.Damage, actor.HP)
type PlayEffect func(widget.EffectId, actor.ActorId)
type SetDisappear func(actor.ActorId)

type SkillToSequenceFunc func(gameSkill.SkillId) EventSequenceId

func ToEventSequenceId(skillId gameSkill.SkillId) EventSequenceId {
	return EventSequenceId(skillId)
}

func CreateServeBattleEventSequence() ServeBattleEventSequenceFunc {
	dict := map[EventSequenceId]*BattleEventSequence{}
	register := func(id gameSkill.SkillId, rows []BattleEventRow) {
		eventId := ToEventSequenceId(id)
		dict[eventId] = &BattleEventSequence{
			Id:   eventId,
			Rows: rows,
		}
	}
	normalAttack := []BattleEventRow{
		&DisplayMessageEvent{
			Frame: 1,
			Text:  text.TextIdLuneAttackDesc,
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
			EmotionType: battleemotion.BattleEmotionDamage,
		},
		&ChangeEmotionEvent{
			Frame:       60,
			EmotionType: battleemotion.BattleEmotionNormal,
		},
	}
	register(gameSkill.SkillIdLuneAttack, normalAttack)
	register(gameSkill.SkillIdNormalTackle, normalAttack)
	luneFire := []BattleEventRow{
		&DisplayMessageEvent{
			Frame: 1,
			Text:  text.TextIdLuneFireDesc,
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
			EmotionType: battleemotion.BattleEmotionDamage,
		},
		&ChangeEmotionEvent{
			Frame:       96,
			EmotionType: battleemotion.BattleEmotionNormal,
		},
	}
	register(gameSkill.SkillIdLuneFireEnemy, luneFire)
	combinationThunder := []BattleEventRow{
		&DisplayMessageEvent{
			Frame: 1,
			Text:  text.TextIdCombinationThunder,
		},
		&PlayEffectEvent{
			Frame:    60,
			EffectId: widget.EffectIdExplode,
		},
		&ChangeEmotionEvent{
			Frame:       126,
			EmotionType: battleemotion.BattleEmotionDamage,
		},
		&ShakeActorAnimationEvent{
			Frame: 126,
		},
		&DisplayDamageEvent{
			Frame: 126,
		},
		&ChangeEmotionEvent{
			Frame:       180,
			EmotionType: battleemotion.BattleEmotionNormal,
		},
	}
	register(gameSkill.SkillIdCombinationThunder, combinationThunder)
	punchingBagBeaten := []BattleEventRow{
		&ChangeEmotionEvent{
			Frame:       1,
			EmotionType: battleemotion.BattleEmotionDamage,
		},
		&EnemyDisappearEvent{
			Frame: 1,
		},
		&DisplayMessageEvent{
			Frame: 1,
			Text:  text.TextIdEnemyBeaten,
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
	Target  actor.ActorId
	Damage  battleSkill.Damage
	AfterHP actor.HP
}

type EventSequenceArgs struct {
	SequenceId EventSequenceId
	Actor      actor.ActorId
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
	textServer text.ServeTextDataFunc,
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
						display.SetText(text.Text.String(), false)
					case *ShakeActorAnimationEvent:
						target := args.Target[0]
						shakeActor(target.Target)
					case *DisplayDamageEvent:
						damageMap := func() map[actor.ActorId][]*DamageInformation {
							result := map[actor.ActorId][]*DamageInformation{}
							for _, target := range args.Target {
								result[target.Target] = append(result[target.Target], target)
							}
							return result
						}()
						for target, damages := range damageMap {
							allDamage := func() battleSkill.Damage {
								var result battleSkill.Damage = 0
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
	Text  text.TextId
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
	EmotionType battleemotion.BattleEmotionType
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
