package text

import "github.com/asragi/yasoba-prototype/toolkit/drawing"

type viewInterface interface {
	update(*drawing.Vector, float64)
	draw(drawing.DrawFunc)
}

type textLine interface {
	Size() *drawing.Vector
	SetDisable(bool)
	IsDisabled() bool
}

type textOptionWithCost struct {
	width    float64
	textLine textLine
	view     viewInterface
}

func newPresenter(
	width float64,
	textLine textLine,
	view viewInterface,
) *textOptionWithCost {
	return &textOptionWithCost{
		width:    width,
		textLine: textLine,
		view:     view,
	}
}

func (t *textOptionWithCost) Draw(drawFunc drawing.DrawFunc) {
	t.view.draw(drawFunc)
}

func (t *textOptionWithCost) Size() *drawing.Vector {
	height := t.textLine.Size().Y
	return &drawing.Vector{
		X: t.width,
		Y: height,
	}
}

func (t *textOptionWithCost) SetDisable(disable bool) {
	t.textLine.SetDisable(disable)
}

func (t *textOptionWithCost) IsDisabled() bool {
	return t.textLine.IsDisabled()
}
