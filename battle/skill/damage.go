package skill

import (
	"math"
	"strconv"

	"github.com/asragi/yasoba-prototype/common/character"
)

// Damage represents damage dealt to a target.
type Damage int

// String renders the damage value as text.
func (d Damage) String() string {
	return strconv.Itoa(int(d))
}

// Apply computes the remaining HP after taking damage.
func (d Damage) Apply(hp character.HP) character.HP {
	return character.HP(math.Max(0, float64(hp)-float64(d)))
}
