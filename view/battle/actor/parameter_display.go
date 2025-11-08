package actor

import (
	"github.com/asragi/yasoba-prototype/common/character"
	"github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/constant"
	"github.com/asragi/yasoba-prototype/view/battle/hp"
	widgetwindow "github.com/asragi/yasoba-prototype/view/common/window"
)

type BattleParameterDisplay struct {
	hpDisplay *hp.BattleHPDisplay
	window    widgetwindow.WindowInterface
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
	newWindow widgetwindow.NewWindowFunc,
	newBattleHPDisplay hp.NewBattleHPDisplayFunc,
) NewBattleParameterDisplayFunc {
	const windowCornerSize = 3
	const size = constant.FaceSize + windowCornerSize*2
	height := windowCornerSize*2 + 13.0
	return func(initialHp character.HP, pivot *drawing.Pivot) *BattleParameterDisplay {
		return &BattleParameterDisplay{
			hpDisplay: newBattleHPDisplay(initialHp),
			window: newWindow(
				&widgetwindow.WindowOption{
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
