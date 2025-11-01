package selection

import (
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/input"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/widget"
)

type NewSelectWindowFunc func(
	*drawing.Vector,
	*drawing.Pivot,
	drawing.Depth,
	[]text.TextId,
	func(int),
	bool,
) *SelectWindow

type NewCursor func(
	relativePosition *drawing.Vector,
	pivot *drawing.Pivot,
	depth drawing.Depth,
) Cursor

func StandByNewSelectWindow(
	newCursor NewCursor,
	newText widget.NewTextFunc,
	textServer text.ServeTextDataFunc,
) NewSelectWindowFunc {
	return func(
		relativePosition *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
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
		size := &drawing.Vector{X: width, Y: float64(lineHeight * count)}
		pivotModification := pivot.ApplyToSize(size)
		cursorPositions := func() []*drawing.Vector {
			positions := make([]*drawing.Vector, len(commands))
			for i := 0; i < count; i++ {
				positions[i] = &drawing.Vector{
					X: relativePosition.X - pivotModification.X,
					Y: relativePosition.Y - pivotModification.Y + float64(lineHeight*i),
				}
			}
			return positions
		}()
		cursor := newCursor(
			cursorPositions[0],
			drawing.PivotTopLeft,
			depth,
		)
		cursorWidth := cursor.Size().X
		texts := func() []textInterface {
			relativePositions := func() []*drawing.Vector {
				var positions []*drawing.Vector
				for i := 0; i < count; i++ {
					positions = append(
						positions, &drawing.Vector{
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
						Pivot:            drawing.PivotTopLeft,
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
			smoother:        input.NewInputSmoother(),
			closeOnSubmit:   closeOnSubmit,
			view:            view,
		}
	}
}
