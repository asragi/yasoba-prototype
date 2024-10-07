package core

import (
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
	Actors []*Actor
}

type BattleEndType int

const (
	BattleEndTypeNone BattleEndType = iota
	BattleEndTypeWin
	BattleEndTypeLose
)

func (s *BattleState) IsBattleShouldBeEnd() BattleEndType {
	if s.IsAllBeaten(ActorSidePlayer) {
		return BattleEndTypeLose
	}
	if s.IsAllBeaten(ActorSideEnemy) {
		return BattleEndTypeWin
	}
	return BattleEndTypeNone
}

func (s *BattleState) IsAllBeaten(side ActorSide) bool {
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

func (s *BattleState) GetOtherSideActors(actionActor *Actor) []*Actor {
	side := actionActor.Side
	var result []*Actor
	for _, actor := range s.Actors {
		if actor.Side == side {
			continue
		}
		result = append(result, actor)
	}
	return result
}

func (s *BattleState) GetMainActor() *Actor {
	for _, actor := range s.Actors {
		if actor.IsMainActor() {
			return actor
		}
	}
	return nil
}

func (s *BattleState) GetSubActor() *Actor {
	for _, actor := range s.Actors {
		if actor.IsSubActor() {
			return actor
		}
	}
	return nil
}

type BattleAction struct {
	SelectedSkill  SkillId
	TargetActorIds []ActorId
}

type ServeBattleState func() *BattleState

func CreateServeBattleState(supplyActor AllActorServer) ServeBattleState {
	return func() *BattleState {
		actors := supplyActor.GetAllActor()
		return &BattleState{
			Actors: actors,
		}
	}
}

type DecideActionFunc func(*Actor, *BattleState) *BattleAction
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
		return func(actor *Actor, state *BattleState) *BattleAction {
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
	actor *Actor,
	actors *BattleState,
) []ActorId

func CreateChoiceSkillTarget(getRandom util.EmitRandomFunc) ChoiceSkillTargetFunc {
	choiceSingleTarget := func(actionActor *Actor, state *BattleState) *Actor {
		random := getRandom()
		otherSideActors := state.GetOtherSideActors(actionActor)
		possibleActors := func() []*Actor {
			var result []*Actor
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
		actor *Actor,
		state *BattleState,
	) []ActorId {
		if skill.TargetType == SkillTargetTypeSingleOther {
			target := choiceSingleTarget(actor, state)
			return []ActorId{target.Id}
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
