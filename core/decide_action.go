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

type DecideActionFunc func(state *BattleState) *BattleAction
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
		return func(state *BattleState) *BattleAction {
			random := getRandom()
			skillIndex := int(random * float64(len(skills)))
			skill := skills[skillIndex]
			targets := choiceSkillTarget(skill, state.Actors)
			return &BattleAction{
				SelectedSkill:  skill.Id,
				TargetActorIds: targets,
			}
		}
	}
}

type ChoiceSkillTargetFunc func(
	skill *SkillData,
	actors []*Actor,
) []ActorId

func CreateChoiceSkillTarget(getRandom util.EmitRandomFunc) ChoiceSkillTargetFunc {
	choiceSingleTarget := func(actor []*Actor) *Actor {
		random := getRandom()
		targetIndex := int(random * float64(len(actor)))
		return actor[targetIndex]
	}
	return func(
		skill *SkillData,
		actors []*Actor,
	) []ActorId {
		if skill.TargetType == SkillTargetTypeSingleOther {
			actor := choiceSingleTarget(actors)
			return []ActorId{actor.Id}
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
