package animation

import "github.com/asragi/yasoba-prototype/common/texture"

// AnimationData describes how to play a sprite-sheet animation.
type AnimationData struct {
	TextureID      texture.ID
	RowCount       int
	ColumnCount    int
	AnimationCount int
	Duration       int
	IsLoop         bool
}
