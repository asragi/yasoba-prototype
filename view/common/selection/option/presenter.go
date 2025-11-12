package option

import (
	"image/color"

	"github.com/asragi/yasoba-prototype/global"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/widget"
)

type textInterface interface {
	SetText(string, bool)
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	SetTextColor(color color.Color)
	Size() *drawing.Vector
}

type newTextFunc func(
	options *widget.TextOptionsNew,
) textInterface

type NewTextItemFunc func(text.TextId) *TextItem

type TextItem struct {
	text     textInterface
	disabled bool
}

func CreateNewTextItem(
	newText widget.NewTextFunc,
	textServer text.ServeTextDataFunc,
) NewTextItemFunc {
	return func(
		textId text.TextId,
	) *TextItem {
		text := newText(
			&widget.TextOptionsNew{
				RelativePosition: &drawing.Vector{X: 0, Y: 0},
				Pivot:            drawing.PivotTopLeft,
				Font:             global.PrimaryFont,
				Speed:            1,
				Depth:            drawing.DepthWindow,
			},
		)
		textData := textServer(textId)
		text.SetText(textData.Text.String(), true)
		return &TextItem{
			text: text,
		}
	}
}

func (i *TextItem) Update(parentPosition *drawing.Vector) {
	i.text.Update(parentPosition)
}

func (i *TextItem) Draw(drawFunc drawing.DrawFunc) {
	i.text.Draw(drawFunc)
}

func (i *TextItem) Size() *drawing.Vector {
	return i.text.Size()
}

func (i *TextItem) SetDisable(disabled bool) {
	i.disabled = disabled
	if disabled {
		i.text.SetTextColor(global.DisableColor)
		return
	}
	i.text.SetTextColor(color.White)
}

func (i *TextItem) IsDisabled() bool {
	return i.disabled
}
