package core

import "github.com/asragi/yasoba-prototype/util"

type PartnerActionForecast struct {
	SelectedSkill  SkillId
	SelectedTarget ActorId
}

type DecidePartnerActionFunc func(*BattleState) *PartnerActionForecast

func CreateDecidePartnerAction(random util.EmitRandomFunc) DecidePartnerActionFunc {
	skillList := []SkillId{
		SkillIdSunnyKick,
		SkillIdSunnyUppercut,
	}
	target := func(enemies []*Actor) ActorId {
		return enemies[0].Id
	}
	return func(s *BattleState) *PartnerActionForecast {
		subActor := s.GetSubActor()
		skill := skillList[int(random()*float64(len(skillList)))]
		return &PartnerActionForecast{
			SelectedSkill:  skill,
			SelectedTarget: target(s.GetOtherSideActors(subActor)),
		}
	}
}
