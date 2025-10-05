package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
)

type DebugOption struct {
	OnSelectBattle func()
}

type DebugScene struct {
	ui debugUI
}

func (s *DebugScene) Update() {
	s.ui.Update()
}

func (s *DebugScene) Draw(drawFunc frontend.DrawFunc) {
	s.ui.Draw(drawFunc)
}

type debugUI struct {
	selectWindow *component.SelectWindow
	input        frontend.InputManager
}

func (ui *debugUI) Update() {
	ui.selectWindow.Update(frontend.VectorZero)
	ui.input.Update()
}

func (ui *debugUI) Draw(drawFunc frontend.DrawFunc) {
	ui.selectWindow.Draw(drawFunc)
}

type CreateDebugScene func(*DebugOption) *DebugScene

func InitializeCreateDebugScene(
	newSelectWindow component.NewSelectWindowFunc,
) CreateDebugScene {
	return func(option *DebugOption) *DebugScene {
		keyboard := &frontend.KeyBoardInput{}
		selectWindow := newSelectWindow(
			&frontend.Vector{X: 48, Y: 48},
			frontend.PivotTopLeft,
			frontend.DepthWindow,
			[]core.TextId{core.TextIdDebugMenuBattle},
			func(int) {
				option.OnSelectBattle()
			},
			true,
		)
		keyboard.Set(selectWindow)
		selectWindow.Open()

		return &DebugScene{
			ui: debugUI{
				selectWindow: selectWindow,
				input:        keyboard,
			},
		}
	}
}
