package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battlesetup"
	"github.com/asragi/yasoba-prototype/characterdata"
	"github.com/asragi/yasoba-prototype/enemydata"
)

// InitializeBattleRequest points to the characters and enemies to prepare.
type InitializeBattleRequest struct {
	MainActorCharacterId characterdata.CharacterId
	SubActorCharacterId  characterdata.CharacterId
	EnemyIds             []enemydata.EnemyId
}

// InitializeBattleResponse reports the actor IDs that were prepared.
type InitializeBattleResponse struct {
	MainActorId actor.ActorId
	SubActorId  actor.ActorId
	EnemyIds    []*battlesetup.EnemyIdPair
}

// InitializeBattleFunc initializes battle actors from the provided request.
type InitializeBattleFunc func(*InitializeBattleRequest) *InitializeBattleResponse

// CreateInitializeBattle prepares actors and triggers partner plan selection.
func CreateInitializeBattle(
	prepareActorService battlesetup.Service,
	decidePartnerPlan func(),
) InitializeBattleFunc {
	return func(option *InitializeBattleRequest) *InitializeBattleResponse {
		prepareResult := prepareActorService(
			&battlesetup.PrepareArgs{
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
