package command

type Id string

type Model interface {
	ID() Id
	Cost() Cost
}

type commandModel struct {
	id   Id
	cost Cost
}

func NewModel(id Id, cost Cost) Model {
	return &commandModel{
		id:   id,
		cost: cost,
	}
}

func (m *commandModel) ID() Id {
	return m.id
}

func (m *commandModel) Cost() Cost {
	return m.cost
}
