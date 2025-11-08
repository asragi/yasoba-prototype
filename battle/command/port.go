package command

type GetCommandModelPort func() []Model
type GetCommandModelFunc func(Id) Model

func CreateGetCommandModel(port GetCommandModelPort) GetCommandModelFunc {
	dict := make(map[Id]Model)
	for _, model := range port() {
		dict[model.ID()] = model
	}
	return func(id Id) Model {
		return dict[id]
	}
}
