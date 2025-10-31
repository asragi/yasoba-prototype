package selection

import (
	"image/color"
	"testing"

	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/text"
	"github.com/asragi/yasoba-prototype/widget"
)

type dummyCursor struct {
	size             *drawing.Vector
	relativePosition *drawing.Vector
}

func newDummyCursor(size *drawing.Vector) func(*drawing.Vector, *drawing.Pivot, drawing.Depth) Cursor {
	return func(relativePosition *drawing.Vector, _ *drawing.Pivot, _ drawing.Depth) Cursor {
		return &dummyCursor{
			size:             size,
			relativePosition: relativePosition,
		}
	}
}

func (c *dummyCursor) Update(*drawing.Vector) {}

func (c *dummyCursor) Draw(drawing.DrawFunc) {}

func (c *dummyCursor) SetRelativePosition(position *drawing.Vector) {
	c.relativePosition = position
}

func (c *dummyCursor) Size() *drawing.Vector {
	return c.size
}

type mockText struct {
	updateCalls []*drawing.Vector
	drawCount   int
}

func (m *mockText) Update(parentPosition *drawing.Vector) {
	m.updateCalls = append(m.updateCalls, parentPosition)
}

func (m *mockText) Draw(drawing.DrawFunc) {
	m.drawCount++
}

func (m *mockText) ForceComplete() {}

func (m *mockText) SetText(string, bool) {}

func (m *mockText) Size() *drawing.Vector {
	return drawing.VectorZero
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
		newDummyCursor(drawing.VectorZero),
		mockTextFactory,
		newDummyTextServer(),
	)
	selectWindow := newSelectWindow(
		drawing.VectorZero,
		drawing.PivotTopLeft,
		drawing.DepthWindow,
		[]text.TextId{text.TextIdBattleCommandAttack},
		func(int) {},
		true,
	)
	selectWindow.Open()

	selectWindow.Draw(func(drawing.DrawArgFunc, drawing.Depth) {})

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
