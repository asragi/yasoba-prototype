package scene

import (
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

type BattleScene struct {
	messageWindow      *component.MessageWindow
	battleSelectWindow *component.BattleSelectWindow
	faceWindow         *component.FaceWindow
}

func (s *BattleScene) Update() {
	s.messageWindow.Update(frontend.VectorZero)
	s.faceWindow.Update(&frontend.Vector{X: 0, Y: 288})
	s.battleSelectWindow.Update(s.faceWindow.GetTopLeftPosition())
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.messageWindow.Shake(2, 10)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		s.battleSelectWindow.MoveCursorDown()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		s.battleSelectWindow.MoveCursorUp()
	}
}

func (s *BattleScene) Draw(drawFunc frontend.DrawFunc) {
	s.messageWindow.Draw(drawFunc)
	s.battleSelectWindow.Draw(drawFunc)
	s.faceWindow.Draw(drawFunc)
}

type NewBattleScene func() *BattleScene

func StandByNewBattleScene(
	newMessageWindow component.NewMessageWindowFunc,
	newBattleSelectWindow component.NewBattleSelectWindowFunc,
	newFaceWindow component.NewFaceWindowFunc,
) NewBattleScene {
	return func() *BattleScene {
		messageWindow := newMessageWindow(
			&frontend.Vector{X: 192, Y: 0},
			&frontend.Vector{X: 292, Y: 62},
			frontend.DepthWindow,
			frontend.PivotTopCenter,
		)
		testString := "あのイーハトーヴォのすきとおった風\n夏でも底に冷たさをもつ青いそら\nうつくしい森で飾られたモリーオ市"
		messageWindow.SetText(testString, false)
		battleSelectWindow := newBattleSelectWindow(
			&frontend.Vector{X: 0, Y: 0},
			frontend.PivotBottomLeft,
			frontend.DepthWindow,
			[]component.BattleCommand{
				component.BattleCommandAttack,
				component.BattleCommandFire,
				/*
					component.BattleCommandThunder,
					component.BattleCommandBarrier,
					component.BattleCommandWind,
				*/
				component.BattleCommandFocus,
				component.BattleCommandDefend,
			},
		)
		battleSelectWindow.Open()
		faceWindow := newFaceWindow(
			&frontend.Vector{X: 0, Y: 0},
			frontend.DepthWindow,
			frontend.PivotBottomLeft,
		)
		return &BattleScene{
			messageWindow:      messageWindow,
			battleSelectWindow: battleSelectWindow,
			faceWindow:         faceWindow,
		}
	}
}
