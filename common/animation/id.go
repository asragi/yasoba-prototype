package animation

// ID identifies sprite-sheet animation data shared across UI layers.
type ID int

const (
	MarshmallowNormal ID = iota
	MarshmallowDamage
	BattleEffectImpact
	BattleEffectFire
	BattleEffectExplode
	LuneNormal
	LuneDamage
	SunnyNormal
	SunnyDamage
	SunnySmile
	SunnyAngry
	SunnyAnnoyed
)
