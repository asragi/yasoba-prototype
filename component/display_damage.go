package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/util"
	"github.com/asragi/yasoba-prototype/widget"
	"image/color"
)

const DamageDisplayDuration = 40
const DamageDisplayPopFrame = 30
const DamageDisplayPopHeight = 10

type DisplayDamage struct {
	text     *widget.Text
	popFrame int
}

type NewDisplayDamageFunc func() *DisplayDamage

func CreateNewDisplayDamage(resource *frontend.ResourceManager) NewDisplayDamageFunc {
	damageTextColor := color.White
	font := resource.GetFont(frontend.MaruMinya)
	return func() *DisplayDamage {
		text := widget.NewText(
			&widget.TextOptions{
				RelativePosition: frontend.VectorZero,
				Pivot:            frontend.PivotCenter,
				TextFace:         font,
				Speed:            4,
				Depth:            frontend.DepthDamageText,
				Color:            damageTextColor,
				EnableOutline:    true,
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
