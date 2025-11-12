package text

import "github.com/asragi/yasoba-prototype/toolkit/drawing"

type item interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
}

type base = item

type costView = item

type view struct {
	baseItem base
	costView costView
}

func newView(
	baseItem base,
	costView costView,
) viewInterface {
	return &view{
		baseItem: baseItem,
		costView: costView,
	}
}

func (v *view) update(parentPosition *drawing.Vector, width float64) {
	v.baseItem.Update(parentPosition)
	v.costView.Update(parentPosition.Add(&drawing.Vector{X: width, Y: 0}))
}

func (v *view) draw(drawFunc drawing.DrawFunc) {
	v.baseItem.Draw(drawFunc)
	v.costView.Draw(drawFunc)
}
