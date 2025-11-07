package damage

import (
	"image/color"

	"github.com/asragi/yasoba-prototype/adapter/ebiten/frontend"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
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
	positionDiff := &drawing.Vector{Y: 33}
	return func() *DisplayDamage {
		text := newText(
			&widget.TextOptionsNew{
				RelativePosition: positionDiff,
				Pivot:            drawing.PivotCenter,
				Font:             frontend.MaruMinya,
				Speed:            4,
				Depth:            drawing.DepthDamageText,
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

func (d *DisplayDamage) DisplayDamage(damage battleSkill.Damage) {
	d.text.SetText(damage.String(), true)
	d.popFrame = 0
}

func (d *DisplayDamage) Update(parentPosition *drawing.Vector) {
	d.popFrame++
	positionDiff := func() *drawing.Vector {
		if d.popFrame < DamageDisplayPopFrame {
			y := DamageDisplayPopHeight * (1 - util.EaseOutBounce(float64(d.popFrame)/DamageDisplayPopFrame))
			return &drawing.Vector{X: 0, Y: -y}
		}
		return drawing.VectorZero
	}()
	d.text.Update(parentPosition.Add(positionDiff))
}

func (d *DisplayDamage) Draw(drawFunc drawing.DrawFunc) {
	if d.popFrame > DamageDisplayDuration {
		return
	}
	d.text.Draw(drawFunc)
}
