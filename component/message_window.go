package component

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type MessageWindow struct {
	text   widget.TextInterface
	window widget.WindowInterface
	shake  *frontend.EmitShake
	isOpen bool
}

func (m *MessageWindow) Shake(amplitude float64, period int) {
	m.shake.Shake(amplitude, period)
}

func (m *MessageWindow) FitToMessage() {
	size := m.text.Size().Add(m.window.GetPadding().Multiply(2))
	m.window.SetSize(size)
}

func (m *MessageWindow) Open() {
	m.isOpen = true
}

func (m *MessageWindow) Close() {
	m.isOpen = false
}

func (m *MessageWindow) IsTextEnd() bool {
	if !m.isOpen {
		return true
	}
	return m.text.CheckIsEnd()
}

func (m *MessageWindow) SetText(textString string, displayAll bool) {
	m.text.SetText(textString, displayAll)
}

func (m *MessageWindow) Update(parentPosition *frontend.Vector) {
	m.shake.Update()
	m.window.Update(parentPosition)
	m.text.Update(m.window.GetPositionUpperLeft())
}

func (m *MessageWindow) Draw(drawFunc frontend.DrawFunc) {
	if !m.isOpen {
		return
	}
	m.window.Draw(drawFunc)
	m.text.Draw(drawFunc)
}

type NewMessageWindowFunc func(
	relativePosition *frontend.Vector,
	size *frontend.Vector,
	depth frontend.Depth,
	pivot *frontend.Pivot,
) *MessageWindow

func StandByNewMessageWindow(
	newText widget.NewTextFunc,
	newWindow widget.NewWindowFunc,
) NewMessageWindowFunc {
	cornerSize := 6
	padding := &frontend.Vector{X: 16, Y: 8}
	speed := 5
	return func(
		relativePosition *frontend.Vector,
		size *frontend.Vector,
		depth frontend.Depth,
		pivot *frontend.Pivot,
	) *MessageWindow {
		window := newWindow(
			&widget.WindowOption{
				Texture:          frontend.TextureWindow,
				CornerSize:       cornerSize,
				RelativePosition: relativePosition,
				Size:             size,
				Depth:            depth,
				Pivot:            pivot,
				Padding:          padding,
			},
		)

		text := newText(
			&widget.TextOptionsNew{
				RelativePosition: window.GetPadding(),
				Pivot:            frontend.PivotTopLeft,
				Font:             frontend.MaruMinya,
				Speed:            speed,
				Depth:            depth,
			},
		)

		return &MessageWindow{
			text:   text,
			window: window,
			shake:  frontend.NewShake(),
		}
	}
}
