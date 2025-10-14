package setup

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/game/enemy"
)

type PrepareArgs struct {
	MainActorCharacterId character.CharacterId
	SubActorCharacterId  character.CharacterId
	EnemyIds             []enemy.EnemyId
}

type EnemyIdPair struct {
	EnemyId enemy.EnemyId
	ActorId actor.ActorId
}

type PrepareResult struct {
	MainActorId actor.ActorId
	SubActorId  actor.ActorId
	EnemyIds    []*EnemyIdPair
}

type Service func(*PrepareArgs) *PrepareResult

type actorInserter interface {
	ClearAll()
	Upsert(*actor.Actor)
}

func characterToActor(character *character.CharacterData, id actor.ActorId) *actor.Actor {
	return &actor.Actor{
		Id:    id,
		MaxHP: character.MaxHP,
		HP:    character.HP,
		ATK:   character.ATK,
		MAG:   character.MAG,
		DEF:   character.DEF,
		SPD:   character.SPD,
		Side:  actor.ActorSidePlayer,
	}
}

func enemyToActor(enemy *enemy.EnemyData, id actor.ActorId) *actor.Actor {
	return &actor.Actor{
		Id:    id,
		MaxHP: enemy.MaxHP,
		HP:    enemy.MaxHP.ToHP(),
		ATK:   enemy.Atk,
		MAG:   enemy.Mag,
		DEF:   enemy.Def,
		SPD:   enemy.Spd,
		Side:  actor.ActorSideEnemy,
	}
}

func NewPrepareService(
	serveCharacter character.ServeCharacterFunc,
	serveEnemy enemy.ServeEnemyData,
	actorServer actorInserter,
) Service {
	const mainActorId = actor.ActorLuneId
	const subActorId = actor.ActorSunnyId
	return func(args *PrepareArgs) *PrepareResult {
		actorServer.ClearAll()
		mainCharacter := serveCharacter(args.MainActorCharacterId)
		mainActor := characterToActor(mainCharacter, mainActorId)
		actorServer.Upsert(mainActor)
		if args.SubActorCharacterId != character.CharacterEmptyId {
			subCharacter := serveCharacter(args.SubActorCharacterId)
			subActor := characterToActor(subCharacter, subActorId)
			actorServer.Upsert(subActor)
		}
		result := make([]*EnemyIdPair, len(args.EnemyIds))
		for i, id := range args.EnemyIds {
			enemyData := serveEnemy(id)
			actorId := actor.ActorId(fmt.Sprintf("%s_%d", enemyData.Id, i))
			result[i] = &EnemyIdPair{
				EnemyId: id,
				ActorId: actorId,
			}
			actorServer.Upsert(enemyToActor(enemyData, actorId))
		}

		return &PrepareResult{
			MainActorId: mainActorId,
			SubActorId:  subActorId,
			EnemyIds:    result,
		}
	}
}
