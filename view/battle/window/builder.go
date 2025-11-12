package window

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/common/character/hero"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/window/line"
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
	newTextItemWithCost line.NewTextWithCostOptionFunc,
	commandPort command.GetCommandModelFunc,
) NewBattleSelectWindowFunc {
	newView := createNewView(newSelectWindow, newTextItemWithCost)
	return func(
		playerMp hero.MP,
		position *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		commands []command.Id,
		outerOnSubmit func(command.Id),
	) *BattleSelectWindow {
		createView := func(onSubmit func(command.Id)) viewInterface {
			costList := func() map[command.Id]command.Cost {
				m := make(map[command.Id]command.Cost, len(commands))
				for _, id := range commands {
					model := commandPort(id)
					if model == nil {
						panic("command model not found: " + string(id))
					}
					m[id] = model.Cost()
				}
				return m
			}()
			return newView(
				position,
				pivot,
				depth,
				commands,
				costList,
				onSubmit,
			)
		}
		return newBattleSelectWindow(playerMp, commands, commandPort, createView, outerOnSubmit)
	}
}
