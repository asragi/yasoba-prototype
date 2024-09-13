package core

import (
	"github.com/asragi/yasoba-prototype/util"
)

type PlayerCommand int

const (
	PlayerCommandAttack PlayerCommand = iota
	PlayerCommandFire
	PlayerCommandThunder
	PlayerCommandBarrier
	PlayerCommandWind
	PlayerCommandFocus
	PlayerCommandDefend
)

func (b *PlayerCommand) ToTextId() TextId {
	switch *b {
	case PlayerCommandAttack:
		return "battle_command_attack"
	case PlayerCommandFire:
		return "battle_command_fire"
	case PlayerCommandThunder:
		return "battle_command_thunder"
	case PlayerCommandBarrier:
		return "battle_command_barrier"
	case PlayerCommandWind:
		return "battle_command_wind"
	case PlayerCommandFocus:
		return "battle_command_focus"
	case PlayerCommandDefend:
		return "battle_command_defend"
	}
	return ""
}

type BattlePlayerCommandResult struct {
	SkillApplyArgs *SelectedAction
}

type ProcessPlayerCommandFunc func(*PostCommandRequest) *BattlePlayerCommandResult

func CreateProcessPlayerCommand(supplyActor ActorSupplier) ProcessPlayerCommandFunc {
	isToEnemy := func(targets []ActorId) bool {
		if len(targets) == 0 {
			return false
		}
		target := supplyActor(targets[0])
		return target.Side == ActorSideEnemy
	}
	return func(command *PostCommandRequest) *BattlePlayerCommandResult {
		decidedSkillId := func() SkillId {
			if isToEnemy(command.TargetId) {
				switch command.Command {
				case PlayerCommandAttack:
					return SkillIdLuneAttack
				case PlayerCommandFire:
					return SkillIdLuneFireEnemy
				default:
					panic("not implemented")
				}
			}
			// TODO: implement
			return SkillIdLuneAttack
		}()
		return &BattlePlayerCommandResult{
			SkillApplyArgs: &SelectedAction{
				Id:       decidedSkillId,
				Actor:    command.ActorId,
				SubActor: ActorEmptyId,
				Target:   command.TargetId,
			},
		}
	}
}

func decideActionOrder(actorServer ActorServer) []ActorId {
	var result []ActorId
	mainActorId := ActorLuneId
	actors := actorServer.GetAllActor()
	actorSet := util.NewSet(actors)
	subActor, err := actorSet.Find(func(a *Actor) bool { return a.Id != mainActorId && a.Side == ActorSidePlayer })
	if err == nil {
		result = append(result, subActor.Id)
	}
	result = append(result, mainActorId)
	enemies := actorSet.Filter(func(a *Actor) bool { return a.Side == ActorSideEnemy })
	enemyIds := util.SetSelect(enemies, func(a *Actor) ActorId { return a.Id })
	result = append(result, enemyIds.ToArray()...)
	return result
}

type InitializeBattleRequest struct {
	MainActorCharacterId CharacterId
	SubActorCharacterId  CharacterId
	EnemyIds             []EnemyId
}

type InitializeBattleResponse struct {
	MainActorId ActorId
	SubActorId  ActorId
	EnemyIds    []*EnemyIdPair
}

type InitializeBattleFunc func(*InitializeBattleRequest) *InitializeBattleResponse

func CreateInitializeBattle(
	prepareActorService PrepareActorService,
) func(*InitializeBattleRequest) *InitializeBattleResponse {
	return func(option *InitializeBattleRequest) *InitializeBattleResponse {
		prepareResult := prepareActorService(
			&PrepareActorArgs{
				MainActorCharacterId: option.MainActorCharacterId,
				SubActorCharacterId:  option.SubActorCharacterId,
				EnemyIds:             option.EnemyIds,
			},
		)
		return &InitializeBattleResponse{
			MainActorId: prepareResult.MainActorId,
			SubActorId:  prepareResult.SubActorId,
			EnemyIds:    prepareResult.EnemyIds,
		}
	}
}

type PostCommandRequest struct {
	ActorId  ActorId
	TargetId []ActorId
	Command  PlayerCommand
}

type PostCommandResponse struct {
	Actions []*SelectedAction
}

type PostCommandFunc func(*PostCommandRequest) *PostCommandResponse

func CreatePostCommand(
	processPlayerCommand ProcessPlayerCommandFunc,
) PostCommandFunc {
	return func(command *PostCommandRequest) *PostCommandResponse {
		result := processPlayerCommand(command)
		return &PostCommandResponse{
			Actions: []*SelectedAction{result.SkillApplyArgs},
		}
	}
}
