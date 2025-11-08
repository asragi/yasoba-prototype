package window

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/selection"
)

type view struct {
	selectWindow *selection.SelectWindow
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

type newViewFunc func(
	*drawing.Vector,
	*drawing.Pivot,
	drawing.Depth,
	[]command.Id,
	func(command.Id),
) *view

func createNewView(
	newSelectWindow selection.NewSelectWindowFunc,
) newViewFunc {
	return func(
		relativePosition *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		commands []command.Id,
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
		marginRight := 10.0
		window := newSelectWindow(
			relativePosition,
			pivot,
			depth,
			marginRight,
			commandTexts,
			onSubmitIndex,
			false,
		)

		return &view{
			selectWindow: window,
		}
	}
}
