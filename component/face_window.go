package component

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type FaceWindow struct {
	face   *widget.Image
	window *widget.Window
}

type NewFaceWindowFunc func(
	*frontend.Vector,
	frontend.Depth,
	*frontend.Pivot,
	frontend.TextureId,
) *FaceWindow

func (f *FaceWindow) Update(parentPosition *frontend.Vector) {
	f.window.Update(parentPosition)
	f.face.Update(f.window.GetPositionCenter())
}

func (f *FaceWindow) Draw(drawFunc frontend.DrawFunc) {
	f.window.Draw(drawFunc)
	f.face.Draw(drawFunc)
}

func (f *FaceWindow) GetTopLeftPosition() *frontend.Vector {
	return f.window.GetPositionUpperLeft()
}

func (f *FaceWindow) GetBottomRightPosition() *frontend.Vector {
	return f.window.GetPositionLowerRight()
}

func (f *FaceWindow) GetCenterPosition() *frontend.Vector {
	return f.window.GetPositionCenter()
}

func StandByNewFaceWindow(resource *frontend.ResourceManager) NewFaceWindowFunc {
	return func(
		relativePosition *frontend.Vector,
		depth frontend.Depth,
		pivot *frontend.Pivot,
		texture frontend.TextureId,
	) *FaceWindow {
		const padding = 6
		face := widget.NewImage(
			frontend.VectorZero,
			frontend.PivotCenter,
			depth,
			resource.GetTexture(texture),
		)
		face.SetScaleBySize(&frontend.Vector{X: 74, Y: 74})
		window := widget.NewWindow(
			&widget.WindowOption{
				Image:            resource.GetTexture(frontend.TextureWindow),
				CornerSize:       6,
				RelativePosition: relativePosition,
				Size:             face.Size().Add(&frontend.Vector{X: padding, Y: padding}),
				Depth:            depth,
				Pivot:            pivot,
			},
		)
		return &FaceWindow{
			face:   face,
			window: window,
		}
	}
}

type BattleEmotionType int

const (
	BattleEmotionNormal BattleEmotionType = iota
	BattleEmotionDamage
)
