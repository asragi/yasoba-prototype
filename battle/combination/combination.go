package combination

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/skill"
)

// Request carries the input data for checking skill combinations.
type Request struct {
	MainActorSkillId skill.SkillId
	MainActorTarget  actor.ActorId
	SubActorSkillId  skill.SkillId
	SubActorTarget   actor.ActorId
}

// Response reports whether a combination occurs and the resulting skill.
type Response struct {
	IsCombination bool
	SkillId       skill.SkillId
	TargetId      actor.ActorId
}

// CheckFunc evaluates whether two skills should trigger a combination.
type CheckFunc func(request *Request) *Response

// CreateCheckCombination returns the default combination checker.

func CreateCheckCombination() CheckFunc {
	combinationDict := map[skill.SkillId]map[skill.SkillId]skill.SkillId{
		skill.SkillIdLuneFireEnemy: {
			skill.SkillIdSunnyUppercut: skill.SkillIdCombinationThunder,
		},
	}
	return func(request *Request) *Response {
		onFailure := func() *Response {
			return &Response{
				IsCombination: false,
			}
		}
		if request.MainActorTarget != request.SubActorTarget {
			return onFailure()
		}
		combination, ok := combinationDict[request.MainActorSkillId]
		if !ok {
			return onFailure()
		}
		skillId, ok := combination[request.SubActorSkillId]
		if !ok {
			return onFailure()
		}
		return &Response{
			IsCombination: true,
			SkillId:       skillId,
			TargetId:      request.MainActorTarget,
		}
	}
}
