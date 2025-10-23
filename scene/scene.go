package scene

import "github.com/asragi/yasoba-prototype/frontend"

type Scene interface {
	Update()
	Draw(frontend.DrawFunc)
}
