package battle

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/combination"
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/battle/decision"
	"github.com/asragi/yasoba-prototype/battle/enemy"
	"github.com/asragi/yasoba-prototype/battle/partner"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/common/character"
	"github.com/asragi/yasoba-prototype/common/character/hero"
)

// PostCommandRequest captures the player's chosen command.
type PostCommandRequest struct {
	ActorId  actor.ActorId
	TargetId []actor.ActorId
	Command  command.Id
}

// ProcessBattleRequest is the input for executing a turn.
type ProcessBattleRequest struct {
	TargetId []actor.ActorId
	Command  command.Id
}

// ProcessBattleResponse contains the set of applied skill results.
type ProcessBattleResponse struct {
	SkillApplyResults []*skill.SkillApplyResult
	IsMainActorBeaten bool
	DidRecoverMp      bool
	CurrentMp         hero.MP
}

// MpStatus reports the current MP and recovers it when needed.
type MpStatus interface {
	Recover()
	CurrentMP() hero.MP
}

// ProcessBattleFunc advances the battle by one set of actions.
type ProcessBattleFunc func(*ProcessBattleRequest) *ProcessBattleResponse

// NewProcessBattleFunc constructs a ProcessBattleFunc for a prepared battle.
type NewProcessBattleFunc func(
	res *InitializeBattleResponse,
	onBattleEnd func(BattleEndType),
	mpStatus MpStatus,
) ProcessBattleFunc

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
		mpStatus MpStatus,
	) ProcessBattleFunc {
		if mpStatus == nil {
			panic("mp status is nil")
		}
		mainActorId := initializeBattleResponse.MainActorId
		subActorId := initializeBattleResponse.SubActorId

		onBattleEnd := func(battleState BattleEndType, applyResult []*skill.SkillApplyResult) *ProcessBattleResponse {
			onBattleEndArg(battleState)
			mainActor := getActor(mainActorId)
			if mainActor == nil {
				panic("main actor not found: " + string(mainActorId))
			}
			isBeaten := mainActor.IsBeaten()
			return &ProcessBattleResponse{
				SkillApplyResults: applyResult,
				IsMainActorBeaten: isBeaten,
				DidRecoverMp:      false,
				CurrentMp:         mpStatus.CurrentMP(),
			}
		}
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
			result := make([]*skill.SkillApplyResult, 0)
			didRecoverMp := false
			mainActor := getActor(mainActorId)
			if mainActor == nil {
				panic("main actor not found: " + string(mainActorId))
			}
			subActor := getActor(subActorId)
			if subActor == nil {
				panic("sub actor not found: " + string(subActorId))
			}
			subActorAlive := !subActor.IsBeaten()

			runPlayerTurn := func(req *ProcessBattleRequest) []*skill.SkillApplyResult {
				actualAction := processPlayerCommand(
					&PostCommandRequest{
						ActorId:  mainActorId,
						TargetId: req.TargetId,
						Command:  req.Command,
					},
				)
				selectedAction := actualAction.SkillApplyArgs
				var partnerPlan *partner.PartnerActionPlan
				if subActorAlive {
					partnerPlan = getPartnerPlan()
				}
				combinationResult := func() *combination.Response {
					if partnerPlan == nil {
						return &combination.Response{
							IsCombination: false,
						}
					}
					return checkCombination(
						&combination.Request{
							MainActorSkillId: selectedAction.Id,
							// TODO: consider multi target
							MainActorTarget: selectedAction.Target[0],
							SubActorSkillId: partnerPlan.SkillId,
							SubActorTarget:  partnerPlan.SelectedTarget,
						},
					)
				}()
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
				return []*skill.SkillApplyResult{skillApply(resultAction)}
			}

			shouldRunPlayer := func(req *ProcessBattleRequest) bool {
				if req == nil {
					return false
				}
				if mainActor.IsBeaten() {
					return false
				}
				return true
			}

			if shouldRunPlayer(request) {
				playerResults := runPlayerTurn(request)
				result = append(result, playerResults...)
				if battleState, battleShouldEnd := checkBattleShouldEnd(); battleShouldEnd {
					return onBattleEnd(battleState, result)
				}
			}

			if battleState, battleShouldEnd := checkBattleShouldEnd(); battleShouldEnd {
				return onBattleEnd(battleState, result)
			}

			actionOrder := decideActionOrder()
			for _, actorId := range actionOrder {
				actionActor := getActor(actorId)
				if actionActor == nil {
					panic("action actor not found: " + string(actorId))
				}
				if actionActor.IsBeaten() {
					continue
				}
				state := getState()
				decideActionFunction, ok := choiceActionList[actorId]
				if !ok || decideActionFunction == nil {
					panic("choice action not found: " + string(actorId))
				}
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

			latestMainActor := getActor(mainActorId)
			if latestMainActor == nil {
				panic("main actor not found: " + string(mainActorId))
			}
			if !latestMainActor.IsBeaten() {
				mpStatus.Recover()
				didRecoverMp = true
			}

			return &ProcessBattleResponse{
				SkillApplyResults: result,
				IsMainActorBeaten: latestMainActor.IsBeaten(),
				DidRecoverMp:      didRecoverMp,
				CurrentMp:         mpStatus.CurrentMP(),
			}
		}
	}
}
