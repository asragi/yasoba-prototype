package battle_decision

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/util"
)

// BattleAction describes an auto-selected battle action.
type BattleAction struct {
	SelectedSkill  core.SkillId
	TargetActorIds []actor.ActorId
}

// DecideActionFunc determines a battle action for a given actor.
type DecideActionFunc func(*actor.Actor, *BattleState) *BattleAction

// NewChoiceRandomActionFunc prepares a random action selector given skill IDs.
type NewChoiceRandomActionFunc func([]core.SkillId) DecideActionFunc

// StandByCreateRandomAction prepares a function that picks random skills.
func StandByCreateRandomAction(
	getRandom util.EmitRandomFunc,
	serveSkillData core.ServeSkillData,
	choiceSkillTarget ChoiceSkillTargetFunc,
) NewChoiceRandomActionFunc {
	return func(skillIds []core.SkillId) DecideActionFunc {
		skills := func() []*core.SkillData {
			var result []*core.SkillData
			for _, id := range skillIds {
				result = append(result, serveSkillData(id))
			}
			return result
		}()
		return func(actionActor *actor.Actor, state *BattleState) *BattleAction {
			random := getRandom()
			skillIndex := int(random * float64(len(skills)))
			skill := skills[skillIndex]
			targets := choiceSkillTarget(skill, actionActor, state)
			return &BattleAction{
				SelectedSkill:  skill.SkillId,
				TargetActorIds: targets,
			}
		}
	}
}

// ChoiceSkillTargetFunc selects targets for a skill.
type ChoiceSkillTargetFunc func(
	skill *core.SkillData,
	actionActor *actor.Actor,
	state *BattleState,
) []actor.ActorId

// CreateChoiceSkillTarget builds a target chooser using the provided RNG.
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
		skill *core.SkillData,
		actionActor *actor.Actor,
		state *BattleState,
	) []actor.ActorId {
		if skill.TargetType == core.SkillTargetTypeSingleOther {
			target := choiceSingleTarget(actionActor, state)
			return []actor.ActorId{target.Id}
		}
		panic("Not implemented")
	}
}

// NewChoiceActionFunc constructs DecideActionFunc from a choice action ID.
type NewChoiceActionFunc func(ChoiceActionId) DecideActionFunc

// CreateNewChoiceAction links choice IDs to specific random action factories.
func CreateNewChoiceAction(newChoiceRandomAction NewChoiceRandomActionFunc) NewChoiceActionFunc {
	dict := map[ChoiceActionId]DecideActionFunc{
		CharacterIdToChoiceActionId(core.CharacterSunnyId): newChoiceRandomAction([]core.SkillId{core.SkillIdNormalTackle}),
		EnemyIdToChoiceActionId(core.EnemyPunchingBagId):   newChoiceRandomAction([]core.SkillId{core.SkillIdNormalTackle}),
	}
	return func(id ChoiceActionId) DecideActionFunc {
		return dict[id]
	}
}
