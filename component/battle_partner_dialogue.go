package component

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattlePartnerDialogue struct {
	window widget.WindowInterface
}

type NewBattlePartnerDialogueFunc func() *BattlePartnerDialogue

func (d *BattlePartnerDialogue) Draw(drawFunc frontend.DrawFunc) {
	d.window.Draw(drawFunc)
}

func (d *BattlePartnerDialogue) Update(parentPosition *frontend.Vector) {
	d.window.Update(parentPosition)
}

func CreateNewBattlePartnerDialogue(newWindow widget.NewWindowFunc) NewBattlePartnerDialogueFunc {
	texture := frontend.TextureWindow
	cornerSize := 6
	padding := &frontend.Vector{X: 16, Y: 8}

	return func() *BattlePartnerDialogue {
		return &BattlePartnerDialogue{
			window: newWindow(
				&widget.WindowOption{
					Texture:          texture,
					CornerSize:       cornerSize,
					RelativePosition: frontend.VectorZero,
					Size:             frontend.VectorOne,
					Depth:            frontend.DepthWindow,
					Pivot:            frontend.PivotBottomRight,
					Padding:          padding,
				},
			),
		}
	}
}
