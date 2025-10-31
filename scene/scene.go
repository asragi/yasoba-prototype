package scene

import "github.com/asragi/yasoba-prototype/drawing"

type Scene interface {
	Update()
	Draw(drawing.DrawFunc)
}
