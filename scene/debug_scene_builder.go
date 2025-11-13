package scene

import (
	ebiteninput "github.com/asragi/yasoba-prototype/adapter/ebiten/input"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/common/selection"
	"github.com/asragi/yasoba-prototype/view/common/selection/option"
)

func InitializeCreateDebugScene(
	newSelectWindow selection.NewSelectWindowFunc,
	newTextItem option.NewTextItemFunc,
) CreateDebugScene {
	return func() *DebugScene {
		scene := &DebugScene{}
		anchor := drawing.Vector{X: 48, Y: 48}
		texts := []text.TextId{
			text.TextIdDebugMenuBattle,
			text.TextIdDebugMenuSequence,
			text.TextIdDebugMenuEffect,
		}
		items := func() []selection.Item {
			result := []selection.Item{}
			for _, textId := range texts {
				item := newTextItem(textId)
				result = append(result, item)
			}
			return result
		}()
		selectWindow := newSelectWindow(
			&anchor,
			drawing.PivotTopLeft,
			drawing.DepthWindow,
			0,
			items,
			scene.handleSelect,
			nil,
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
