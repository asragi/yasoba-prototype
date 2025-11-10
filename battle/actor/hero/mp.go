package hero

import (
	"github.com/asragi/yasoba-prototype/battle/command"
	"github.com/asragi/yasoba-prototype/common/character/hero"
)

type MPPort interface {
	CurrentMP() hero.MP
	Recover()
	Consume(cost command.Cost)
	MaxMP() hero.MaxMP
}
