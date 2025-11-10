package command

import "github.com/asragi/yasoba-prototype/common/character/hero"

type commonMPManager interface {
	InitialMP() hero.InitialMP
	RecoverMP() hero.RecoverMP
	MaxMP() hero.MaxMP
}

type battleMPManagerImpl struct {
	commonMPManager commonMPManager
	currentMP       hero.MP
}

func NewMpManager(commonManager commonMPManager) *battleMPManagerImpl {
	return &battleMPManagerImpl{
		commonMPManager: commonManager,
	}
}

func (m *battleMPManagerImpl) InitialMP() hero.InitialMP {
	return m.commonMPManager.InitialMP()
}

func (m *battleMPManagerImpl) CurrentMP() hero.MP {
	return m.currentMP
}

func (m *battleMPManagerImpl) Recover() {
	commonManager := m.commonMPManager
	recoverMp := commonManager.RecoverMP()
	mapMp := commonManager.MaxMP()
	m.currentMP = m.currentMP.Recover(recoverMp, mapMp)
}

func (m *battleMPManagerImpl) Consume(cost Cost) {
	m.currentMP = cost.Consume(m.currentMP)
}

func (m *battleMPManagerImpl) MaxMP() hero.MaxMP {
	return m.commonMPManager.MaxMP()
}
