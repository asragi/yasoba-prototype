package actor

// ActorSide indicates whether an actor belongs to the player or enemy.
type ActorSide int

const (
	ActorSidePlayer ActorSide = iota
	ActorSideEnemy
)

// Invert returns the opposite side.
func (s ActorSide) Invert() ActorSide {
	if s == ActorSidePlayer {
		return ActorSideEnemy
	}
	return ActorSidePlayer
}
