package line

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
)

type NewTextWithCostOptionFunc func(
	width float64,
	textId text.TextId,
	cost command.Cost,
) *textOptionWithCost
type newCostDisplayFunc func(cost int) interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
}
type newTextLineFunc func(text.TextId) interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	Size() *drawing.Vector
	SetDisable(bool)
	IsDisabled() bool
}

func CreateNewTextWithCostOption(
	newCostDisplay newCostDisplayFunc,
	newTextLine newTextLineFunc,
) NewTextWithCostOptionFunc {
	return func(
		width float64,
		textId text.TextId,
		cost command.Cost,
	) *textOptionWithCost {
		textLine := newTextLine(textId)
		costDisplay := newCostDisplay(cost.ToInt())
		view := newView(
			textLine,
			textLine.Size().Y,
			costDisplay,
		)
		return newPresenter(
			width,
			textLine,
			view,
		)
	}
}
