package window

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/selection"
)

type NewBattleSelectWindowFunc func(
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
		position *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		commands []command.Id,
		onSubmit func(command.Id),
	) *BattleSelectWindow {
		view := newView(
			position,
			pivot,
			depth,
			commands,
			onSubmit,
		)
		return newBattleSelectWindow(commands, commandPort, view)
	}
}
