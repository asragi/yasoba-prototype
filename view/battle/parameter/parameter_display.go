package parameter

import (
	"github.com/asragi/yasoba-prototype/common/character"
	"github.com/asragi/yasoba-prototype/common/character/hero"
	"github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/constant"
	"github.com/asragi/yasoba-prototype/view/battle/hp"
	"github.com/asragi/yasoba-prototype/view/common/window"
)

type mpDisplay interface {
	SetCurrentMP(hero.MP)
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
}

type BattleParameterDisplay struct {
	mpDisplay mpDisplay
	hpDisplay *hp.BattleHPDisplay
	window    window.WindowInterface
}

func (d *BattleParameterDisplay) SetHP(hp character.HP) {
	d.hpDisplay.SetHP(hp)
}

func (d *BattleParameterDisplay) GetHeight() float64 {
	return d.window.Size().Y
}

func (d *BattleParameterDisplay) Update(parentPosition *drawing.Vector) {
	d.window.Update(parentPosition)
	d.hpDisplay.Update(d.window.GetPositionLowerRight())
	if d.mpDisplay == nil {
		// mpDisplay は optional なので nil でも panic しない
		return
	}
	d.mpDisplay.Update(d.window.GetPositionCenter())
}

func (d *BattleParameterDisplay) Draw(drawFunc drawing.DrawFunc) {
	d.window.Draw(drawFunc)
	d.hpDisplay.Draw(drawFunc)
	if d.mpDisplay == nil {
		return
	}
	d.mpDisplay.Draw(drawFunc)
}

// mpDisplay is optional
type NewBattleParameterDisplayFunc func(character.HP, *drawing.Pivot, mpDisplay) *BattleParameterDisplay

func CreateNewBattleParameterDisplay(
	newWindow window.NewWindowFunc,
	newBattleHPDisplay hp.NewBattleHPDisplayFunc,
) NewBattleParameterDisplayFunc {
	const windowCornerSize = 3
	const size = constant.FaceSize + windowCornerSize*2
	height := windowCornerSize*2 + 13.0
	return func(
		initialHp character.HP,
		pivot *drawing.Pivot,
		mpDisplay mpDisplay) *BattleParameterDisplay {
		return &BattleParameterDisplay{
			hpDisplay: newBattleHPDisplay(initialHp),
			mpDisplay: mpDisplay,
			window: newWindow(
				&window.WindowOption{
					Texture:          texture.Window,
					CornerSize:       windowCornerSize,
					RelativePosition: &drawing.Vector{X: 0, Y: 0},
					Size:             &drawing.Vector{X: size, Y: height},
					Depth:            drawing.DepthPlayer,
					Pivot:            pivot,
				},
			),
		}
	}
}
