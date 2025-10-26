package selection

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/widget"
)

type NewSelectWindowFunc func(
	*frontend.Vector,
	*frontend.Pivot,
	frontend.Depth,
	[]text.TextId,
	func(int),
	bool,
) *SelectWindow

type NewCursor func(
	relativePosition *frontend.Vector,
	pivot *frontend.Pivot,
	depth frontend.Depth,
) Cursor

func StandByNewSelectWindow(
	newCursor NewCursor,
	newText widget.NewTextFunc,
	textServer text.ServeTextDataFunc,
) NewSelectWindowFunc {
	return func(
		relativePosition *frontend.Vector,
		pivot *frontend.Pivot,
		depth frontend.Depth,
		commands []text.TextId,
		onSubmit func(int),
		closeOnSubmit bool,
	) *SelectWindow {
		// TODO: Use actual values
		const lineHeight = 16
		const width = 64
		const marginX = 4
		const offsetY = -1
		count := len(commands)
		size := &frontend.Vector{X: width, Y: float64(lineHeight * count)}
		pivotModification := pivot.ApplyToSize(size)
		cursorPositions := func() []*frontend.Vector {
			positions := make([]*frontend.Vector, len(commands))
			for i := 0; i < count; i++ {
				positions[i] = &frontend.Vector{
					X: relativePosition.X - pivotModification.X,
					Y: relativePosition.Y - pivotModification.Y + float64(lineHeight*i),
				}
			}
			return positions
		}()
		cursor := newCursor(
			cursorPositions[0],
			frontend.PivotTopLeft,
			depth,
		)
		cursorWidth := cursor.Size().X
		texts := func() []textInterface {
			relativePositions := func() []*frontend.Vector {
				var positions []*frontend.Vector
				for i := 0; i < count; i++ {
					positions = append(
						positions, &frontend.Vector{
							X: relativePosition.X - pivotModification.X + cursorWidth + marginX,
							Y: relativePosition.Y - pivotModification.Y + float64(lineHeight*i) + offsetY,
						},
					)
				}
				return positions
			}()
			var texts []textInterface
			for i, command := range commands {
				text := newText(
					&widget.TextOptionsNew{
						RelativePosition: relativePositions[i],
						Pivot:            frontend.PivotTopLeft,
						Font:             frontend.MaruMinya,
						Speed:            1,
						Depth:            depth,
					},
				)
				text.SetText(textServer(command).Text.String(), true)
				texts = append(texts, text)
			}
			return texts
		}()

		view := newSelectWindowView(texts, cursor)
		view.setCursorRelativePosition(cursorPositions[0])

		return &SelectWindow{
			cursorPositions: cursorPositions,
			index:           0,
			isActive:        false,
			isOpen:          false,
			onSubmit:        onSubmit,
			smoother:        frontend.NewInputSmoother(),
			closeOnSubmit:   closeOnSubmit,
			view:            view,
		}
	}
}
