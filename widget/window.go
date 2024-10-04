package widget

import (
	"errors"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type WindowOption struct {
	Image            *ebiten.Image
	CornerSize       int
	RelativePosition *frontend.Vector
	Size             *frontend.Vector
	Depth            frontend.Depth
	Pivot            *frontend.Pivot
	Padding          *frontend.Vector
}

func (o *WindowOption) Validation() error {
	if o.Image == nil {
		return errors.New("image is required")
	}
	if o.CornerSize <= 0 {
		return errors.New("corner Size must be greater than 0")
	}
	if o.Size == nil {
		return errors.New("size is required")
	}
	if o.Depth == frontend.Zero {
		return errors.New("depth is required")
	}
	return nil
}

type windowRect struct {
	x0 int
	y0 int
	x1 int
	y1 int
}

type Window struct {
	image               *ebiten.Image
	relativePosition    *frontend.Vector
	parentPosition      *frontend.Vector
	size                *frontend.Vector
	pivot               *frontend.Pivot
	GetContentUpperLeft func() *frontend.Vector
	corners             []*windowRect
	sides               []*windowRect
	cornerPosition      []*frontend.Vector
	sidePosition        []*frontend.Vector
	sideScale           []*frontend.Vector
	cornerSize          int
	depth               frontend.Depth
}

func (w *Window) GetPositionUpperLeft() *frontend.Vector {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	return &frontend.Vector{
		X: w.relativePosition.X + w.parentPosition.X - pivotDiff.X,
		Y: w.relativePosition.Y + w.parentPosition.Y - pivotDiff.Y,
	}
}

func (w *Window) GetPositionCenter() *frontend.Vector {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	return &frontend.Vector{
		X: w.relativePosition.X + w.parentPosition.X - pivotDiff.X + w.size.X/2,
		Y: w.relativePosition.Y + w.parentPosition.Y - pivotDiff.Y + w.size.Y/2,
	}
}

func (w *Window) GetPositionLowerRight() *frontend.Vector {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	return &frontend.Vector{
		X: w.relativePosition.X + w.parentPosition.X - pivotDiff.X + w.size.X,
		Y: w.relativePosition.Y + w.parentPosition.Y - pivotDiff.Y + w.size.Y,
	}
}

func (w *Window) Update(passedPosition *frontend.Vector) {
	w.parentPosition = passedPosition
}

func (w *Window) Draw(drawFunc frontend.DrawFunc) {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	textureWidth := w.image.Bounds().Dx()
	textureHeight := w.image.Bounds().Dy()
	targetXSize := w.size.X - float64(w.cornerSize*2)
	targetYSize := w.size.Y - float64(w.cornerSize*2)
	for i, v := range w.corners {
		op := &ebiten.DrawImageOptions{}
		x := w.cornerPosition[i].X + w.relativePosition.X - pivotDiff.X + w.parentPosition.X
		y := w.cornerPosition[i].Y + w.relativePosition.Y - pivotDiff.Y + w.parentPosition.Y
		op.GeoM.Translate(x, y)
		subImage := w.image.SubImage(image.Rect(v.x0, v.y0, v.x1, v.y1)).(*ebiten.Image)
		drawFunc(
			func(screen *ebiten.Image) {
				screen.DrawImage(subImage, op)
			}, w.depth,
		)
	}
	for i, v := range w.sides {
		op := &ebiten.DrawImageOptions{}
		x := w.sidePosition[i].X + w.relativePosition.X - pivotDiff.X + w.parentPosition.X
		y := w.sidePosition[i].Y + w.relativePosition.Y - pivotDiff.Y + w.parentPosition.Y
		op.GeoM.Scale(w.sideScale[i].X, w.sideScale[i].Y)
		op.GeoM.Translate(x, y)
		subImage := w.image.SubImage(image.Rect(v.x0, v.y0, v.x1, v.y1)).(*ebiten.Image)
		drawFunc(
			func(screen *ebiten.Image) {
				screen.DrawImage(subImage, op)
			}, w.depth,
		)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(
		targetXSize/float64(textureWidth-w.cornerSize*2),
		targetYSize/float64(textureHeight-w.cornerSize*2),
	)
	op.GeoM.Translate(
		w.relativePosition.X-pivotDiff.X+w.parentPosition.X+float64(w.cornerSize),
		w.relativePosition.Y-pivotDiff.Y+w.parentPosition.Y+float64(w.cornerSize),
	)
	subImage := w.image.SubImage(
		image.Rect(
			w.cornerSize,
			w.cornerSize,
			textureWidth-w.cornerSize,
			textureHeight-w.cornerSize,
		),
	).(*ebiten.Image)
	drawFunc(
		func(screen *ebiten.Image) {
			screen.DrawImage(subImage, op)
		}, w.depth,
	)
}

func (w *Window) Size() *frontend.Vector {
	return w.size
}

func NewWindow(option *WindowOption) *Window {
	parentPosition := frontend.VectorZero
	relativePosition := option.RelativePosition

	if err := option.Validation(); err != nil {
		panic(err)
	}
	img := option.Image
	textureWidth := img.Bounds().Dx()
	textureHeight := img.Bounds().Dy()
	corners := []*windowRect{
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
	sides := []*windowRect{
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
	cornerPosition := []*frontend.Vector{
		{0, 0},
		{option.Size.X - float64(option.CornerSize), 0},
		{0, option.Size.Y - float64(option.CornerSize)},
		{option.Size.X - float64(option.CornerSize), option.Size.Y - float64(option.CornerSize)},
	}
	sidePosition := []*frontend.Vector{
		{float64(option.CornerSize), 0},
		{option.Size.X - float64(option.CornerSize), float64(option.CornerSize)},
		{0, float64(option.CornerSize)},
		{float64(option.CornerSize), option.Size.Y - float64(option.CornerSize)},
	}
	sideXSize := float64(textureWidth - option.CornerSize*2)
	targetXSize := option.Size.X - float64(option.CornerSize*2)
	sideYSize := float64(textureHeight - option.CornerSize*2)
	targetYSize := option.Size.Y - float64(option.CornerSize*2)
	sideScales := []*frontend.Vector{
		{targetXSize / sideXSize, 1},
		{1, targetYSize / sideYSize},
		{1, targetYSize / sideYSize},
		{targetXSize / sideXSize, 1},
	}

	getContentPosition := func() *frontend.Vector {
		return option.Padding
		/*
			return &frontend.Vector{
				X: relativePosition.X - pivotDiff.X + parentPosition.X + option.Padding.X,
				Y: relativePosition.Y - pivotDiff.Y + parentPosition.Y + option.Padding.Y,
			}
		*/
	}

	return &Window{
		image:               img,
		relativePosition:    relativePosition,
		parentPosition:      parentPosition,
		size:                option.Size,
		pivot:               option.Pivot,
		GetContentUpperLeft: getContentPosition,
		corners:             corners,
		sides:               sides,
		cornerPosition:      cornerPosition,
		sidePosition:        sidePosition,
		sideScale:           sideScales,
		cornerSize:          option.CornerSize,
		depth:               option.Depth,
	}
}
