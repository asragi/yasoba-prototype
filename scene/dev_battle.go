package scene

type battlePhase int

const (
	battlePhaseStart battlePhase = iota
	battlePhaseAction
	battlePhaseCommand
	battlePhaseEnd
)

type BattleOption struct {
	OnEnd OnEndBattle
}

type BattleSequence struct {
	frame           int
	phase           battlePhase
	openingSequence func()
	onEnd           OnEndBattle
}

type BattleResult struct{}

type OnEndBattle func(BattleResult)

type BattleDrawer interface {
}

func NewBattleSequencer(option BattleOption, drawer BattleDrawer) *BattleSequence {
	sequencer := &BattleSequence{
		frame: 0,
		onEnd: option.OnEnd,
	}
	openingSequence := newBattleOpeningPhaseSequence(sequencer.ToAction)
	sequencer.openingSequence = openingSequence
	return sequencer
}

func (s *BattleSequence) ToAction() {
	s.phase = battlePhaseAction
}

func (s *BattleSequence) Update() {
	if s.phase == battlePhaseStart {
		s.openingSequence()
		return
	}
	if s.phase == battlePhaseAction {
		return
	}
	if s.phase == battlePhaseCommand {
		return
	}
}

func newBattleOpeningPhaseSequence(onEnd func()) func() {
	innerFrame := 0
	displayInitialInformation := func() {}

	sequence := make(map[int]func())
	sequence[0] = displayInitialInformation
	sequence[60] = onEnd

	return func() {
		innerFrame++
		if _, ok := sequence[innerFrame]; !ok {
			return
		}
		sequence[innerFrame]()
	}
}
