package selection

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/frontend"
)

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

func (w *selectWindowView) update(parentPosition *frontend.Vector) {
	for _, text := range w.texts {
		text.Update(parentPosition)
	}
	w.cursor.Update(parentPosition)
}

func (w *selectWindowView) draw(drawFunc frontend.DrawFunc) {
	w.cursor.Draw(drawFunc)
	fmt.Printf("text: %+v\n", w.texts)
	for _, text := range w.texts {
		text.Draw(drawFunc)
	}
}

func (w *selectWindowView) setCursorRelativePosition(position *frontend.Vector) {
	w.cursor.SetRelativePosition(position)
}
