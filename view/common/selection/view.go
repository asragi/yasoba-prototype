package selection

import (
	"github.com/asragi/yasoba-prototype/global"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/window"
	"github.com/asragi/yasoba-prototype/widget"
)

type windowInterface interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	SetSize(*drawing.Vector)
	GetPadding() *drawing.Vector
	GetPositionUpperLeft() *drawing.Vector
}

type selectWindowView struct {
	cursorPositions []*drawing.Vector
	window          windowInterface
	items           []Item
	cursor          Cursor
	padding         *drawing.Vector
	lineHeight      float64
	lineMargin      float64
}

func NewSelectWindowView(
	relativePosition *drawing.Vector,
	options []Item,
	pivot *drawing.Pivot,
	depth drawing.Depth,
	marginRight float64,
	newWindow window.NewWindowFunc,
	newText widget.NewTextFunc,
	textServer text.ServeTextDataFunc,
	newCursor NewCursor,
) viewInterface {
	count := len(options)
	// TODO: Use actual values
	lineHeight := options[0].Size().Y
	lineMargin := global.LineMargin
	//const marginX = 4
	//const offsetY = -1
	cursorWidth := 4.0
	padding := window.WindowOptionDefaultPadding
	// size := &drawing.Vector{X: width, Y: float64(lineHeight * count)}
	// pivotModification := pivot.ApplyToSize(size)
	cursorPositions := func() []*drawing.Vector {
		positions := make([]*drawing.Vector, count)
		for i := 0; i < count; i++ {
			positions[i] = &drawing.Vector{
				// TODO: たまたま2.0で割るといい感じなだけ
				X: padding.X / 3.0,
				Y: padding.Y + (lineHeight+lineMargin)*float64(i),
			}
		}
		return positions
	}()
	cursor := newCursor(
		cursorPositions[0],
		drawing.PivotTopLeft,
		depth,
	)
	contentSize := func() *drawing.Vector {
		width := func() float64 {
			var maxWidth float64
			for _, option := range options {
				if option.Size().X > maxWidth {
					maxWidth = option.Size().X
				}
			}
			return maxWidth
		}()
		height := lineHeight * float64(count)
		return &drawing.Vector{X: width, Y: height}
	}()
	contentSizeWithPadding := contentSize.Add(padding.Multiply(2)).Add(&drawing.Vector{X: cursorWidth, Y: 0}).Add(&drawing.Vector{X: marginRight, Y: 0})
	sizeWithMargin := contentSizeWithPadding.Add(&drawing.Vector{X: 0, Y: lineMargin * float64(count-1)})
	window := newWindow(&window.WindowOption{
		Texture:          window.WindowOptionDefaultTexture,
		CornerSize:       window.WindowOptionDefaultCornerSize,
		RelativePosition: relativePosition,
		Size:             sizeWithMargin,
		Depth:            depth,
		Pivot:            pivot,
		Padding:          padding,
	})
	cursor.SetRelativePosition(cursorPositions[0])
	return &selectWindowView{
		cursorPositions: cursorPositions,
		window:          window,
		items:           options,
		cursor:          cursor,
		padding:         padding,
		lineHeight:      lineHeight,
		lineMargin:      lineMargin,
	}
}

func (w *selectWindowView) update(parentPosition *drawing.Vector) {
	w.window.Update(parentPosition)
	windowPos := w.window.GetPositionUpperLeft()
	for i, text := range w.items {
		// text.Update(parentPosition)
		text.Update(windowPos.Add(w.padding).Add(&drawing.Vector{X: 0, Y: (w.lineHeight + w.lineMargin) * float64(i)}))
	}
	w.cursor.Update(windowPos)
}

func (w *selectWindowView) draw(drawFunc drawing.DrawFunc) {
	w.window.Draw(drawFunc)
	w.cursor.Draw(drawFunc)
	for _, item := range w.items {
		item.Draw(drawFunc)
	}
}

func (w *selectWindowView) onChangeIndex(index int) {
	w.cursor.SetRelativePosition(w.cursorPositions[index])
}
