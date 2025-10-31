package actor

import (
	"github.com/asragi/yasoba-prototype/component/battle/emotion"
	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/widget"
)

type (
	emotionQueue interface {
		Current() emotion.BattleEmotionType
		Enqueue(emotion.BattleEmotionType)
		Apply(func(emotion.BattleEmotionType) *widget.Animation) *widget.Animation
	}

	animation interface {
		Update(*drawing.Vector)
		Draw(drawing.DrawFunc)
		SetScaleBySize(*drawing.Vector)
	}

	window interface {
		Update(*drawing.Vector)
		Draw(drawing.DrawFunc)
		GetPositionCenter() *drawing.Vector
		GetPositionTopCenter() *drawing.Vector
		GetPositionUpperLeft() *drawing.Vector
		GetPositionLowerRight() *drawing.Vector
	}

	resourceProvider interface {
		GetAnimationData(frontend.AnimationId) *frontend.AnimationData
		GetTexture(frontend.TextureId) drawing.Image
	}

	newWindowFunc    func(*widget.WindowOption) window
	newAnimationFunc func(
		*drawing.Vector,
		*drawing.Pivot,
		drawing.Depth,
		drawing.Image,
		*frontend.AnimationData,
	) animation
	newEmotionQueueFunc func(emotion.BattleEmotionType) emotionQueue
)

type FaceWindow struct {
	emotion emotionQueue
	face    map[emotion.BattleEmotionType]animation
	window  window
}

type NewFaceWindowFunc func(
	*drawing.Vector,
	drawing.Depth,
	*drawing.Pivot,
	character.CharacterId,
) *FaceWindow

func (f *FaceWindow) getCurrentAnimation() animation {
	return f.face[f.emotion.Current()]
}

func (f *FaceWindow) SetEmotion(emotion emotion.BattleEmotionType) {
	f.emotion.Enqueue(emotion)
}

func (f *FaceWindow) Update(parentPosition *drawing.Vector) {
	f.window.Update(parentPosition)
	anim := f.emotion.Apply(func(emotion emotion.BattleEmotionType) *widget.Animation {
		return f.face[emotion].(*widget.Animation)
	})
	anim.Update(f.window.GetPositionCenter())
}

func (f *FaceWindow) Draw(drawFunc drawing.DrawFunc) {
	f.window.Draw(drawFunc)
	f.getCurrentAnimation().Draw(drawFunc)
}

func (f *FaceWindow) GetTopCenterPosition() *drawing.Vector {
	return f.window.GetPositionTopCenter()
}

func (f *FaceWindow) GetTopLeftPosition() *drawing.Vector {
	return f.window.GetPositionUpperLeft()
}

func (f *FaceWindow) GetBottomRightPosition() *drawing.Vector {
	return f.window.GetPositionLowerRight()
}

func (f *FaceWindow) GetCenterPosition() *drawing.Vector {
	return f.window.GetPositionCenter()
}

func StandByNewFaceWindow(
	resource resourceProvider,
	newWindow widget.NewWindowFunc,
) NewFaceWindowFunc {
	return standByNewFaceWindow(
		resource,
		func(option *widget.WindowOption) window {
			return newWindow(option)
		},
		func(
			relativePosition *drawing.Vector,
			pivot *drawing.Pivot,
			depth drawing.Depth,
			texture drawing.Image,
			data *frontend.AnimationData,
		) animation {
			return widget.NewAnimation(
				relativePosition,
				pivot,
				depth,
				texture,
				data,
			)
		},
		func(initial emotion.BattleEmotionType) emotionQueue {
			queue := emotion.NewQueued(initial)
			return &queue
		},
	)
}

func standByNewFaceWindow(
	resource resourceProvider,
	newWindow newWindowFunc,
	newAnimation newAnimationFunc,
	newEmotionQueue newEmotionQueueFunc,
) NewFaceWindowFunc {
	getAllEmotion := createGetAllEmotionFunc()
	allEmotion := getAllEmotion()
	return func(
		relativePosition *drawing.Vector,
		depth drawing.Depth,
		pivot *drawing.Pivot,
		characterId character.CharacterId,
	) *FaceWindow {
		const padding = 6
		const faceSize = 74
		animationMap := func() map[emotion.BattleEmotionType]animation {
			result := map[emotion.BattleEmotionType]animation{}
			for emotion, animationId := range allEmotion[characterId] {
				animationData := resource.GetAnimationData(animationId)
				texture := resource.GetTexture(animationData.TextureId)
				animation := newAnimation(
					drawing.VectorZero,
					drawing.PivotCenter,
					depth,
					texture,
					animationData,
				)
				animation.SetScaleBySize(&drawing.Vector{X: faceSize, Y: faceSize})
				result[emotion] = animation
			}
			return result
		}()
		window := newWindow(
			&widget.WindowOption{
				Texture:          frontend.TextureWindow,
				CornerSize:       6,
				RelativePosition: relativePosition,
				Size:             drawing.NewVector(faceSize, faceSize).Add(drawing.NewVector(padding, padding)),
				Depth:            depth,
				Pivot:            pivot,
			},
		)
		return &FaceWindow{
			emotion: newEmotionQueue(emotion.BattleEmotionNormal),
			face:    animationMap,
			window:  window,
		}
	}
}

type getAllEmotionFunc func() map[character.CharacterId]map[emotion.BattleEmotionType]frontend.AnimationId

func createGetAllEmotionFunc() getAllEmotionFunc {
	dict := map[character.CharacterId]map[emotion.BattleEmotionType]frontend.AnimationId{
		character.CharacterLuneId: {
			emotion.BattleEmotionNormal: frontend.AnimationIdLuneNormal,
			emotion.BattleEmotionDamage: frontend.AnimationIdLuneDamage,
		},
		character.CharacterSunnyId: {
			emotion.BattleEmotionNormal:  frontend.AnimationIdSunnyNormal,
			emotion.BattleEmotionDamage:  frontend.AnimationIdSunnyDamage,
			emotion.BattleEmotionSmile:   frontend.AnimationIdSunnySmile,
			emotion.BattleEmotionAngry:   frontend.AnimationIdSunnyAngry,
			emotion.BattleEmotionAnnoyed: frontend.AnimationIdSunnyAnnoyed,
		},
	}
	return func() map[character.CharacterId]map[emotion.BattleEmotionType]frontend.AnimationId {
		return dict
	}
}
