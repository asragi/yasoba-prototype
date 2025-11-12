package window

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/common/character/hero"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/window/cost"
	"github.com/asragi/yasoba-prototype/view/common/selection"
	"github.com/asragi/yasoba-prototype/view/common/selection/option"
)

type NewBattleSelectWindowFunc func(
	hero.MP,
	*drawing.Vector,
	*drawing.Pivot,
	drawing.Depth,
	[]command.Id,
	func(command.Id),
) *BattleSelectWindow

func StandByNewBattleSelectWindow(
	newSelectWindow selection.NewSelectWindowFunc,
	newTextItem option.NewTextItemFunc,
	newCostDisplay cost.NewCostDisplayFunc,
	commandPort command.GetCommandModelFunc,
) NewBattleSelectWindowFunc {
	newTextItemWrapper := func(textId text.TextId) itemInterface {
		return newTextItem(textId)
	}
	newView := createNewView(newSelectWindow, newTextItemWrapper)
	return func(
		playerMp hero.MP,
		position *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		commands []command.Id,
		outerOnSubmit func(command.Id),
	) *BattleSelectWindow {
		createView := func(onSubmit func(command.Id)) viewInterface {
			return newView(
				position,
				pivot,
				depth,
				commands,
				onSubmit,
			)
		}
		return newBattleSelectWindow(playerMp, commands, commandPort, createView, outerOnSubmit)
	}
}
