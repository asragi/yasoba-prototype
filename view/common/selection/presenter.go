package selection

import (
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/toolkit/input"
)

type viewInterface interface {
	update(*drawing.Vector)
	draw(drawing.DrawFunc)
	onChangeIndex(int)
}

type smoother interface {
	Update()
	Do(input.SmoothKey) bool
}

type Cursor interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	SetRelativePosition(*drawing.Vector)
	Size() *drawing.Vector
}

type SelectWindow struct {
	items         []Item
	index         int
	isActive      bool
	isOpen        bool
	onSubmit      func(int)
	closeOnSubmit bool
	smoother      smoother
	view          viewInterface
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

func (w *SelectWindow) OnInputUp() {
	if !w.smoother.Do(input.SmoothKeyUp) {
		return
	}
	count := len(w.items)
	w.index = (w.index - 1 + count) % count
	w.view.onChangeIndex(w.index)
}

func (w *SelectWindow) OnInputDown() {
	if !w.smoother.Do(input.SmoothKeyDown) {
		return
	}
	count := len(w.items)
	w.index = (w.index + 1) % count
	w.view.onChangeIndex(w.index)
}

func (w *SelectWindow) OnInputSubmit() {
	if w.items[w.index].IsDisabled() {
		return
	}
	w.onSubmit(w.index)
	if !w.closeOnSubmit {
		return
	}
	w.Close()
}
