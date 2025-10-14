package actor

import "strconv"

type ActorId string

const (
	ActorEmptyId ActorId = "empty"
	ActorLuneId          = ActorId("lune")
	ActorSunnyId         = ActorId("sunny")
)

type ActorSide int

const (
	ActorSidePlayer ActorSide = iota
	ActorSideEnemy
)

func (s ActorSide) Invert() ActorSide {
	if s == ActorSidePlayer {
		return ActorSideEnemy
	}
	return ActorSidePlayer
}

type MaxHP int

func (h MaxHP) ToHP() HP {
	return HP(h)
}

type HP int

func (h HP) String() string {
	return strconv.Itoa(int(h))
}

type HPRatio float64

func (h HP) Ratio(max MaxHP) HPRatio {
	if max <= 0 {
		return 0
	}
	return HPRatio(float64(h) / float64(max))
}

func (r HPRatio) Float64() float64 {
	return float64(r)
}

type ATK int
type MAG int
type DEF int
type SPD int

// Actor is a parameter set for a character in a battle.
type Actor struct {
	Id    ActorId
	MaxHP MaxHP
	HP    HP
	ATK   ATK
	MAG   MAG
	DEF   DEF
	SPD   SPD
	Side  ActorSide
}

func (a *Actor) IsMainActor() bool {
	return a.Id == ActorLuneId
}

func (a *Actor) IsSubActor() bool {
	return a.Id == ActorSunnyId
}

func (a *Actor) IsEnemy() bool {
	return a.Side == ActorSideEnemy
}

func (a *Actor) IsBeaten() bool {
	return a.HP <= 0
}

type ActorSupplier func(ActorId) *Actor
type UpdateActorFunc func(*Actor)

type ActorServer interface {
	Get(ActorId) *Actor
	GetAllActor() []*Actor
	Upsert(*Actor)
	ClearAll()
}

type InMemoryActorServer struct {
	actors map[ActorId]*Actor
}

func NewInMemoryActorServer() *InMemoryActorServer {
	actors := map[ActorId]*Actor{}
	return &InMemoryActorServer{
		actors: actors,
	}
}

func (s *InMemoryActorServer) Get(id ActorId) *Actor {
	return s.actors[id]
}

func (s *InMemoryActorServer) Upsert(actor *Actor) {
	s.actors[actor.Id] = actor
}

func (s *InMemoryActorServer) GetAllActor() []*Actor {
	result := make([]*Actor, 0, len(s.actors))
	for _, actor := range s.actors {
		result = append(result, actor)
	}
	return result
}

func (s *InMemoryActorServer) ClearAll() {
	s.actors = map[ActorId]*Actor{}
}
