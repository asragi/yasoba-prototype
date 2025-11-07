package message

import (
	commonfont "github.com/asragi/yasoba-prototype/common/font"
	"github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/shake"
	widgetwindow "github.com/asragi/yasoba-prototype/view/common/window"
	"github.com/asragi/yasoba-prototype/widget"
)

type MessageWindow struct {
	text   widget.TextInterface
	window widgetwindow.WindowInterface
	shake  *shake.EmitShake
	isOpen bool
}

func (m *MessageWindow) Shake(amplitude float64, period int) {
	m.shake.Shake(amplitude, period)
}

func (m *MessageWindow) FitToMessage() {
	// TODO: content size だけを渡して padding の処理は window 側に任せるべき
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

func (m *MessageWindow) Update(parentPosition *drawing.Vector) {
	m.shake.Update()
	m.window.Update(parentPosition)
	m.text.Update(m.window.GetPositionUpperLeft())
}

func (m *MessageWindow) Draw(drawFunc drawing.DrawFunc) {
	if !m.isOpen {
		return
	}
	m.window.Draw(drawFunc)
	m.text.Draw(drawFunc)
}

type NewMessageWindowFunc func(
	relativePosition *drawing.Vector,
	size *drawing.Vector,
	depth drawing.Depth,
	pivot *drawing.Pivot,
) *MessageWindow

func StandByNewMessageWindow(
	newText widget.NewTextFunc,
	newWindow widgetwindow.NewWindowFunc,
) NewMessageWindowFunc {
	cornerSize := 6
	padding := &drawing.Vector{X: 16, Y: 8}
	speed := 5
	return func(
		relativePosition *drawing.Vector,
		size *drawing.Vector,
		depth drawing.Depth,
		pivot *drawing.Pivot,
	) *MessageWindow {
		window := newWindow(
			&widgetwindow.WindowOption{
				Texture:          texture.Window,
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
				Pivot:            drawing.PivotTopLeft,
				Font:             commonfont.MaruMinya,
				Speed:            speed,
				Depth:            depth,
			},
		)

		return &MessageWindow{
			text:   text,
			window: window,
			shake:  shake.NewShake(),
		}
	}
}
