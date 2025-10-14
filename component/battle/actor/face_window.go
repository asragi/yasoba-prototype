package actor

import (
	battleemotion "github.com/asragi/yasoba-prototype/component/battle/emotion"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/widget"
)

type FaceWindow struct {
	emotion battleemotion.Queued
	face    map[battleemotion.BattleEmotionType]*widget.Animation
	window  widget.WindowInterface
}

type NewFaceWindowFunc func(
	*frontend.Vector,
	frontend.Depth,
	*frontend.Pivot,
	character.CharacterId,
) *FaceWindow

func (f *FaceWindow) getCurrentAnimation() *widget.Animation {
	return f.face[f.emotion.Current()]
}

func (f *FaceWindow) SetEmotion(emotion battleemotion.BattleEmotionType) {
	f.emotion.Enqueue(emotion)
}

func (f *FaceWindow) Update(parentPosition *frontend.Vector) {
	f.window.Update(parentPosition)
	animation := f.emotion.Apply(func(emotion battleemotion.BattleEmotionType) *widget.Animation {
		return f.face[emotion]
	})
	animation.Update(f.window.GetPositionCenter())
}

func (f *FaceWindow) Draw(drawFunc frontend.DrawFunc) {
	f.window.Draw(drawFunc)
	f.getCurrentAnimation().Draw(drawFunc)
}

func (f *FaceWindow) GetTopCenterPosition() *frontend.Vector {
	return f.window.GetPositionTopCenter()
}

func (f *FaceWindow) GetTopLeftPosition() *frontend.Vector {
	return f.window.GetPositionUpperLeft()
}

func (f *FaceWindow) GetBottomRightPosition() *frontend.Vector {
	return f.window.GetPositionLowerRight()
}

func (f *FaceWindow) GetCenterPosition() *frontend.Vector {
	return f.window.GetPositionCenter()
}

func StandByNewFaceWindow(
	resource *frontend.ResourceManager,
	newWindow widget.NewWindowFunc,
) NewFaceWindowFunc {
	getAllEmotion := createGetAllEmotionFunc()
	allEmotion := getAllEmotion()
	return func(
		relativePosition *frontend.Vector,
		depth frontend.Depth,
		pivot *frontend.Pivot,
		characterId character.CharacterId,
	) *FaceWindow {
		const padding = 6
		const faceSize = 74
		animationMap := func() map[battleemotion.BattleEmotionType]*widget.Animation {
			result := map[battleemotion.BattleEmotionType]*widget.Animation{}
			for emotion, animationId := range allEmotion[characterId] {
				animationData := resource.GetAnimationData(animationId)
				texture := resource.GetTexture(animationData.TextureId)
				animation := widget.NewAnimation(
					frontend.VectorZero,
					frontend.PivotCenter,
					depth,
					texture,
					animationData,
				)
				animation.SetScaleBySize(&frontend.Vector{X: faceSize, Y: faceSize})
				result[emotion] = animation
			}
			return result
		}()
		window := newWindow(
			&widget.WindowOption{
				Texture:          frontend.TextureWindow,
				CornerSize:       6,
				RelativePosition: relativePosition,
				Size:             frontend.NewVectorShort(faceSize).Add(&frontend.Vector{X: padding, Y: padding}),
				Depth:            depth,
				Pivot:            pivot,
			},
		)
		return &FaceWindow{
			emotion: battleemotion.NewQueued(battleemotion.BattleEmotionNormal),
			face:    animationMap,
			window:  window,
		}
	}
}

type getAllEmotionFunc func() map[character.CharacterId]map[battleemotion.BattleEmotionType]frontend.AnimationId

func createGetAllEmotionFunc() getAllEmotionFunc {
	dict := map[character.CharacterId]map[battleemotion.BattleEmotionType]frontend.AnimationId{
		character.CharacterLuneId: {
			battleemotion.BattleEmotionNormal: frontend.AnimationIdLuneNormal,
			battleemotion.BattleEmotionDamage: frontend.AnimationIdLuneDamage,
		},
		character.CharacterSunnyId: {
			battleemotion.BattleEmotionNormal:  frontend.AnimationIdSunnyNormal,
			battleemotion.BattleEmotionDamage:  frontend.AnimationIdSunnyDamage,
			battleemotion.BattleEmotionSmile:   frontend.AnimationIdSunnySmile,
			battleemotion.BattleEmotionAngry:   frontend.AnimationIdSunnyAngry,
			battleemotion.BattleEmotionAnnoyed: frontend.AnimationIdSunnyAnnoyed,
		},
	}
	return func() map[character.CharacterId]map[battleemotion.BattleEmotionType]frontend.AnimationId {
		return dict
	}
}
