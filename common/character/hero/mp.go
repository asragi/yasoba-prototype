package hero

type MP int

type MaxMP MP

func (max MaxMP) Clamp(m MP) MP {
	if m < 0 {
		return 0
	}
	if m > MP(max) {
		return MP(max)
	}
	return m
}

type RecoverMP MP

func (m MP) Recover(recover RecoverMP, maxMP MaxMP) MP {
	return maxMP.Clamp(m + MP(recover))
}

type InitialMP MP

func (m InitialMP) ToMP() MP {
	return MP(m)
}
