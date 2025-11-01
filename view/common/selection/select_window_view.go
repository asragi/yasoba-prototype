package selection

import "github.com/asragi/yasoba-prototype/toolkit/drawing"

type selectWindowView struct {
	texts  []textInterface
	cursor Cursor
}

func newSelectWindowView(texts []textInterface, cursor Cursor) selectWindowViewInterface {
	return &selectWindowView{
		texts:  texts,
		cursor: cursor,
	}
}

func (w *selectWindowView) update(parentPosition *drawing.Vector) {
	for _, text := range w.texts {
		text.Update(parentPosition)
	}
	w.cursor.Update(parentPosition)
}

func (w *selectWindowView) draw(drawFunc drawing.DrawFunc) {
	w.cursor.Draw(drawFunc)
	for _, text := range w.texts {
		text.Draw(drawFunc)
	}
}

func (w *selectWindowView) setCursorRelativePosition(position *drawing.Vector) {
	w.cursor.SetRelativePosition(position)
}
