package debug

import (
	"fmt"
	"runtime"

	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

func fps() float64 {
	return ebiten.ActualFPS()
}

func tps() float64 {
	return ebiten.ActualTPS()
}

type Debug struct {
	Update func()
	Draw   func(frontend.DrawFunc)
}

func CreateDrawParameters(newText widget.NewTextFunc) *Debug {
	textFPS := newText(&widget.TextOptionsNew{
		RelativePosition: &frontend.Vector{X: 0, Y: 0},
		Font:             frontend.MaruMinya,
		Scale:            1,
		Depth:            frontend.DepthDebug,
		Pivot:            frontend.PivotTopLeft,
	})
	textTPS := newText(&widget.TextOptionsNew{
		RelativePosition: &frontend.Vector{X: 0, Y: 16},
		Font:             frontend.MaruMinya,
		Scale:            1,
		Depth:            frontend.DepthDebug,
		Pivot:            frontend.PivotTopLeft,
	})
	textMemory := newText(&widget.TextOptionsNew{
		RelativePosition: &frontend.Vector{X: 0, Y: 32},
		Font:             frontend.MaruMinya,
		Scale:            1,
		Depth:            frontend.DepthDebug,
		Pivot:            frontend.PivotTopLeft,
	})
	textSys := newText(&widget.TextOptionsNew{
		RelativePosition: &frontend.Vector{X: 0, Y: 48},
		Font:             frontend.MaruMinya,
		Scale:            1,
		Depth:            frontend.DepthDebug,
		Pivot:            frontend.PivotTopLeft,
	})

	position := frontend.VectorZero

	update := func() {
		var memory runtime.MemStats
		runtime.ReadMemStats(&memory)

		textFPS.Update(position)
		textFPS.SetText(fmt.Sprintf("FPS: %f", fps()), true)
		textTPS.Update(position)
		textTPS.SetText(fmt.Sprintf("TPS: %f", tps()), true)
		textMemory.Update(position)
		textMemory.SetText(fmt.Sprintf("Alloc: %d MB", memory.Alloc/1024/1024), true)
		textSys.Update(position)
		textSys.SetText(fmt.Sprintf("Sys: %d MB", memory.Sys/1024/1024), true)
	}
	draw := func(draw frontend.DrawFunc) {
		textFPS.Draw(draw)
		textTPS.Draw(draw)
		textMemory.Draw(draw)
		textSys.Draw(draw)
	}
	return &Debug{
		Update: update,
		Draw:   draw,
	}
}
