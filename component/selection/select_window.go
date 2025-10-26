package selection

import "github.com/asragi/yasoba-prototype/frontend"

type selectWindowViewInterface interface {
	update(*frontend.Vector)
	draw(frontend.DrawFunc)
	setCursorRelativePosition(*frontend.Vector)
}

type smoother interface {
	Update()
	Do(frontend.SmoothKey) bool
}

type textInterface interface {
	Update(*frontend.Vector)
	Draw(frontend.DrawFunc)
}

type Cursor interface {
	Update(*frontend.Vector)
	Draw(frontend.DrawFunc)
	SetRelativePosition(*frontend.Vector)
	Size() *frontend.Vector
}

type SelectWindow struct {
	cursorPositions []*frontend.Vector
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

func (w *SelectWindow) Update(parentPosition *frontend.Vector) {
	w.smoother.Update()
	w.view.update(parentPosition)
}

func (w *SelectWindow) Draw(drawFunc frontend.DrawFunc) {
	if !w.isOpen {
		return
	}
	w.view.draw(drawFunc)
}

func (w *SelectWindow) calculateCursorPosition() *frontend.Vector {
	return w.cursorPositions[w.index]
}

func (w *SelectWindow) OnInputUp() {
	if !w.smoother.Do(frontend.SmoothKeyUp) {
		return
	}
	count := len(w.cursorPositions)
	w.index = (w.index - 1 + count) % count
	w.view.setCursorRelativePosition(w.calculateCursorPosition())
}

func (w *SelectWindow) OnInputDown() {
	if !w.smoother.Do(frontend.SmoothKeyDown) {
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
