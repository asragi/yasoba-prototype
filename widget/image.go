package widget

import (
	"image"

	"github.com/asragi/yasoba-prototype/drawing"
)

type newEmptyTextureFunc func() drawing.Image

type Image struct {
	relativePosition *drawing.Vector
	parentPosition   *drawing.Vector
	pivot            *drawing.Pivot
	image            drawing.Image
	depth            drawing.Depth
	scale            *drawing.Vector
	rect             *image.Rectangle
	shader           *drawing.Shader
	newEmptyTexture  newEmptyTextureFunc
}

func (i *Image) Update(passedPosition *drawing.Vector) {
	i.parentPosition = passedPosition
	if i.shader == nil {
		return
	}
	i.shader.Update()
}

func (i *Image) Draw(drawFunc drawing.DrawFunc) {
	if i.parentPosition == nil {
		return
	}
	op := drawing.NewDrawOptions()
	pivotModification := i.pivot.ApplyToSize(i.Size())
	op.SetScale(
		i.scale.X,
		i.scale.Y,
	)
	op.Translate(
		i.parentPosition.X+i.relativePosition.X-pivotModification.X,
		i.parentPosition.Y+i.relativePosition.Y-pivotModification.Y,
	)
	imageToDraw := i.image.SubImage(*i.rect)

	drawFunc(
		func(screen drawing.Image) {
			if i.shader == nil {
				screen.DrawImage(imageToDraw, op)
				return
			}
			w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
			if i.newEmptyTexture == nil {
				panic("widget: shader render target factory is not set")
			}
			renderTarget := i.newEmptyTexture()
			targetBounds := renderTarget.Bounds()
			if targetBounds.Dx() != imageToDraw.Bounds().Dx() || targetBounds.Dy() != imageToDraw.Bounds().Dy() {
				// panic("widget: shader render target size mismatch")
			}
			renderTarget.DrawImage(imageToDraw, op)
			shaderOption := drawing.NewDrawRectShaderOptions()
			shaderOption.SetImage(0, renderTarget)
			shaderOption.SetUniforms(i.shader.GetUniforms())
			screen.DrawRectShader(w, h, i.shader, shaderOption)
		}, i.depth,
	)
}

func (i *Image) SetRect(rect image.Rectangle) {
	i.rect = &rect
}

func (i *Image) Size() *drawing.Vector {
	return &drawing.Vector{
		X: float64(i.rect.Dx()) * i.scale.X,
		Y: float64(i.rect.Dy()) * i.scale.Y,
	}
}

func (i *Image) TextureSize() *drawing.Vector {
	return &drawing.Vector{
		X: float64(i.image.Bounds().Dx()),
		Y: float64(i.image.Bounds().Dy()),
	}
}

func (i *Image) SetScaleBySize(size *drawing.Vector) {
	i.scale = &drawing.Vector{
		X: size.X / float64(i.image.Bounds().Dx()),
		Y: size.Y / float64(i.image.Bounds().Dy()),
	}
}

func (i *Image) SetRelativePosition(position *drawing.Vector) {
	i.relativePosition = position
}

func (i *Image) SetShader(shader *drawing.Shader) {
	i.shader = shader
	shader.Reset()
}

func (i *Image) SetRenderTargetFactory(factory newEmptyTextureFunc) {
	i.newEmptyTexture = factory
}

func (i *Image) SetShaderUniforms(key string, value interface{}) {
	i.shader.SetUniforms(key, value)
}

func NewImage(
	relativePosition *drawing.Vector,
	pivot *drawing.Pivot,
	depth drawing.Depth,
	imageData drawing.Image,
) *Image {
	rect := image.Rect(0, 0, imageData.Bounds().Dx(), imageData.Bounds().Dy())
	return &Image{
		relativePosition: relativePosition,
		pivot:            pivot,
		image:            imageData,
		depth:            depth,
		scale:            &drawing.Vector{X: 1, Y: 1},
		rect:             &rect,
		shader:           nil,
	}
}
