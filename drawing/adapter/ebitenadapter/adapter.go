package ebitenadapter

import (
	"fmt"

	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/hajimehoshi/ebiten/v2"
)

type Adapter struct {
	manager *drawing.Manager
}

func New(manager *drawing.Manager) *Adapter {
	if manager == nil {
		panic("ebitenadapter: manager must not be nil")
	}
	return &Adapter{manager: manager}
}

func (a *Adapter) Draw(fn func(*ebiten.Image), depth int) {
	a.manager.Draw(func(surface drawing.Surface) {
		img, ok := surface.(*ebiten.Image)
		if !ok {
			panic(fmt.Sprintf("ebitenadapter: unexpected surface type %T", surface))
		}
		fn(img)
	}, depth)
}

func (a *Adapter) Flush(screen *ebiten.Image) {
	a.manager.DrawEnd(screen)
}
