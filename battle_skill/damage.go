package battle_skill

import (
	"math"
	"strconv"

	"github.com/asragi/yasoba-prototype/actor"
)

// Damage represents damage dealt to a target.
type Damage int

// String renders the damage value as text.
func (d Damage) String() string {
	return strconv.Itoa(int(d))
}

// Apply computes the remaining HP after taking damage.
func (d Damage) Apply(hp actor.HP) actor.HP {
	return actor.HP(math.Max(0, float64(hp)-float64(d)))
}
