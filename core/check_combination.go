package core

type CheckCombinationRequest struct {
	MainActorSkillId SkillId
	MainActorTarget  ActorId
	SubActorSkillId  SkillId
	SubActorTarget   ActorId
}

type CheckCombinationResponse struct {
	IsCombination bool
	SkillId       SkillId
	TargetId      ActorId
}
type CheckCombinationFunc func(request *CheckCombinationRequest) *CheckCombinationResponse

func CreateCheckCombination() CheckCombinationFunc {
	combinationDict := map[SkillId]map[SkillId]SkillId{
		SkillIdLuneFireEnemy: {
			SkillIdSunnyUppercut: SkillIdCombinationThunder,
		},
	}
	return func(request *CheckCombinationRequest) *CheckCombinationResponse {
		onFailure := func() *CheckCombinationResponse {
			return &CheckCombinationResponse{
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
		return &CheckCombinationResponse{
			IsCombination: true,
			SkillId:       skillId,
			TargetId:      request.MainActorTarget,
		}
	}
}
