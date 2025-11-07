package actor

import (
	"github.com/asragi/yasoba-prototype/battle/actor"
	"github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/common/character"
	"github.com/asragi/yasoba-prototype/common/emotion"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/view/battle/damage"
)

type BattleActorDisplay struct {
	faceWindow       *FaceWindow
	displayDamage    *damage.DisplayDamage
	parameterDisplay *BattleParameterDisplay
}

func (d *BattleActorDisplay) SetDamage(damage skill.Damage, afterHP character.HP) {
	d.displayDamage.DisplayDamage(damage)
	d.parameterDisplay.hpDisplay.SetHP(afterHP)
}

func (d *BattleActorDisplay) Update(bottomLeftPosition *drawing.Vector) {
	d.faceWindow.Update(bottomLeftPosition)
	d.displayDamage.Update(d.faceWindow.GetCenterPosition())
	d.parameterDisplay.Update(bottomLeftPosition)
}

func (d *BattleActorDisplay) Draw(drawFunc drawing.DrawFunc) {
	d.faceWindow.Draw(drawFunc)
	d.displayDamage.Draw(drawFunc)
	d.parameterDisplay.Draw(drawFunc)
}

func (d *BattleActorDisplay) GetMainCharacterPosition() *drawing.Vector {
	return d.faceWindow.GetCenterPosition()
}

func (d *BattleActorDisplay) GetMainCharacterTopLeftPosition() *drawing.Vector {
	return d.faceWindow.GetTopLeftPosition()
}

func (d *BattleActorDisplay) SetEmotion(value emotion.EmotionType) {
	d.faceWindow.SetEmotion(value)
}

type NewBattleActorDisplayFunc func(*actor.Actor) *BattleActorDisplay

func CreateNewBattleActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage damage.NewDisplayDamageFunc,
	newParameterDisplay NewBattleParameterDisplayFunc,
) NewBattleActorDisplayFunc {
	return func(actor *actor.Actor) *BattleActorDisplay {
		initialHp := actor.HP
		parameter := newParameterDisplay(initialHp, drawing.PivotBottomLeft)

		return &BattleActorDisplay{
			faceWindow: newFaceWindow(
				&drawing.Vector{X: 0, Y: -parameter.GetHeight()},
				drawing.DepthPlayer,
				drawing.PivotBottomLeft,
				character.CharacterLuneId,
			),
			displayDamage:    newDisplayDamage(),
			parameterDisplay: parameter,
		}
	}
}
