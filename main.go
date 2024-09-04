package main

import (
	"bytes"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/fonts"
	"github.com/asragi/yasoba-prototype/images"
	"github.com/asragi/yasoba-prototype/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"log"
)

var (
	drawing      *Drawing
	fontSource   *text.GoTextFaceSource
	windowSource *ebiten.Image
	window       *widgets.Window
)

func init() {
	drawing = NewDrawing()
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MaruMinya))
	if err != nil {
		log.Fatal(err)
	}
	fontSource = s
	w, _, err := image.Decode(bytes.NewReader(images.Window))
	if err != nil {
		log.Fatal(err)
	}
	windowSource = ebiten.NewImageFromImage(w)
	window = widgets.NewWindow(
		&widgets.WindowOption{
			Image:      windowSource,
			CornerSize: 6,
			Size: &core.Vector{
				X: 300,
				Y: 200,
			},
		},
	)
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	op.Filter = ebiten.FilterLinear
	drawing.Draw(
		func() {
			text.Draw(
				screen, "こんにちは世界", &text.GoTextFace{
					Source: fontSource,
					Size:   12,
				}, op,
			)
		}, DepthDebug,
	)
	drawing.Draw(func() { window.Draw(screen) }, DepthWindow)
	drawing.DrawEnd()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 384, 288
}

func main() {
	ebiten.SetWindowSize(768, 576)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
