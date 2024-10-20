package core

import (
	"github.com/asragi/yasoba-prototype/util"
)

// PlayerCommand is a command that the player can select in the battle.
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

type DecideActionOrderFunc func() []ActorId

type AllActorServer interface {
	GetAllActor() []*Actor
}

func CreateDecideActionOrder(actorServer AllActorServer) DecideActionOrderFunc {
	return func() []ActorId {
		var result []ActorId
		actors := actorServer.GetAllActor()
		actorSet := util.NewSet(actors)
		subActor, err := actorSet.Find(func(a *Actor) bool { return a.IsSubActor() })
		if err == nil {
			result = append(result, subActor.Id)
		}
		enemies := actorSet.Filter(func(a *Actor) bool { return a.Side == ActorSideEnemy })
		enemyIds := util.SetSelect(enemies, func(a *Actor) ActorId { return a.Id })
		result = append(result, enemyIds.ToArray()...)
		return result
	}
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
	decidePartnerPlan decidePartnerPlanFunc,
) func(*InitializeBattleRequest) *InitializeBattleResponse {
	return func(option *InitializeBattleRequest) *InitializeBattleResponse {
		prepareResult := prepareActorService(
			&PrepareActorArgs{
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

type PostCommandRequest struct {
	ActorId  ActorId
	TargetId []ActorId
	Command  PlayerCommand
}

type ProcessBattleRequest struct {
	TargetId []ActorId
	Command  PlayerCommand
}

type ProcessBattleResponse struct {
	SkillApplyResults []*SkillApplyResult
}

type ProcessBattleFunc func(*ProcessBattleRequest) *ProcessBattleResponse
type NewProcessBattleFunc func(res *InitializeBattleResponse, onBattleEnd func(BattleEndType)) ProcessBattleFunc

func StandByCreateProcessBattle(
	getActor ActorSupplier,
	getState ServeBattleState,
	processPlayerCommand ProcessPlayerCommandFunc,
	getPartnerPlan GetPartnerPlanFunc,
	checkCombination CheckCombinationFunc,
	skillApply SkillApplyFunc,
	decideActionOrder DecideActionOrderFunc,
	newChoiceAction NewChoiceActionFunc,
) NewProcessBattleFunc {
	checkBattleShouldEnd := func() (b BattleEndType, shouldEnd bool) {
		state := getState()
		if state.IsAllBeaten(ActorSidePlayer) {
			return BattleEndTypeLose, true
		}
		if state.IsAllBeaten(ActorSideEnemy) {
			return BattleEndTypeWin, true
		}
		return BattleEndTypeNone, false
	}

	return func(
		initializeBattleResponse *InitializeBattleResponse,
		onBattleEndArg func(BattleEndType),
	) ProcessBattleFunc {
		onBattleEnd := func(battleState BattleEndType, applyResult []*SkillApplyResult) *ProcessBattleResponse {
			onBattleEndArg(battleState)
			return &ProcessBattleResponse{
				SkillApplyResults: applyResult,
			}
		}
		mainActorId := initializeBattleResponse.MainActorId
		subActorId := initializeBattleResponse.SubActorId
		actorIdToEnemy := func() map[ActorId]EnemyId {
			result := make(map[ActorId]EnemyId)
			for _, pair := range initializeBattleResponse.EnemyIds {
				result[pair.ActorId] = pair.EnemyId
			}
			return result
		}()
		choiceActionList := func() map[ActorId]DecideActionFunc {
			result := make(map[ActorId]DecideActionFunc)
			for key, value := range actorIdToEnemy {
				result[key] = newChoiceAction(EnemyIdToChoiceActionId(value))
			}
			result[subActorId] = newChoiceAction(CharacterIdToChoiceActionId(CharacterSunnyId))
			return result
		}()

		return func(request *ProcessBattleRequest) *ProcessBattleResponse {
			actualAction := processPlayerCommand(
				&PostCommandRequest{
					ActorId:  mainActorId,
					TargetId: request.TargetId,
					Command:  request.Command,
				},
			)
			selectedAction := actualAction.SkillApplyArgs
			partnerPlan := getPartnerPlan()
			combinationResult := checkCombination(
				&CheckCombinationRequest{
					MainActorSkillId: selectedAction.Id,
					// TODO: consider multi target
					MainActorTarget: selectedAction.Target[0],
					SubActorSkillId: partnerPlan.SkillId,
					SubActorTarget:  partnerPlan.SelectedTarget,
				},
			)
			resultAction := func() *SelectedAction {
				if combinationResult.IsCombination {
					return &SelectedAction{
						Id:       combinationResult.SkillId,
						Actor:    selectedAction.Actor,
						SubActor: subActorId,
						Target:   []ActorId{combinationResult.TargetId},
					}
				}
				return selectedAction
			}()

			mainActorApplyResult := skillApply(resultAction)
			result := []*SkillApplyResult{mainActorApplyResult}
			if battleState, battleShouldEnd := checkBattleShouldEnd(); battleShouldEnd {
				return onBattleEnd(battleState, result)
			}

			actionOrder := decideActionOrder()
			for _, actorId := range actionOrder {
				actor := getActor(actorId)
				if actor.IsBeaten() {
					continue
				}
				state := getState()
				decideActionFunction := choiceActionList[actorId]
				decidedAction := decideActionFunction(actor, state)
				applyResult := skillApply(
					&SelectedAction{
						Id:       decidedAction.SelectedSkill,
						Actor:    actorId,
						SubActor: ActorEmptyId,
						Target:   decidedAction.TargetActorIds,
					},
				)
				result = append(result, applyResult)

				if battleState, battleShouldEnd := checkBattleShouldEnd(); battleShouldEnd {
					return onBattleEnd(battleState, result)
				}
			}

			return &ProcessBattleResponse{
				SkillApplyResults: result,
			}
		}
	}
}
