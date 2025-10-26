package scene

import (
	"github.com/asragi/yasoba-prototype/component/selection"
	"github.com/asragi/yasoba-prototype/frontend"
)

type debugUI struct {
	selectWindow *selection.SelectWindow
	input        frontend.InputManager
	anchor       frontend.Vector
	layout       debugUILayout
}

type debugUILayout struct {
	selectAnchor frontend.Vector
}

func (ui *debugUI) Update() {
	ui.layout = computeDebugUILayout(&ui.anchor)
	ui.selectWindow.Update(&ui.layout.selectAnchor)
	ui.input.Update()
}

func (ui *debugUI) Draw(drawFunc frontend.DrawFunc) {
	ui.selectWindow.Draw(drawFunc)
}

func computeDebugUILayout(anchor *frontend.Vector) debugUILayout {
	return debugUILayout{
		selectAnchor: *anchor,
	}
}
