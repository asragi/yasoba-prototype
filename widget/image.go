package widget

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
)

type Image struct {
	relativePosition *frontend.Vector
	parentPosition   *frontend.Vector
	pivot            *frontend.Pivot
	image            *ebiten.Image
	depth            frontend.Depth
}

func (i *Image) Update(passedPosition *frontend.Vector) {
	i.parentPosition = passedPosition
}

func (i *Image) Draw(drawFunc frontend.DrawFunc) {
	if i.parentPosition == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	pivotModification := i.pivot.ApplyToSize(i.Size())
	op.GeoM.Translate(
		i.parentPosition.X+i.relativePosition.X-pivotModification.X,
		i.parentPosition.Y+i.relativePosition.Y-pivotModification.Y,
	)
	drawFunc(
		func(screen *ebiten.Image) {
			screen.DrawImage(i.image, op)
		}, i.depth,
	)
}

func (i *Image) Size() *frontend.Vector {
	return &frontend.Vector{
		X: float64(i.image.Bounds().Dx()),
		Y: float64(i.image.Bounds().Dy()),
	}
}

func (i *Image) SetRelativePosition(position *frontend.Vector) {
	i.relativePosition = position
}

func NewImage(
	relativePosition *frontend.Vector,
	pivot *frontend.Pivot,
	depth frontend.Depth,
	image *ebiten.Image,
) *Image {
	return &Image{
		relativePosition: relativePosition,
		pivot:            pivot,
		image:            image,
		depth:            depth,
	}
}
