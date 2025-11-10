package mp

import (
	"github.com/asragi/yasoba-prototype/common/character/hero"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
)

type newViewFunc func(
	initialMp int,
	maxMp int,
) viewInterface

type viewInterface interface {
	update(*drawing.Vector, int)
	draw(drawing.DrawFunc)
}

type mpDisplay struct {
	currentMP hero.MP
	view      viewInterface
}

func (m *mpDisplay) SetCurrentMP(mp hero.MP) {
	m.currentMP = mp
}

func (m *mpDisplay) Update(parentPosition *drawing.Vector) {
	currentMpCount := m.currentMP.ToInt()
	m.view.update(parentPosition, currentMpCount)
}

func (m *mpDisplay) Draw(drawFunc drawing.DrawFunc) {
	m.view.draw(drawFunc)
}

func createNewPresenter(
	newView newViewFunc,
) NewMPDisplayFunc {
	return func(
		initialMP hero.InitialMP,
		maxMP hero.MaxMP) *mpDisplay {
		view := newView(
			initialMP.ToInt(),
			maxMP.ToInt(),
		)
		return &mpDisplay{
			currentMP: initialMP.ToMP(),
			view:      view,
		}
	}
}
