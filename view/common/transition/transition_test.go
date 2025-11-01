package transition

import (
	"math"
	"testing"

	"github.com/asragi/yasoba-prototype/toolkit/drawing"
)

func TestTransitionUpdateAndDraw(t *testing.T) {
	const maxFrame = 4

	rates := make([]float64, 0, maxFrame+1)
	tr := New(maxFrame, func(drawFunc drawing.DrawFunc, rate float64) {
		rates = append(rates, rate)
	}, InitialStateTransparent)
	drawFunc := func(fn drawing.DrawArgFunc, depth drawing.Depth) {}

	tr.FadeOut()
	rates = rates[:0]
	for i := 1; i <= maxFrame; i++ {
		tr.Update()
		tr.Draw(drawFunc)

		if len(rates) == 0 {
			t.Fatalf("expected Draw to record a rate during fade out")
		}
		expected := float64(i) / float64(maxFrame)
		got := rates[len(rates)-1]
		if math.Abs(got-expected) > 1e-9 {
			t.Fatalf("fade out step %d: expected rate %.2f, got %.2f", i, expected, got)
		}
	}

	tr.Update()
	tr.Draw(drawFunc)
	if got := rates[len(rates)-1]; math.Abs(got-1.0) > 1e-9 {
		t.Fatalf("expected rate to stay at 1 after fade out completes, got %.2f", got)
	}

	tr.FadeIn()
	rates = rates[:0]
	for i := maxFrame - 1; i >= 0; i-- {
		tr.Update()
		tr.Draw(drawFunc)

		if len(rates) == 0 {
			t.Fatalf("expected Draw to record a rate during fade in")
		}
		expected := float64(i) / float64(maxFrame)
		got := rates[len(rates)-1]
		if math.Abs(got-expected) > 1e-9 {
			t.Fatalf("fade in step %d: expected rate %.2f, got %.2f", maxFrame-i, expected, got)
		}
	}

	tr.Update()
	tr.Draw(drawFunc)
	if got := rates[len(rates)-1]; math.Abs(got-0.0) > 1e-9 {
		t.Fatalf("expected rate to stay at 0 after fade in completes, got %.2f", got)
	}
}

func TestTransitionInitialStateOpaque(t *testing.T) {
	const maxFrame = 4

	rates := make([]float64, 0, 1)
	tr := New(maxFrame, func(drawFunc drawing.DrawFunc, rate float64) {
		rates = append(rates, rate)
	}, InitialStateOpaque)

	drawFunc := func(fn drawing.DrawArgFunc, depth drawing.Depth) {}

	tr.Draw(drawFunc)

	if len(rates) == 0 {
		t.Fatalf("expected rate to be recorded")
	}
	if got := rates[0]; math.Abs(got-1.0) > 1e-9 {
		t.Fatalf("expected initial opaque rate 1.0, got %.2f", got)
	}

	tr.FadeIn()
	tr.Update()
	rates = rates[:0]
	tr.Draw(drawFunc)
	if len(rates) == 0 {
		t.Fatalf("expected rate to be recorded after fade in update")
	}
	if got := rates[0]; got >= 1.0 {
		t.Fatalf("expected rate to decrease after fade in, got %.2f", got)
	}
}
