package drawing

type Depth int

const (
	Zero Depth = iota
	DepthEnemy
	DepthPlayer
	DepthBattleStatusText
	DepthEffect
	DepthDamageText
	DepthWindow
	DepthDebug
	DepthTransition
)

var AllDepths = []Depth{
	DepthEnemy,
	DepthPlayer,
	DepthBattleStatusText,
	DepthEffect,
	DepthDamageText,
	DepthWindow,
	DepthDebug,
	DepthTransition,
}
