package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	battleactor "github.com/asragi/yasoba-prototype/component/battle/actor"
	battledialogue "github.com/asragi/yasoba-prototype/component/battle/dialogue"
	battleenemy "github.com/asragi/yasoba-prototype/component/battle/enemy"
	battleselect "github.com/asragi/yasoba-prototype/component/battle/window"
	"github.com/asragi/yasoba-prototype/component/selection"
	"github.com/asragi/yasoba-prototype/component/shake"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/input"
	"github.com/asragi/yasoba-prototype/widget"
)

type battleUI struct {
	messageWindow      *component.MessageWindow
	battleSelectWindow *battleselect.BattleSelectWindow
	actorDisplay       *battleactor.BattleActorDisplay
	subActorDisplay    *battleactor.BattleSubActorDisplay
	subActorDialog     *battledialogue.BattlePartnerDialogue
	targetSelectWindow *selection.SelectWindow
	input              input.InputManager
	battleEnemyDisplay *battleenemy.BattleEnemyDisplay
	effectManager      *widget.EffectManager
	shake              *shake.EmitShake
	layout             battleUILayout
}

type battleUILayout struct {
	messageWindowPosition drawing.Vector
	actorDisplayPosition  drawing.Vector
	subActorDisplayPos    drawing.Vector
	enemyDisplayPosition  drawing.Vector
	subActorDialogAnchor  drawing.Vector
	selectAnchor          drawing.Vector
}

func (ui *battleUI) Update() {
	delta := updateBattleShake(ui.shake)
	ui.layout = computeBattleUILayout(ui, delta)

	ui.messageWindow.Update(&ui.layout.messageWindowPosition)
	ui.actorDisplay.Update(&ui.layout.actorDisplayPosition)
	ui.subActorDisplay.Update(&ui.layout.subActorDisplayPos)
	ui.subActorDialog.Update(&ui.layout.subActorDialogAnchor)
	ui.battleEnemyDisplay.Update(&ui.layout.enemyDisplayPosition)
	ui.battleSelectWindow.Update(&ui.layout.selectAnchor)
	ui.targetSelectWindow.Update(&ui.layout.selectAnchor)

	ui.input.Update()
	ui.effectManager.Update()
}

func updateBattleShake(shake *shake.EmitShake) *drawing.Vector {
	shake.Update()
	return shake.Delta()
}

func (ui *battleUI) Draw(drawFunc drawing.DrawFunc) {
	ui.messageWindow.Draw(drawFunc)
	ui.battleSelectWindow.Draw(drawFunc)
	ui.targetSelectWindow.Draw(drawFunc)
	ui.battleEnemyDisplay.Draw(drawFunc)
	ui.subActorDisplay.Draw(drawFunc)
	ui.subActorDialog.Draw(drawFunc)
	ui.actorDisplay.Draw(drawFunc)
	ui.effectManager.Draw(drawFunc)
}

func computeBattleUILayout(ui *battleUI, delta *drawing.Vector) battleUILayout {
	actorAnchor := ui.actorDisplay.GetMainCharacterTopLeftPosition()
	subActorAnchor := ui.subActorDisplay.GetTopCenterPosition()

	return battleUILayout{
		messageWindowPosition: drawing.Vector{X: delta.X, Y: delta.Y},
		actorDisplayPosition:  drawing.Vector{X: delta.X, Y: delta.Y + 288},
		subActorDisplayPos:    drawing.Vector{X: delta.X + 384, Y: delta.Y + 288},
		enemyDisplayPosition:  drawing.Vector{X: delta.X + 192, Y: delta.Y + 144},
		subActorDialogAnchor:  drawing.Vector{X: subActorAnchor.X + delta.X, Y: subActorAnchor.Y + delta.Y},
		selectAnchor:          drawing.Vector{X: actorAnchor.X + delta.X, Y: actorAnchor.Y + delta.Y},
	}
}
