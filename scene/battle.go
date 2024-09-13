package scene

import (
	"fmt"
	"github.com/asragi/yasoba-prototype/component"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game"
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
	faceSubWindow      *component.FaceWindow
	enemyData          []*core.EnemyIdPair
}

func (s *BattleScene) Update() {
	s.messageWindow.Update(frontend.VectorZero)
	s.faceWindow.Update(&frontend.Vector{X: 0, Y: 288})
	s.faceSubWindow.Update(&frontend.Vector{X: 384, Y: 288})
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
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		s.battleSelectWindow.OnSubmit()
	}
}

func (s *BattleScene) Draw(drawFunc frontend.DrawFunc) {
	s.messageWindow.Draw(drawFunc)
	s.battleSelectWindow.Draw(drawFunc)
	s.faceWindow.Draw(drawFunc)
	s.faceSubWindow.Draw(drawFunc)
}

type BattleOption struct {
	OnEnd           OnEndBattle
	BattleSettingId game.BattleSettingId
}

type NewBattleScene func(*BattleOption) *BattleScene

func StandByNewBattleScene(
	newMessageWindow component.NewMessageWindowFunc,
	newBattleSelectWindow component.NewBattleSelectWindowFunc,
	newFaceWindow component.NewFaceWindowFunc,
	initializeBattle core.InitializeBattleFunc,
	postCommand core.PostCommandFunc,
	getBattleSetting game.ServeBattleSetting,
) NewBattleScene {
	return func(option *BattleOption) *BattleScene {
		battleSetting := getBattleSetting(option.BattleSettingId)
		enemyIds := func() []core.EnemyId {
			ids := make([]core.EnemyId, len(battleSetting.Enemies))
			for i, set := range battleSetting.Enemies {
				ids[i] = set.EnemyId
			}
			return ids
		}()
		initializeRequest := &core.InitializeBattleRequest{
			// TODO: variables must be provided by args
			MainActorCharacterId: core.CharacterLuneId,
			SubActorCharacterId:  core.CharacterSunnyId,
			EnemyIds:             enemyIds,
		}
		battleResponse := initializeBattle(initializeRequest)

		messageWindow := newMessageWindow(
			&frontend.Vector{X: 192, Y: 0},
			&frontend.Vector{X: 292, Y: 62},
			frontend.DepthWindow,
			frontend.PivotTopCenter,
		)
		testString := "あのイーハトーヴォのすきとおった風\n夏でも底に冷たさをもつ青いそら\nうつくしい森で飾られたモリーオ市"
		messageWindow.SetText(testString, false)

		var battleSelectWindow *component.BattleSelectWindow
		onSubmit := func(command core.PlayerCommand) {
			response := postCommand(
				&core.PostCommandRequest{
					ActorId:  core.ActorLuneId,
					TargetId: nil,
					Command:  command,
				},
			)
			fmt.Printf("response: %+v\n", response)
			battleSelectWindow.Close()
		}
		battleSelectWindow = newBattleSelectWindow(
			&frontend.Vector{X: 0, Y: 0},
			frontend.PivotBottomLeft,
			frontend.DepthWindow,
			[]core.PlayerCommand{
				core.PlayerCommandAttack,
				core.PlayerCommandFire,
				core.PlayerCommandFocus,
				core.PlayerCommandDefend,
			},
			onSubmit,
		)
		battleSelectWindow.Open()

		faceWindow := newFaceWindow(
			&frontend.Vector{X: 0, Y: 0},
			frontend.DepthWindow,
			frontend.PivotBottomLeft,
			frontend.TextureFaceLuneNormal,
		)
		faceSubWindow := newFaceWindow(
			&frontend.Vector{X: 0, Y: 0},
			frontend.DepthWindow,
			frontend.PivotBottomRight,
			frontend.TextureFaceSunnyNormal,
		)
		return &BattleScene{
			messageWindow:      messageWindow,
			battleSelectWindow: battleSelectWindow,
			faceWindow:         faceWindow,
			faceSubWindow:      faceSubWindow,
			enemyData:          battleResponse.EnemyIds,
		}
	}
}
