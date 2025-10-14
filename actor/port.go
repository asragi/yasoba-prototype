package actor

// ActorSupplier provides actor data by ID.
type ActorSupplier func(ActorId) *Actor

// UpdateActorFunc persists updates to an actor.
type UpdateActorFunc func(*Actor)

// ActorServer defines the behavior required to manage actor state.
type ActorServer interface {
	Get(ActorId) *Actor
	GetAllActor() []*Actor
	Upsert(*Actor)
	ClearAll()
}
