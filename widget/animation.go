package widget

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Animation struct {
	image *Image
	frame int
	data  *frontend.AnimationData
}

func NewAnimation(
	relativePosition *frontend.Vector,
	pivot *frontend.Pivot,
	depth frontend.Depth,
	image *ebiten.Image,
	data *frontend.AnimationData,
) *Animation {
	return &Animation{
		image: NewImage(
			relativePosition,
			pivot,
			depth,
			image,
		),
		frame: 0,
		data:  data,
	}
}

func (a *Animation) setRect() {
	textureSize := a.image.TextureSize()
	width := textureSize.X / float64(a.data.ColumnCount)
	height := textureSize.Y / float64(a.data.RowCount)
	target := func() int {
		if a.data.IsLoop {
			return (a.frame / a.data.Duration) % a.data.AnimationCount
		}
		return util.ClampInt(a.frame/a.data.Duration, 0, a.data.AnimationCount-1)
	}()
	row := target / a.data.ColumnCount
	column := target % a.data.ColumnCount
	a.image.SetRect(
		image.Rect(
			int(width*float64(column)),
			int(height*float64(row)),
			int(width*float64(column+1)),
			int(height*float64(row+1)),
		),
	)
}

func (a *Animation) Update(passedPosition *frontend.Vector) {
	a.frame++
	a.setRect()
	a.image.Update(passedPosition)
}

func (a *Animation) Draw(drawFunc frontend.DrawFunc) {
	a.image.Draw(drawFunc)
}

func (a *Animation) Reset() {
	a.frame = 0
}

func (a *Animation) IsEnd() bool {
	if a.data.IsLoop {
		return false
	}
	return a.frame >= a.data.Duration*a.data.AnimationCount
}

func (a *Animation) SetShader(shader *frontend.Shader) {
	a.image.SetShader(shader)
}

func (a *Animation) SetScaleBySize(size *frontend.Vector) {
	width := size.X * float64(a.data.ColumnCount)
	height := size.Y * float64(a.data.RowCount)
	a.image.SetScaleBySize(&frontend.Vector{X: width, Y: height})
}
