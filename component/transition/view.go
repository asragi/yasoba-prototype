package transition

import (
	"image/color"

	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
)

type View struct {
	overlay OverlayImage
	depth   frontend.Depth
}

type NewTransitionViewFunc func() *View

type OverlayImage interface {
	Fill(color color.Color)
	Draw(target *ebiten.Image, op *ebiten.DrawImageOptions)
}

type CreateImageFunc func(width, height int) OverlayImage

func CreateNewView(
	width, height int, depth frontend.Depth, createImage CreateImageFunc,
) NewTransitionViewFunc {
	if width <= 0 {
		panic("transition: width must be positive")
	}
	if height <= 0 {
		panic("transition: height must be positive")
	}
	if createImage == nil {
		panic("transition: createImage must not be nil")
	}
	return func() *View {
		img := createImage(width, height)
		if img == nil {
			panic("transition: createImage returned nil")
		}
		img.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})
		return &View{
			overlay: img,
			depth:   depth,
		}
	}
}

func (v *View) Draw(drawFunc frontend.DrawFunc, rate float64) {
	if v == nil || v.overlay == nil {
		return
	}
	if drawFunc == nil {
		return
	}
	alpha := clampRate(rate)
	if alpha <= 0 {
		return
	}

	drawFunc(func(screen *ebiten.Image) {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(float32(alpha))
		v.overlay.Draw(screen, op)
	}, v.depth)
}

func clampRate(rate float64) float64 {
	if rate < 0 {
		return 0
	}
	if rate > 1 {
		return 1
	}
	return rate
}

type ebitenOverlayImage struct {
	img *ebiten.Image
}

func (e *ebitenOverlayImage) Fill(c color.Color) {
	if e.img == nil {
		return
	}
	e.img.Fill(c)
}

func (e *ebitenOverlayImage) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	if e.img == nil || target == nil || op == nil {
		return
	}
	target.DrawImage(e.img, op)
}

func NewEbitenImageFactory() CreateImageFunc {
	return func(width, height int) OverlayImage {
		return &ebitenOverlayImage{
			img: ebiten.NewImage(width, height),
		}
	}
}
