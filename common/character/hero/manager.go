package hero

type mpManager struct {
	initialMP InitialMP
	recovery  RecoverMP
	maxMP     MaxMP
}

func NewCommonMPManager(initialMP InitialMP, recovery RecoverMP, maxMP MaxMP) *mpManager {
	return &mpManager{
		initialMP: initialMP,
		recovery:  recovery,
		maxMP:     maxMP,
	}
}

func (m *mpManager) InitialMP() InitialMP {
	return m.initialMP
}

func (m *mpManager) RecoverMP() RecoverMP {
	return m.recovery
}

func (m *mpManager) MaxMP() MaxMP {
	return m.maxMP
}
