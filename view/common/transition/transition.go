package transition

import "github.com/asragi/yasoba-prototype/toolkit/drawing"

type transitionViewDraw func(drawFunc drawing.DrawFunc, rate float64)

type mode int

const (
	ModeNone mode = iota
	ModeFadeOut
	ModeFadeIn
)

type InitialState int

const (
	InitialStateTransparent InitialState = iota
	InitialStateOpaque
)

type Transition struct {
	frame    int
	maxFrame int
	mode     mode
	view     transitionViewDraw
}

func (t *Transition) FadeOut() {
	t.mode = ModeFadeOut
}

func (t *Transition) FadeIn() {
	t.mode = ModeFadeIn
}

func (t *Transition) fadeOutUpdate() {
	if t.mode != ModeFadeOut {
		return
	}
	if t.frame >= t.maxFrame {
		t.mode = ModeNone
		return
	}
	t.frame++
}

func (t *Transition) fadeInUpdate() {
	if t.mode != ModeFadeIn {
		return
	}
	if t.frame <= 0 {
		t.mode = ModeNone
		return
	}
	t.frame--
}

func (t *Transition) Update() {
	if t.mode == ModeNone {
		return
	}
	t.fadeOutUpdate()
	t.fadeInUpdate()
}

func (t *Transition) Draw(drawFunc drawing.DrawFunc) {
	if drawFunc == nil {
		return
	}
	t.view(drawFunc, float64(t.frame)/float64(t.maxFrame))
}

func New(maxFrame int, view transitionViewDraw, initial InitialState) *Transition {
	if maxFrame <= 0 {
		panic("transition: maxFrame must be positive")
	}
	tr := &Transition{
		maxFrame: maxFrame,
		mode:     ModeNone,
		view:     view,
	}
	tr.setInitialState(initial)
	return tr
}

func (t *Transition) setInitialState(initial InitialState) {
	if initial == InitialStateOpaque {
		t.frame = t.maxFrame
		return
	}
	t.frame = 0
}
