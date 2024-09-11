package core

type ActorId string

const (
	ActorLuneId  = ActorId("lune")
	ActorSunnyId = ActorId("sunny")
)

type Actor struct {
	Id          ActorId
	CharacterId CharacterId
	MaxHP       MaxHP
	HP          HP
	ATK         ATK
	MAG         MAG
	DEF         DEF
	SPD         SPD
}

type ServeActorCharacterFunc func(CharacterId) *Actor

func CreateServeActorCharacter(
	serveCharacter ServeCharacterFunc,
) ServeActorCharacterFunc {
	return func(id CharacterId) *Actor {
		character := serveCharacter(id)
		return &Actor{
			CharacterId: id,
			MaxHP:       character.MaxHP,
			HP:          character.HP,
			ATK:         character.ATK,
			MAG:         character.MAG,
			DEF:         character.DEF,
			SPD:         character.SPD,
		}
	}
}

type ActorServer interface {
	Get(ActorId) *Actor
	Upsert(*Actor)
}

type InMemoryActorServer struct {
	actors map[ActorId]*Actor
}

func NewInMemoryActorServer() *InMemoryActorServer {
	actors := map[ActorId]*Actor{
		ActorLuneId: {
			Id:          ActorLuneId,
			CharacterId: CharacterLuneId,
			MaxHP:       40,
			HP:          40,
			ATK:         5,
			MAG:         30,
			DEF:         6,
			SPD:         6,
		},
	}
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
