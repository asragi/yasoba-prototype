package widgets

import (
	"errors"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type WindowOption struct {
	Image            *ebiten.Image
	CornerSize       int
	RelativePosition *core.Vector
	Size             *core.Vector
	Depth            core.Depth
	Pivot            *core.Pivot
	Padding          *core.Vector
}

func (o *WindowOption) Validation() error {
	if o.Image == nil {
		return errors.New("image is required")
	}
	if o.CornerSize <= 0 {
		return errors.New("corner size must be greater than 0")
	}
	if o.Size == nil {
		return errors.New("size is required")
	}
	if o.Depth == core.Zero {
		return errors.New("depth is required")
	}
	return nil
}

type Window struct {
	GetContentUpperLeft func() *core.Vector
	MoveTo              func(relativePosition *core.Vector)
	Update              func(parentPosition *core.Vector)
	Draw                func(core.DrawFunc)
}

func NewWindow(option *WindowOption) *Window {
	var parentPosition *core.Vector
	relativePosition := option.RelativePosition

	if err := option.Validation(); err != nil {
		panic(err)
	}
	img := option.Image
	type rect struct {
		x0 int
		y0 int
		x1 int
		y1 int
	}
	pivotDiff := option.Pivot.ApplyToSize(option.Size)
	textureWidth := img.Bounds().Dx()
	textureHeight := img.Bounds().Dy()
	corners := []*rect{
		{0, 0, option.CornerSize, option.CornerSize},
		{img.Bounds().Dx() - option.CornerSize, 0, img.Bounds().Dx(), option.CornerSize},
		{0, img.Bounds().Dy() - option.CornerSize, option.CornerSize, img.Bounds().Dy()},
		{
			img.Bounds().Dx() - option.CornerSize,
			img.Bounds().Dy() - option.CornerSize,
			img.Bounds().Dx(),
			img.Bounds().Dy(),
		},
	}
	sides := []*rect{
		{option.CornerSize, 0, img.Bounds().Dx() - option.CornerSize, option.CornerSize},
		{
			img.Bounds().Dx() - option.CornerSize,
			option.CornerSize,
			img.Bounds().Dx(),
			img.Bounds().Dy() - option.CornerSize,
		},
		{0, option.CornerSize, option.CornerSize, img.Bounds().Dy() - option.CornerSize},
		{
			option.CornerSize,
			img.Bounds().Dy() - option.CornerSize,
			img.Bounds().Dx() - option.CornerSize,
			img.Bounds().Dy(),
		},
	}
	cornerPosition := []*core.Vector{
		{0, 0},
		{option.Size.X - float64(option.CornerSize), 0},
		{0, option.Size.Y - float64(option.CornerSize)},
		{option.Size.X - float64(option.CornerSize), option.Size.Y - float64(option.CornerSize)},
	}
	sidePosition := []*core.Vector{
		{float64(option.CornerSize), 0},
		{option.Size.X - float64(option.CornerSize), float64(option.CornerSize)},
		{0, float64(option.CornerSize)},
		{float64(option.CornerSize), option.Size.Y - float64(option.CornerSize)},
	}
	sideXSize := float64(textureWidth - option.CornerSize*2)
	targetXSize := option.Size.X - float64(option.CornerSize*2)
	sideYSize := float64(textureHeight - option.CornerSize*2)
	targetYSize := option.Size.Y - float64(option.CornerSize*2)
	sideScales := []*core.Vector{
		{targetXSize / sideXSize, 1},
		{1, targetYSize / sideYSize},
		{1, targetYSize / sideYSize},
		{targetXSize / sideXSize, 1},
	}

	moveTo := func(passedPosition *core.Vector) {
		relativePosition = passedPosition
	}

	update := func(passedPosition *core.Vector) {
		parentPosition = passedPosition
	}

	drawWindow := func(drawing core.DrawFunc) {
		for i, v := range corners {
			op := &ebiten.DrawImageOptions{}
			x := cornerPosition[i].X + relativePosition.X - pivotDiff.X + parentPosition.X
			y := cornerPosition[i].Y + relativePosition.Y - pivotDiff.Y + parentPosition.Y
			op.GeoM.Translate(x, y)
			subImage := img.SubImage(image.Rect(v.x0, v.y0, v.x1, v.y1)).(*ebiten.Image)
			drawing(
				func(screen *ebiten.Image) {
					screen.DrawImage(subImage, op)
				}, option.Depth,
			)
		}
		for i, v := range sides {
			op := &ebiten.DrawImageOptions{}
			x := sidePosition[i].X + relativePosition.X - pivotDiff.X + parentPosition.X
			y := sidePosition[i].Y + relativePosition.Y - pivotDiff.Y + parentPosition.Y
			op.GeoM.Scale(sideScales[i].X, sideScales[i].Y)
			op.GeoM.Translate(x, y)
			subImage := img.SubImage(image.Rect(v.x0, v.y0, v.x1, v.y1)).(*ebiten.Image)
			drawing(
				func(screen *ebiten.Image) {
					screen.DrawImage(subImage, op)
				}, option.Depth,
			)
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(
			targetXSize/float64(textureWidth-option.CornerSize*2),
			targetYSize/float64(textureHeight-option.CornerSize*2),
		)
		op.GeoM.Translate(
			relativePosition.X-pivotDiff.X+parentPosition.X+float64(option.CornerSize),
			relativePosition.Y-pivotDiff.Y+parentPosition.Y+float64(option.CornerSize),
		)
		subImage := img.SubImage(
			image.Rect(
				option.CornerSize,
				option.CornerSize,
				textureWidth-option.CornerSize,
				textureHeight-option.CornerSize,
			),
		).(*ebiten.Image)
		drawing(
			func(screen *ebiten.Image) {
				screen.DrawImage(subImage, op)
			}, option.Depth,
		)
	}

	getContentPosition := func() *core.Vector {
		return &core.Vector{
			X: relativePosition.X - pivotDiff.X + parentPosition.X + option.Padding.X,
			Y: relativePosition.Y - pivotDiff.Y + parentPosition.Y + option.Padding.Y,
		}
	}

	return &Window{
		GetContentUpperLeft: getContentPosition,
		MoveTo:              moveTo,
		Update:              update,
		Draw:                drawWindow,
	}
}
