package core

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/util"
)

// ChoiceActionId is issued for each one enemy
type ChoiceActionId string

func CharacterIdToChoiceActionId(id CharacterId) ChoiceActionId {
	return ChoiceActionId(id)
}

func EnemyIdToChoiceActionId(id EnemyId) ChoiceActionId {
	return ChoiceActionId(id)
}

type BattleState struct {
	Actors []*actor.Actor
}

type BattleEndType int

const (
	BattleEndTypeNone BattleEndType = iota
	BattleEndTypeWin
	BattleEndTypeLose
)

func (s *BattleState) IsBattleShouldBeEnd() BattleEndType {
	if s.IsAllBeaten(actor.ActorSidePlayer) {
		return BattleEndTypeLose
	}
	if s.IsAllBeaten(actor.ActorSideEnemy) {
		return BattleEndTypeWin
	}
	return BattleEndTypeNone
}

func (s *BattleState) IsAllBeaten(side actor.ActorSide) bool {
	for _, actor := range s.Actors {
		if actor.Side != side {
			continue
		}
		if !actor.IsBeaten() {
			return false
		}
	}
	return true
}

func (s *BattleState) GetOtherSideActors(actionActor *actor.Actor) []*actor.Actor {
	side := actionActor.Side
	var result []*actor.Actor
	for _, actor := range s.Actors {
		if actor.Side == side {
			continue
		}
		result = append(result, actor)
	}
	return result
}

func (s *BattleState) GetMainActor() *actor.Actor {
	for _, actor := range s.Actors {
		if actor.IsMainActor() {
			return actor
		}
	}
	return nil
}

func (s *BattleState) GetSubActor() *actor.Actor {
	for _, actor := range s.Actors {
		if actor.IsSubActor() {
			return actor
		}
	}
	return nil
}

type BattleAction struct {
	SelectedSkill  SkillId
	TargetActorIds []actor.ActorId
}

type ServeBattleState func() *BattleState

type allActorServer interface {
	GetAllActor() []*actor.Actor
}

func CreateServeBattleState(supplyActor allActorServer) ServeBattleState {
	return func() *BattleState {
		actors := supplyActor.GetAllActor()
		return &BattleState{
			Actors: actors,
		}
	}
}

type DecideActionFunc func(*actor.Actor, *BattleState) *BattleAction
type NewChoiceRandomActionFunc func([]SkillId) DecideActionFunc

func StandByCreateRandomAction(
	getRandom util.EmitRandomFunc,
	serveSkillData ServeSkillData,
	choiceSkillTarget ChoiceSkillTargetFunc,
) NewChoiceRandomActionFunc {
	return func(skillIds []SkillId) DecideActionFunc {
		skills := func() []*SkillData {
			var result []*SkillData
			for _, id := range skillIds {
				result = append(result, serveSkillData(id))
			}
			return result
		}()
		return func(actor *actor.Actor, state *BattleState) *BattleAction {
			random := getRandom()
			skillIndex := int(random * float64(len(skills)))
			skill := skills[skillIndex]
			targets := choiceSkillTarget(skill, actor, state)
			return &BattleAction{
				SelectedSkill:  skill.SkillId,
				TargetActorIds: targets,
			}
		}
	}
}

type ChoiceSkillTargetFunc func(
	skill *SkillData,
	actor *actor.Actor,
	actors *BattleState,
) []actor.ActorId

func CreateChoiceSkillTarget(getRandom util.EmitRandomFunc) ChoiceSkillTargetFunc {
	choiceSingleTarget := func(actionActor *actor.Actor, state *BattleState) *actor.Actor {
		random := getRandom()
		otherSideActors := state.GetOtherSideActors(actionActor)
		possibleActors := func() []*actor.Actor {
			var result []*actor.Actor
			for _, actor := range otherSideActors {
				if actor.IsBeaten() {
					continue
				}
				result = append(result, actor)
			}
			return result
		}()
		targetIndex := int(random * float64(len(possibleActors)))
		return possibleActors[targetIndex]
	}
	return func(
		skill *SkillData,
		actionActor *actor.Actor,
		state *BattleState,
	) []actor.ActorId {
		if skill.TargetType == SkillTargetTypeSingleOther {
			target := choiceSingleTarget(actionActor, state)
			return []actor.ActorId{target.Id}
		}
		panic("Not implemented")
	}
}

type NewChoiceActionFunc func(ChoiceActionId) DecideActionFunc

func CreateNewChoiceAction(newChoiceRandomAction NewChoiceRandomActionFunc) NewChoiceActionFunc {
	dict := map[ChoiceActionId]DecideActionFunc{
		CharacterIdToChoiceActionId(CharacterSunnyId): newChoiceRandomAction([]SkillId{SkillIdNormalTackle}),
		EnemyIdToChoiceActionId(EnemyPunchingBagId):   newChoiceRandomAction([]SkillId{SkillIdNormalTackle}),
	}
	return func(id ChoiceActionId) DecideActionFunc {
		return dict[id]
	}
}
