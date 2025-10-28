package adapter

import (
	"image"

	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/hajimehoshi/ebiten/v2"
)

// implements drawing.Image using ebiten.Image
type EbitenImage struct {
	img *ebiten.Image
}

func NewEbitenImage(
	img *ebiten.Image,
) *EbitenImage {
	return &EbitenImage{
		img: img,
	}
}

func (i *EbitenImage) Bounds() image.Rectangle {
	return i.img.Bounds()
}

func (i *EbitenImage) SubImage(rect image.Rectangle) drawing.Image {
	img, ok := i.img.SubImage(rect).(*ebiten.Image)
	if !ok {
		panic("drawing/adapter: img is not EbitenImage")
	}
	return NewEbitenImage(img)
}

func (i *EbitenImage) DrawImage(img drawing.Image, options *drawing.DrawOptions) {
	ebitenImg, ok := img.(*EbitenImage)
	if !ok {
		panic("drawing/adapter: img is not EbitenImage")
	}
	i.img.DrawImage(ebitenImg.img, i.toEbitenOptions(options))
}

func (i *EbitenImage) DrawRectShader(
	width, height int,
	shader *drawing.Shader,
	options *drawing.DrawRectShaderOptions,
) {
	ebitenShaderOptions := i.toEbitenShaderOptions(options)
	i.img.DrawRectShader(
		width,
		height,
		shader.GetShader(),
		ebitenShaderOptions,
	)
}

func (_ *EbitenImage) toEbitenOptions(options *drawing.DrawOptions) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(options.Scale.X, options.Scale.Y)
	op.GeoM.Translate(options.Position.X, options.Position.Y)
	op.ColorScale.ScaleAlpha(float32(options.Opacity))
	return op
}

func (_ *EbitenImage) toEbitenShaderOptions(
	options *drawing.DrawRectShaderOptions,
) *ebiten.DrawRectShaderOptions {
	op := &ebiten.DrawRectShaderOptions{}
	for index, img := range options.Images {
		if img == nil {
			continue
		}
		ebitenImg, ok := img.(*EbitenImage)
		if !ok {
			panic("drawing/adapter: img is not EbitenImage")
		}
		op.Images[index] = ebitenImg.img
	}
	for key, value := range options.Uniforms {
		op.Uniforms[key] = value
	}
	return op
}
