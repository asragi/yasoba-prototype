package component

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type MessageWindow struct {
	text   *widget.Text
	window *widget.Window
}

func (m *MessageWindow) Update(parentPosition *frontend.Vector) {
	m.window.Update(parentPosition)
	m.text.Update(m.window.GetContentUpperLeft())
}

func (m *MessageWindow) Draw(drawFunc frontend.DrawFunc) {
	m.window.Draw(drawFunc)
	m.text.Draw(drawFunc)
}

func NewMessageWindow(
	relativePosition *frontend.Vector,
	size *frontend.Vector,
	depth frontend.Depth,
	pivot *frontend.Pivot,
) *MessageWindow {
	window := widget.NewWindow(
		&widget.WindowOption{
			Image:            nil,
			CornerSize:       6,
			RelativePosition: relativePosition,
			Size:             size,
			Depth:            depth,
			Pivot:            pivot,
			Padding:          &frontend.Vector{X: 16, Y: 8},
		},
	)

	text := widget.NewText(
		"",
		&widget.TextOptions{
			RelativePosition: window.GetContentUpperLeft(),
			Pivot:            frontend.PivotTopLeft,
			TextFace:         nil,
			DisplayAll:       false,
			Speed:            6,
			Depth:            depth,
		},
	)

	return &MessageWindow{
		text:   text,
		window: window,
	}
}
