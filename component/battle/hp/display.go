package hp

import (
	"image/color"

	"github.com/asragi/yasoba-prototype/actor"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleHPDisplay struct {
	text widget.TextInterface
}

func (d *BattleHPDisplay) Update(parentPosition *drawing.Vector) {
	d.text.Update(parentPosition)
}

func (d *BattleHPDisplay) Draw(drawFunc drawing.DrawFunc) {
	d.text.Draw(drawFunc)
}

func (d *BattleHPDisplay) SetHP(afterHp actor.HP) {
	d.text.SetText(afterHp.String(), false)
}

type NewBattleHPDisplayFunc func(actor.HP) *BattleHPDisplay

func CreateNewBattleHPDisplay(
	font frontend.FontId,
	newText widget.NewTextFunc,
) NewBattleHPDisplayFunc {
	const margin float64 = 4
	return func(initialHp actor.HP) *BattleHPDisplay {
		text := newText(
			&widget.TextOptionsNew{
				RelativePosition: &drawing.Vector{X: -margin, Y: -margin},
				Pivot:            drawing.PivotBottomRight,
				Font:             font,
				Speed:            4,
				Depth:            drawing.DepthDebug,
				Color:            color.White,
				EnableOutline:    true,
				Scale:            1,
			},
		)

		text.SetText(initialHp.String(), true)
		return &BattleHPDisplay{
			text: text,
		}
	}
}
