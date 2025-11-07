package actor

import (
	commonanimation "github.com/asragi/yasoba-prototype/common/animation"
	"github.com/asragi/yasoba-prototype/common/character"
	commonemotion "github.com/asragi/yasoba-prototype/common/emotion"
	commontexture "github.com/asragi/yasoba-prototype/common/texture"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	viewemotion "github.com/asragi/yasoba-prototype/view/battle/emotion"
	viewanimation "github.com/asragi/yasoba-prototype/view/common/animation"
	"github.com/asragi/yasoba-prototype/widget"
)

type (
	emotionQueue interface {
		Current() commonemotion.EmotionType
		Enqueue(commonemotion.EmotionType)
		Apply(func(commonemotion.EmotionType) *viewanimation.Animation) *viewanimation.Animation
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
		GetAnimationData(commonanimation.ID) *commonanimation.AnimationData
		GetTexture(commontexture.ID) drawing.Image
	}

	newWindowFunc    func(*widget.WindowOption) window
	newAnimationFunc func(
		*drawing.Vector,
		*drawing.Pivot,
		drawing.Depth,
		drawing.Image,
		*commonanimation.AnimationData,
	) *viewanimation.Animation
	newEmotionQueueFunc func(commonemotion.EmotionType) emotionQueue
)

type FaceWindow struct {
	emotion emotionQueue
	face    map[commonemotion.EmotionType]*viewanimation.Animation
	window  window
}

type NewFaceWindowFunc func(
	*drawing.Vector,
	drawing.Depth,
	*drawing.Pivot,
	character.CharacterId,
) *FaceWindow

func (f *FaceWindow) getCurrentAnimation() *viewanimation.Animation {
	return f.face[f.emotion.Current()]
}

func (f *FaceWindow) SetEmotion(value commonemotion.EmotionType) {
	f.emotion.Enqueue(value)
}

func (f *FaceWindow) Update(parentPosition *drawing.Vector) {
	f.window.Update(parentPosition)
	animation := f.emotion.Apply(func(value commonemotion.EmotionType) *viewanimation.Animation {
		return f.face[value]
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
			data *commonanimation.AnimationData,
		) *viewanimation.Animation {
			return viewanimation.New(
				func(
					relativePosition *drawing.Vector,
					pivot *drawing.Pivot,
					depth drawing.Depth,
					image drawing.Image,
				) viewanimation.Sprite {
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
		func(initial commonemotion.EmotionType) emotionQueue {
			queue := viewemotion.NewQueued(initial)
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
		animationMap := func() map[commonemotion.EmotionType]*viewanimation.Animation {
			result := map[commonemotion.EmotionType]*viewanimation.Animation{}
			for emotion, animationId := range allEmotion[characterId] {
				animationData := resource.GetAnimationData(animationId)
				texture := resource.GetTexture(animationData.TextureID)
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
				Texture:          commontexture.Window,
				CornerSize:       6,
				RelativePosition: relativePosition,
				Size:             drawing.NewVector(faceSize, faceSize).Add(drawing.NewVector(padding, padding)),
				Depth:            depth,
				Pivot:            pivot,
			},
		)
		return &FaceWindow{
			emotion: newEmotionQueue(commonemotion.EmotionNormal),
			face:    animationMap,
			window:  window,
		}
	}
}

type getAllEmotionFunc func() map[character.CharacterId]map[commonemotion.EmotionType]commonanimation.ID

func createGetAllEmotionFunc() getAllEmotionFunc {
	dict := map[character.CharacterId]map[commonemotion.EmotionType]commonanimation.ID{
		character.CharacterLuneId: {
			commonemotion.EmotionNormal: commonanimation.LuneNormal,
			commonemotion.EmotionDamage: commonanimation.LuneDamage,
		},
		character.CharacterSunnyId: {
			commonemotion.EmotionNormal:  commonanimation.SunnyNormal,
			commonemotion.EmotionDamage:  commonanimation.SunnyDamage,
			commonemotion.EmotionSmile:   commonanimation.SunnySmile,
			commonemotion.EmotionAngry:   commonanimation.SunnyAngry,
			commonemotion.EmotionAnnoyed: commonanimation.SunnyAnnoyed,
		},
	}
	return func() map[character.CharacterId]map[commonemotion.EmotionType]commonanimation.ID {
		return dict
	}
}
