package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/core"
)

// InitializeBattleRequest points to the characters and enemies to prepare.
type InitializeBattleRequest struct {
	MainActorCharacterId core.CharacterId
	SubActorCharacterId  core.CharacterId
	EnemyIds             []core.EnemyId
}

// InitializeBattleResponse reports the actor IDs that were prepared.
type InitializeBattleResponse struct {
	MainActorId actor.ActorId
	SubActorId  actor.ActorId
	EnemyIds    []*core.EnemyIdPair
}

// InitializeBattleFunc initializes battle actors from the provided request.
type InitializeBattleFunc func(*InitializeBattleRequest) *InitializeBattleResponse

// CreateInitializeBattle prepares actors and triggers partner plan selection.
func CreateInitializeBattle(
	prepareActorService core.PrepareActorService,
	decidePartnerPlan func(),
) InitializeBattleFunc {
	return func(option *InitializeBattleRequest) *InitializeBattleResponse {
		prepareResult := prepareActorService(
			&core.PrepareActorArgs{
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
