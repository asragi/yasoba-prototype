package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle/setup"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/game/enemy"
)

// InitializeBattleRequest points to the characters and enemies to prepare.
type InitializeBattleRequest struct {
	MainActorCharacterId character.CharacterId
	SubActorCharacterId  character.CharacterId
	EnemyIds             []enemy.EnemyId
}

// InitializeBattleResponse reports the actor IDs that were prepared.
type InitializeBattleResponse struct {
	MainActorId actor.ActorId
	SubActorId  actor.ActorId
	EnemyIds    []*setup.EnemyIdPair
}

// InitializeBattleFunc initializes battle actors from the provided request.
type InitializeBattleFunc func(*InitializeBattleRequest) *InitializeBattleResponse

// CreateInitializeBattle prepares actors and triggers partner plan selection.
func CreateInitializeBattle(
	prepareActorService setup.Service,
	decidePartnerPlan func(),
) InitializeBattleFunc {
	return func(option *InitializeBattleRequest) *InitializeBattleResponse {
		prepareResult := prepareActorService(
			&setup.PrepareArgs{
				MainActorCharacterId: option.MainActorCharacterId,
				SubActorCharacterId:  option.SubActorCharacterId,
				EnemyIds:             option.EnemyIds,
			},
		)
		decidePartnerPlan()
		return &InitializeBattleResponse{
			MainActorId: prepareResult.MainActorId,
			SubActorId:  prepareResult.SubActorId,
			EnemyIds:    prepareResult.EnemyIds,
		}
	}
}
