package selection

import (
	"testing"

	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/widget"
)

type dummyCursor struct {
	size             *frontend.Vector
	relativePosition *frontend.Vector
}

func newDummyCursor(size *frontend.Vector) func(*frontend.Vector, *frontend.Pivot, frontend.Depth) Cursor {
	return func(relativePosition *frontend.Vector, _ *frontend.Pivot, _ frontend.Depth) Cursor {
		return &dummyCursor{
			size:             size,
			relativePosition: relativePosition,
		}
	}
}

func (c *dummyCursor) Update(*frontend.Vector) {}

func (c *dummyCursor) Draw(frontend.DrawFunc) {}

func (c *dummyCursor) SetRelativePosition(position *frontend.Vector) {
	c.relativePosition = position
}

func (c *dummyCursor) Size() *frontend.Vector {
	return c.size
}

func newDummyTextServer() text.ServeTextDataFunc {
	return func(id text.TextId) *text.Data {
		return &text.Data{
			Id:   id,
			Text: text.String("dummy"),
		}
	}
}

func newNoopDrawFunc() frontend.DrawFunc {
	return func(fn frontend.DrawArgFunc, depth frontend.Depth) {}
}

func TestSelectWindow_DrawBeforeUpdate(t *testing.T) {
	resourceManager, err := frontend.CreateResourceManager()
	if err != nil {
		t.Fatalf("failed to create resource manager: %v", err)
	}
	newText := widget.CreateNewText(resourceManager)
	newSelectWindow := StandByNewSelectWindow(
		newDummyCursor(frontend.VectorZero),
		newText,
		newDummyTextServer(),
	)
	selectWindow := newSelectWindow(
		frontend.VectorZero,
		frontend.PivotTopLeft,
		frontend.DepthWindow,
		[]text.TextId{text.TextIdBattleCommandAttack},
		func(int) {},
		true,
	)
	selectWindow.Open()

	selectWindow.Draw(newNoopDrawFunc())
}
