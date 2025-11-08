package command

// Model exposes read-only access to a command definition.
type Model interface {
	ID() Id
	Cost() Cost
}

type commandModel struct {
	id   Id
	cost Cost
}

// NewModel creates a command model from its identifier and MP cost.
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
