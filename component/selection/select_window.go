package selection

import (
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/input"
)

type selectWindowViewInterface interface {
	update(*drawing.Vector)
	draw(drawing.DrawFunc)
	setCursorRelativePosition(*drawing.Vector)
}

type smoother interface {
	Update()
	Do(input.SmoothKey) bool
}

type textInterface interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
}

type Cursor interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	SetRelativePosition(*drawing.Vector)
	Size() *drawing.Vector
}

type SelectWindow struct {
	cursorPositions []*drawing.Vector
	index           int
	isActive        bool
	isOpen          bool
	onSubmit        func(int)
	closeOnSubmit   bool
	smoother        smoother
	view            selectWindowViewInterface
}

func (w *SelectWindow) OnInputCancel() {}

func (w *SelectWindow) OnInputSubButton() {}

func (w *SelectWindow) OnInputLeft() {}

func (w *SelectWindow) OnInputRight() {}

func (w *SelectWindow) Open() {
	w.isOpen = true
}

func (w *SelectWindow) Close() {
	w.isOpen = false
}

func (w *SelectWindow) Update(parentPosition *drawing.Vector) {
	w.smoother.Update()
	w.view.update(parentPosition)
}

func (w *SelectWindow) Draw(drawFunc drawing.DrawFunc) {
	if !w.isOpen {
		return
	}
	w.view.draw(drawFunc)
}

func (w *SelectWindow) calculateCursorPosition() *drawing.Vector {
	return w.cursorPositions[w.index]
}

func (w *SelectWindow) OnInputUp() {
	if !w.smoother.Do(input.SmoothKeyUp) {
		return
	}
	count := len(w.cursorPositions)
	w.index = (w.index - 1 + count) % count
	w.view.setCursorRelativePosition(w.calculateCursorPosition())
}

func (w *SelectWindow) OnInputDown() {
	if !w.smoother.Do(input.SmoothKeyDown) {
		return
	}
	count := len(w.cursorPositions)
	w.index = (w.index + 1) % count
	w.view.setCursorRelativePosition(w.calculateCursorPosition())
}

func (w *SelectWindow) OnInputSubmit() {
	w.onSubmit(w.index)
	if !w.closeOnSubmit {
		return
	}
	w.Close()
}
