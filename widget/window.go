package widget

import (
	"errors"
	"image"
	"math"

	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
)

type windowRect struct {
	x0 int
	y0 int
	x1 int
	y1 int
}

type WindowInterface interface {
	PositionUpdater
	Drawer
	Size() *drawing.Vector
	GetPositionUpperLeft() *drawing.Vector
	GetPositionTopCenter() *drawing.Vector
	GetPositionCenter() *drawing.Vector
	GetPositionLowerRight() *drawing.Vector
	GetPadding() *drawing.Vector
	SetSize(size *drawing.Vector)
}

type Window struct {
	screenWidth      int
	screenHeight     int
	image            drawing.Image
	relativePosition *drawing.Vector
	parentPosition   *drawing.Vector
	size             *drawing.Vector
	pivot            *drawing.Pivot
	corners          []*windowRect
	sides            []*windowRect
	cornerPosition   []*drawing.Vector
	sidePosition     []*drawing.Vector
	sideScale        []*drawing.Vector
	cornerSize       int
	depth            drawing.Depth
	padding          *drawing.Vector
}

// Windowの枠を含めた全体のサイズを指定する
func (w *Window) SetSize(size *drawing.Vector) {
	w.size = size
	w.cornerPosition = calculateCornerPosition(size, float64(w.cornerSize))
	w.sidePosition = calculateSidePosition(size, float64(w.cornerSize))
	w.sideScale = calculateSideScale(
		size,
		float64(w.cornerSize),
		float64(w.image.Bounds().Dx()),
		float64(w.image.Bounds().Dy()),
	)
}

func (w *Window) GetPositionUpperLeft() *drawing.Vector {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	return &drawing.Vector{
		X: w.relativePosition.X + w.parentPosition.X - pivotDiff.X,
		Y: w.relativePosition.Y + w.parentPosition.Y - pivotDiff.Y,
	}
}

func (w *Window) GetPositionTopCenter() *drawing.Vector {
	return w.GetPositionUpperLeft().Add(&drawing.Vector{X: w.size.X / 2, Y: 0})
}

func (w *Window) GetPositionCenter() *drawing.Vector {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	return &drawing.Vector{
		X: w.relativePosition.X + w.parentPosition.X - pivotDiff.X + w.size.X/2,
		Y: w.relativePosition.Y + w.parentPosition.Y - pivotDiff.Y + w.size.Y/2,
	}
}

func (w *Window) GetPositionLowerRight() *drawing.Vector {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	return &drawing.Vector{
		X: w.relativePosition.X + w.parentPosition.X - pivotDiff.X + w.size.X,
		Y: w.relativePosition.Y + w.parentPosition.Y - pivotDiff.Y + w.size.Y,
	}
}

// 画面からはみ出す分を計算し修正に必要なVectorを返す
func (w *Window) calculateInWindowPosition(passedParentPosition *drawing.Vector) *drawing.Vector {
	// TODO: 左や上に飛び出す場合を想定していない
	// TODO: フラグではみ出しを許容するかどうかを変えたい
	pivotDiff := w.pivot.ApplyToSize(w.size)
	xDiff := math.Max(passedParentPosition.X+w.relativePosition.X+w.size.X-pivotDiff.X-float64(w.screenWidth), 0)
	yDiff := math.Max(passedParentPosition.Y+w.relativePosition.Y+w.size.Y-pivotDiff.Y-float64(w.screenHeight), 0)
	return &drawing.Vector{X: xDiff, Y: yDiff}
}

func (w *Window) GetPadding() *drawing.Vector {
	return w.padding
}

func (w *Window) Update(passedPosition *drawing.Vector) {
	inWindowPosition := w.calculateInWindowPosition(passedPosition)
	w.parentPosition = passedPosition.Sub(inWindowPosition)
}

