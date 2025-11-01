package actor

// InMemoryActorServer is a simple in-memory implementation of ActorServer.
type InMemoryActorServer struct {
	actors map[ActorId]*Actor
}

// NewInMemoryActorServer creates a new in-memory actor server.
func NewInMemoryActorServer() *InMemoryActorServer {
	return &InMemoryActorServer{
		actors: map[ActorId]*Actor{},
	}
}

// Get returns the actor with the specified ID.
func (s *InMemoryActorServer) Get(id ActorId) *Actor {
	return s.actors[id]
}

// Upsert adds or replaces the provided actor.
func (s *InMemoryActorServer) Upsert(actor *Actor) {
	s.actors[actor.Id] = actor
}

// GetAllActor returns every actor currently stored.
func (s *InMemoryActorServer) GetAllActor() []*Actor {
	result := make([]*Actor, 0, len(s.actors))
	for _, actor := range s.actors {
		result = append(result, actor)
	}
	return result
}

// ClearAll removes all stored actors.
func (s *InMemoryActorServer) ClearAll() {
	s.actors = map[ActorId]*Actor{}
}
