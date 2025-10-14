package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/core"
)

// PostCommandRequest captures the player's chosen command.
type PostCommandRequest struct {
	ActorId  actor.ActorId
	TargetId []actor.ActorId
	Command  PlayerCommand
}

// ProcessBattleRequest is the input for executing a turn.
type ProcessBattleRequest struct {
	TargetId []actor.ActorId
	Command  PlayerCommand
}

// ProcessBattleResponse contains the set of applied skill results.
type ProcessBattleResponse struct {
	SkillApplyResults []*core.SkillApplyResult
}

// ProcessBattleFunc advances the battle by one set of actions.
type ProcessBattleFunc func(*ProcessBattleRequest) *ProcessBattleResponse

// NewProcessBattleFunc constructs a ProcessBattleFunc for a prepared battle.
type NewProcessBattleFunc func(res *InitializeBattleResponse, onBattleEnd func(core.BattleEndType)) ProcessBattleFunc

// StandByCreateProcessBattle wires all dependencies to handle battle progression.
func StandByCreateProcessBattle(
	getActor actor.ActorSupplier,
	getState core.ServeBattleState,
	processPlayerCommand ProcessPlayerCommandFunc,
	getPartnerPlan core.GetPartnerPlanFunc,
	checkCombination core.CheckCombinationFunc,
	skillApply core.SkillApplyFunc,
	decideActionOrder DecideActionOrderFunc,
	newChoiceAction core.NewChoiceActionFunc,
) NewProcessBattleFunc {
	checkBattleShouldEnd := func() (core.BattleEndType, bool) {
		state := getState()
		if state.IsAllBeaten(actor.ActorSidePlayer) {
			return core.BattleEndTypeLose, true
		}
		if state.IsAllBeaten(actor.ActorSideEnemy) {
			return core.BattleEndTypeWin, true
		}
		return core.BattleEndTypeNone, false
	}

	return func(
		initializeBattleResponse *InitializeBattleResponse,
		onBattleEndArg func(core.BattleEndType),
	) ProcessBattleFunc {
		onBattleEnd := func(battleState core.BattleEndType, applyResult []*core.SkillApplyResult) *ProcessBattleResponse {
			onBattleEndArg(battleState)
			return &ProcessBattleResponse{
				SkillApplyResults: applyResult,
			}
		}
		mainActorId := initializeBattleResponse.MainActorId
		subActorId := initializeBattleResponse.SubActorId
		actorIdToEnemy := func() map[actor.ActorId]core.EnemyId {
			result := make(map[actor.ActorId]core.EnemyId)
			for _, pair := range initializeBattleResponse.EnemyIds {
				result[pair.ActorId] = pair.EnemyId
			}
			return result
		}()
		choiceActionList := func() map[actor.ActorId]core.DecideActionFunc {
			result := make(map[actor.ActorId]core.DecideActionFunc)
			for key, value := range actorIdToEnemy {
				result[key] = newChoiceAction(core.EnemyIdToChoiceActionId(value))
			}
			result[subActorId] = newChoiceAction(core.CharacterIdToChoiceActionId(core.CharacterSunnyId))
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
				&core.CheckCombinationRequest{
					MainActorSkillId: selectedAction.Id,
					// TODO: consider multi target
					MainActorTarget: selectedAction.Target[0],
					SubActorSkillId: partnerPlan.SkillId,
					SubActorTarget:  partnerPlan.SelectedTarget,
				},
			)
			resultAction := func() *core.SelectedAction {
				if combinationResult.IsCombination {
					return &core.SelectedAction{
						Id:       combinationResult.SkillId,
						Actor:    selectedAction.Actor,
						SubActor: subActorId,
						Target:   []actor.ActorId{combinationResult.TargetId},
					}
				}
				return selectedAction
			}()

			mainActorApplyResult := skillApply(resultAction)
			result := []*core.SkillApplyResult{mainActorApplyResult}
			if battleState, battleShouldEnd := checkBattleShouldEnd(); battleShouldEnd {
				return onBattleEnd(battleState, result)
			}

			actionOrder := decideActionOrder()
			for _, actorId := range actionOrder {
				actionActor := getActor(actorId)
				if actionActor.IsBeaten() {
					continue
				}
				state := getState()
				decideActionFunction := choiceActionList[actorId]
				decidedAction := decideActionFunction(actionActor, state)
				applyResult := skillApply(
					&core.SelectedAction{
						Id:       decidedAction.SelectedSkill,
						Actor:    actorId,
						SubActor: actor.ActorEmptyId,
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
