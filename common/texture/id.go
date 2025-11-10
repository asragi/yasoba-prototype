package texture

// ID identifies a texture asset that can be used across UI layers.
type ID int

const (
	Window ID = iota
	Cursor
	FaceLuneNormal
	FaceLuneDamage
	FaceSunnyNormal
	FaceSunnyDamage
	FaceSunnySmile
	FaceSunnyAngry
	FaceSunnyAnnoyed
	MarshmallowNormal
	MarshmallowDamage
	BattleEffectImpact
	BattleEffectFire
	BattleEffectExplode
	MPIcon
	MPIconEmpty
)
