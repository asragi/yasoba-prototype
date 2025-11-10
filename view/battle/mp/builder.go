package mp

import "github.com/asragi/yasoba-prototype/common/character/hero"

type NewMPDisplayFunc func(
	initialMP hero.InitialMP,
	maxMP hero.MaxMP,
) *mpDisplay

func CreateNewMPDisplay(
	newImage ImageFactory,
	getImageData GetImageFunc,
) NewMPDisplayFunc {
	return func(initialMP hero.InitialMP, maxMP hero.MaxMP) *mpDisplay {
		newView := func(
			initialMp int,
			maxMp int,
		) viewInterface {
			return newView(
				newImage,
				getImageData,
				maxMp,
			)
		}
		presenter := createNewPresenter(
			newView,
		)
		return presenter(
			initialMP,
			maxMP,
		)
	}
}
