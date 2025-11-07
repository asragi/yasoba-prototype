package selection

import (
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/toolkit/input"
	"github.com/asragi/yasoba-prototype/view/common/window"
	"github.com/asragi/yasoba-prototype/widget"
)

type NewSelectWindowFunc func(
	relativePosition *drawing.Vector,
	pivot *drawing.Pivot,
	depth drawing.Depth,
	texts []text.TextId,
	onSubmit func(int),
	closeOnSubmit bool,
) *SelectWindow

type NewCursor func(
	relativePosition *drawing.Vector,
	pivot *drawing.Pivot,
	depth drawing.Depth,
) Cursor

type NewSelectWindowViewFunc func(
	relativePosition *drawing.Vector,
	commands []text.TextId,
	pivot *drawing.Pivot,
	depth drawing.Depth,
	newWindow window.NewWindowFunc,
	newText widget.NewTextFunc,
	textServer text.ServeTextDataFunc,
	newCursor NewCursor,
) selectWindowViewInterface

func StandByNewSelectWindow(
	newCursor NewCursor,
	newText widget.NewTextFunc,
	newWindow window.NewWindowFunc,
	textServer text.ServeTextDataFunc,
	newView NewSelectWindowViewFunc,
) NewSelectWindowFunc {
	return func(
		relativePosition *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		commands []text.TextId,
		onSubmit func(int),
		closeOnSubmit bool,
	) *SelectWindow {

		view := newView(
			relativePosition,
			commands,
			pivot,
			depth,
			newWindow,
			newText,
			textServer,
			newCursor,
		)

		return &SelectWindow{
			indexSize:     len(commands),
			index:         0,
			isActive:      false,
			isOpen:        false,
			onSubmit:      onSubmit,
			smoother:      input.NewInputSmoother(),
			closeOnSubmit: closeOnSubmit,
			view:          view,
		}
	}
}
