package component

import (
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/widget"
)

type VariableMessageWindowInterface interface {
	widget.PositionUpdater
	widget.Drawer
	SetActive(bool)
	SetText(text.TextId)
}

type VariableMessageWindow struct {
	isActive  bool
	text      widget.TextInterface
	window    widget.WindowInterface
	serveText text.ServeTextDataFunc
}

func (w *VariableMessageWindow) Draw(drawFunc drawing.DrawFunc) {
	if !w.isActive {
		return
	}
	w.window.Draw(drawFunc)
	w.text.Draw(drawFunc)
}

func (w *VariableMessageWindow) Update(parentPosition *drawing.Vector) {
	w.window.Update(parentPosition)
	w.text.Update(w.window.GetPositionUpperLeft().Add(w.window.GetPadding()))
}

func (w *VariableMessageWindow) SetActive(isActive bool) {
	w.isActive = isActive
}

func (w *VariableMessageWindow) SetText(textId text.TextId) {
	text := w.serveText(textId)
	w.text.SetText(text.Text.String(), false)
	padding := w.window.GetPadding().Multiply(2)
	w.window.SetSize(w.text.Size().Add(padding))
}

type NewVariableMessageWindowFunc func(
	*drawing.Vector,
	drawing.Depth,
	*drawing.Pivot,
) VariableMessageWindowInterface

func StandByNewVariableMessageWindow(
	newWindow widget.NewWindowFunc,
	newText widget.NewTextFunc,
	serveTextData text.ServeTextDataFunc,
) NewVariableMessageWindowFunc {
	padding := &drawing.Vector{X: 16, Y: 8}
	windowTexture := frontend.TextureWindow
	font := frontend.MaruMinya
	speed := 5
	margin := 5.0
	return func(
		relativePosition *drawing.Vector,
		depth drawing.Depth,
		pivot *drawing.Pivot,
	) VariableMessageWindowInterface {
		text := newText(
			&widget.TextOptionsNew{
				Font:             font,
				Speed:            speed,
				RelativePosition: drawing.VectorZero,
				Depth:            depth,
				Pivot:            drawing.PivotTopLeft,
			},
		)

		window := newWindow(
			&widget.WindowOption{
				Texture:          windowTexture,
				CornerSize:       6,
				RelativePosition: relativePosition.Add(&drawing.Vector{X: 0, Y: -margin}),
				Size:             drawing.VectorOne,
				Depth:            depth,
				Pivot:            pivot,
				Padding:          padding,
			},
		)

		return &VariableMessageWindow{
			text:      text,
			window:    window,
			serveText: serveTextData,
		}
	}
}
