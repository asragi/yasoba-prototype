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

type NewBattleHPDisplayFunc func() *BattleHPDisplay

func CreateNewBattleHPDisplay(
	font frontend.FontId,
	resource *frontend.ResourceManager,
) NewBattleHPDisplayFunc {
	return func() *BattleHPDisplay {
		text := widget.NewText(
			&widget.TextOptions{
				RelativePosition: &frontend.Vector{X: 0, Y: 0},
				Pivot:            frontend.PivotBottomRight,
				TextFace:         resource.GetFont(font),
				Speed:            4,
				Depth:            frontend.DepthDamageText,
				Color:            color.White,
				EnableOutline:    true,
				Scale:            1,
				XSpacing:         8,
			},
		)

		return &BattleHPDisplay{
			text: text,
		}
	}
}
