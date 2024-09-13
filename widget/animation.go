package widget

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	image *Image
}

func NewAnimation(
	relativePosition *frontend.Vector,
	pivot *frontend.Pivot,
	depth frontend.Depth,
	image *ebiten.Image,
) *Animation {
	return &Animation{
		image: NewImage(
			relativePosition,
			pivot,
			depth,
			image,
		),
	}
}

func (a *Animation) Update(passedPosition *frontend.Vector) {
	a.image.Update(passedPosition)
}

func (a *Animation) Draw(drawFunc frontend.DrawFunc) {
	a.Draw(drawFunc)
}
