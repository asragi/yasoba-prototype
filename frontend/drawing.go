package frontend

import (
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/drawing/adapter/ebitenadapter"
	"github.com/hajimehoshi/ebiten/v2"
)

type DrawArgFunc func(*ebiten.Image)
type DrawFunc func(DrawArgFunc, Depth)
type DrawEnd func(*ebiten.Image)

type Drawing struct {
	adapter *ebitenadapter.Adapter
}

func NewDrawing() *Drawing {
	const maxDrawSize = 256

	depthOrder := make([]int, len(AllDepths))
	for i, depth := range AllDepths {
		depthOrder[i] = int(depth)
	}

	manager := drawing.NewManager(depthOrder, maxDrawSize)
	adapter := ebitenadapter.New(manager)

	return &Drawing{adapter: adapter}
}

func (d *Drawing) Draw(fn DrawArgFunc, depth Depth) {
	d.adapter.Draw(fn, int(depth))
}

func (d *Drawing) DrawEnd(screen *ebiten.Image) {
	d.adapter.Flush(screen)
}
