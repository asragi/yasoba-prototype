package frontend

type Depth int

const (
	Zero Depth = iota
	DepthEnemy
	DepthWindow
	DepthDebug
)

var AllDepths = []Depth{
	DepthEnemy,
	DepthWindow,
	DepthDebug,
}
