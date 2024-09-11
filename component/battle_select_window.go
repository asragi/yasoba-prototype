package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type BattleCommand int

const (
	BattleCommandAttack BattleCommand = iota
	BattleCommandFire
	BattleCommandThunder
	BattleCommandBarrier
	BattleCommandWind
	BattleCommandFocus
	BattleCommandDefend
)

func (b *BattleCommand) ToTextId() core.TextId {
	switch *b {
	case BattleCommandAttack:
		return "battle_command_attack"
	case BattleCommandFire:
		return "battle_command_fire"
	case BattleCommandThunder:
		return "battle_command_thunder"
	case BattleCommandBarrier:
		return "battle_command_barrier"
	case BattleCommandWind:
		return "battle_command_wind"
	case BattleCommandFocus:
		return "battle_command_focus"
	case BattleCommandDefend:
		return "battle_command_defend"
	}
	return ""
}

type BattleSelectWindow struct {
	commands []BattleCommand
	texts    []*widget.Text
	index    int
	isActive bool
	isOpen   bool
}

func (w *BattleSelectWindow) Open() {
	w.isOpen = true
}

func (w *BattleSelectWindow) Update(parentPosition *frontend.Vector) {
	for _, text := range w.texts {
		text.Update(parentPosition)
	}
}

func (w *BattleSelectWindow) Draw(drawFunc frontend.DrawFunc) {
	if !w.isOpen {
		return
	}
	for _, text := range w.texts {
		text.Draw(drawFunc)
	}
}

type NewBattleSelectWindowFunc func(
	*frontend.Vector,
	*frontend.Pivot,
	frontend.Depth,
	[]BattleCommand,
) *BattleSelectWindow

func StandByNewBattleSelectWindow(
	resource *frontend.ResourceManager,
	textServer core.ServeTextDataFunc,
) NewBattleSelectWindowFunc {
	font := resource.GetFont(frontend.MaruMinya)
	return func(
		relativePosition *frontend.Vector,
		pivot *frontend.Pivot,
		depth frontend.Depth,
		commands []BattleCommand,
	) *BattleSelectWindow {
		texts := func() []*widget.Text {
			// TODO: Use actual values
			const lineHeight = 16
			const width = 64
			count := len(commands)
			size := &frontend.Vector{X: width, Y: float64(lineHeight * count)}
			pivotModification := pivot.ApplyToSize(size)
			relativePositions := func() []*frontend.Vector {
				var positions []*frontend.Vector
				for i := 0; i < count; i++ {
					positions = append(
						positions, &frontend.Vector{
							X: relativePosition.X - pivotModification.X,
							Y: relativePosition.Y - pivotModification.Y + float64(lineHeight*i),
						},
					)
				}
				return positions
			}()
			var texts []*widget.Text
			for i, command := range commands {
				text := widget.NewText(
					&widget.TextOptions{
						RelativePosition: relativePositions[i],
						Pivot:            frontend.PivotTopLeft,
						TextFace:         font,
						Speed:            1,
						Depth:            depth,
					},
				)
				text.SetText(textServer(command.ToTextId()).Text, true)
				texts = append(texts, text)
			}
			return texts
		}()

		return &BattleSelectWindow{
			texts: texts,
		}
	}
}
