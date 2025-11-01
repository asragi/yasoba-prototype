package scene

import "github.com/asragi/yasoba-prototype/toolkit/drawing"

type Scene interface {
	Update()
	Draw(drawing.DrawFunc)
}
