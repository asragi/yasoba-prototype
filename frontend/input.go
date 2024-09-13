package frontend

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

type KeyBoardInput struct {
	receiver InputReceiver
}

func (k *KeyBoardInput) Update() {
	checkButton := func() {
		if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
			k.receiver.OnInputSubmit()
			return
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyX) {
			k.receiver.OnInputCancel()
			return
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			k.receiver.OnInputSubButton()
			return
		}
	}

	checkButton()
	if inpututil.KeyPressDuration(ebiten.KeyUp) >= 1 {
		k.receiver.OnInputUp()
	}
	if inpututil.KeyPressDuration(ebiten.KeyDown) >= 1 {
		k.receiver.OnInputDown()
	}
	if inpututil.KeyPressDuration(ebiten.KeyLeft) >= 1 {
		k.receiver.OnInputLeft()
	}
	if inpututil.KeyPressDuration(ebiten.KeyRight) >= 1 {
		k.receiver.OnInputRight()
	}
}

func (k *KeyBoardInput) Set(receiver InputReceiver) {
	k.receiver = receiver
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
