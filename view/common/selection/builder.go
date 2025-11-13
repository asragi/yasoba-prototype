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
	marginRight float64,
	texts []Item,
	onSubmit func(int),
	onCancel func(),
	closeOnSubmit bool,
) *SelectWindow

type NewCursor func(
	relativePosition *drawing.Vector,
	pivot *drawing.Pivot,
	depth drawing.Depth,
) Cursor

type Item interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	Size() *drawing.Vector
	IsDisabled() bool
}

type NewSelectWindowViewFunc func(
	relativePosition *drawing.Vector,
	commands []Item,
	pivot *drawing.Pivot,
	depth drawing.Depth,
	marginRight float64,
	newWindow window.NewWindowFunc,
	newText widget.NewTextFunc,
	textServer text.ServeTextDataFunc,
	newCursor NewCursor,
) viewInterface

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
		marginRight float64,
		commands []Item,
		onSubmit func(int),
		onCancel func(),
		closeOnSubmit bool,
	) *SelectWindow {

		view := newView(
			relativePosition,
			commands,
			pivot,
			depth,
			marginRight,
			newWindow,
			newText,
			textServer,
			newCursor,
		)

		return &SelectWindow{
			items:         commands,
			index:         0,
			isActive:      false,
			isOpen:        false,
			onSubmit:      onSubmit,
			onCancel:      onCancel,
			smoother:      input.NewInputSmoother(),
			closeOnSubmit: closeOnSubmit,
			view:          view,
		}
	}
}
