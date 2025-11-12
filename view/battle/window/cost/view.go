package cost

import (
	"github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/global"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
)

type Image interface {
	Update(*drawing.Vector)
	Draw(drawing.DrawFunc)
}

type GetImageDataFunc func(texture.ID) drawing.Image

type ImageFactory func(
	*drawing.Vector,
	*drawing.Pivot,
	drawing.Depth,
	drawing.Image,
) Image

type view struct {
	shouldRender   bool
	icons          []Image
	parentPosition *drawing.Vector
}

func newView(
	newImage ImageFactory,
	getImageData GetImageDataFunc,
	cost int,
) viewInterface {
	margin := global.MPMargin
	if cost <= 0 {
		return &view{
			shouldRender: false,
		}
	}
	iconImgData := getImageData(texture.MPIcon)
	icons := make([]Image, cost)
	for i := 0; i < cost; i++ {
		relative := &drawing.Vector{
			X: float64(i * (iconImgData.Bounds().Dx() + int(margin.X))),
			Y: margin.Y,
		}
		icons[i] = newImage(
			relative,
			drawing.PivotCenterLeft,
			drawing.DepthDebug,
			iconImgData,
		)
	}
	return &view{
		shouldRender: true,
		icons:        icons,
	}
}

func (v *view) update(parentPosition *drawing.Vector) {
	if !v.shouldRender {
		return
	}
	v.parentPosition = parentPosition
	for _, icon := range v.icons {
		icon.Update(parentPosition)
	}
}

func (v *view) draw(drawFunc drawing.DrawFunc) {
	if !v.shouldRender {
		return
	}
	for _, icon := range v.icons {
		icon.Draw(drawFunc)
	}
}
