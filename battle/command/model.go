package command

type Model interface {
	ID() string
	Cost() Cost
}

type commandModel struct {
	id   string
	cost Cost
}

func NewModel(id string, cost Cost) Model {
	return &commandModel{
		id:   id,
		cost: cost,
	}
}

func (m *commandModel) ID() string {
	return m.id
}

func (m *commandModel) Cost() Cost {
	return m.cost
}
