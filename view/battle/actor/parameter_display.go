package actor

import (
	"github.com/asragi/yasoba-prototype/common/character"
	"github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/hp"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleParameterDisplay struct {
	hpDisplay *hp.BattleHPDisplay
	window    widget.WindowInterface
}

func (d *BattleParameterDisplay) GetHeight() float64 {
	return d.window.Size().Y
}

func (d *BattleParameterDisplay) Update(parentPosition *drawing.Vector) {
	d.window.Update(parentPosition)
	d.hpDisplay.Update(d.window.GetPositionLowerRight())
}

func (d *BattleParameterDisplay) Draw(drawFunc drawing.DrawFunc) {
	d.window.Draw(drawFunc)
	d.hpDisplay.Draw(drawFunc)
}

type NewBattleParameterDisplayFunc func(character.HP, *drawing.Pivot) *BattleParameterDisplay

func CreateNewBattleParameterDisplay(
	newWindow widget.NewWindowFunc,
	newBattleHPDisplay hp.NewBattleHPDisplayFunc,
) NewBattleParameterDisplayFunc {
	const windowCornerSize = 3
	const faceSize = 80
	height := windowCornerSize*2 + 13.0
	return func(initialHp character.HP, pivot *drawing.Pivot) *BattleParameterDisplay {
		return &BattleParameterDisplay{
			hpDisplay: newBattleHPDisplay(initialHp),
			window: newWindow(
				&widget.WindowOption{
					Texture:          texture.Window,
					CornerSize:       windowCornerSize,
					RelativePosition: &drawing.Vector{X: 0, Y: 0},
					Size:             &drawing.Vector{X: faceSize, Y: height},
					Depth:            drawing.DepthPlayer,
					Pivot:            pivot,
				},
			),
		}
	}
}
