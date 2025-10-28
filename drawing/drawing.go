package drawing

import (
	"image"
)

type Image interface {
	Bounds() image.Rectangle
	SubImage(image.Rectangle) Image
	DrawImage(Image, *DrawOptions)
	DrawRectShader(int, int, *Shader, *DrawRectShaderOptions)
}

type NewEmptyImageFunc func(width, height int) Image

type DrawRectShaderOptions struct {
	Images   []Image
	Uniforms map[string]interface{}
}

func NewDrawRectShaderOptions() *DrawRectShaderOptions {
	images := make([]Image, 4)
	return &DrawRectShaderOptions{
		Images:   images,
		Uniforms: map[string]interface{}{},
	}
}

func (o *DrawRectShaderOptions) SetImage(index int, img Image) {
	if index < 0 || index >= len(o.Images) {
		panic("drawing: image index out of range")
	}
	o.Images[index] = img
}

func (o *DrawRectShaderOptions) SetUniforms(uniforms map[string]interface{}) {
	o.Uniforms = uniforms
}

type DrawOptions struct {
	Opacity  float64
	Position *Vector
	Scale    *Vector
}

func NewDrawOptions() *DrawOptions {
	return &DrawOptions{
		Opacity:  1.0,
		Position: &Vector{X: 0, Y: 0},
		Scale:    &Vector{X: 1, Y: 1},
	}
}

func (o *DrawOptions) Translate(x, y float64) {
	o.Position.X += x
	o.Position.Y += y
}

func (o *DrawOptions) SetScale(sx, sy float64) {
	o.Scale.X = sx
	o.Scale.Y = sy
}

func (o *DrawOptions) SetOpacity(opacity float64) {
	o.Opacity = opacity
}

type DrawArgFunc func(Image)
type DrawFunc func(DrawArgFunc, Depth)
type DrawEndFunc func(Image)
