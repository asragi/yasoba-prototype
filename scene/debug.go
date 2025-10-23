package scene

import "github.com/asragi/yasoba-prototype/frontend"

type DebugScene struct {
	ui debugUI
}

func (s *DebugScene) Update() {
	s.ui.Update()
}

func (s *DebugScene) Draw(drawFunc frontend.DrawFunc) {
	s.ui.Draw(drawFunc)
}

type CreateDebugScene func() *DebugScene
