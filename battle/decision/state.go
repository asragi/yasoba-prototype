package decision

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/game/enemy"
)

// ChoiceActionId is issued for each one enemy.
type ChoiceActionId string

// CharacterIdToChoiceActionId converts a character ID to a choice action ID.
func CharacterIdToChoiceActionId(id character.CharacterId) ChoiceActionId {
	return ChoiceActionId(id)
}

// EnemyIdToChoiceActionId converts an enemy ID to a choice action ID.
func EnemyIdToChoiceActionId(id enemy.EnemyId) ChoiceActionId {
	return ChoiceActionId(id)
}

// BattleEndType expresses the termination state of a battle.
type BattleEndType int

const (
	BattleEndTypeNone BattleEndType = iota
	BattleEndTypeWin
	BattleEndTypeLose
)

// BattleState represents the participating actors in a battle.
type BattleState struct {
	Actors []*actor.Actor
}

// IsBattleShouldBeEnd reports whether the battle should end, and why.
func (s *BattleState) IsBattleShouldBeEnd() BattleEndType {
	if s.IsAllBeaten(actor.ActorSidePlayer) {
		return BattleEndTypeLose
	}
	if s.IsAllBeaten(actor.ActorSideEnemy) {
		return BattleEndTypeWin
	}
	return BattleEndTypeNone
}

// IsAllBeaten checks whether all actors on the given side are defeated.
func (s *BattleState) IsAllBeaten(side actor.ActorSide) bool {
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

// GetOtherSideActors retrieves actors on the opposite side.
func (s *BattleState) GetOtherSideActors(actionActor *actor.Actor) []*actor.Actor {
	side := actionActor.Side
	var result []*actor.Actor
	for _, actor := range s.Actors {
		if actor.Side == side {
			continue
		}
		result = append(result, actor)
	}
	return result
}

// GetMainActor returns the primary player actor.
func (s *BattleState) GetMainActor() *actor.Actor {
	for _, actor := range s.Actors {
		if actor.IsMainActor() {
			return actor
		}
	}
	return nil
}

// GetSubActor returns the partner actor.
func (s *BattleState) GetSubActor() *actor.Actor {
	for _, actor := range s.Actors {
		if actor.IsSubActor() {
			return actor
		}
	}
	return nil
}

// ServeBattleState supplies the current battle state snapshot.
type ServeBattleState func() *BattleState

type allActorServer interface {
	GetAllActor() []*actor.Actor
}

// CreateServeBattleState constructs a state supplier backed by an actor server.
func CreateServeBattleState(supplyActor allActorServer) ServeBattleState {
	return func() *BattleState {
		actors := supplyActor.GetAllActor()
		return &BattleState{
			Actors: actors,
		}
	}
}
