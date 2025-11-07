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

type FontId int

const (
	MaruMinya FontId = iota
)

type ResourceManager struct {
	textureDict   map[commontexture.ID]*ebiten.Image
	fontDict      map[FontId]*text.GoTextFace
	animationDict map[commonanimation.ID]*commonanimation.AnimationData
	shaderDict    map[drawing.ShaderId]*ebiten.Shader
}

type ResourceManagerInterface interface {
	GetTexture(id commontexture.ID) drawing.Image
	GetFont(id FontId) *text.GoTextFace
	GetAnimationData(id commonanimation.ID) *commonanimation.AnimationData
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

func (r *ResourceManager) GetAnimationData(id commonanimation.ID) *commonanimation.AnimationData {
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
		commontexture.Window:              load.Window,
		commontexture.Cursor:              load.Cursor,
		commontexture.FaceLuneNormal:      load.FaceLuneNormal,
		commontexture.FaceLuneDamage:      load.FaceLuneDamage,
		commontexture.FaceSunnyNormal:     load.FaceSunnyNormal,
		commontexture.FaceSunnyDamage:     load.FaceSunnyDamage,
		commontexture.FaceSunnySmile:      load.FaceSunnySmile,
		commontexture.FaceSunnyAngry:      load.FaceSunnyAngry,
		commontexture.FaceSunnyAnnoyed:    load.FaceSunnyAnnoyed,
		commontexture.MarshmallowNormal:   load.MarshmallowNormal,
		commontexture.MarshmallowDamage:   load.MarshmallowDamage,
		commontexture.BattleEffectImpact:  load.BattleEffectImpact,
		commontexture.BattleEffectFire:    load.BattleEffectFire,
		commontexture.BattleEffectExplode: load.BattleEffectExplode,
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
	animationDict := map[commonanimation.ID]*commonanimation.AnimationData{
		commonanimation.LuneNormal: {
			TextureID:      commontexture.FaceLuneNormal,
			RowCount:       1,
			ColumnCount:    2,
			AnimationCount: 2,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.LuneDamage: {
			TextureID:      commontexture.FaceLuneDamage,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.SunnyNormal: {
			TextureID:      commontexture.FaceSunnyNormal,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.SunnyDamage: {
			TextureID:      commontexture.FaceSunnyDamage,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.SunnySmile: {
			TextureID:      commontexture.FaceSunnySmile,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.SunnyAngry: {
			TextureID:      commontexture.FaceSunnyAngry,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.SunnyAnnoyed: {
			TextureID:      commontexture.FaceSunnyAnnoyed,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.MarshmallowNormal: {
			TextureID:      commontexture.MarshmallowNormal,
			RowCount:       1,
			ColumnCount:    2,
			AnimationCount: 2,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.MarshmallowDamage: {
			TextureID:      commontexture.MarshmallowDamage,
			RowCount:       1,
			ColumnCount:    1,
			AnimationCount: 1,
			Duration:       20,
			IsLoop:         true,
		},
		commonanimation.BattleEffectImpact: {
			TextureID:      commontexture.BattleEffectImpact,
			RowCount:       4,
			ColumnCount:    4,
			AnimationCount: 16,
			Duration:       4,
			IsLoop:         false,
		},
		commonanimation.BattleEffectFire: {
			TextureID:      commontexture.BattleEffectFire,
			RowCount:       5,
			ColumnCount:    6,
			AnimationCount: 25,
			Duration:       4,
			IsLoop:         false,
		},
		commonanimation.BattleEffectExplode: {
			TextureID:      commontexture.BattleEffectExplode,
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
