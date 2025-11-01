package actor

import "github.com/asragi/yasoba-prototype/common/character"

// Actor is a parameter set for a character in a battle.
type Actor struct {
	Id    ActorId
	MaxHP character.MaxHP
	HP    character.HP
	ATK   character.ATK
	MAG   character.MAG
	DEF   character.DEF
	SPD   character.SPD
	Side  ActorSide
}

// IsMainActor reports whether the actor is the main player-controlled character.
func (a *Actor) IsMainActor() bool {
	return a.Id == ActorLuneId
}

// IsSubActor reports whether the actor is the partner character.
func (a *Actor) IsSubActor() bool {
	return a.Id == ActorSunnyId
}

// IsEnemy reports whether the actor belongs to the enemy side.
func (a *Actor) IsEnemy() bool {
	return a.Side == ActorSideEnemy
}

// IsBeaten reports whether the actor has been defeated.
func (a *Actor) IsBeaten() bool {
	return a.HP <= 0
}
