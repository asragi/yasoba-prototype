package window

import (
	"github.com/asragi/yasoba-prototype/battle"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/selection"
)

// SelectWindowに対する薄いwrapper
type BattleSelectWindow struct {
	commands     []battle.PlayerCommand
	onSubmit     func(battle.PlayerCommand)
	selectWindow *selection.SelectWindow
}

func (w *BattleSelectWindow) OnInputSubmit() {
	w.selectWindow.OnInputSubmit()
}

func (w *BattleSelectWindow) OnInputCancel() {}

func (w *BattleSelectWindow) OnInputSubButton() {}

func (w *BattleSelectWindow) OnInputLeft() {}

func (w *BattleSelectWindow) OnInputRight() {}

func (w *BattleSelectWindow) Open() {
	w.selectWindow.Open()
}

func (w *BattleSelectWindow) Close() {
	w.selectWindow.Close()
}

func (w *BattleSelectWindow) Update(parentPosition *drawing.Vector) {
	w.selectWindow.Update(parentPosition)
}

func (w *BattleSelectWindow) Draw(drawFunc drawing.DrawFunc) {
	w.selectWindow.Draw(drawFunc)
}

func (w *BattleSelectWindow) OnInputUp() {
	w.selectWindow.OnInputUp()
}

func (w *BattleSelectWindow) OnInputDown() {
	w.selectWindow.OnInputDown()
}

type NewBattleSelectWindowFunc func(
	*drawing.Vector,
	*drawing.Pivot,
	drawing.Depth,
	[]battle.PlayerCommand,
	func(battle.PlayerCommand),
) *BattleSelectWindow

func StandByNewBattleSelectWindow(
	newSelectWindow selection.NewSelectWindowFunc,
) NewBattleSelectWindowFunc {
	return func(
		relativePosition *drawing.Vector,
		pivot *drawing.Pivot,
		depth drawing.Depth,
		commands []battle.PlayerCommand,
		onSubmit func(battle.PlayerCommand),
	) *BattleSelectWindow {
		onSubmitIndex := func(index int) {
			onSubmit(commands[index])
		}
		commandTexts := func() []text.TextId {
			texts := make([]text.TextId, len(commands))
			for i, command := range commands {
				texts[i] = command.ToTextId()
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

		return &BattleSelectWindow{
			commands:     commands,
			onSubmit:     onSubmit,
			selectWindow: window,
		}
	}
}
