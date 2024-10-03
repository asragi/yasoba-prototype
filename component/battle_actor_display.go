package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
)

type BattleActorDisplay struct {
	faceWindow    *FaceWindow
	displayDamage *DisplayDamage
	hpDisplay     *BattleHPDisplay
}

func (d *BattleActorDisplay) SetDamage(damage core.Damage) {
	d.displayDamage.DisplayDamage(damage)
}

func (d *BattleActorDisplay) Update(
	bottomLeftPosition *frontend.Vector,
) {
	d.faceWindow.Update(bottomLeftPosition)
	d.displayDamage.Update(d.faceWindow.GetCenterPosition())
	d.hpDisplay.Update(d.faceWindow.GetTopLeftPosition())
}

func (d *BattleActorDisplay) Draw(
	drawFunc frontend.DrawFunc,
) {
	d.faceWindow.Draw(drawFunc)
	d.displayDamage.Draw(drawFunc)
	d.hpDisplay.Draw(drawFunc)
}

func (d *BattleActorDisplay) GetMainCharacterPosition() *frontend.Vector {
	return d.faceWindow.GetCenterPosition()
}

func (d *BattleActorDisplay) GetMainCharacterTopLeftPosition() *frontend.Vector {
	return d.faceWindow.GetTopLeftPosition()
}

type NewBattleActorDisplayFunc func() *BattleActorDisplay

func CreateNewBattleActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage NewDisplayDamageFunc,
	newBattleHPDisplay NewBattleHPDisplayFunc,
) NewBattleActorDisplayFunc {
	return func() *BattleActorDisplay {
		return &BattleActorDisplay{
			faceWindow: newFaceWindow(
				&frontend.Vector{X: 0, Y: 0},
				frontend.DepthPlayer,
				frontend.PivotBottomLeft,
				frontend.TextureFaceLuneNormal,
			),
			displayDamage: newDisplayDamage(),
			hpDisplay:     newBattleHPDisplay(),
		}
	}
}

type BattleSubActorDisplay struct {
	faceWindow    *FaceWindow
	displayDamage *DisplayDamage
	shake         *frontend.EmitShake
}

type NewBattleSubActorDisplayFunc func() *BattleSubActorDisplay

func CreateNewBattleSubActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage NewDisplayDamageFunc,
) NewBattleSubActorDisplayFunc {
	return func() *BattleSubActorDisplay {
		return &BattleSubActorDisplay{
			faceWindow: newFaceWindow(
				&frontend.Vector{X: 0, Y: 0},
				frontend.DepthPlayer,
				frontend.PivotBottomRight,
				frontend.TextureFaceSunnyNormal,
			),
			displayDamage: newDisplayDamage(),
			shake:         frontend.NewShake(),
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
	d.faceWindow.Update(bottomRightPosition.Add(d.shake.Delta()))
	d.displayDamage.Update(d.faceWindow.GetCenterPosition())
}

func (d *BattleSubActorDisplay) Draw(
	drawFunc frontend.DrawFunc,
) {
	d.faceWindow.Draw(drawFunc)
	d.displayDamage.Draw(drawFunc)
}

func (d *BattleSubActorDisplay) GetCenterPosition() *frontend.Vector {
	return d.faceWindow.GetCenterPosition()
}
