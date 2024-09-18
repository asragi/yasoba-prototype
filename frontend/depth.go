package frontend

type Depth int

const (
	Zero Depth = iota
	DepthEnemy
	DepthEffect
	DepthWindow
	DepthDebug
)

var AllDepths = []Depth{
	DepthEnemy,
	DepthEffect,
	DepthWindow,
	DepthDebug,
}
