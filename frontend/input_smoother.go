package frontend

type SmoothKey int

const (
	SmoothKeyNone SmoothKey = iota
	SmoothKeyUp
	SmoothKeyDown
	SmoothKeyLeft
	SmoothKeyRight
)

type InputSmoother struct {
	frame     int
	beforeKey SmoothKey
	inputted  bool
}

func NewInputSmoother() *InputSmoother {
	return &InputSmoother{}
}

func (s *InputSmoother) Do(input SmoothKey) bool {
	const MaxFrame = 12
	const Duration = 8
	s.inputted = true
	if input != s.beforeKey {
		s.frame = 0
		s.beforeKey = input
		return true
	}
	s.frame++
	if s.frame <= MaxFrame {
		return false
	}
	if s.frame%Duration != 0 {
		return false
	}
	return true
}

func (s *InputSmoother) Update() {
	if s.inputted {
		s.inputted = false
		return
	}
	s.frame = 0
	s.beforeKey = SmoothKeyNone
}
