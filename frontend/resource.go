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
)

type FontId int

const (
	MaruMinya FontId = iota
)

type ResourceManager struct {
	textureDict map[TextureId]*ebiten.Image
	fontDict    map[FontId]*text.GoTextFace
}

func (r *ResourceManager) GetTexture(id TextureId) *ebiten.Image {
	return r.textureDict[id]
}

func (r *ResourceManager) GetFont(id FontId) *text.GoTextFace {
	return r.fontDict[id]
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

	fontDict := map[FontId]*text.GoTextFace{}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(font.MaruMinya))
	if err != nil {
		return handleError(err)
	}
	fontDict[MaruMinya] = &text.GoTextFace{Source: s, Size: 12}
	return &ResourceManager{
		textureDict: textureDict,
		fontDict:    fontDict,
	}, nil
}
