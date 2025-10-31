package scene

import "github.com/asragi/yasoba-prototype/drawing"

type DebugScene struct {
	ui             debugUI
	onSelectBattle func()
}

func (s *DebugScene) Update() {
	s.ui.Update()
}

func (s *DebugScene) Draw(drawFunc drawing.DrawFunc) {
	s.ui.Draw(drawFunc)
}

func (s *DebugScene) SetOnSelectBattle(f func()) {
	s.onSelectBattle = f
}

func (s *DebugScene) handleSelect(index int) {
	if index == 0 && s.onSelectBattle != nil {
		s.onSelectBattle()
	}
}

type CreateDebugScene func() *DebugScene
