package battle

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/battle/combination"
	"github.com/asragi/yasoba-prototype/battle/decision"
	"github.com/asragi/yasoba-prototype/battle/partner"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/game/enemy"
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
	SkillApplyResults []*skill.SkillApplyResult
}

// ProcessBattleFunc advances the battle by one set of actions.
type ProcessBattleFunc func(*ProcessBattleRequest) *ProcessBattleResponse

// NewProcessBattleFunc constructs a ProcessBattleFunc for a prepared battle.
type NewProcessBattleFunc func(res *InitializeBattleResponse, onBattleEnd func(BattleEndType)) ProcessBattleFunc

// StandByCreateProcessBattle wires all dependencies to handle battle progression.
func StandByCreateProcessBattle(
	getActor actor.ActorSupplier,
	getState decision.ServeBattleState,
	processPlayerCommand ProcessPlayerCommandFunc,
	getPartnerPlan partner.GetPartnerPlanFunc,
	checkCombination combination.CheckFunc,
	skillApply skill.SkillApplyFunc,
	decideActionOrder DecideActionOrderFunc,
	newChoiceAction decision.NewChoiceActionFunc,
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
		onBattleEnd := func(battleState BattleEndType, applyResult []*skill.SkillApplyResult) *ProcessBattleResponse {
			onBattleEndArg(battleState)
			return &ProcessBattleResponse{
				SkillApplyResults: applyResult,
			}
		}
		mainActorId := initializeBattleResponse.MainActorId
		subActorId := initializeBattleResponse.SubActorId
		actorIdToEnemy := func() map[actor.ActorId]enemy.EnemyId {
			result := make(map[actor.ActorId]enemy.EnemyId)
			for _, pair := range initializeBattleResponse.EnemyIds {
				result[pair.ActorId] = pair.EnemyId
			}
			return result
		}()
		choiceActionList := func() map[actor.ActorId]decision.DecideActionFunc {
			result := make(map[actor.ActorId]decision.DecideActionFunc)
			for key, value := range actorIdToEnemy {
				result[key] = newChoiceAction(decision.EnemyIdToChoiceActionId(value))
			}
			result[subActorId] = newChoiceAction(decision.CharacterIdToChoiceActionId(character.CharacterSunnyId))
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
				&combination.Request{
					MainActorSkillId: selectedAction.Id,
					// TODO: consider multi target
					MainActorTarget: selectedAction.Target[0],
					SubActorSkillId: partnerPlan.SkillId,
					SubActorTarget:  partnerPlan.SelectedTarget,
				},
			)
			resultAction := func() *skill.SelectedAction {
				if combinationResult.IsCombination {
					return &skill.SelectedAction{
						Id:       combinationResult.SkillId,
						Actor:    selectedAction.Actor,
						SubActor: subActorId,
						Target:   []actor.ActorId{combinationResult.TargetId},
					}
				}
				return selectedAction
			}()

			mainActorApplyResult := skillApply(resultAction)
			result := []*skill.SkillApplyResult{mainActorApplyResult}
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
					&skill.SelectedAction{
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
