package selection

import (
	"github.com/asragi/yasoba-prototype/common/font"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/constant"
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

type textInterface interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	Size() *drawing.Vector
}

type selectWindowView struct {
	cursorPositions []*drawing.Vector
	window          windowInterface
	texts           []textInterface
	cursor          Cursor
}

func NewSelectWindowView(
	relativePosition *drawing.Vector,
	commands []text.TextId,
	pivot *drawing.Pivot,
	depth drawing.Depth,
	marginRight float64,
	newWindow window.NewWindowFunc,
	newText widget.NewTextFunc,
	textServer text.ServeTextDataFunc,
	newCursor NewCursor,
) selectWindowViewInterface {
	count := len(commands)
	// TODO: Use actual values
	lineHeight := constant.CommandLineHeight
	//const marginX = 4
	//const offsetY = -1
	cursorWidth := 4.0
	padding := window.WindowOptionDefaultPadding
	// size := &drawing.Vector{X: width, Y: float64(lineHeight * count)}
	// pivotModification := pivot.ApplyToSize(size)
	cursorPositions := func() []*drawing.Vector {
		positions := make([]*drawing.Vector, len(commands))
		for i := 0; i < count; i++ {
			positions[i] = &drawing.Vector{
				// TODO: たまたま2.0で割るといい感じなだけ
				X: padding.X / 2.0,
				Y: padding.Y + float64(lineHeight*i),
			}
		}
		return positions
	}()
	cursor := newCursor(
		cursorPositions[0],
		drawing.PivotTopLeft,
		depth,
	)
	texts := func() []textInterface {
		relativePositions := func() []*drawing.Vector {
			var positions []*drawing.Vector
			for i := 0; i < count; i++ {
				positions = append(
					positions, &drawing.Vector{
						X: cursorWidth + padding.X,           //relativePosition.X - pivotModification.X + cursorWidth + marginX,
						Y: float64(lineHeight*i) + padding.Y, //relativePosition.Y - pivotModification.Y + float64(lineHeight*i) + offsetY,
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
					Font:             font.MaruMinya,
					Speed:            1,
					Depth:            depth,
				},
			)
			text.SetText(textServer(command).Text.String(), true)
			texts = append(texts, text)
		}
		return texts
	}()
	contentSize := func() *drawing.Vector {
		width := func() float64 {
			var maxWidth float64
			for _, text := range texts {
				if text.Size().X > maxWidth {
					maxWidth = text.Size().X
				}
			}
			return maxWidth
		}()
		height := float64(lineHeight * len(texts))
		return &drawing.Vector{X: width, Y: height}
	}()
	contentSizeWithPadding := contentSize.Add(padding.Multiply(2)).Add(&drawing.Vector{X: cursorWidth, Y: 0}).Add(&drawing.Vector{X: marginRight, Y: 0})
	window := newWindow(&window.WindowOption{
		Texture:          window.WindowOptionDefaultTexture,
		CornerSize:       window.WindowOptionDefaultCornerSize,
		RelativePosition: relativePosition,
		Size:             contentSizeWithPadding,
		Depth:            depth,
		Pivot:            pivot,
		Padding:          padding,
	})
	cursor.SetRelativePosition(cursorPositions[0])
	return &selectWindowView{
		cursorPositions: cursorPositions,
		window:          window,
		texts:           texts,
		cursor:          cursor,
	}
}

func (w *selectWindowView) update(parentPosition *drawing.Vector) {
	w.window.Update(parentPosition)
	windowPos := w.window.GetPositionUpperLeft()
	for _, text := range w.texts {
		// text.Update(parentPosition)
		text.Update(windowPos)
	}
	w.cursor.Update(windowPos)
}

func (w *selectWindowView) draw(drawFunc drawing.DrawFunc) {
	w.window.Draw(drawFunc)
	w.cursor.Draw(drawFunc)
	for _, text := range w.texts {
		text.Draw(drawFunc)
	}
}

func (w *selectWindowView) onChangeIndex(index int) {
	w.cursor.SetRelativePosition(w.cursorPositions[index])
}
