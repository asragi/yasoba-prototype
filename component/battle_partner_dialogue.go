package component

import (
	"github.com/asragi/yasoba-prototype/frontend"
)

type PartnerDialogueMessageWindow interface {
	Open()
	Close()
	SetText(string, bool)
	FitToMessage()
	Draw(frontend.DrawFunc)
	Update(parentPosition *frontend.Vector)
}

type BattlePartnerDialogue struct {
	window        PartnerDialogueMessageWindow
	newWindowFunc NewMessageWindowFunc
}

type NewBattlePartnerDialogueFunc func() *BattlePartnerDialogue

func (d *BattlePartnerDialogue) Close() {
	if d.window == nil {
		return
	}
	d.window.Close()
}

func (d *BattlePartnerDialogue) Open() {
	d.window = d.newWindow()
	d.window.Open()
}

func (d *BattlePartnerDialogue) SetText(textString string, displayAll bool) {
	d.window.SetText(textString, displayAll)
	d.window.FitToMessage()
}

func (d *BattlePartnerDialogue) Draw(drawFunc frontend.DrawFunc) {
	if d.window == nil {
		return
	}
	d.window.Draw(drawFunc)
}

func (d *BattlePartnerDialogue) Update(parentPosition *frontend.Vector) {
	if d.window == nil {
		return
	}
	d.window.Update(parentPosition)
}

func (d *BattlePartnerDialogue) newWindow() PartnerDialogueMessageWindow {
	window := d.newWindowFunc(
		frontend.VectorZero,
		frontend.VectorZero,
		frontend.DepthWindow,
		frontend.PivotBottomRight,
	)
	window.FitToMessage()
	return window
}

func CreateNewBattlePartnerDialogue(newWindow NewMessageWindowFunc) NewBattlePartnerDialogueFunc {
	return func() *BattlePartnerDialogue {
		return &BattlePartnerDialogue{
			newWindowFunc: newWindow,
		}
	}
}
