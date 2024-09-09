package widget

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"strings"
)

type Text struct {
	ForceComplete func()
	Update        func(parentPosition *frontend.Vector)
	Draw          func(frontend.DrawFunc)
}

type TextOptions struct {
	RelativePosition *frontend.Vector
	Pivot            *frontend.Pivot
	TextFace         *text.GoTextFace
	DisplayAll       bool
	Speed            int
	Depth            frontend.Depth
}

func NewText(textString string, options *TextOptions) *Text {
	// TODO: characterSizeX should be calculated from font size
	characterSizeX := 14
	lineHeight := 16
	charactersSet, sizes, textSize := func() ([][]string, []int, int) {
		texts := strings.Split(textString, "\n")
		textSize := 0
		characters := make([][]string, len(texts))
		sizes := make([]int, len(texts))
		for i, t := range texts {
			characters[i] = util.SplitString(t)
			sizes[i] = len(characters[i])
			textSize += sizes[i]
		}
		return characters, sizes, textSize
	}()
	drawText := func(
		characters []string,
		currentIndex int,
		parentPosition *frontend.Vector,
		line int,
		drawFunc frontend.DrawFunc,
	) {
		characterPosition := func() []*frontend.Vector {
			result := make([]*frontend.Vector, textSize)
			for i := 0; i < textSize; i++ {
				result[i] = &frontend.Vector{
					X: options.RelativePosition.X + float64(i*characterSizeX),
					Y: options.RelativePosition.Y,
				}
			}
			return result
		}()
		for i := 0; i < currentIndex; i++ {
			op := &text.DrawOptions{}
			x := characterPosition[i].X + parentPosition.X
			y := characterPosition[i].Y + parentPosition.Y + float64(line*lineHeight)
			op.GeoM.Translate(x, y)
			targetCharacter := characters[i]
			drawFunc(
				func(screen *ebiten.Image) {
					text.Draw(screen, targetCharacter, options.TextFace, op)
				}, options.Depth,
			)
		}
	}
	currentIndex := 0
	frameCounter := 0
	parentPosition := &frontend.Vector{}
	if options.DisplayAll {
		currentIndex = textSize
	}

	forceComplete := func() {
		currentIndex = textSize
	}

	update := func(passedParentPosition *frontend.Vector) {
		frameCounter++
		currentIndex = util.ClampInt(frameCounter/options.Speed, currentIndex, textSize)
		parentPosition = passedParentPosition
	}

	draw := func(drawFunc frontend.DrawFunc) {
		for i := 0; i < len(charactersSet); i++ {
			tmpCurrentIndex := currentIndex
			for j := 0; j < i; j++ {
				tmpCurrentIndex -= sizes[j]
			}
			tmpCurrentIndex = util.ClampInt(tmpCurrentIndex, 0, sizes[i])
			drawText(charactersSet[i], tmpCurrentIndex, parentPosition, i, drawFunc)
		}
	}

	return &Text{
		ForceComplete: forceComplete,
		Update:        update,
		Draw:          draw,
	}
}
