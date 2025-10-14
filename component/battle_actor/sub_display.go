package battle_actor

import (
	"github.com/asragi/yasoba-prototype/actor"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/character"
)

type BattleSubActorDisplay struct {
	faceWindow       *component.FaceWindow
	displayDamage    *component.DisplayDamage
	parameterDisplay *BattleParameterDisplay
	shake            *frontend.EmitShake
}

type NewBattleSubActorDisplayFunc func(*actor.Actor) *BattleSubActorDisplay

func CreateNewBattleSubActorDisplay(
	newFaceWindow component.NewFaceWindowFunc,
	newDisplayDamage component.NewDisplayDamageFunc,
	newParameterDisplay NewBattleParameterDisplayFunc,
) NewBattleSubActorDisplayFunc {
	return func(actor *actor.Actor) *BattleSubActorDisplay {
		parameterDisplay := newParameterDisplay(actor.HP, frontend.PivotBottomRight)
		height := parameterDisplay.GetHeight()
		return &BattleSubActorDisplay{
			faceWindow: newFaceWindow(
				&frontend.Vector{X: 0, Y: -height},
				frontend.DepthPlayer,
				frontend.PivotBottomRight,
				character.CharacterSunnyId,
			),
			displayDamage:    newDisplayDamage(),
			parameterDisplay: parameterDisplay,
			shake:            frontend.NewShake(),
		}
	}
}

func (d *BattleSubActorDisplay) SetDamage(damage battleSkill.Damage, afterHP actor.HP) {
	d.displayDamage.DisplayDamage(damage)
	d.parameterDisplay.hpDisplay.SetHP(afterHP)
}

func (d *BattleSubActorDisplay) Shake() {
	d.shake.Shake(
		frontend.ShakeDefaultAmplitude,
		frontend.ShakeDefaultPeriod,
	)
}

func (d *BattleSubActorDisplay) Update(bottomRightPosition *frontend.Vector) {
	d.shake.Update()
	delta := d.shake.Delta()
	d.faceWindow.Update(bottomRightPosition.Add(delta))
	d.parameterDisplay.Update(bottomRightPosition.Add(delta))
	d.displayDamage.Update(d.faceWindow.GetCenterPosition())
}

func (d *BattleSubActorDisplay) Draw(drawFunc frontend.DrawFunc) {
	d.faceWindow.Draw(drawFunc)
	d.parameterDisplay.Draw(drawFunc)
	d.displayDamage.Draw(drawFunc)
}

func (d *BattleSubActorDisplay) GetTopCenterPosition() *frontend.Vector {
	return d.faceWindow.GetTopCenterPosition()
}

func (d *BattleSubActorDisplay) GetCenterPosition() *frontend.Vector {
	return d.faceWindow.GetCenterPosition()
}

func (d *BattleSubActorDisplay) SetEmotion(emotion component.BattleEmotionType) {
	d.faceWindow.SetEmotion(emotion)
}
