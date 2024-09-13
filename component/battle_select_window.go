package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
)

type BattleCommand int

type BattleSelectWindow struct {
	commands     []core.PlayerCommand
	onSubmit     func(core.PlayerCommand)
	selectWindow *SelectWindow
}

func (w *BattleSelectWindow) Open() {
	w.selectWindow.Open()
}

func (w *BattleSelectWindow) Close() {
	w.selectWindow.Close()
}

func (w *BattleSelectWindow) Update(parentPosition *frontend.Vector) {
	w.selectWindow.Update(parentPosition)
}

func (w *BattleSelectWindow) Draw(drawFunc frontend.DrawFunc) {
	w.selectWindow.Draw(drawFunc)
}

func (w *BattleSelectWindow) MoveCursorUp() {
	w.selectWindow.MoveCursorUp()
}

func (w *BattleSelectWindow) MoveCursorDown() {
	w.selectWindow.MoveCursorDown()
}

func (w *BattleSelectWindow) OnSubmit() {
	w.selectWindow.OnSubmit()
}

type NewBattleSelectWindowFunc func(
	*frontend.Vector,
	*frontend.Pivot,
	frontend.Depth,
	[]core.PlayerCommand,
	func(core.PlayerCommand),
) *BattleSelectWindow

func StandByNewBattleSelectWindow(
	newSelectWindow NewSelectWindowFunc,
) NewBattleSelectWindowFunc {
	return func(
		relativePosition *frontend.Vector,
		pivot *frontend.Pivot,
		depth frontend.Depth,
		commands []core.PlayerCommand,
		onSubmit func(core.PlayerCommand),
	) *BattleSelectWindow {
		onSubmitIndex := func(index int) {
			onSubmit(commands[index])
		}
		commandTexts := func() []core.TextId {
			texts := make([]core.TextId, len(commands))
			for i, command := range commands {
				texts[i] = command.ToTextId()
			}
			return texts
		}()
		window := newSelectWindow(
			relativePosition,
			pivot,
			depth,
			commandTexts,
			onSubmitIndex,
		)

		return &BattleSelectWindow{
			commands:     commands,
			onSubmit:     onSubmit,
			selectWindow: window,
		}
	}
}
