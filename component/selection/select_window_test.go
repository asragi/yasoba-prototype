package selection

import (
	"image/color"
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

type mockText struct {
	updateCalls []*frontend.Vector
	drawCount   int
}

func (m *mockText) Update(parentPosition *frontend.Vector) {
	m.updateCalls = append(m.updateCalls, parentPosition)
}

func (m *mockText) Draw(frontend.DrawFunc) {
	m.drawCount++
}

func (m *mockText) ForceComplete() {}

func (m *mockText) SetText(string, bool) {}

func (m *mockText) Size() *frontend.Vector {
	return frontend.VectorZero
}

func (m *mockText) CheckIsEnd() bool {
	return true
}

func (m *mockText) SetTextColor(color.Color) {}

func newDummyTextServer() text.ServeTextDataFunc {
	return func(id text.TextId) *text.Data {
		return &text.Data{
			Id:   id,
			Text: text.String("dummy"),
		}
	}
}

func TestSelectWindow_DrawBeforeUpdate(t *testing.T) {
	var createdTexts []*mockText
	mockTextFactory := func(*widget.TextOptionsNew) widget.TextInterface {
		text := &mockText{}
		createdTexts = append(createdTexts, text)
		return text
	}

	newSelectWindow := StandByNewSelectWindow(
		newDummyCursor(frontend.VectorZero),
		mockTextFactory,
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

	selectWindow.Draw(func(frontend.DrawArgFunc, frontend.Depth) {})

	if len(createdTexts) == 0 {
		t.Fatalf("expected mock text to be created")
	}
	if createdTexts[0].drawCount == 0 {
		t.Errorf("expected mock text Draw to be called at least once")
	}
	if len(createdTexts[0].updateCalls) != 0 {
		t.Errorf("expected Update not to be called before Draw, got %d calls", len(createdTexts[0].updateCalls))
	}
}
