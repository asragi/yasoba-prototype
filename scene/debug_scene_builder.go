package scene

import (
	"github.com/asragi/yasoba-prototype/component/selection"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/input/adapter"
	"github.com/asragi/yasoba-prototype/text"
)

func InitializeCreateDebugScene(newSelectWindow selection.NewSelectWindowFunc) CreateDebugScene {
	return func() *DebugScene {
		scene := &DebugScene{}
		anchor := drawing.Vector{X: 48, Y: 48}
		selectWindow := newSelectWindow(
			&anchor,
			drawing.PivotTopLeft,
			drawing.DepthWindow,
			[]text.TextId{
				text.TextIdDebugMenuBattle,
				text.TextIdDebugMenuSequence,
				text.TextIdDebugMenuEffect,
			},
			scene.handleSelect,
			false,
		)
		selectWindow.Open()

		inputManager := &adapter.KeyBoardInput{}
		inputManager.Set(selectWindow)

		scene.ui = debugUI{
			selectWindow: selectWindow,
			input:        inputManager,
			anchor:       anchor,
		}
		return scene
	}
}
