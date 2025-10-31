package actor

import (
	"github.com/asragi/yasoba-prototype/actor"
	battleSkill "github.com/asragi/yasoba-prototype/battle/skill"
	"github.com/asragi/yasoba-prototype/component"
	battleemotion "github.com/asragi/yasoba-prototype/component/battle/emotion"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/game/character"
)

type BattleActorDisplay struct {
	faceWindow       *FaceWindow
	displayDamage    *component.DisplayDamage
	parameterDisplay *BattleParameterDisplay
}

func (d *BattleActorDisplay) SetDamage(damage battleSkill.Damage, afterHP actor.HP) {
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

func (d *BattleActorDisplay) SetEmotion(emotion battleemotion.BattleEmotionType) {
	d.faceWindow.SetEmotion(emotion)
}

type NewBattleActorDisplayFunc func(*actor.Actor) *BattleActorDisplay

func CreateNewBattleActorDisplay(
	newFaceWindow NewFaceWindowFunc,
	newDisplayDamage component.NewDisplayDamageFunc,
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
