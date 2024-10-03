package frontend

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
)

var AllDepths = []Depth{
	DepthEnemy,
	DepthPlayer,
	DepthBattleStatusText,
	DepthEffect,
	DepthDamageText,
	DepthWindow,
	DepthDebug,
}
