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
	drawing      *core.Drawing
	fontSource   *text.GoTextFaceSource
	windowSource *ebiten.Image
	window       *widgets.Window
)

func init() {
	drawing = core.NewDrawing()
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
			Image:            windowSource,
			CornerSize:       6,
			RelativePosition: &core.Vector{X: 384, Y: 288},
			Size: &core.Vector{
				X: 300,
				Y: 200,
			},
			Depth: core.DepthWindow,
			Pivot: core.PivotBottomRight,
		},
	)
}

type Game struct{}

func (g *Game) Update() error {
	window.Update(&core.Vector{X: 0, Y: 0})
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	op.Filter = ebiten.FilterLinear
	drawing.Draw(
		func(screen *ebiten.Image) {
			text.Draw(
				screen, "こんにちは世界", &text.GoTextFace{
					Source: fontSource,
					Size:   12,
				}, op,
			)
		}, core.DepthDebug,
	)
	window.Draw(drawing.Draw)
	drawing.DrawEnd(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 384, 288
}

func main() {
	ebiten.SetWindowSize(768, 576)
	ebiten.SetWindowTitle("yasoba-prototype")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
