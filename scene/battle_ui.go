package scene

import (
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/toolkit/input"
	"github.com/asragi/yasoba-prototype/view/battle/actor"
	"github.com/asragi/yasoba-prototype/view/battle/dialogue"
	"github.com/asragi/yasoba-prototype/view/battle/enemy"
	"github.com/asragi/yasoba-prototype/view/battle/window"
	"github.com/asragi/yasoba-prototype/view/common/constant"
	"github.com/asragi/yasoba-prototype/view/common/message"
	"github.com/asragi/yasoba-prototype/view/common/selection"
	"github.com/asragi/yasoba-prototype/view/common/shake"
	transitionview "github.com/asragi/yasoba-prototype/view/common/transition"
	"github.com/asragi/yasoba-prototype/widget"
)

type battleUI struct {
	messageWindow      *message.MessageWindow
	battleSelectWindow *window.BattleSelectWindow
	actorDisplay       *actor.BattleActorDisplay
	subActorDisplay    *actor.BattleSubActorDisplay
	subActorDialog     *dialogue.BattlePartnerDialogue
	targetSelectWindow *selection.SelectWindow
	input              input.Manager
	battleEnemyDisplay *enemy.BattleEnemyDisplay
	effectManager      *widget.EffectManager
	shake              *shake.EmitShake
	layout             battleUILayout
	transition         *transitionview.Transition
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
	ui.transition.Update()
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
	ui.transition.Draw(drawFunc)
}

func computeBattleUILayout(ui *battleUI, delta *drawing.Vector) battleUILayout {
	actorAnchor := ui.actorDisplay.GetMainCharacterTopLeftPosition()
	subActorAnchor := ui.subActorDisplay.GetTopCenterPosition()
	screenWidth := float64(constant.GameWidth)
	screenHeight := float64(constant.GameHeight)

	return battleUILayout{
		messageWindowPosition: drawing.Vector{X: delta.X, Y: delta.Y},
		actorDisplayPosition:  drawing.Vector{X: delta.X, Y: delta.Y + screenHeight},
		subActorDisplayPos:    drawing.Vector{X: delta.X + screenWidth, Y: delta.Y + screenHeight},
		// TODO: enemy positionはbattle settingから得る
		enemyDisplayPosition: drawing.Vector{X: delta.X + screenWidth/2, Y: delta.Y + screenHeight/2},
		subActorDialogAnchor: drawing.Vector{X: subActorAnchor.X + delta.X, Y: subActorAnchor.Y + delta.Y},
		selectAnchor:         drawing.Vector{X: actorAnchor.X + delta.X, Y: actorAnchor.Y + delta.Y},
	}
}
