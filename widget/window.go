package widget

import (
	"errors"
	"image"
	"math"

	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
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
	Size() *frontend.Vector
	GetPositionUpperLeft() *frontend.Vector
	GetPositionTopCenter() *frontend.Vector
	GetPositionCenter() *frontend.Vector
	GetPositionLowerRight() *frontend.Vector
	GetPadding() *frontend.Vector
	SetSize(size *frontend.Vector)
}

type Window struct {
	screenWidth      int
	screenHeight     int
	image            *ebiten.Image
	relativePosition *frontend.Vector
	parentPosition   *frontend.Vector
	size             *frontend.Vector
	pivot            *frontend.Pivot
	corners          []*windowRect
	sides            []*windowRect
	cornerPosition   []*frontend.Vector
	sidePosition     []*frontend.Vector
	sideScale        []*frontend.Vector
	cornerSize       int
	depth            frontend.Depth
	padding          *frontend.Vector
}

// Windowの枠を含めた全体のサイズを指定する
func (w *Window) SetSize(size *frontend.Vector) {
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

func (w *Window) GetPositionUpperLeft() *frontend.Vector {
	pivotDiff := w.pivot.ApplyToSize(w.size)
	return &frontend.Vector{
		X: w.relativePosition.X + w.parentPosition.X - pivotDiff.X,
		Y: w.relativePosition.Y + w.parentPosition.Y - pivotDiff.Y,
	}
}

func (w *Window) GetPositionTopCenter() *frontend.Vector {
	return w.GetPositionUpperLeft().Add(&frontend.Vector{X: w.size.X / 2, Y: 0})
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

// 画面からはみ出す分を計算し修正に必要なVectorを返す
func (w *Window) calculateInWindowPosition(passedParentPosition *frontend.Vector) *frontend.Vector {
	// TODO: 左や上に飛び出す場合を想定していない
	// TODO: フラグではみ出しを許容するかどうかを変えたい
	pivotDiff := w.pivot.ApplyToSize(w.size)
	xDiff := math.Max(passedParentPosition.X+w.relativePosition.X+w.size.X-pivotDiff.X-float64(w.screenWidth), 0)
	yDiff := math.Max(passedParentPosition.Y+w.relativePosition.Y+w.size.Y-pivotDiff.Y-float64(w.screenHeight), 0)
	return &frontend.Vector{X: xDiff, Y: yDiff}
}

func (w *Window) GetPadding() *frontend.Vector {
	return w.padding
}

func (w *Window) Update(passedPosition *frontend.Vector) {
	inWindowPosition := w.calculateInWindowPosition(passedPosition)
	w.parentPosition = passedPosition.Sub(inWindowPosition)
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

func calculateCornerPosition(
	size *frontend.Vector,
	cornerSize float64,
) []*frontend.Vector {
	return []*frontend.Vector{
		{0, 0},
		{size.X - cornerSize, 0},
		{0, size.Y - cornerSize},
		{size.X - cornerSize, size.Y - cornerSize},
	}
}

func calculateSidePosition(
	size *frontend.Vector,
	cornerSize float64,
) []*frontend.Vector {
	return []*frontend.Vector{
		{cornerSize, 0},
		{size.X - cornerSize, cornerSize},
		{0, cornerSize},
		{cornerSize, size.Y - cornerSize},
	}
}

func calculateSideScale(
	size *frontend.Vector,
	cornerSize float64,
	textureWidth float64,
	textureHeight float64,
) []*frontend.Vector {
	sideXSize := textureWidth - cornerSize*2
	targetXSize := size.X - cornerSize*2
	sideYSize := textureHeight - cornerSize*2
	targetYSize := size.Y - cornerSize*2
	return []*frontend.Vector{
		{targetXSize / sideXSize, 1},
		{1, targetYSize / sideYSize},
		{1, targetYSize / sideYSize},
		{targetXSize / sideXSize, 1},
	}
}

type WindowOption struct {
	Texture          frontend.TextureId
	CornerSize       int
	RelativePosition *frontend.Vector
	Size             *frontend.Vector
	Depth            frontend.Depth
	Pivot            *frontend.Pivot
	Padding          *frontend.Vector
}

type NewWindowFunc func(*WindowOption) WindowInterface

func (o *WindowOption) Validation() error {
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

func CreateNewWindow(
	resource *frontend.ResourceManager,
	screenWidth, screenHeight int,
) NewWindowFunc {
	return func(option *WindowOption) WindowInterface {
		parentPosition := frontend.VectorZero
		relativePosition := option.RelativePosition

		if err := option.Validation(); err != nil {
			panic(err)
		}
		img := resource.GetTexture(option.Texture)
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
