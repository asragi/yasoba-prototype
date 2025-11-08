package input

type Receiver interface {
	OnInputSubmit()
	OnInputCancel()
	OnInputSubButton()
	OnInputUp()
	OnInputDown()
	OnInputLeft()
	OnInputRight()
}

type Manager interface {
	Update()
	Set(Receiver)
}

type EmptyReceiver struct{}

var EmptyInstance = &EmptyReceiver{}

func (i *EmptyReceiver) OnInputSubmit()    {}
func (i *EmptyReceiver) OnInputCancel()    {}
func (i *EmptyReceiver) OnInputSubButton() {}
func (i *EmptyReceiver) OnInputUp()        {}
func (i *EmptyReceiver) OnInputDown()      {}
func (i *EmptyReceiver) OnInputLeft()      {}
func (i *EmptyReceiver) OnInputRight()     {}
