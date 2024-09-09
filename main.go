package main

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

var (
	drawing  *frontend.Drawing
	window   *widget.Window
	testText *widget.Text
)

func init() {
	drawing = frontend.NewDrawing()
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		log.Fatal(err)
	}
	window = widget.NewWindow(
		&widget.WindowOption{
			Image:            resource.GetTexture(frontend.TextureWindow),
			CornerSize:       6,
			RelativePosition: &frontend.Vector{X: 192, Y: 0},
			Size: &frontend.Vector{
				X: 292,
				Y: 62,
			},
			Depth:   frontend.DepthWindow,
			Pivot:   frontend.PivotTopCenter,
			Padding: &frontend.Vector{X: 16, Y: 8},
		},
	)
	testText = widget.NewText(
		"あのイーハトーヴォのすきとおった風\n夏でも底に冷たさをもつ青いそら\nうつくしい森で飾られたモリーオ市",
		&widget.TextOptions{
			RelativePosition: &frontend.Vector{X: 0, Y: 0},
			Pivot:            frontend.PivotTopLeft,
			TextFace:         resource.GetFont(frontend.MaruMinya),
			DisplayAll:       false,
			Speed:            6,
			Depth:            frontend.DepthWindow,
		},
	)
}

type Game struct{}

func (g *Game) Update() error {
	window.Update(&frontend.Vector{X: 0, Y: 0})
	testText.Update(window.GetContentUpperLeft())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	op.Filter = ebiten.FilterLinear
	window.Draw(drawing.Draw)
	testText.Draw(drawing.Draw)
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
