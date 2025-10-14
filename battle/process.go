package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle_combination"
	"github.com/asragi/yasoba-prototype/battle_decision"
	"github.com/asragi/yasoba-prototype/battle_partner"
	"github.com/asragi/yasoba-prototype/battle_skill"
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
	SkillApplyResults []*battle_skill.SkillApplyResult
}

// ProcessBattleFunc advances the battle by one set of actions.
type ProcessBattleFunc func(*ProcessBattleRequest) *ProcessBattleResponse

// NewProcessBattleFunc constructs a ProcessBattleFunc for a prepared battle.
type NewProcessBattleFunc func(res *InitializeBattleResponse, onBattleEnd func(BattleEndType)) ProcessBattleFunc

// StandByCreateProcessBattle wires all dependencies to handle battle progression.
func StandByCreateProcessBattle(
	getActor actor.ActorSupplier,
	getState battle_decision.ServeBattleState,
	processPlayerCommand ProcessPlayerCommandFunc,
	getPartnerPlan battle_partner.GetPartnerPlanFunc,
	checkCombination battle_combination.CheckFunc,
	skillApply battle_skill.SkillApplyFunc,
	decideActionOrder DecideActionOrderFunc,
	newChoiceAction battle_decision.NewChoiceActionFunc,
) NewProcessBattleFunc {
	checkBattleShouldEnd := func() (BattleEndType, bool) {
		state := getState()
		if state.IsAllBeaten(actor.ActorSidePlayer) {
			return BattleEndTypeLose, true
		}
		if state.IsAllBeaten(actor.ActorSideEnemy) {
			return BattleEndTypeWin, true
		}
		return BattleEndTypeNone, false
	}

	return func(
		initializeBattleResponse *InitializeBattleResponse,
		onBattleEndArg func(BattleEndType),
	) ProcessBattleFunc {
		onBattleEnd := func(battleState BattleEndType, applyResult []*battle_skill.SkillApplyResult) *ProcessBattleResponse {
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
		choiceActionList := func() map[actor.ActorId]battle_decision.DecideActionFunc {
			result := make(map[actor.ActorId]battle_decision.DecideActionFunc)
			for key, value := range actorIdToEnemy {
				result[key] = newChoiceAction(battle_decision.EnemyIdToChoiceActionId(value))
			}
			result[subActorId] = newChoiceAction(battle_decision.CharacterIdToChoiceActionId(core.CharacterSunnyId))
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
				&battle_combination.Request{
					MainActorSkillId: selectedAction.Id,
					// TODO: consider multi target
					MainActorTarget: selectedAction.Target[0],
					SubActorSkillId: partnerPlan.SkillId,
					SubActorTarget:  partnerPlan.SelectedTarget,
				},
			)
			resultAction := func() *battle_skill.SelectedAction {
				if combinationResult.IsCombination {
					return &battle_skill.SelectedAction{
						Id:       combinationResult.SkillId,
						Actor:    selectedAction.Actor,
						SubActor: subActorId,
						Target:   []actor.ActorId{combinationResult.TargetId},
					}
				}
				return selectedAction
			}()

			mainActorApplyResult := skillApply(resultAction)
			result := []*battle_skill.SkillApplyResult{mainActorApplyResult}
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
					&battle_skill.SelectedAction{
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
