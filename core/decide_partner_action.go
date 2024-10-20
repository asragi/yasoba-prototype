package core

import "github.com/asragi/yasoba-prototype/util"

type PartnerActionPlan struct {
	SkillId        SkillId
	SelectedTarget ActorId
}

type decidePartnerPlanFunc func()

func createDecidePartnerAction(random util.EmitRandomFunc, state *BattleState) *PartnerActionPlan {
	skillList := []SkillId{
		SkillIdSunnyKick,
		SkillIdSunnyUppercut,
	}
	target := func(enemies []*Actor) ActorId {
		return enemies[0].Id
	}
	return func() *PartnerActionPlan {
		subActor := state.GetSubActor()
		skill := skillList[int(random()*float64(len(skillList)))]
		return &PartnerActionPlan{
			SkillId:        skill,
			SelectedTarget: target(state.GetOtherSideActors(subActor)),
		}
	}()
}

type GetPartnerPlanFunc func() *PartnerActionPlan

type PartnerActionServer struct {
	StoredPlan *PartnerActionPlan
	random     util.EmitRandomFunc
	serveState ServeBattleState
}

type NewPartnerActionServer func() *PartnerActionServer

func StandByNewPartnerActionServer(
	random util.EmitRandomFunc,
	serveState ServeBattleState,
) NewPartnerActionServer {
	return func() *PartnerActionServer {
		return &PartnerActionServer{
			random:     random,
			serveState: serveState,
		}
	}
}

func (s *PartnerActionServer) DecidePlan() {
	state := s.serveState()
	s.StoredPlan = createDecidePartnerAction(s.random, state)
}

func (s *PartnerActionServer) GetPlan() *PartnerActionPlan {
	return s.StoredPlan
}
