package input

type InputReceiver interface {
	OnInputSubmit()
	OnInputCancel()
	OnInputSubButton()
	OnInputUp()
	OnInputDown()
	OnInputLeft()
	OnInputRight()
}

type InputManager interface {
	Update()
	Set(InputReceiver)
}

type InputReceiverEmpty struct{}

var InputReceiverEmptyInstance = &InputReceiverEmpty{}

func (i *InputReceiverEmpty) OnInputSubmit()    {}
func (i *InputReceiverEmpty) OnInputCancel()    {}
func (i *InputReceiverEmpty) OnInputSubButton() {}
func (i *InputReceiverEmpty) OnInputUp()        {}
func (i *InputReceiverEmpty) OnInputDown()      {}
func (i *InputReceiverEmpty) OnInputLeft()      {}
func (i *InputReceiverEmpty) OnInputRight()     {}
