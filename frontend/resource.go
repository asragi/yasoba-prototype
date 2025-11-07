package frontend

import (
	"bytes"
	"fmt"
	"image"

	"github.com/asragi/yasoba-prototype/adapter/ebiten/drawing/adapter"
	"github.com/asragi/yasoba-prototype/assets/font"
	load "github.com/asragi/yasoba-prototype/assets/image"
	commonanimation "github.com/asragi/yasoba-prototype/common/animation"
	commontexture "github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	TextureWindow commontexture.ID = iota
	TextureCursor
	TextureFaceLuneNormal
	TextureFaceLuneDamage
	TextureFaceSunnyNormal
	TextureFaceSunnyDamage
	TextureFaceSunnySmile
	TextureFaceSunnyAngry
	TextureFaceSunnyAnnoyed
	TextureMarshmallowNormal
	TextureMarshmallowDamage
	TextureBattleEffectImpact
	TextureBattleEffectFire
	TextureBattleEffectExplode
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
	AnimationBattleEffectExplode
	AnimationIdLuneNormal
	AnimationIdLuneDamage
	AnimationIdSunnyNormal
	AnimationIdSunnyDamage
	AnimationIdSunnySmile
	AnimationIdSunnyAngry
	AnimationIdSunnyAnnoyed
)

type ResourceManager struct {
	textureDict   map[commontexture.ID]*ebiten.Image
	fontDict      map[FontId]*text.GoTextFace
	animationDict map[AnimationId]*commonanimation.AnimationData
	shaderDict    map[drawing.ShaderId]*ebiten.Shader
}

type ResourceManagerInterface interface {
	GetTexture(id commontexture.ID) drawing.Image
	GetFont(id FontId) *text.GoTextFace
	GetAnimationData(id AnimationId) *commonanimation.AnimationData
	GetShader(id drawing.ShaderId) *drawing.Shader
	NewEmptyImage(width, height int) drawing.Image
}

func (r *ResourceManager) GetTexture(id commontexture.ID) drawing.Image {
	t, ok := r.textureDict[id]
	if !ok {
		panic(fmt.Sprintf("texture not found: %d", id))
	}
	return adapter.NewEbitenImage(t)
}

func (r *ResourceManager) GetFont(id FontId) *text.GoTextFace {
	return r.fontDict[id]
}

func (r *ResourceManager) GetAnimationData(id AnimationId) *commonanimation.AnimationData {
	data, ok := r.animationDict[id]
	if !ok {
		panic(fmt.Sprintf("animation data not found: %d", id))
	}
	return data
}

func (r *ResourceManager) GetShader(id drawing.ShaderId) *drawing.Shader {
	s, ok := r.shaderDict[id]
	if !ok {
		panic(fmt.Sprintf("shader not found: %d", id))
	}
	return drawing.NewShader(s)
}

func (r *ResourceManager) NewEmptyImage(width, height int) drawing.Image {
	return adapter.NewEbitenImage(ebiten.NewImage(width, height))
}

func CreateResourceManager() (*ResourceManager, error) {
	handleError := func(err error) (*ResourceManager, error) {
		return nil, fmt.Errorf("failed to create resource manager: %w", err)
	}
	textureDict := map[commontexture.ID]*ebiten.Image{}
	loadTexture := func(data []byte, id commontexture.ID) error {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
		textureDict[id] = ebiten.NewImageFromImage(img)
		return nil
	}
	shaderDict := map[drawing.ShaderId]*ebiten.Shader{}
	loadShader := func(data []byte, id drawing.ShaderId) error {
		shader, err := ebiten.NewShader(data)
		if err != nil {
			return err
		}
		shaderDict[id] = shader
		return nil
	}
	// TODO: この辺の処理go:generateとかで自動生成したいね
	imageLoadMap := map[commontexture.ID][]byte{
		TextureWindow:              load.Window,
		TextureCursor:              load.Cursor,
		TextureFaceLuneNormal:      load.FaceLuneNormal,
		TextureFaceLuneDamage:      load.FaceLuneDamage,
		TextureFaceSunnyNormal:     load.FaceSunnyNormal,
		TextureFaceSunnyDamage:     load.FaceSunnyDamage,
		TextureFaceSunnySmile:      load.FaceSunnySmile,
		TextureFaceSunnyAngry:      load.FaceSunnyAngry,
		TextureFaceSunnyAnnoyed:    load.FaceSunnyAnnoyed,
		TextureMarshmallowNormal:   load.MarshmallowNormal,
		TextureMarshmallowDamage:   load.MarshmallowDamage,
		TextureBattleEffectImpact:  load.BattleEffectImpact,
		TextureBattleEffectFire:    load.BattleEffectFire,
		TextureBattleEffectExplode: load.BattleEffectExplode,
	}

	for id, data := range imageLoadMap {
		if err := loadTexture(data, id); err != nil {
			return handleError(err)
		}
	}

	fontDict := map[FontId]*text.GoTextFace{}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(font.MaruMinya))
	if err != nil {
		return handleError(err)
	}
	fontDict[MaruMinya] = &text.GoTextFace{Source: s, Size: 12}

	// TODO: 外部ファイルとかから動的に読み込みたい
	animationDict := map[AnimationId]*commonanimation.AnimationData{
		AnimationIdLuneNormal: {
			TextureID:      TextureFaceLuneNormal,
			RowCount:       1,
			ColumnCount:    2,
			AnimationCount: 2,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationIdLuneDamage: {
			TextureID:      TextureFaceLuneDamage,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationIdSunnyNormal: {
			TextureID:      TextureFaceSunnyNormal,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationIdSunnyDamage: {
			TextureID:      TextureFaceSunnyDamage,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationIdSunnySmile: {
			TextureID:      TextureFaceSunnySmile,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationIdSunnyAngry: {
			TextureID:      TextureFaceSunnyAngry,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationIdSunnyAnnoyed: {
			TextureID:      TextureFaceSunnyAnnoyed,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationMarshmallowNormal: {
			TextureID:      TextureMarshmallowNormal,
			RowCount:       1,
			ColumnCount:    2,
			AnimationCount: 2,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationMarshmallowDamage: {
			TextureID:      TextureMarshmallowDamage,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		AnimationBattleEffectImpact: {
			TextureID:      TextureBattleEffectImpact,
			RowCount:       4,
			ColumnCount:    4,
			AnimationCount: 16,
			Duration:       4,
			IsLoop:         false,
		},
		AnimationBattleEffectFire: {
			TextureID:      TextureBattleEffectFire,
			RowCount:       5,
			ColumnCount:    6,
			AnimationCount: 25,
			Duration:       4,
			IsLoop:         false,
		},
		AnimationBattleEffectExplode: {
			TextureID:      TextureBattleEffectExplode,
			RowCount:       5,
			ColumnCount:    6,
			AnimationCount: 30,
			Duration:       4,
			IsLoop:         false,
		},
	}
	if err = loadShader(load.DisappearShader, drawing.ShaderDisappear); err != nil {
		return handleError(err)
	}
	return &ResourceManager{
		textureDict:   textureDict,
		fontDict:      fontDict,
		animationDict: animationDict,
		shaderDict:    shaderDict,
	}, nil
}
