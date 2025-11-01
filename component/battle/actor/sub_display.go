package actor

import (
	"github.com/asragi/yasoba-prototype/actor"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/component"
	battleemotion "github.com/asragi/yasoba-prototype/component/battle/emotion"
	componentshake "github.com/asragi/yasoba-prototype/component/shake"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/game/character"
)

type BattleSubActorDisplay struct {
	faceWindow       *FaceWindow
	displayDamage    *component.DisplayDamage
	parameterDisplay *BattleParameterDisplay
	shake            *componentshake.EmitShake
}

type NewBattleSubActorDisplayFunc func(*actor.Actor) *BattleSubActorDisplay

func CreateNewBattleSubActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage component.NewDisplayDamageFunc,
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
			shake:            componentshake.NewShake(),
		}
	}
}

func (d *BattleSubActorDisplay) SetDamage(damage battleSkill.Damage, afterHP actor.HP) {
	d.displayDamage.DisplayDamage(damage)
	d.parameterDisplay.hpDisplay.SetHP(afterHP)
}

func (d *BattleSubActorDisplay) Shake() {
	d.shake.Shake(
		componentshake.ShakeDefaultAmplitude,
		componentshake.ShakeDefaultPeriod,
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

func (d *BattleSubActorDisplay) SetEmotion(emotion battleemotion.BattleEmotionType) {
	d.faceWindow.SetEmotion(emotion)
}
