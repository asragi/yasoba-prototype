package scene

import (
	"github.com/asragi/yasoba-prototype/component/selection"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/text"
)

func InitializeCreateDebugScene(newSelectWindow selection.NewSelectWindowFunc) CreateDebugScene {
	return func() *DebugScene {
		scene := &DebugScene{}
		anchor := frontend.Vector{X: 48, Y: 48}
		selectWindow := newSelectWindow(
			&anchor,
			frontend.PivotTopLeft,
			frontend.DepthWindow,
			[]text.TextId{
				text.TextIdDebugMenuBattle,
				text.TextIdDebugMenuSequence,
				text.TextIdDebugMenuEffect,
			},
			scene.handleSelect,
			false,
		)
		selectWindow.Open()

		input := &frontend.KeyBoardInput{}
		input.Set(selectWindow)

		scene.ui = debugUI{
			selectWindow: selectWindow,
			input:        input,
			anchor:       anchor,
		}
		return scene
	}
}
