package actor

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/common/character"
	"github.com/asragi/yasoba-prototype/common/emotion"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/damage"
	"github.com/asragi/yasoba-prototype/view/common/shake"
)

type BattleSubActorDisplay struct {
	faceWindow       *FaceWindow
	displayDamage    *damage.DisplayDamage
	parameterDisplay *BattleParameterDisplay
	shake            *shake.EmitShake
}

type NewBattleSubActorDisplayFunc func(*actor.Actor) *BattleSubActorDisplay

func CreateNewBattleSubActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage damage.NewDisplayDamageFunc,
	newParameterDisplay NewBattleParameterDisplayFunc,
) NewBattleSubActorDisplayFunc {
	return func(actor *actor.Actor) *BattleSubActorDisplay {
		parameterDisplay := newParameterDisplay(actor.HP, drawing.PivotBottomRight)
		height := parameterDisplay.GetHeight()
		return &BattleSubActorDisplay{
			faceWindow: newFaceWindow(
				&drawing.Vector{X: 0, Y: -height},
				drawing.DepthPlayer,
				drawing.PivotBottomRight,
				character.CharacterSunnyId,
			),
			displayDamage:    newDisplayDamage(),
			parameterDisplay: parameterDisplay,
			shake:            shake.NewShake(),
		}
	}
}

func (d *BattleSubActorDisplay) SetDamage(damage skill.Damage, afterHP character.HP) {
	d.displayDamage.DisplayDamage(damage)
	d.parameterDisplay.hpDisplay.SetHP(afterHP)
}

func (d *BattleSubActorDisplay) Shake() {
	d.shake.Shake(
		shake.ShakeDefaultAmplitude,
		shake.ShakeDefaultPeriod,
	)
}

func (d *BattleSubActorDisplay) Update(bottomRightPosition *drawing.Vector) {
	d.shake.Update()
	delta := d.shake.Delta()
	d.faceWindow.Update(bottomRightPosition.Add(delta))
	d.parameterDisplay.Update(bottomRightPosition.Add(delta))
	d.displayDamage.Update(d.faceWindow.GetCenterPosition())
}

func (d *BattleSubActorDisplay) Draw(drawFunc drawing.DrawFunc) {
	d.faceWindow.Draw(drawFunc)
	d.parameterDisplay.Draw(drawFunc)
	d.displayDamage.Draw(drawFunc)
}

func (d *BattleSubActorDisplay) GetTopCenterPosition() *drawing.Vector {
	return d.faceWindow.GetTopCenterPosition()
}

func (d *BattleSubActorDisplay) GetCenterPosition() *drawing.Vector {
	return d.faceWindow.GetCenterPosition()
}

func (d *BattleSubActorDisplay) SetEmotion(value emotion.EmotionType) {
	d.faceWindow.SetEmotion(value)
}
