package mp

import (
	"github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
)

// widget.Image の代替インターフェース
type Image interface {
	Update(parentPosition *drawing.Vector)
	Draw(drawing.DrawFunc)
}

type GetImageFunc func(texture.ID) drawing.Image

type view struct {
	icons          []Image
	emptyIcons     []Image
	parentPosition *drawing.Vector
	maxMp          int
	currentMp      int
}

type ImageFactory func(
	*drawing.Vector,
	*drawing.Pivot,
	drawing.Depth,
	drawing.Image,
) Image

var pivot = drawing.PivotCenterLeft
var depth = drawing.DepthDebug
var margin = 4

func newView(
	newImage ImageFactory,
	getImageData GetImageFunc,
	maxMp int,
) *view {
	iconImgData := getImageData(texture.MPIcon)
	emptyIconData := getImageData(texture.MPIconEmpty)
	icons := make([]Image, maxMp)
	emptyIcons := make([]Image, maxMp)
	for i := 0; i < maxMp; i++ {
		relative := &drawing.Vector{
			X: float64(i) * float64(iconImgData.Bounds().Dx()+margin),
			Y: 0,
		}
		icons[i] = newImage(
			relative,
			pivot,
			depth,
			iconImgData,
		)
		emptyIcons[i] = newImage(
			relative,
			pivot,
			depth,
			emptyIconData,
		)
	}
	return &view{
		emptyIcons:     emptyIcons,
		icons:          icons,
		parentPosition: nil,
		maxMp:          maxMp,
	}
}

func (m *view) update(parentPosition *drawing.Vector, currentMp int) {
	m.parentPosition = parentPosition
	m.currentMp = currentMp
	for i := 0; i < m.maxMp; i++ {
		m.icons[i].Update(parentPosition)
		m.emptyIcons[i].Update(parentPosition)
	}
}

func (m *view) draw(drawFunc drawing.DrawFunc) {
	if m.parentPosition == nil {
		return
	}
	for i := 0; i < m.maxMp; i++ {
		if i < m.currentMp {
			m.icons[i].Draw(drawFunc)
			continue
		}
		m.emptyIcons[i].Draw(drawFunc)
	}
}
