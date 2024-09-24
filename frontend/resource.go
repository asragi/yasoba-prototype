package frontend

import (
	"bytes"
	"fmt"
	"github.com/asragi/yasoba-prototype/font"
	load "github.com/asragi/yasoba-prototype/image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
)

type TextureId int

const (
	TextureWindow TextureId = iota
	TextureCursor
	TextureFaceLuneNormal
	TextureFaceSunnyNormal
	TextureMarshmallowNormal
	TextureMarshmallowDamage
	TextureBattleEffectImpact
	TextureBattleEffectFire
)

type FontId int

const (
	MaruMinya FontId = iota
)

type AnimationId int

const (
	AnimationMarshmallowNormal AnimationId = iota
	AnimationMarshmallowDamage
	AnimationBattleEffectImpact
	AnimationBattleEffectFire
)

type ResourceManager struct {
	textureDict   map[TextureId]*ebiten.Image
	fontDict      map[FontId]*text.GoTextFace
	animationDict map[AnimationId]*AnimationData
	shaderDict    map[ShaderId]*ebiten.Shader
}

func (r *ResourceManager) GetTexture(id TextureId) *ebiten.Image {
	t, ok := r.textureDict[id]
	if !ok {
		panic(fmt.Sprintf("texture not found: %d", id))
	}
	return t
}

func (r *ResourceManager) GetFont(id FontId) *text.GoTextFace {
	return r.fontDict[id]
}

func (r *ResourceManager) GetAnimationData(id AnimationId) *AnimationData {
	return r.animationDict[id]
}

func (r *ResourceManager) GetShader(id ShaderId) *Shader {
	s, ok := r.shaderDict[id]
	if !ok {
		panic(fmt.Sprintf("shader not found: %d", id))
	}
	return NewShader(s)
}

func CreateResourceManager() (*ResourceManager, error) {
	handleError := func(err error) (*ResourceManager, error) {
		return nil, fmt.Errorf("failed to create resource manager: %w", err)
	}
	textureDict := map[TextureId]*ebiten.Image{}
	loadTexture := func(data []byte, id TextureId) error {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
		textureDict[id] = ebiten.NewImageFromImage(img)
		return nil
	}
	shaderDict := map[ShaderId]*ebiten.Shader{}
	loadShader := func(data []byte, id ShaderId) error {
		shader, err := ebiten.NewShader(data)
		if err != nil {
			return err
		}
		shaderDict[id] = shader
		return nil
	}
	// TODO: この辺の処理go:generateとかで自動生成したいね
	if err := loadTexture(load.Window, TextureWindow); err != nil {
		return handleError(err)
	}
	if err := loadTexture(load.Cursor, TextureCursor); err != nil {
		return handleError(err)
	}
	if err := loadTexture(load.FaceLuneNormal, TextureFaceLuneNormal); err != nil {
		return handleError(err)
	}
	if err := loadTexture(load.FaceSunnyNormal, TextureFaceSunnyNormal); err != nil {
		return handleError(err)
	}
	if err := loadTexture(load.MarshmallowNormal, TextureMarshmallowNormal); err != nil {
		return handleError(err)
	}
	if err := loadTexture(load.MarshmallowDamage, TextureMarshmallowDamage); err != nil {
		return handleError(err)
	}
	if err := loadTexture(load.BattleEffectImpact, TextureBattleEffectImpact); err != nil {
		return handleError(err)
	}
	if err := loadTexture(load.BattleEffectFire, TextureBattleEffectFire); err != nil {
		return handleError(err)
	}

	fontDict := map[FontId]*text.GoTextFace{}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(font.MaruMinya))
	if err != nil {
		return handleError(err)
	}
	fontDict[MaruMinya] = &text.GoTextFace{Source: s, Size: 12}

	animationDict := map[AnimationId]*AnimationData{
		AnimationMarshmallowNormal: {
			TextureId:      TextureMarshmallowNormal,
			RowCount:       1,
			ColumnCount:    2,
			AnimationCount: 2,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationMarshmallowDamage: {
			TextureId:      TextureMarshmallowDamage,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationBattleEffectImpact: {
			TextureId:      TextureBattleEffectImpact,
			RowCount:       4,
			ColumnCount:    4,
			AnimationCount: 16,
			Duration:       4,
			IsLoop:         false,
		},
		AnimationBattleEffectFire: {
			TextureId:      TextureBattleEffectFire,
			RowCount:       5,
			ColumnCount:    6,
			AnimationCount: 25,
			Duration:       4,
			IsLoop:         false,
		},
	}
	if err = loadShader(load.DisappearShader, ShaderDisappear); err != nil {
		return handleError(err)
	}
	return &ResourceManager{
		textureDict:   textureDict,
		fontDict:      fontDict,
		animationDict: animationDict,
		shaderDict:    shaderDict,
	}, nil
}

type AnimationData struct {
	TextureId      TextureId
	RowCount       int
	ColumnCount    int
	AnimationCount int
	Duration       int
	IsLoop         bool
}
