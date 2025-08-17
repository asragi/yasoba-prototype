package frontend

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Drawing struct {
	Draw    DrawFunc
	DrawEnd DrawEnd
}

type DrawArgFunc func(*ebiten.Image)
type DrawFunc func(DrawArgFunc, Depth)
type DrawEnd func(*ebiten.Image)

func NewDrawing() *Drawing {
	const MaxDrawSize = 256
	drawIndex := make(map[Depth]int)
	drawMap := make(map[Depth][]DrawArgFunc)
	for i := 0; i < len(AllDepths); i++ {
		depth := AllDepths[i]
		drawMap[depth] = make([]DrawArgFunc, MaxDrawSize)
	}
	draw := func(d DrawArgFunc, depth Depth) {
		if drawIndex[depth] >= MaxDrawSize {
			panic(fmt.Sprintf("drawIndex[depth] >= MaxDrawSize: %d >= %d", drawIndex[depth], MaxDrawSize))
		}
		drawMap[depth][drawIndex[depth]] = d
		drawIndex[depth]++
	}
	drawEnd := func(screen *ebiten.Image) {
		for i := 0; i < len(AllDepths); i++ {
			depth := AllDepths[i]
			for j := 0; j < drawIndex[depth]; j++ {
				drawMap[depth][j](screen)
				drawMap[depth][j] = nil
			}
			drawIndex[depth] = 0
		}
	}
	return &Drawing{
		Draw:    draw,
		DrawEnd: drawEnd,
	}
}
