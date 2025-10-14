package battle_partner

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle_decision"
	"github.com/asragi/yasoba-prototype/skilldata"
	"github.com/asragi/yasoba-prototype/util"
)

// PartnerActionPlan describes the partner's planned action.
type PartnerActionPlan struct {
	SkillId        skilldata.SkillId
	SelectedTarget actor.ActorId
}

type decidePartnerPlanFunc func()

func createDecidePartnerAction(random util.EmitRandomFunc, state *battle_decision.BattleState) *PartnerActionPlan {
	skillList := []skilldata.SkillId{
		skilldata.SkillIdSunnyKick,
		skilldata.SkillIdSunnyUppercut,
	}
	target := func(enemies []*actor.Actor) actor.ActorId {
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

// GetPartnerPlanFunc retrieves the stored partner plan.
type GetPartnerPlanFunc func() *PartnerActionPlan

// PartnerActionServer stores and updates partner action plans.
type PartnerActionServer struct {
	StoredPlan *PartnerActionPlan
	random     util.EmitRandomFunc
	serveState battle_decision.ServeBattleState
}

// NewPartnerActionServer creates new partner action servers lazily.
type NewPartnerActionServer func() *PartnerActionServer

// StandByNewPartnerActionServer prepares a constructor for PartnerActionServer.
func StandByNewPartnerActionServer(
	random util.EmitRandomFunc,
	serveState battle_decision.ServeBattleState,
) NewPartnerActionServer {
	return func() *PartnerActionServer {
		return &PartnerActionServer{
			random:     random,
			serveState: serveState,
		}
	}
}

// DecidePlan creates and stores a new partner plan.
func (s *PartnerActionServer) DecidePlan() {
	state := s.serveState()
	s.StoredPlan = createDecidePartnerAction(s.random, state)
}

// GetPlan returns the stored partner plan.
func (s *PartnerActionServer) GetPlan() *PartnerActionPlan {
	return s.StoredPlan
}
