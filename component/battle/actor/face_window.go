package actor

import (
	"github.com/asragi/yasoba-prototype/component/battle/emotion"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/game/character"
	"github.com/asragi/yasoba-prototype/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type (
	emotionQueue interface {
		Current() emotion.BattleEmotionType
		Enqueue(emotion.BattleEmotionType)
		Apply(func(emotion.BattleEmotionType) *widget.Animation) *widget.Animation
	}

	animation interface {
		Update(*frontend.Vector)
		Draw(frontend.DrawFunc)
		SetScaleBySize(*frontend.Vector)
	}

	window interface {
		Update(*frontend.Vector)
		Draw(frontend.DrawFunc)
		GetPositionCenter() *frontend.Vector
		GetPositionTopCenter() *frontend.Vector
		GetPositionUpperLeft() *frontend.Vector
		GetPositionLowerRight() *frontend.Vector
	}

	resourceProvider interface {
		GetAnimationData(frontend.AnimationId) *frontend.AnimationData
		GetTexture(frontend.TextureId) *ebiten.Image
	}

	newWindowFunc    func(*widget.WindowOption) window
	newAnimationFunc func(
		*frontend.Vector,
		*frontend.Pivot,
		frontend.Depth,
		*ebiten.Image,
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
	*frontend.Vector,
	frontend.Depth,
	*frontend.Pivot,
	character.CharacterId,
) *FaceWindow

func (f *FaceWindow) getCurrentAnimation() animation {
	return f.face[f.emotion.Current()]
}

func (f *FaceWindow) SetEmotion(emotion emotion.BattleEmotionType) {
	f.emotion.Enqueue(emotion)
}

func (f *FaceWindow) Update(parentPosition *frontend.Vector) {
	f.window.Update(parentPosition)
	anim := f.emotion.Apply(func(emotion emotion.BattleEmotionType) *widget.Animation {
		return f.face[emotion].(*widget.Animation)
	})
	anim.Update(f.window.GetPositionCenter())
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
	return standByNewFaceWindow(
		resource,
		func(option *widget.WindowOption) window {
			return newWindow(option)
		},
		func(
			relativePosition *frontend.Vector,
			pivot *frontend.Pivot,
			depth frontend.Depth,
			texture *ebiten.Image,
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
		relativePosition *frontend.Vector,
		depth frontend.Depth,
		pivot *frontend.Pivot,
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
