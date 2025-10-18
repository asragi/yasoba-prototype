package drawing

import "fmt"

type Surface interface{}
type DrawArgFunc func(Surface)
type DrawFunc func(DrawArgFunc, int)

type Manager struct {
	depthOrder []int
	drawIndex  map[int]int
	drawMap    map[int][]DrawArgFunc
	maxSize    int
}

func NewManager(depthOrder []int, maxSize int) *Manager {
	if maxSize <= 0 {
		panic("drawing: maxSize must be positive")
	}

	drawIndex := make(map[int]int, len(depthOrder))
	drawMap := make(map[int][]DrawArgFunc, len(depthOrder))
	for _, depth := range depthOrder {
		drawMap[depth] = make([]DrawArgFunc, maxSize)
	}

	return &Manager{
		depthOrder: append([]int(nil), depthOrder...),
		drawIndex:  drawIndex,
		drawMap:    drawMap,
		maxSize:    maxSize,
	}
}

func (m *Manager) Draw(fn DrawArgFunc, depth int) {
	if _, ok := m.drawMap[depth]; !ok {
		panic(fmt.Sprintf("drawing: depth %d is not registered", depth))
	}
	if m.drawIndex[depth] >= m.maxSize {
		panic(fmt.Sprintf("drawing: drawIndex[%d] >= maxSize: %d >= %d", depth, m.drawIndex[depth], m.maxSize))
	}
	m.drawMap[depth][m.drawIndex[depth]] = fn
	m.drawIndex[depth]++
}

func (m *Manager) DrawEnd(surface Surface) {
	for _, depth := range m.depthOrder {
		callbacks := m.drawMap[depth]
		for i := 0; i < m.drawIndex[depth]; i++ {
			callbacks[i](surface)
			callbacks[i] = nil
		}
		m.drawIndex[depth] = 0
	}
}
