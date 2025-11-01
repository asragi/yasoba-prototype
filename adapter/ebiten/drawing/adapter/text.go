package adapter

import (
	"image/color"

	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func NewDrawTextFunc() drawing.DrawTextFunc {
	return func(img drawing.Image, content string, face drawing.TextFace, options *drawing.TextDrawOptions) {
		target, ok := img.(*EbitenImage)
		if !ok {
			panic("drawing/adapter: img is not EbitenImage")
		}
		textFace, ok := face.(*text.GoTextFace)
		if !ok {
			panic("drawing/adapter: face is not *text.GoTextFace")
		}
		drawOptions := &text.DrawOptions{}
		if options != nil {
			drawOptions.GeoM.Scale(options.Scale.X, options.Scale.Y)
			drawOptions.GeoM.Translate(options.Position.X, options.Position.Y)
			applyColor(&drawOptions.ColorScale, options.Color)
		}
		text.Draw(target.img, content, textFace, drawOptions)
	}
}

func applyColor(scale interface {
	Reset()
	ScaleWithColor(color.Color)
}, clr color.Color) {
	if clr == nil {
		return
	}
	scale.Reset()
	scale.ScaleWithColor(clr)
}
