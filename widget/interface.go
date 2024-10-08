package widget

import "github.com/asragi/yasoba-prototype/frontend"

type PositionUpdater interface {
	Update(parentPosition *frontend.Vector)
}

type Drawer interface {
	Draw(drawFunc frontend.DrawFunc)
}
