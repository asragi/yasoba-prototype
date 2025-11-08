package input

import (
	"github.com/asragi/yasoba-prototype/toolkit/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyBoardInput struct {
	receiver input.Receiver
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

func (k *KeyBoardInput) Set(receiver input.Receiver) {
	k.receiver = receiver
}
