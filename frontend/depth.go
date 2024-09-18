package frontend

type Depth int

const (
	Zero Depth = iota
	DepthEnemy
	DepthEffect
	DepthDamageText
	DepthWindow
	DepthDebug
)

var AllDepths = []Depth{
	DepthEnemy,
	DepthEffect,
	DepthDamageText,
	DepthWindow,
	DepthDebug,
}
