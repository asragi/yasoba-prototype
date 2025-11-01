package actor

import (
	anim "github.com/asragi/yasoba-prototype/animation"
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
		Apply(func(emotion.BattleEmotionType) *anim.Animation) *anim.Animation
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
	) *anim.Animation
	newEmotionQueueFunc func(emotion.BattleEmotionType) emotionQueue
)

type FaceWindow struct {
	emotion emotionQueue
	face    map[emotion.BattleEmotionType]*anim.Animation
	window  window
}

type NewFaceWindowFunc func(
	*drawing.Vector,
	drawing.Depth,
	*drawing.Pivot,
	character.CharacterId,
) *FaceWindow

func (f *FaceWindow) getCurrentAnimation() *anim.Animation {
	return f.face[f.emotion.Current()]
}

func (f *FaceWindow) SetEmotion(emotion emotion.BattleEmotionType) {
	f.emotion.Enqueue(emotion)
}

func (f *FaceWindow) Update(parentPosition *drawing.Vector) {
	f.window.Update(parentPosition)
	animation := f.emotion.Apply(func(emotion emotion.BattleEmotionType) *anim.Animation {
		return f.face[emotion]
	})
	animation.Update(f.window.GetPositionCenter())
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
		) *anim.Animation {
			return anim.New(
				func(
					relativePosition *drawing.Vector,
					pivot *drawing.Pivot,
					depth drawing.Depth,
					image drawing.Image,
				) anim.Sprite {
					return widget.NewImage(relativePosition, pivot, depth, image)
				},
				relativePosition,
				pivot,
				depth,
				texture,
				data,
				nil,
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
		animationMap := func() map[emotion.BattleEmotionType]*anim.Animation {
			result := map[emotion.BattleEmotionType]*anim.Animation{}
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
