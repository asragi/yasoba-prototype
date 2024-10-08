package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
	"image/color"
)

type BattleHPDisplay struct {
	text widget.TextInterface
}

func (d *BattleHPDisplay) Update(parentPosition *frontend.Vector) {
	d.text.Update(parentPosition)
}

func (d *BattleHPDisplay) Draw(drawFunc frontend.DrawFunc) {
	d.text.Draw(drawFunc)
}

func (d *BattleHPDisplay) SetHP(afterHp core.HP) {
	d.text.SetText(afterHp.String(), false)
}

type NewBattleHPDisplayFunc func(core.HP) *BattleHPDisplay

func CreateNewBattleHPDisplay(
	font frontend.FontId,
	newText widget.NewTextFunc,
) NewBattleHPDisplayFunc {
	const margin float64 = 4
	return func(initialHp core.HP) *BattleHPDisplay {
		text := newText(
			&widget.TextOptionsNew{
				RelativePosition: &frontend.Vector{X: -margin, Y: -margin},
				Pivot:            frontend.PivotBottomRight,
				Font:             font,
				Speed:            4,
				Depth:            frontend.DepthDebug,
				Color:            color.White,
				EnableOutline:    true,
				Scale:            1,
				XSpacing:         7,
			},
		)

		text.SetText(initialHp.String(), true)
		return &BattleHPDisplay{
			text: text,
		}
	}
}
