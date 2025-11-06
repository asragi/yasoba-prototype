package sequence

import "github.com/asragi/yasoba-prototype/battle"

type SwitchToBattleScene func(battle.BattleId)

type SwitchToBattleSceneModel struct {
	battleId battle.BattleId
}

type SwitchToBattleSceneDataPort func(EventID) *SwitchToBattleSceneModel

type createSwitchToBattleSceneEvent func(EventID) *eventUnit

func produceCreateSwitchToBattleSceneEventToUnit(
	switchToBattleSceneDataPort SwitchToBattleSceneDataPort,
	switchToBattle SwitchToBattleScene,
) createSwitchToBattleSceneEvent {
	if switchToBattle == nil {
		panic("sequence: switchToBattle is nil")
	}
	return func(eventID EventID) *eventUnit {
		model := switchToBattleSceneDataPort(eventID)
		if model == nil {
			panic("switch to battle scene model not found: " + string(eventID))
		}
		return &eventUnit{
			start: func() {
				switchToBattle(model.battleId)
			},
			checkIsEnd: func() IsEnd { return true },
			reset:      func() {},
		}
	}
}

func NewSwitchToBattleSceneModel(battleId battle.BattleId) *SwitchToBattleSceneModel {
	return &SwitchToBattleSceneModel{battleId: battleId}
}
