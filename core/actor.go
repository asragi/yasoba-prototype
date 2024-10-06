package core

import "fmt"

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

type ServeActorCharacterFunc func(CharacterId) *Actor

func characterToActor(character *CharacterData, id ActorId) *Actor {
	return &Actor{
		Id:    id,
		MaxHP: character.MaxHP,
		HP:    character.HP,
		ATK:   character.ATK,
		MAG:   character.MAG,
		DEF:   character.DEF,
		SPD:   character.SPD,
		Side:  ActorSidePlayer,
	}
}

func enemyToActor(enemy *EnemyData, id ActorId) *Actor {
	return &Actor{
		Id:    id,
		MaxHP: enemy.MaxHP,
		HP:    enemy.MaxHP.ToHP(),
		ATK:   enemy.Atk,
		MAG:   enemy.Mag,
		DEF:   enemy.Def,
		SPD:   enemy.Spd,
		Side:  ActorSideEnemy,
	}
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
	var result []*Actor
	for _, actor := range s.actors {
		result = append(result, actor)
	}
	return result
}

func (s *InMemoryActorServer) ClearAll() {
	s.actors = map[ActorId]*Actor{}
}

type PrepareActorArgs struct {
	MainActorCharacterId CharacterId
	SubActorCharacterId  CharacterId
	EnemyIds             []EnemyId
}

type EnemyIdPair struct {
	EnemyId EnemyId
	ActorId ActorId
}

type PrepareActorResult struct {
	MainActorId ActorId
	SubActorId  ActorId
	EnemyIds    []*EnemyIdPair
}

type PrepareActorService func(*PrepareActorArgs) *PrepareActorResult

type actorInserter interface {
	ClearAll()
	Upsert(actor *Actor)
}

func CreatePrepareActorService(
	serveCharacter ServeCharacterFunc,
	serveEnemy ServeEnemyData,
	actorServer actorInserter,
) PrepareActorService {
	const MainActorId = ActorLuneId
	const SubActorId = ActorSunnyId
	return func(args *PrepareActorArgs) *PrepareActorResult {
		actorServer.ClearAll()
		mainCharacter := serveCharacter(args.MainActorCharacterId)
		mainActor := characterToActor(mainCharacter, MainActorId)
		actorServer.Upsert(mainActor)
		if args.SubActorCharacterId != CharacterEmptyId {
			subCharacter := serveCharacter(args.SubActorCharacterId)
			subActor := characterToActor(subCharacter, SubActorId)
			actorServer.Upsert(subActor)
		}
		result := make([]*EnemyIdPair, len(args.EnemyIds))
		for i, id := range args.EnemyIds {
			enemyData := serveEnemy(id)
			actorId := ActorId(fmt.Sprintf("%s_%d", enemyData.Id, i))
			result[i] = &EnemyIdPair{
				EnemyId: id,
				ActorId: actorId,
			}
			actorServer.Upsert(enemyToActor(serveEnemy(id), actorId))
		}

		return &PrepareActorResult{
			MainActorId: MainActorId,
			SubActorId:  SubActorId,
			EnemyIds:    result,
		}
	}
}
