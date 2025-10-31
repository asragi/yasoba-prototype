package scene

import (
	"github.com/asragi/yasoba-prototype/component/selection"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
)

type debugUI struct {
	selectWindow *selection.SelectWindow
	input        frontend.InputManager
	anchor       drawing.Vector
	layout       debugUILayout
}

type debugUILayout struct {
	selectAnchor drawing.Vector
}

func (ui *debugUI) Update() {
	ui.layout = computeDebugUILayout(&ui.anchor)
	ui.selectWindow.Update(&ui.layout.selectAnchor)
	ui.input.Update()
}

func (ui *debugUI) Draw(drawFunc drawing.DrawFunc) {
	ui.selectWindow.Draw(drawFunc)
}

func computeDebugUILayout(anchor *drawing.Vector) debugUILayout {
	return debugUILayout{
		selectAnchor: *anchor,
	}
}
