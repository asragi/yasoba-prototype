package core

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/actor"
)

type PrepareActorArgs struct {
	MainActorCharacterId CharacterId
	SubActorCharacterId  CharacterId
	EnemyIds             []EnemyId
}

type EnemyIdPair struct {
	EnemyId EnemyId
	ActorId actor.ActorId
}

type PrepareActorResult struct {
	MainActorId actor.ActorId
	SubActorId  actor.ActorId
	EnemyIds    []*EnemyIdPair
}

type PrepareActorService func(*PrepareActorArgs) *PrepareActorResult

type actorInserter interface {
	ClearAll()
	Upsert(actor *actor.Actor)
}

func characterToActor(character *CharacterData, id actor.ActorId) *actor.Actor {
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

func enemyToActor(enemy *EnemyData, id actor.ActorId) *actor.Actor {
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

func CreatePrepareActorService(
	serveCharacter ServeCharacterFunc,
	serveEnemy ServeEnemyData,
	actorServer actorInserter,
) PrepareActorService {
	const MainActorId = actor.ActorLuneId
	const SubActorId = actor.ActorSunnyId
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
			actorId := actor.ActorId(fmt.Sprintf("%s_%d", enemyData.Id, i))
			result[i] = &EnemyIdPair{
				EnemyId: id,
				ActorId: actorId,
			}
			actorServer.Upsert(enemyToActor(enemyData, actorId))
		}

		return &PrepareActorResult{
			MainActorId: MainActorId,
			SubActorId:  SubActorId,
			EnemyIds:    result,
		}
	}
}
