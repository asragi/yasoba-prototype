package scene

import (
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/toolkit/input"
	"github.com/asragi/yasoba-prototype/view/common/selection"
)

type debugUI struct {
	selectWindow *selection.SelectWindow
	input        input.Manager
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
