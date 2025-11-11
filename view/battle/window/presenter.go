package window

import (
	"image/color"

	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/common/character/hero"
	"github.com/asragi/yasoba-prototype/global"
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
	SetCommandColor(index int, color color.Color)
}

// SelectWindowに対する薄いwrapper
type BattleSelectWindow struct {
	view       viewInterface
	commands   map[command.Id]command.Cost
	commandIds []command.Id
	playerMp   hero.MP
}

func newBattleSelectWindow(
	playerMp hero.MP,
	commands []command.Id,
	commandPort command.GetCommandModelFunc,
	createView func(onSubmit func(command.Id)) viewInterface,
	outerOnSubmit func(command.Id),
) *BattleSelectWindow {
	commandsDict := func() map[command.Id]command.Cost {
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
	window := &BattleSelectWindow{
		commands:   commandsDict,
		commandIds: commands,
		playerMp:   playerMp,
	}
	onSubmit := func(id command.Id) {
		if !window.canSelectCommand(id) {
			return
		}
		outerOnSubmit(id)
	}
	window.view = createView(onSubmit)
	window.updateCommandColors()
	return window
}

func (w *BattleSelectWindow) canSelectCommand(id command.Id) bool {
	cost := w.commands[id]
	return cost.HasEnoughMP(w.playerMp)
}

func (w *BattleSelectWindow) OnInputSubmit() {
	w.view.OnInputSubmit()
}

func (w *BattleSelectWindow) OnInputCancel() {}

func (w *BattleSelectWindow) OnInputSubButton() {}

func (w *BattleSelectWindow) OnInputLeft() {}

func (w *BattleSelectWindow) OnInputRight() {}

func (w *BattleSelectWindow) Open(mp hero.MP) {
	w.playerMp = mp
	w.updateCommandColors()
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

func (w *BattleSelectWindow) updateCommandColors() {
	for idx, id := range w.commandIds {
		var textColor color.Color = color.White
		if !w.canSelectCommand(id) {
			textColor = global.DisableColor
		}
		w.view.SetCommandColor(idx, textColor)
	}
}
