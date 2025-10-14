package component

import (
	"github.com/asragi/yasoba-prototype/actor"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleActorDisplay struct {
	faceWindow       *FaceWindow
	displayDamage    *DisplayDamage
	parameterDisplay *BattleParameterDisplay
}

func (d *BattleActorDisplay) SetDamage(damage battleSkill.Damage, afterHP actor.HP) {
	d.displayDamage.DisplayDamage(damage)
	d.parameterDisplay.hpDisplay.SetHP(afterHP)
}

func (d *BattleActorDisplay) Update(
	bottomLeftPosition *frontend.Vector,
) {
	d.faceWindow.Update(bottomLeftPosition)
	d.displayDamage.Update(d.faceWindow.GetCenterPosition())
	d.parameterDisplay.Update(bottomLeftPosition)
}

func (d *BattleActorDisplay) Draw(
	drawFunc frontend.DrawFunc,
) {
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

func (d *BattleActorDisplay) SetEmotion(emotion BattleEmotionType) {
	d.faceWindow.SetEmotion(emotion)
}

type NewBattleActorDisplayFunc func(*actor.Actor) *BattleActorDisplay

func CreateNewBattleActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage NewDisplayDamageFunc,
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

type BattleSubActorDisplay struct {
	faceWindow       *FaceWindow
	displayDamage    *DisplayDamage
	parameterDisplay *BattleParameterDisplay
	shake            *frontend.EmitShake
}

type NewBattleSubActorDisplayFunc func(*actor.Actor) *BattleSubActorDisplay

func CreateNewBattleSubActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage NewDisplayDamageFunc,
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

func (d *BattleSubActorDisplay) Update(
	bottomRightPosition *frontend.Vector,
) {
	d.shake.Update()
	delta := d.shake.Delta()
	d.faceWindow.Update(bottomRightPosition.Add(delta))
	d.parameterDisplay.Update(bottomRightPosition.Add(delta))
	d.displayDamage.Update(d.faceWindow.GetCenterPosition())
}

func (d *BattleSubActorDisplay) Draw(
	drawFunc frontend.DrawFunc,
) {
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

func (d *BattleSubActorDisplay) SetEmotion(emotion BattleEmotionType) {
	d.faceWindow.SetEmotion(emotion)
}

type BattleParameterDisplay struct {
	hpDisplay *BattleHPDisplay
	window    widget.WindowInterface
}

func (d *BattleParameterDisplay) GetHeight() float64 {
	return d.window.Size().Y
}

func (d *BattleParameterDisplay) Update(
	parentPosition *frontend.Vector,
) {
	d.window.Update(parentPosition)
	d.hpDisplay.Update(d.window.GetPositionLowerRight())
}

func (d *BattleParameterDisplay) Draw(
	drawFunc frontend.DrawFunc,
) {
	d.window.Draw(drawFunc)
	d.hpDisplay.Draw(drawFunc)
}

type NewBattleParameterDisplayFunc func(actor.HP, *frontend.Pivot) *BattleParameterDisplay

func CreateNewBattleParameterDisplay(
	newWindow widget.NewWindowFunc,
	newBattleHPDisplay NewBattleHPDisplayFunc,
) NewBattleParameterDisplayFunc {
	const windowCornerSize = 3
	const faceSize = 80
	height := windowCornerSize*2 + 13.0
	return func(initialHp actor.HP, pivot *frontend.Pivot) *BattleParameterDisplay {
		return &BattleParameterDisplay{
			hpDisplay: newBattleHPDisplay(initialHp),
			window: newWindow(
				&widget.WindowOption{
					Texture:          frontend.TextureWindow,
					CornerSize:       windowCornerSize,
					RelativePosition: &frontend.Vector{X: 0, Y: 0},
					Size:             &frontend.Vector{X: faceSize, Y: height},
					Depth:            frontend.DepthPlayer,
					Pivot:            pivot,
				},
			),
		}
	}
}
