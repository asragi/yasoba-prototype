package window

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/global"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/window/line"
	"github.com/asragi/yasoba-prototype/view/common/selection"
)

type itemInterface interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	Size() *drawing.Vector
	IsDisabled() bool
	SetDisable(bool)
}

type view struct {
	selectWindow *selection.SelectWindow
	items        []itemInterface
}

func (v *view) Open() {
	v.selectWindow.Open()
}

func (v *view) Close() {
	v.selectWindow.Close()
}

func (v *view) Update(parentPosition *drawing.Vector) {
	v.selectWindow.Update(parentPosition)
}

func (v *view) Draw(drawFunc drawing.DrawFunc) {
	v.selectWindow.Draw(drawFunc)
}

func (v *view) OnInputUp() {
	v.selectWindow.OnInputUp()
}

func (v *view) OnInputDown() {
	v.selectWindow.OnInputDown()
}

func (v *view) OnInputSubmit() {
	v.selectWindow.OnInputSubmit()
}

func (v *view) SetDisable(index int, disable bool) {
	v.items[index].SetDisable(disable)
}

type newViewFunc func(
	*drawing.Vector,
	*drawing.Pivot,
	drawing.Depth,
	[]command.Id,
	map[command.Id]command.Cost,
	func(command.Id),
) *view

func createNewView(
	newSelectWindow selection.NewSelectWindowFunc,
	newTextWithCost line.NewTextWithCostOptionFunc,
) newViewFunc {
	return func(
		relativePosition *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		commands []command.Id,
		costList map[command.Id]command.Cost,
		onSubmit func(command.Id),
	) *view {
		onSubmitIndex := func(index int) {
			onSubmit(commands[index])
		}
		commandTexts := func() []text.TextId {
			texts := make([]text.TextId, len(commands))
			for i, cmd := range commands {
				texts[i] = cmd.ToTextId()
			}
			return texts
		}()
		items := func() []itemInterface {
			items := make([]itemInterface, len(commandTexts))
			for i, textId := range commandTexts {
				items[i] = newTextWithCost(global.SelectionWindowWidth, textId, costList[commands[i]])
			}
			return items
		}()
		itemsToArg := func() []selection.Item {
			result := make([]selection.Item, len(items))
			for i, item := range items {
				result[i] = item
			}
			return result
		}()
		marginRight := 10.0
		window := newSelectWindow(
			relativePosition,
			pivot,
			depth,
			marginRight,
			itemsToArg,
			onSubmitIndex,
			nil,
			false,
		)

		return &view{
			items:        items,
			selectWindow: window,
		}
	}
}
