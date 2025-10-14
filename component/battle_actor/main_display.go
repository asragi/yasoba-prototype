package battle_actor

import (
	"github.com/asragi/yasoba-prototype/actor"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/character"
)

type BattleActorDisplay struct {
	faceWindow       *component.FaceWindow
	displayDamage    *component.DisplayDamage
	parameterDisplay *BattleParameterDisplay
}

func (d *BattleActorDisplay) SetDamage(damage battleSkill.Damage, afterHP actor.HP) {
	d.displayDamage.DisplayDamage(damage)
	d.parameterDisplay.hpDisplay.SetHP(afterHP)
}

func (d *BattleActorDisplay) Update(bottomLeftPosition *frontend.Vector) {
	d.faceWindow.Update(bottomLeftPosition)
	d.displayDamage.Update(d.faceWindow.GetCenterPosition())
	d.parameterDisplay.Update(bottomLeftPosition)
}

func (d *BattleActorDisplay) Draw(drawFunc frontend.DrawFunc) {
	d.faceWindow.Draw(drawFunc)
	d.displayDamage.Draw(drawFunc)
	d.parameterDisplay.Draw(drawFunc)
}

func (d *BattleActorDisplay) GetMainCharacterPosition() *frontend.Vector {
	return d.faceWindow.GetCenterPosition()
}

func (d *BattleActorDisplay) GetMainCharacterTopLeftPosition() *frontend.Vector {
	return d.faceWindow.GetTopLeftPosition()
}

func (d *BattleActorDisplay) SetEmotion(emotion component.BattleEmotionType) {
	d.faceWindow.SetEmotion(emotion)
}

type NewBattleActorDisplayFunc func(*actor.Actor) *BattleActorDisplay

func CreateNewBattleActorDisplay(
	newFaceWindow component.NewFaceWindowFunc,
	newDisplayDamage component.NewDisplayDamageFunc,
	newParameterDisplay NewBattleParameterDisplayFunc,
) NewBattleActorDisplayFunc {
	return func(actor *actor.Actor) *BattleActorDisplay {
		initialHp := actor.HP
		parameter := newParameterDisplay(initialHp, frontend.PivotBottomLeft)

		return &BattleActorDisplay{
			faceWindow: newFaceWindow(
				&frontend.Vector{X: 0, Y: -parameter.GetHeight()},
				frontend.DepthPlayer,
				frontend.PivotBottomLeft,
				character.CharacterLuneId,
			),
			displayDamage:    newDisplayDamage(),
			parameterDisplay: parameter,
		}
	}
}
