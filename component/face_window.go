package component

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/widget"
)

type FaceWindow struct {
	emotion BattleEmotionType
	face    map[BattleEmotionType]*widget.Animation
	window  widget.WindowInterface
}

type NewFaceWindowFunc func(
	*frontend.Vector,
	frontend.Depth,
	*frontend.Pivot,
	core.CharacterId,
) *FaceWindow

func (f *FaceWindow) getCurrentAnimation() *widget.Animation {
	return f.face[f.emotion]
}

func (f *FaceWindow) SetEmotion(emotion BattleEmotionType) {
	f.emotion = emotion
}

func (f *FaceWindow) Update(parentPosition *frontend.Vector) {
	f.window.Update(parentPosition)
	f.getCurrentAnimation().Update(f.window.GetPositionCenter())
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
		characterId core.CharacterId,
	) *FaceWindow {
		const padding = 6
		const faceSize = 74
		animationMap := func() map[BattleEmotionType]*widget.Animation {
			result := map[BattleEmotionType]*widget.Animation{}
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
			face:   animationMap,
			window: window,
		}
	}
}

type BattleEmotionType int

const (
	BattleEmotionNormal BattleEmotionType = iota
	BattleEmotionDamage
	BattleEmotionSmile
	BattleEmotionAngry
	BattleEmotionAnnoyed
)

type getAllEmotionFunc func() map[core.CharacterId]map[BattleEmotionType]frontend.AnimationId

func createGetAllEmotionFunc() getAllEmotionFunc {
	dict := map[core.CharacterId]map[BattleEmotionType]frontend.AnimationId{
		core.CharacterLuneId: {
			BattleEmotionNormal: frontend.AnimationIdLuneNormal,
			BattleEmotionDamage: frontend.AnimationIdLuneDamage,
		},
		core.CharacterSunnyId: {
			BattleEmotionNormal:  frontend.AnimationIdSunnyNormal,
			BattleEmotionDamage:  frontend.AnimationIdSunnyDamage,
			BattleEmotionSmile:   frontend.AnimationIdSunnySmile,
			BattleEmotionAngry:   frontend.AnimationIdSunnyAngry,
			BattleEmotionAnnoyed: frontend.AnimationIdSunnyAnnoyed,
		},
	}
	return func() map[core.CharacterId]map[BattleEmotionType]frontend.AnimationId {
		return dict
	}
}