func (w *Window) Draw(drawFunc drawing.DrawFunc) {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	textureWidth := w.image.Bounds().Dx()
	textureHeight := w.image.Bounds().Dy()
	targetXSize := w.size.X - float64(w.cornerSize*2)
	targetYSize := w.size.Y - float64(w.cornerSize*2)
	for i, v := range w.corners {
		op := &drawing.DrawOptions{}
		x := w.cornerPosition[i].X + w.relativePosition.X - pivotDiff.X + w.parentPosition.X
		y := w.cornerPosition[i].Y + w.relativePosition.Y - pivotDiff.Y + w.parentPosition.Y
		op.Translate(x, y)
		subImage := w.image.SubImage(image.Rect(v.x0, v.y0, v.x1, v.y1))
		drawFunc(
			func(screen drawing.Image) {
				screen.DrawImage(subImage, op)
			}, w.depth,
		)
	}
	for i, v := range w.sides {
		op := drawing.NewDrawOptions()
		x := w.sidePosition[i].X + w.relativePosition.X - pivotDiff.X + w.parentPosition.X
		y := w.sidePosition[i].Y + w.relativePosition.Y - pivotDiff.Y + w.parentPosition.Y
		op.SetScale(w.sideScale[i].X, w.sideScale[i].Y)
		op.Translate(x, y)
		subImage := w.image.SubImage(image.Rect(v.x0, v.y0, v.x1, v.y1))
		drawFunc(
			func(screen drawing.Image) {
				screen.DrawImage(subImage, op)
			}, w.depth,
		)
	}
	op := drawing.NewDrawOptions()
	op.SetScale(
		targetXSize/float64(textureWidth-w.cornerSize*2),
		targetYSize/float64(textureHeight-w.cornerSize*2),
	)
	op.Translate(
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
	)
	drawFunc(
		func(screen drawing.Image) {
			screen.DrawImage(subImage, op)
		}, w.depth,
	)
}

func (w *Window) Size() *drawing.Vector {
	return w.size
}

func calculateCornerPosition(
	size *drawing.Vector,
	cornerSize float64,
) []*drawing.Vector {
	return []*drawing.Vector{
		{0, 0},
		{size.X - cornerSize, 0},
		{0, size.Y - cornerSize},
		{size.X - cornerSize, size.Y - cornerSize},
	}
}

func calculateSidePosition(
	size *drawing.Vector,
	cornerSize float64,
) []*drawing.Vector {
	return []*drawing.Vector{
		{cornerSize, 0},
		{size.X - cornerSize, cornerSize},
		{0, cornerSize},
		{cornerSize, size.Y - cornerSize},
	}
}

func calculateSideScale(
	size *drawing.Vector,
	cornerSize float64,
	textureWidth float64,
	textureHeight float64,
) []*drawing.Vector {
	sideXSize := textureWidth - cornerSize*2
	targetXSize := size.X - cornerSize*2
	sideYSize := textureHeight - cornerSize*2
	targetYSize := size.Y - cornerSize*2
	return []*drawing.Vector{
		{targetXSize / sideXSize, 1},
		{1, targetYSize / sideYSize},
		{1, targetYSize / sideYSize},
		{targetXSize / sideXSize, 1},
	}
}

type WindowOption struct {
	Texture          frontend.TextureId
	CornerSize       int
	RelativePosition *drawing.Vector
	Size             *drawing.Vector
	Depth            drawing.Depth
	Pivot            *drawing.Pivot
	Padding          *drawing.Vector
}

type NewWindowFunc func(*WindowOption) WindowInterface

func (o *WindowOption) Validation() error {
	if o.CornerSize <= 0 {
		return errors.New("corner Size must be greater than 0")
	}
	if o.Size == nil {
		return errors.New("size is required")
	}
	if o.Depth == drawing.Zero {
		return errors.New("depth is required")
	}
	return nil
}

type GetImageFunc func(frontend.TextureId) drawing.Image

func CreateNewWindow(
	getImage GetImageFunc,
	screenWidth, screenHeight int,
) NewWindowFunc {
	return func(option *WindowOption) WindowInterface {
		parentPosition := drawing.VectorZero
		relativePosition := option.RelativePosition

		if err := option.Validation(); err != nil {
			panic(err)
		}
		img := getImage(option.Texture)
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
		cornerPosition := calculateCornerPosition(option.Size, float64(option.CornerSize))
		sidePosition := calculateSidePosition(option.Size, float64(option.CornerSize))
		sideScales := calculateSideScale(
			option.Size,
			float64(option.CornerSize),
			float64(textureWidth),
			float64(textureHeight),
		)

		return &Window{
			screenWidth:      screenWidth,
			screenHeight:     screenHeight,
			image:            img,
			relativePosition: relativePosition,
			parentPosition:   parentPosition,
			size:             option.Size,
			pivot:            option.Pivot,
			padding:          option.Padding,
			corners:          corners,
			sides:            sides,
			cornerPosition:   cornerPosition,
			sidePosition:     sidePosition,
			sideScale:        sideScales,
			cornerSize:       option.CornerSize,
			depth:            option.Depth,
		}
	}
}
