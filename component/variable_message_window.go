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
}

type VariableMessageWindow struct {
	isActive  bool
	text      *widget.Text
	window    widget.WindowInterface
	serveText core.ServeTextDataFunc
}

func (w *VariableMessageWindow) Draw(drawFunc frontend.DrawFunc) {
	if !w.isActive {
		return
	}
	w.text.Draw(drawFunc)
	w.window.Draw(drawFunc)
}

func (w *VariableMessageWindow) Update(parentPosition *frontend.Vector) {
	w.text.Update(parentPosition)
	w.window.Update(parentPosition)
}

func (w *VariableMessageWindow) SetActive(isActive bool) {
	w.isActive = isActive
}

func (w *VariableMessageWindow) SetText(textId core.TextId) {
	text := w.serveText(textId)
	w.text.SetText(text.Text, true)
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
	return func(
		relativePosition *frontend.Vector,
		depth frontend.Depth,
		pivot *frontend.Pivot,
	) VariableMessageWindowInterface {
		text := newText(
			&widget.TextOptionsNew{
				Font:             font,
				Speed:            speed,
				RelativePosition: relativePosition,
				Depth:            depth,
				Pivot:            pivot,
			},
		)

		window := newWindow(
			&widget.WindowOption{
				Texture:          windowTexture,
				CornerSize:       6,
				RelativePosition: relativePosition,
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
