package transition

import "github.com/asragi/yasoba-prototype/drawing"

type transitionViewDraw func(drawFunc drawing.DrawFunc, rate float64)

type mode int

const (
	ModeNone mode = iota
	ModeFadeOut
	ModeFadeIn
)

type transition struct {
	frame    int
	maxFrame int
	mode     mode
	view     transitionViewDraw
}

func (t *transition) FadeOut() {
	t.mode = ModeFadeOut
}

func (t *transition) FadeIn() {
	t.mode = ModeFadeIn
}

func (t *transition) fadeOutUpdate() {
	if t.mode != ModeFadeOut {
		return
	}
	if t.frame >= t.maxFrame {
		t.mode = ModeNone
		return
	}
	t.frame++
}

func (t *transition) fadeInUpdate() {
	if t.mode != ModeFadeIn {
		return
	}
	if t.frame <= 0 {
		t.mode = ModeNone
		return
	}
	t.frame--
}

func (t *transition) Update() {
	if t.mode == ModeNone {
		return
	}
	t.fadeOutUpdate()
	t.fadeInUpdate()
}

func (t *transition) Draw(drawFunc drawing.DrawFunc) {
	if drawFunc == nil {
		return
	}
	t.view(drawFunc, float64(t.frame)/float64(t.maxFrame))
}

func New(maxFrame int, view transitionViewDraw) *transition {
	return &transition{
		frame:    0,
		maxFrame: maxFrame,
		mode:     ModeNone,
		view:     view,
	}
}
