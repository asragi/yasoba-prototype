package actor

// ActorId uniquely identifies an actor within a battle.
type ActorId string

const (
	ActorEmptyId ActorId = "empty"
	ActorLuneId          = ActorId("lune")
	ActorSunnyId         = ActorId("sunny")
)
