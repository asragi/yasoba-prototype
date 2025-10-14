package battle_combination

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/core"
)

// Request carries the input data for checking skill combinations.
type Request struct {
	MainActorSkillId core.SkillId
	MainActorTarget  actor.ActorId
	SubActorSkillId  core.SkillId
	SubActorTarget   actor.ActorId
}

// Response reports whether a combination occurs and the resulting skill.
type Response struct {
	IsCombination bool
	SkillId       core.SkillId
	TargetId      actor.ActorId
}

// CheckFunc evaluates whether two skills should trigger a combination.
type CheckFunc func(request *Request) *Response

// CreateCheckCombination returns the default combination checker.
func CreateCheckCombination() CheckFunc {
	combinationDict := map[core.SkillId]map[core.SkillId]core.SkillId{
		core.SkillIdLuneFireEnemy: {
			core.SkillIdSunnyUppercut: core.SkillIdCombinationThunder,
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
