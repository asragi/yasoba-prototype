package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type battleUI struct {
	messageWindow      *component.MessageWindow
	battleSelectWindow *component.BattleSelectWindow
	actorDisplay       *component.BattleActorDisplay
	subActorDisplay    *component.BattleSubActorDisplay
	subActorDialog     *component.BattlePartnerDialogue
	targetSelectWindow *component.SelectWindow
	input              frontend.InputManager
	battleEnemyDisplay *component.BattleEnemyDisplay
	effectManager      *widget.EffectManager
	shake              *frontend.EmitShake
	layout             battleUILayout
}

type battleUILayout struct {
	messageWindowPosition frontend.Vector
	actorDisplayPosition  frontend.Vector
	subActorDisplayPos    frontend.Vector
	enemyDisplayPosition  frontend.Vector
	subActorDialogAnchor  frontend.Vector
	selectAnchor          frontend.Vector
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

func updateBattleShake(shake *frontend.EmitShake) *frontend.Vector {
	shake.Update()
	return shake.Delta()
}

func (ui *battleUI) Draw(drawFunc frontend.DrawFunc) {
	ui.messageWindow.Draw(drawFunc)
	ui.battleSelectWindow.Draw(drawFunc)
	ui.targetSelectWindow.Draw(drawFunc)
	ui.battleEnemyDisplay.Draw(drawFunc)
	ui.subActorDisplay.Draw(drawFunc)
	ui.subActorDialog.Draw(drawFunc)
	ui.actorDisplay.Draw(drawFunc)
	ui.effectManager.Draw(drawFunc)
}

func computeBattleUILayout(ui *battleUI, delta *frontend.Vector) battleUILayout {
	actorAnchor := ui.actorDisplay.GetMainCharacterTopLeftPosition()
	subActorAnchor := ui.subActorDisplay.GetTopCenterPosition()

	return battleUILayout{
		messageWindowPosition: frontend.Vector{X: delta.X, Y: delta.Y},
		actorDisplayPosition:  frontend.Vector{X: delta.X, Y: delta.Y + 288},
		subActorDisplayPos:    frontend.Vector{X: delta.X + 384, Y: delta.Y + 288},
		enemyDisplayPosition:  frontend.Vector{X: delta.X + 192, Y: delta.Y + 144},
		subActorDialogAnchor:  frontend.Vector{X: subActorAnchor.X + delta.X, Y: subActorAnchor.Y + delta.Y},
		selectAnchor:          frontend.Vector{X: actorAnchor.X + delta.X, Y: actorAnchor.Y + delta.Y},
	}
}
