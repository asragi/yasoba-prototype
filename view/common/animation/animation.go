package animation

import (
	"image"

	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/util"
)

type Sprite interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
	SetRect(image.Rectangle)
	TextureSize() *drawing.Vector
	SetShader(*drawing.Shader)
	SetRenderTargetFactory(func() drawing.Image)
	SetScaleBySize(*drawing.Vector)
}

type SpriteFactory func(
	relativePosition *drawing.Vector,
	pivot *drawing.Pivot,
	depth drawing.Depth,
	image drawing.Image,
) Sprite

type Animation struct {
	sprite Sprite
	frame  int
	data   *frontend.AnimationData
}

func New(
	newSprite SpriteFactory,
	relativePosition *drawing.Vector,
	pivot *drawing.Pivot,
	depth drawing.Depth,
	image drawing.Image,
	data *frontend.AnimationData,
	renderTargetFactory func() drawing.Image,
) *Animation {
	sprite := newSprite(
		relativePosition,
		pivot,
		depth,
		image,
	)
	if renderTargetFactory != nil {
		sprite.SetRenderTargetFactory(renderTargetFactory)
	}
	return &Animation{
		sprite: sprite,
		frame:  0,
		data:   data,
	}
}

func (a *Animation) setRect() {
	textureSize := a.sprite.TextureSize()
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
	a.sprite.SetRect(
		image.Rect(
			int(width*float64(column)),
			int(height*float64(row)),
			int(width*float64(column+1)),
			int(height*float64(row+1)),
		),
	)
}

func (a *Animation) Update(passedPosition *drawing.Vector) {
	a.frame++
	a.setRect()
	a.sprite.Update(passedPosition)
}

func (a *Animation) Draw(drawFunc drawing.DrawFunc) {
	a.sprite.Draw(drawFunc)
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

func (a *Animation) SetShader(shader *drawing.Shader) {
	a.sprite.SetShader(shader)
}

func (a *Animation) SetRenderTargetFactory(factory func() drawing.Image) {
	a.sprite.SetRenderTargetFactory(factory)
}

func (a *Animation) SetScaleBySize(size *drawing.Vector) {
	width := size.X * float64(a.data.ColumnCount)
	height := size.Y * float64(a.data.RowCount)
	a.sprite.SetScaleBySize(&drawing.Vector{X: width, Y: height})
}
