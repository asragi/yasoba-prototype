package app

import (
	ebiteninput "github.com/asragi/yasoba-prototype/adapter/ebiten/input"
	"github.com/asragi/yasoba-prototype/toolkit/input"
)

func makeInputManager() input.InputManager {
	return &ebiteninput.KeyBoardInput{}
}
