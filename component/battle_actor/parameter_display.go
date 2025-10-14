package battle_actor

import (
	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleParameterDisplay struct {
	hpDisplay *component.BattleHPDisplay
	window    widget.WindowInterface
}

func (d *BattleParameterDisplay) GetHeight() float64 {
	return d.window.Size().Y
}

func (d *BattleParameterDisplay) Update(parentPosition *frontend.Vector) {
	d.window.Update(parentPosition)
	d.hpDisplay.Update(d.window.GetPositionLowerRight())
}

func (d *BattleParameterDisplay) Draw(drawFunc frontend.DrawFunc) {
	d.window.Draw(drawFunc)
	d.hpDisplay.Draw(drawFunc)
}

type NewBattleParameterDisplayFunc func(actor.HP, *frontend.Pivot) *BattleParameterDisplay

func CreateNewBattleParameterDisplay(
	newWindow widget.NewWindowFunc,
	newBattleHPDisplay component.NewBattleHPDisplayFunc,
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
