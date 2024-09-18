package image

import _ "embed"

var (
	//go:embed cursor.png
	Cursor []byte
	//go:embed window.png
	Window []byte
	//go:embed face_lune_normal.png
	FaceLuneNormal []byte
	//go:embed face_sunny_normal.png
	FaceSunnyNormal []byte
	//go:embed marshmallow_normal.png
	MarshmallowNormal []byte
	//go:embed marshmallow_damage.png
	MarshmallowDamage []byte
	//go:embed battle_effect_impact.png
	BattleEffectImpact []byte
)
