package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/text"
)

func InitializeCreateDebugScene(newSelectWindow component.NewSelectWindowFunc) CreateDebugScene {
	return func() *DebugScene {
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
			func(int) {},
			false,
		)
		selectWindow.Open()

		input := &frontend.KeyBoardInput{}
		input.Set(selectWindow)

		return &DebugScene{
			ui: debugUI{
				selectWindow: selectWindow,
				input:        input,
				anchor:       anchor,
			},
		}
	}
}
