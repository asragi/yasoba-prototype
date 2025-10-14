package dialogue

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/frontend"
)

type PartnerDialogueMessageWindow interface {
	Open()
	Close()
	SetText(string, bool)
	FitToMessage()
	Draw(frontend.DrawFunc)
	Update(parentPosition *frontend.Vector)
	IsTextEnd() bool
}

type makeMessageWindowFunc func(
	*frontend.Vector,
	*frontend.Vector,
	frontend.Depth,
	*frontend.Pivot,
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
		frontend.PivotBottomCenter,
	)
	window.FitToMessage()
	return window
}

func CreateNewBattlePartnerDialogue(newWindow component.NewMessageWindowFunc) NewBattlePartnerDialogueFunc {
	return func() *BattlePartnerDialogue {
		return &BattlePartnerDialogue{
			newWindowFunc: func(
				relativePosition *frontend.Vector,
				size *frontend.Vector,
				depth frontend.Depth,
				pivot *frontend.Pivot,
			) PartnerDialogueMessageWindow {
				return newWindow(relativePosition, size, depth, pivot)
			},
		}
	}
}
