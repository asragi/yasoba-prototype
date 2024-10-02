package frontend

import "math"

type PositionDelta interface {
	Delta() *Vector
}

type EmitShake struct {
	amplitude float64
	period    int
	frame     int
}

const (
	ShakeDefaultAmplitude = 3
	ShakeDefaultPeriod    = 12
)

func NewShake() *EmitShake {
	return &EmitShake{
		amplitude: 0,
		period:    0,
		frame:     9999,
	}
}

func (e *EmitShake) Delta() *Vector {
	if e.frame > e.period {
		return VectorZero
	}
	x := e.amplitude * math.Cos(float64(e.frame)*2.1)
	y := e.amplitude * math.Sin(float64(e.frame))
	return &Vector{X: x, Y: y}
}

func (e *EmitShake) Shake(amplitude float64, period int) {
	e.frame = 0
	e.amplitude = amplitude
	e.period = period
}

func (e *EmitShake) Update() {
	e.frame++
}
