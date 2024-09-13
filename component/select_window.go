package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type SelectWindow struct {
	texts           []*widget.Text
	cursor          *widget.Image
	cursorPositions []*frontend.Vector
	index           int
	isActive        bool
	isOpen          bool
	onSubmit        func(int)
	smoother        *frontend.InputSmoother
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
	for _, text := range w.texts {
		text.Update(parentPosition)
	}
	w.cursor.Update(parentPosition)
}

func (w *SelectWindow) Draw(drawFunc frontend.DrawFunc) {
	if !w.isOpen {
		return
	}
	w.cursor.Draw(drawFunc)
	for _, text := range w.texts {
		text.Draw(drawFunc)
	}
}

func (w *SelectWindow) calculateCursorPosition() *frontend.Vector {
	return w.cursorPositions[w.index]
}

func (w *SelectWindow) OnInputUp() {
	if !w.smoother.Do(frontend.SmoothKeyUp) {
		return
	}
	w.index = (w.index - 1 + len(w.texts)) % len(w.texts)
	w.cursor.SetRelativePosition(w.calculateCursorPosition())
}

func (w *SelectWindow) OnInputDown() {
	if !w.smoother.Do(frontend.SmoothKeyDown) {
		return
	}
	w.index = (w.index + 1) % len(w.texts)
	w.cursor.SetRelativePosition(w.calculateCursorPosition())
}

func (w *SelectWindow) OnInputSubmit() {
	w.onSubmit(w.index)
}

type NewSelectWindowFunc func(
	*frontend.Vector,
	*frontend.Pivot,
	frontend.Depth,
	[]core.TextId,
	func(int),
) *SelectWindow

func StandByNewSelectWindow(
	resource *frontend.ResourceManager,
	textServer core.ServeTextDataFunc,
) NewSelectWindowFunc {
	font := resource.GetFont(frontend.MaruMinya)
	return func(
		relativePosition *frontend.Vector,
		pivot *frontend.Pivot,
		depth frontend.Depth,
		commands []core.TextId,
		onSubmit func(int),
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
		cursor := widget.NewImage(
			cursorPositions[0],
			frontend.PivotTopLeft,
			depth,
			resource.GetTexture(frontend.TextureCursor),
		)
		cursorWidth := cursor.Size().X
		texts := func() []*widget.Text {
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
			var texts []*widget.Text
			for i, command := range commands {
				text := widget.NewText(
					&widget.TextOptions{
						RelativePosition: relativePositions[i],
						Pivot:            frontend.PivotTopLeft,
						TextFace:         font,
						Speed:            1,
						Depth:            depth,
					},
				)
				text.SetText(textServer(command).Text, true)
				texts = append(texts, text)
			}
			return texts
		}()

		return &SelectWindow{
			texts:           texts,
			cursor:          cursor,
			cursorPositions: cursorPositions,
			index:           0,
			isActive:        false,
			isOpen:          false,
			onSubmit:        onSubmit,
			smoother:        frontend.NewInputSmoother(),
		}
	}
}
