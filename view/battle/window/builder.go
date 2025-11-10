package window

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/common/character/hero"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/selection"
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
	commandPort command.GetCommandModelFunc,
) NewBattleSelectWindowFunc {
	newView := createNewView(newSelectWindow)
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
