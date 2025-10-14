package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/util"
)

// DecideActionOrderFunc lists actor IDs in the order they should act.
type DecideActionOrderFunc func() []actor.ActorId

// AllActorServer abstracts actor retrieval for action ordering.
type AllActorServer interface {
	GetAllActor() []*actor.Actor
}

// CreateDecideActionOrder creates a function that determines turn ordering.
func CreateDecideActionOrder(actorServer AllActorServer) DecideActionOrderFunc {
	return func() []actor.ActorId {
		result := make([]actor.ActorId, 0)
		actors := actorServer.GetAllActor()
		actorSet := util.NewSet(actors)
		subActor, err := actorSet.Find(func(a *actor.Actor) bool { return a.IsSubActor() })
		if err == nil {
			result = append(result, subActor.Id)
		}
		enemies := actorSet.Filter(func(a *actor.Actor) bool { return a.Side == actor.ActorSideEnemy })
		enemyIds := util.SetSelect(enemies, func(a *actor.Actor) actor.ActorId { return a.Id })
		result = append(result, enemyIds.ToArray()...)
		return result
	}
}
