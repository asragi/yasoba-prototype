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
	refresh(int, int)
}

type mpDisplay struct {
	currentMP hero.MP
	view      viewInterface
}

func (m *mpDisplay) Refresh(mp hero.MP) {
	if m.currentMP == mp {
		return
	}
	m.view.refresh(m.currentMP.ToInt(), mp.ToInt())
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
