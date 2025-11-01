package dialogue

import (
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/message"
)

type PartnerDialogueMessageWindow interface {
	Open()
	Close()
	SetText(string, bool)
	FitToMessage()
	Draw(drawing.DrawFunc)
	Update(parentPosition *drawing.Vector)
	IsTextEnd() bool
}

type makeMessageWindowFunc func(
	*drawing.Vector,
	*drawing.Vector,
	drawing.Depth,
	*drawing.Pivot,
) PartnerDialogueMessageWindow

type BattlePartnerDialogue struct {
	window        PartnerDialogueMessageWindow
	newWindowFunc makeMessageWindowFunc
	textWait      int
}

const defaultTextWait = 60

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

func (d *BattlePartnerDialogue) IsTextEnd() bool {
	if d.window == nil {
		return true
	}
	isEnd := d.window.IsTextEnd()
	if !isEnd {
		d.textWait = defaultTextWait
	}
	d.textWait--
	return d.textWait <= 0
}

func (d *BattlePartnerDialogue) SetText(textString string, displayAll bool) {
	d.window.SetText(textString, displayAll)
	d.window.FitToMessage()
}

func (d *BattlePartnerDialogue) Draw(drawFunc drawing.DrawFunc) {
	if d.window == nil {
		return
	}
	d.window.Draw(drawFunc)
}

func (d *BattlePartnerDialogue) Update(parentPosition *drawing.Vector) {
	if d.window == nil {
		return
	}
	d.window.Update(parentPosition)
}

func (d *BattlePartnerDialogue) newWindow() PartnerDialogueMessageWindow {
	window := d.newWindowFunc(
		drawing.VectorZero,
		drawing.VectorZero,
		drawing.DepthWindow,
		drawing.PivotBottomCenter,
	)
	window.FitToMessage()
	return window
}

func CreateNewBattlePartnerDialogue(newWindow message.NewMessageWindowFunc) NewBattlePartnerDialogueFunc {
	return func() *BattlePartnerDialogue {
		return &BattlePartnerDialogue{
			newWindowFunc: func(
				relativePosition *drawing.Vector,
				size *drawing.Vector,
				depth drawing.Depth,
				pivot *drawing.Pivot,
			) PartnerDialogueMessageWindow {
				return newWindow(relativePosition, size, depth, pivot)
			},
		}
	}
}
