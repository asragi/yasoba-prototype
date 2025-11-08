package window

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
)

type viewInterface interface {
	Open()
	Close()
	Update(parentPosition *drawing.Vector)
	Draw(drawFunc drawing.DrawFunc)
	OnInputUp()
	OnInputDown()
	OnInputSubmit()
}

// SelectWindowに対する薄いwrapper
type BattleSelectWindow struct {
	view     viewInterface
	commands map[command.Id]command.Model
}

func newBattleSelectWindow(
	commands []command.Id,
	commandPort command.GetCommandModelFunc,
	view viewInterface,
) *BattleSelectWindow {
	commandsDict := func() map[command.Id]command.Model {
		m := make(map[command.Id]command.Model, len(commands))
		for _, id := range commands {
			model := commandPort(id)
			if model == nil {
				panic("command model not found: " + string(id))
			}
			m[id] = model
		}
		return m
	}()
	return &BattleSelectWindow{
		view:     view,
		commands: commandsDict,
	}
}

func (w *BattleSelectWindow) OnInputSubmit() {
	w.view.OnInputSubmit()
}

func (w *BattleSelectWindow) OnInputCancel() {}

func (w *BattleSelectWindow) OnInputSubButton() {}

func (w *BattleSelectWindow) OnInputLeft() {}

func (w *BattleSelectWindow) OnInputRight() {}

func (w *BattleSelectWindow) Open() {
	w.view.Open()
}

func (w *BattleSelectWindow) Close() {
	w.view.Close()
}

func (w *BattleSelectWindow) Update(parentPosition *drawing.Vector) {
	w.view.Update(parentPosition)
}

func (w *BattleSelectWindow) Draw(drawFunc drawing.DrawFunc) {
	w.view.Draw(drawFunc)
}

func (w *BattleSelectWindow) OnInputUp() {
	w.view.OnInputUp()
}

func (w *BattleSelectWindow) OnInputDown() {
	w.view.OnInputDown()
}
