package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleActorDisplay struct {
	faceWindow       *FaceWindow
	displayDamage    *DisplayDamage
	parameterDisplay *BattleParameterDisplay
}

func (d *BattleActorDisplay) SetDamage(damage core.Damage) {
	d.displayDamage.DisplayDamage(damage)
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

type NewBattleActorDisplayFunc func(*core.Actor) *BattleActorDisplay

func CreateNewBattleActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage NewDisplayDamageFunc,
	newParameterDisplay NewBattleParameterDisplayFunc,
) NewBattleActorDisplayFunc {
	return func(actor *core.Actor) *BattleActorDisplay {
		initialHp := actor.HP
		parameter := newParameterDisplay(initialHp, frontend.PivotBottomLeft)

		return &BattleActorDisplay{
			faceWindow: newFaceWindow(
				&frontend.Vector{X: 0, Y: -parameter.GetHeight()},
				frontend.DepthPlayer,
				frontend.PivotBottomLeft,
				frontend.TextureFaceLuneNormal,
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

type NewBattleSubActorDisplayFunc func(*core.Actor) *BattleSubActorDisplay

func CreateNewBattleSubActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage NewDisplayDamageFunc,
	newParameterDisplay NewBattleParameterDisplayFunc,
) NewBattleSubActorDisplayFunc {
	return func(actor *core.Actor) *BattleSubActorDisplay {
		parameterDisplay := newParameterDisplay(actor.HP, frontend.PivotBottomRight)
		height := parameterDisplay.GetHeight()
		return &BattleSubActorDisplay{
			faceWindow: newFaceWindow(
				&frontend.Vector{X: 0, Y: -height},
				frontend.DepthPlayer,
				frontend.PivotBottomRight,
				frontend.TextureFaceSunnyNormal,
			),
			displayDamage:    newDisplayDamage(),
			parameterDisplay: parameterDisplay,
			shake:            frontend.NewShake(),
		}
	}
}

func (d *BattleSubActorDisplay) SetDamage(damage core.Damage) {
	d.displayDamage.DisplayDamage(damage)
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

func (d *BattleSubActorDisplay) GetCenterPosition() *frontend.Vector {
	return d.faceWindow.GetCenterPosition()
}

type BattleParameterDisplay struct {
	hpDisplay *BattleHPDisplay
	window    *widget.Window
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

type NewBattleParameterDisplayFunc func(core.HP, *frontend.Pivot) *BattleParameterDisplay

func CreateNewBattleParameterDisplay(
	resource *frontend.ResourceManager,
	newBattleHPDisplay NewBattleHPDisplayFunc,
) NewBattleParameterDisplayFunc {
	img := resource.GetTexture(frontend.TextureWindow)
	const windowCornerSize = 3
	const faceSize = 80
	height := windowCornerSize*2 + 13.0
	return func(initialHp core.HP, pivot *frontend.Pivot) *BattleParameterDisplay {
		return &BattleParameterDisplay{
			hpDisplay: newBattleHPDisplay(initialHp),
			window: widget.NewWindow(
				&widget.WindowOption{
					Image:            img,
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
