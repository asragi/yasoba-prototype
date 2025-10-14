package core

import (
	"testing"

	"github.com/asragi/yasoba-prototype/actor"
)

func TestHPRatioDelegatesToMaxHP(t *testing.T) {
	hp := actor.HP(90)
	max := actor.MaxHP(150)
	got := hp.Ratio(max)
	want := actor.HPRatio(0.6)
	if got != want {
		t.Fatalf("hp.Ratio() = %v, want %v", got, want)
	}
}

func TestHPRatioZeroMax(t *testing.T) {
	hp := actor.HP(50)
	max := actor.MaxHP(0)
	if hp.Ratio(max) != 0 {
		t.Fatalf("hp.Ratio() with zero max should be 0")
	}
}

func TestHPRatioFloat64(t *testing.T) {
	r := actor.HPRatio(0.75)
	if r.Float64() != 0.75 {
		t.Fatalf("HPRatio.Float64() = %v, want 0.75", r.Float64())
	}
}
