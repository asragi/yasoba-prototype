package component

import "github.com/asragi/yasoba-prototype/frontend"

type BattleActorDisplay struct {
	mainActorDisplay *FaceWindow
	subActorDisplay  *FaceWindow
}

func (d *BattleActorDisplay) Update(
	bottomLeftPosition,
	bottomRightPosition *frontend.Vector,
) {
	d.mainActorDisplay.Update(bottomLeftPosition)
	d.subActorDisplay.Update(bottomRightPosition)
}

func (d *BattleActorDisplay) Draw(
	drawFunc frontend.DrawFunc,
) {
	d.mainActorDisplay.Draw(drawFunc)
	d.subActorDisplay.Draw(drawFunc)
}

func (d *BattleActorDisplay) GetMainCharacterPosition() *frontend.Vector {
	return d.mainActorDisplay.GetCenterPosition()
}

func (d *BattleActorDisplay) GetSubCharacterPosition() *frontend.Vector {
	return d.subActorDisplay.GetCenterPosition()
}

func (d *BattleActorDisplay) GetMainCharacterTopLeftPosition() *frontend.Vector {
	return d.mainActorDisplay.GetTopLeftPosition()
}

type NewBattleActorDisplayFunc func() *BattleActorDisplay

func CreateNewBattleActorDisplay(
	newFaceWindow NewFaceWindowFunc,
) NewBattleActorDisplayFunc {
	return func() *BattleActorDisplay {
		return &BattleActorDisplay{
			mainActorDisplay: newFaceWindow(
				&frontend.Vector{X: 0, Y: 0},
				frontend.DepthPlayer,
				frontend.PivotBottomLeft,
				frontend.TextureFaceLuneNormal,
			),
			subActorDisplay: newFaceWindow(
				&frontend.Vector{X: 0, Y: 0},
				frontend.DepthPlayer,
				frontend.PivotBottomRight,
				frontend.TextureFaceSunnyNormal,
			),
		}
	}
}
