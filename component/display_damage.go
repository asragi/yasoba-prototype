package component

import (
	"image/color"

	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/util"
	"github.com/asragi/yasoba-prototype/widget"
)

const DamageDisplayDuration = 40
const DamageDisplayPopFrame = 30
const DamageDisplayPopHeight = 10

type DisplayDamage struct {
	text     widget.TextInterface
	popFrame int
}

type NewDisplayDamageFunc func() *DisplayDamage

func CreateNewDisplayDamage(newText widget.NewTextFunc) NewDisplayDamageFunc {
	damageTextColor := color.White
	positionDiff := &frontend.Vector{Y: 33}
	return func() *DisplayDamage {
		text := newText(
			&widget.TextOptionsNew{
				RelativePosition: positionDiff,
				Pivot:            frontend.PivotCenter,
				Font:             frontend.MaruMinya,
				Speed:            4,
				Depth:            frontend.DepthDamageText,
				Color:            damageTextColor,
				EnableOutline:    true,
				Scale:            2,
			},
		)

		return &DisplayDamage{
			text:     text,
			popFrame: DamageDisplayDuration,
		}
	}
}

func (d *DisplayDamage) DisplayDamage(damage core.Damage) {
	d.text.SetText(damage.String(), true)
	d.popFrame = 0
}

func (d *DisplayDamage) Update(parentPosition *frontend.Vector) {
	d.popFrame++
	positionDiff := func() *frontend.Vector {
		if d.popFrame < DamageDisplayPopFrame {
			y := DamageDisplayPopHeight * (1 - util.EaseOutBounce(float64(d.popFrame)/DamageDisplayPopFrame))
			return &frontend.Vector{X: 0, Y: -y}
		}
		return frontend.VectorZero
	}()
	d.text.Update(parentPosition.Add(positionDiff))
}

func (d *DisplayDamage) Draw(drawFunc frontend.DrawFunc) {
	if d.popFrame > DamageDisplayDuration {
		return
	}
	d.text.Draw(drawFunc)
}
