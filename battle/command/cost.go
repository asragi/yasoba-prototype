package command

import "github.com/asragi/yasoba-prototype/common/character/hero"

type Cost int

func (c Cost) HasEnoughMP(mp hero.MP) bool {
	return hero.MP(c) <= mp
}

func (c Cost) Consume(mp hero.MP) hero.MP {
	if mp < hero.MP(c) {
		return 0
	}
	return mp - hero.MP(c)
}
