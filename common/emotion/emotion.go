package emotion

// EmotionType represents the emotional state used across battle UI and logic.
type EmotionType int

const (
	EmotionNormal EmotionType = iota
	EmotionDamage
	EmotionSmile
	EmotionAngry
	EmotionAnnoyed
)
