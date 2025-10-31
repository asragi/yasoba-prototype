package widget

import (
	"github.com/asragi/yasoba-prototype/drawing"
)

type PositionUpdater interface {
	Update(parentPosition *drawing.Vector)
}

type Drawer interface {
	Draw(drawFunc drawing.DrawFunc)
}
