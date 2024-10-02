package frontend

type Depth int

const (
	Zero Depth = iota
	DepthEnemy
	DepthPlayer
	DepthEffect
	DepthDamageText
	DepthWindow
	DepthDebug
)

var AllDepths = []Depth{
	DepthEnemy,
	DepthPlayer,
	DepthEffect,
	DepthDamageText,
	DepthWindow,
	DepthDebug,
}
