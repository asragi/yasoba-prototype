package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type VariableMessageWindowInterface interface {
	widget.PositionUpdater
	widget.Drawer
	SetActive(bool)
	SetText(core.TextId)
}

type VariableMessageWindow struct {
	isActive  bool
	text      widget.TextInterface
	window    widget.WindowInterface
	serveText core.ServeTextDataFunc
}

func (w *VariableMessageWindow) Draw(drawFunc frontend.DrawFunc) {
	if !w.isActive {
		return
	}
	w.window.Draw(drawFunc)
	w.text.Draw(drawFunc)
}

func (w *VariableMessageWindow) Update(parentPosition *frontend.Vector) {
	w.window.Update(parentPosition)
	w.text.Update(w.window.GetPositionUpperLeft().Add(w.window.GetPadding()))
}

func (w *VariableMessageWindow) SetActive(isActive bool) {
	w.isActive = isActive
}

func (w *VariableMessageWindow) SetText(textId core.TextId) {
	text := w.serveText(textId)
	w.text.SetText(text.Text, false)
	padding := w.window.GetPadding().Multiply(2)
	w.window.SetSize(w.text.Size().Add(padding))
}

type NewVariableMessageWindowFunc func(
	*frontend.Vector,
	frontend.Depth,
	*frontend.Pivot,
) VariableMessageWindowInterface

func StandByNewVariableMessageWindow(
	newWindow widget.NewWindowFunc,
	newText widget.NewTextFunc,
	serveTextData core.ServeTextDataFunc,
) NewVariableMessageWindowFunc {
	padding := &frontend.Vector{X: 16, Y: 8}
	windowTexture := frontend.TextureWindow
	font := frontend.MaruMinya
	speed := 5
	margin := 5.0
	return func(
		relativePosition *frontend.Vector,
		depth frontend.Depth,
		pivot *frontend.Pivot,
	) VariableMessageWindowInterface {
		text := newText(
			&widget.TextOptionsNew{
				Font:             font,
				Speed:            speed,
				RelativePosition: frontend.VectorZero,
				Depth:            depth,
				Pivot:            frontend.PivotTopLeft,
			},
		)

		window := newWindow(
			&widget.WindowOption{
				Texture:          windowTexture,
				CornerSize:       6,
				RelativePosition: relativePosition.Add(&frontend.Vector{X: 0, Y: -margin}),
				Size:             frontend.VectorOne,
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
