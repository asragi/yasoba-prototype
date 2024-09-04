package main

type Drawing struct {
	Draw    DrawFunc
	DrawEnd DrawEnd
}

type DrawArgFunc func()
type DrawFunc func(DrawArgFunc, Depth)
type DrawEnd func()

func NewDrawing() *Drawing {
	const MaxDrawSize = 128
	drawIndex := make(map[Depth]int)
	drawMap := make(map[Depth][]DrawArgFunc)
	for i := 0; i < len(AllDepths); i++ {
		depth := AllDepths[i]
		drawMap[depth] = make([]DrawArgFunc, MaxDrawSize)
	}
	draw := func(d DrawArgFunc, depth Depth) {
		if drawIndex[depth] >= MaxDrawSize {
			// TODO: export error log
			return
		}
		drawMap[depth][drawIndex[depth]] = d
		drawIndex[depth]++
	}
	drawEnd := func() {
		for i := 0; i < len(AllDepths); i++ {
			depth := AllDepths[i]
			for j := 0; j < drawIndex[depth]; j++ {
				drawMap[depth][j]()
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
