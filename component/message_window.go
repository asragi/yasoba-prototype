package component

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type MessageWindow struct {
	text   *widget.Text
	window *widget.Window
}

func (m *MessageWindow) SetText(textString string, displayAll bool) {
	m.text.SetText(textString, displayAll)
}

func (m *MessageWindow) Update(parentPosition *frontend.Vector) {
	m.window.Update(parentPosition)
	m.text.Update(m.window.GetPositionUpperLeft())
}

func (m *MessageWindow) Draw(drawFunc frontend.DrawFunc) {
	m.window.Draw(drawFunc)
	m.text.Draw(drawFunc)
}

type NewMessageWindowFunc func(
	*frontend.Vector,
	*frontend.Vector,
	frontend.Depth,
	*frontend.Pivot,
) *MessageWindow

func StandByNewMessageWindow(resource *frontend.ResourceManager) NewMessageWindowFunc {
	image := resource.GetTexture(frontend.TextureWindow)
	cornerSize := 6
	padding := &frontend.Vector{X: 16, Y: 8}
	font := resource.GetFont(frontend.MaruMinya)
	speed := 6
	return func(
		relativePosition *frontend.Vector,
		size *frontend.Vector,
		depth frontend.Depth,
		pivot *frontend.Pivot,
	) *MessageWindow {
		window := widget.NewWindow(
			&widget.WindowOption{
				Image:            image,
				CornerSize:       cornerSize,
				RelativePosition: relativePosition,
				Size:             size,
				Depth:            depth,
				Pivot:            pivot,
				Padding:          padding,
			},
		)

		text := widget.NewText(
			&widget.TextOptions{
				RelativePosition: window.GetContentUpperLeft(),
				Pivot:            frontend.PivotTopLeft,
				TextFace:         font,
				Speed:            speed,
				Depth:            depth,
			},
		)

		return &MessageWindow{
			text:   text,
			window: window,
		}
	}
}
