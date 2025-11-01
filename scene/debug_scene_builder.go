package scene

import (
	ebiteninput "github.com/asragi/yasoba-prototype/adapter/ebiten/input"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/selection"
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

		inputManager := &ebiteninput.KeyBoardInput{}
		inputManager.Set(selectWindow)

		scene.ui = debugUI{
			selectWindow: selectWindow,
			input:        inputManager,
			anchor:       anchor,
		}
		return scene
	}
}
