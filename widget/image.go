package widget

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type texture interface {
	Bounds() image.Rectangle
	SubImage(image.Rectangle) texture
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
	AsImage() *ebiten.Image
}

type ebitenTexture struct {
	img *ebiten.Image
}

func newEbitenTexture(img *ebiten.Image) texture {
	if img == nil {
		panic("ebiten texture must not be nil")
	}
	return &ebitenTexture{img: img}
}

func newEmptyEbitenTexture(width, height int) texture {
	return newEbitenTexture(ebiten.NewImage(width, height))
}

func (t *ebitenTexture) Bounds() image.Rectangle {
	return t.img.Bounds()
}

func (t *ebitenTexture) SubImage(rect image.Rectangle) texture {
	return newEbitenTexture(t.img.SubImage(rect).(*ebiten.Image))
}

func (t *ebitenTexture) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	target.DrawImage(t.img, op)
}

func (t *ebitenTexture) AsImage() *ebiten.Image {
	return t.img
}

type Image struct {
	relativePosition *frontend.Vector
	parentPosition   *frontend.Vector
	pivot            *frontend.Pivot
	image            texture
	depth            frontend.Depth
	scale            *frontend.Vector
	rect             *image.Rectangle
	shader           *frontend.Shader
}

func (i *Image) Update(passedPosition *frontend.Vector) {
	i.parentPosition = passedPosition
	if i.shader == nil {
		return
	}
	i.shader.Update()
}

func (i *Image) Draw(drawFunc frontend.DrawFunc) {
	if i.parentPosition == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	pivotModification := i.pivot.ApplyToSize(i.Size())
	op.GeoM.Scale(
		i.scale.X,
		i.scale.Y,
	)
	op.GeoM.Translate(
		i.parentPosition.X+i.relativePosition.X-pivotModification.X,
		i.parentPosition.Y+i.relativePosition.Y-pivotModification.Y,
	)
	imageToDraw := i.image.SubImage(*i.rect)

	drawFunc(
		func(screen *ebiten.Image) {
			if i.shader == nil {
				imageToDraw.Draw(screen, op)
				return
			}
			w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
			renderTarget := newEmptyEbitenTexture(w, h)
			imageToDraw.Draw(renderTarget.AsImage(), op)
			shaderOption := &ebiten.DrawRectShaderOptions{}
			shaderOption.Images[0] = renderTarget.AsImage()
			shaderOption.Uniforms = i.shader.GetUniforms()
			screen.DrawRectShader(w, h, i.shader.GetShader(), shaderOption)
		}, i.depth,
	)
}

func (i *Image) SetRect(rect image.Rectangle) {
	i.rect = &rect
}

func (i *Image) Size() *frontend.Vector {
	return &frontend.Vector{
		X: float64(i.rect.Dx()) * i.scale.X,
		Y: float64(i.rect.Dy()) * i.scale.Y,
	}
}

func (i *Image) TextureSize() *frontend.Vector {
	return &frontend.Vector{
		X: float64(i.image.Bounds().Dx()),
		Y: float64(i.image.Bounds().Dy()),
	}
}

func (i *Image) SetScaleBySize(size *frontend.Vector) {
	i.scale = &frontend.Vector{
		X: size.X / float64(i.image.Bounds().Dx()),
		Y: size.Y / float64(i.image.Bounds().Dy()),
	}
}

func (i *Image) SetRelativePosition(position *frontend.Vector) {
	i.relativePosition = position
}

func (i *Image) SetShader(shader *frontend.Shader) {
	i.shader = shader
	shader.Reset()
}

func (i *Image) SetShaderUniforms(key string, value interface{}) {
	i.shader.SetUniforms(key, value)
}

type ImageOption struct {
	customDraw CustomDrawFunc
}

type CustomDrawArgs struct {
	relativePosition *frontend.Vector
	parentPosition   *frontend.Vector
	pivot            *frontend.Pivot
	image            texture
	depth            frontend.Depth
	scale            *frontend.Vector
	rect             *image.Rectangle
	size             *frontend.Vector
	drawFunc         frontend.DrawFunc
}

type CustomDrawFunc func(*CustomDrawArgs)

func NewImage(
	relativePosition *frontend.Vector,
	pivot *frontend.Pivot,
	depth frontend.Depth,
	imageData *ebiten.Image,
) *Image {
	rect := image.Rect(0, 0, imageData.Bounds().Dx(), imageData.Bounds().Dy())
	return &Image{
		relativePosition: relativePosition,
		pivot:            pivot,
		image:            newEbitenTexture(imageData),
		depth:            depth,
		scale:            &frontend.Vector{X: 1, Y: 1},
		rect:             &rect,
		shader:           nil,
	}
}
