package character

import "strconv"

// MaxHP represents the maximum HP for an actor.
type MaxHP int

// ToHP converts the maximum HP to a usable HP value.
func (h MaxHP) ToHP() HP {
	return HP(h)
}

// HP tracks an actor's current HP.
type HP int

// String returns the HP as a decimal string.
func (h HP) String() string {
	return strconv.Itoa(int(h))
}

// HPRatio represents the ratio of current HP to maximum HP.
type HPRatio float64

// Ratio calculates the current HP ratio against the provided max HP.
func (h HP) Ratio(max MaxHP) HPRatio {
	if max <= 0 {
		return 0
	}
	return HPRatio(float64(h) / float64(max))
}

// Float64 exposes the ratio as a float64 value.
func (r HPRatio) Float64() float64 {
	return float64(r)
}

type (
	ATK int
	MAG int
	DEF int
	SPD int
)
