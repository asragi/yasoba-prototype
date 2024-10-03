package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
	"image/color"
)

type BattleHPDisplay struct {
	text *widget.Text
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
	resource *frontend.ResourceManager,
) NewBattleHPDisplayFunc {
	const margin float64 = 4
	return func(initialHp core.HP) *BattleHPDisplay {
		text := widget.NewText(
			&widget.TextOptions{
				RelativePosition: &frontend.Vector{X: -margin, Y: -margin},
				Pivot:            frontend.PivotBottomRight,
				TextFace:         resource.GetFont(font),
				Speed:            4,
				Depth:            frontend.DepthDebug,
				Color:            color.White,
				EnableOutline:    true,
				Scale:            1,
				XSpacing:         8,
			},
		)

		text.SetText(initialHp.String(), true)

		return &BattleHPDisplay{
			text: text,
		}
	}
}
