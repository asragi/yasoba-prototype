package main

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

var (
	drawing            *frontend.Drawing
	messageWindow      *component.MessageWindow
	battleSelectWindow *component.BattleSelectWindow
	faceWindow         *component.FaceWindow
)

func init() {
	drawing = frontend.NewDrawing()
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		log.Fatal(err)
	}
	textServer := core.CreateServeTextData()
	newMessageWindow := component.StandByNewMessageWindow(resource)
	messageWindow = newMessageWindow(
		&frontend.Vector{X: 192, Y: 0},
		&frontend.Vector{X: 292, Y: 62},
		frontend.DepthWindow,
		frontend.PivotTopCenter,
	)
	testString := "あのイーハトーヴォのすきとおった風\n夏でも底に冷たさをもつ青いそら\nうつくしい森で飾られたモリーオ市"
	messageWindow.SetText(testString, false)
	newBattleSelectWindow := component.StandByNewBattleSelectWindow(resource, textServer)
	battleSelectWindow = newBattleSelectWindow(
		&frontend.Vector{X: 0, Y: 0},
		frontend.PivotBottomLeft,
		frontend.DepthWindow,
		[]component.BattleCommand{
			component.BattleCommandAttack,
			component.BattleCommandFire,
			/*
				component.BattleCommandThunder,
				component.BattleCommandBarrier,
				component.BattleCommandWind,
			*/
			component.BattleCommandFocus,
			component.BattleCommandDefend,
		},
	)
	battleSelectWindow.Open()
	newFaceWindow := component.StandByNewFaceWindow(resource)
	faceWindow = newFaceWindow(
		&frontend.Vector{X: 0, Y: 0},
		frontend.DepthWindow,
		frontend.PivotBottomLeft,
	)
}

type Game struct{}

func (g *Game) Update() error {
	messageWindow.Update(frontend.VectorZero)
	faceWindow.Update(&frontend.Vector{X: 0, Y: 288})
	battleSelectWindow.Update(faceWindow.GetTopLeftPosition())
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		messageWindow.Shake(2, 10)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		battleSelectWindow.MoveCursorDown()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		battleSelectWindow.MoveCursorUp()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	messageWindow.Draw(drawing.Draw)
	battleSelectWindow.Draw(drawing.Draw)
	faceWindow.Draw(drawing.Draw)
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
