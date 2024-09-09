package frontend

type Depth int

const (
	Zero Depth = iota
	DepthWindow

	DepthDebug
)

var AllDepths = []Depth{
	DepthWindow,
	DepthDebug,
}
