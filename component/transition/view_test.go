package transition

import (
	"image/color"
	"testing"

	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
)

type mockOverlay struct {
	fillColor  color.Color
	drawCount  int
	lastTarget *ebiten.Image
	lastOp     *ebiten.DrawImageOptions
}

func (m *mockOverlay) Fill(c color.Color) {
	m.fillColor = c
}

func (m *mockOverlay) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	m.drawCount++
	m.lastTarget = target
	m.lastOp = op
}

func TestCreateNewView(t *testing.T) {
	const (
		width  = 4
		height = 3
		depth  = frontend.DepthWindow
	)

	var overlay *mockOverlay
	create := func(w, h int) OverlayImage {
		if w != width || h != height {
			t.Fatalf("unexpected size: got (%d,%d)", w, h)
		}
		overlay = &mockOverlay{}
		return overlay
	}

	newView := CreateNewView(width, height, depth, create)
	view := newView()
	if view == nil {
		t.Fatal("CreateNewView returned nil view")
	}
	if overlay == nil {
		t.Fatal("createImage was not invoked")
	}
	if view.depth != depth {
		t.Fatalf("unexpected depth: want %v, got %v", depth, view.depth)
	}

	wantColor := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	if got := toRGBA(overlay.fillColor); got != wantColor {
		t.Fatalf("overlay not filled black: want %v, got %v", wantColor, got)
	}
}

func TestCreateNewViewRejectsInvalidSize(t *testing.T) {
	create := func(w, h int) OverlayImage {
		return &mockOverlay{}
	}

	testCases := map[string]func(){
		"width zero": func() {
			CreateNewView(0, 1, frontend.DepthWindow, create)
		},
		"height zero": func() {
			CreateNewView(1, 0, frontend.DepthWindow, create)
		},
		"creator nil": func() {
			CreateNewView(1, 1, frontend.DepthWindow, nil)
		},
	}

	for name, fn := range testCases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("expected panic on invalid argument")
				}
			}()
			fn()
		})
	}
}

func TestViewDrawSkipsWhenTransparent(t *testing.T) {
	view := &View{
		overlay: &mockOverlay{},
		depth:   frontend.DepthWindow,
	}
	var called bool
	draw := func(fn frontend.DrawArgFunc, depth frontend.Depth) {
		called = true
	}

	view.Draw(draw, 0.0)
	if called {
		t.Fatal("expected no draw when rate <= 0")
	}
	if view.overlay.(*mockOverlay).drawCount != 0 {
		t.Fatal("overlay should not be drawn at zero rate")
	}

	view.Draw(draw, -0.5)
	if view.overlay.(*mockOverlay).drawCount != 0 {
		t.Fatal("overlay should not be drawn at negative rate")
	}
}

func TestViewDrawAppliesAlpha(t *testing.T) {
	overlay := &mockOverlay{}
	view := &View{
		overlay: overlay,
		depth:   frontend.DepthWindow,
	}

	var drawnDepth frontend.Depth
	draw := func(fn frontend.DrawArgFunc, depth frontend.Depth) {
		drawnDepth = depth
		fn(nil)
	}

	view.Draw(draw, 0.5)
	if overlay.drawCount != 1 {
		t.Fatalf("expected overlay drawn once, got %d", overlay.drawCount)
	}
	if drawnDepth != frontend.DepthWindow {
		t.Fatalf("unexpected depth: want %v, got %v", frontend.DepthWindow, drawnDepth)
	}
	if overlay.lastOp == nil {
		t.Fatal("expected draw options to be recorded")
	}
	const epsilon = 1e-6
	if got := overlay.lastOp.ColorScale.A(); absFloat(float64(got-0.5)) > epsilon {
		t.Fatalf("unexpected alpha scale: want 0.5, got %f", got)
	}

	view.Draw(draw, 1.5)
	if overlay.drawCount != 2 {
		t.Fatalf("expected second draw, got %d", overlay.drawCount)
	}
	if got := overlay.lastOp.ColorScale.A(); absFloat(float64(got-1.0)) > epsilon {
		t.Fatalf("expected alpha scale clamped to 1, got %f", got)
	}
}

func toRGBA(c color.Color) color.RGBA {
	if c == nil {
		return color.RGBA{}
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
