package cost

import "github.com/asragi/yasoba-prototype/toolkit/drawing"

type viewInterface interface {
	update(*drawing.Vector)
	draw(drawing.DrawFunc)
}

type CostDisplay struct {
	view viewInterface
}

type NewCostDisplayFunc func(
	cost int,
) *CostDisplay

func CreateNewCostDisplay(
	newImage ImageFactory,
	getImageData GetImageDataFunc,
	cost int,
) NewCostDisplayFunc {
	return func(cost int) *CostDisplay {
		return &CostDisplay{
			view: newView(
				newImage,
				getImageData,
				cost,
			),
		}
	}
}

func (c *CostDisplay) Update(parentPosition *drawing.Vector) {
	c.view.update(parentPosition)
}

func (c *CostDisplay) Draw(drawFunc drawing.DrawFunc) {
	c.view.draw(drawFunc)
}
