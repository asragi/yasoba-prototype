package widget

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"strings"
)

type Text struct {
	currentIndex   int
	characterSet   [][]string
	sizes          []int
	textSize       int
	frameCounter   int
	options        *TextOptions
	parentPosition *frontend.Vector
}

func (t *Text) ForceComplete() {
	t.currentIndex = t.textSize
}

func (t *Text) SetText(textString string, displayAll bool) {
	charactersSet, sizes, textSize := func(textString string) ([][]string, []int, int) {
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
	}(textString)
	t.frameCounter = 0
	t.currentIndex = func() int {
		if displayAll {
			return textSize
		}
		return 0
	}()
	t.characterSet = charactersSet
	t.sizes = sizes
	t.textSize = textSize
}

func (t *Text) Update(parentPosition *frontend.Vector) {
	t.frameCounter++
	t.currentIndex = util.ClampInt(t.frameCounter/t.options.Speed, t.currentIndex, t.textSize)
	t.parentPosition = parentPosition
}

func (t *Text) Draw(drawFunc frontend.DrawFunc) {
	for i := 0; i < len(t.characterSet); i++ {
		tmpCurrentIndex := t.currentIndex
		for j := 0; j < i; j++ {
			tmpCurrentIndex -= t.sizes[j]
		}
		tmpCurrentIndex = util.ClampInt(tmpCurrentIndex, 0, t.sizes[i])
		t.drawText(t.characterSet[i], tmpCurrentIndex, t.parentPosition, i, drawFunc)
	}
}

func (t *Text) drawText(
	characters []string,
	currentIndex int,
	parentPosition *frontend.Vector,
	line int,
	drawFunc frontend.DrawFunc,
) {
	// TODO: characterSizeX should be calculated from font size
	const characterSizeX = 14
	const lineHeight = 16
	characterPosition := func() []*frontend.Vector {
		result := make([]*frontend.Vector, t.textSize)
		for i := 0; i < t.textSize; i++ {
			result[i] = &frontend.Vector{
				X: t.options.RelativePosition.X + float64(i*characterSizeX),
				Y: t.options.RelativePosition.Y,
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
				text.Draw(screen, targetCharacter, t.options.TextFace, op)
			}, t.options.Depth,
		)
	}
}

type TextOptions struct {
	RelativePosition *frontend.Vector
	Pivot            *frontend.Pivot
	TextFace         *text.GoTextFace
	Speed            int
	Depth            frontend.Depth
}

func NewText(options *TextOptions) *Text {
	return &Text{
		currentIndex: 0,
		characterSet: nil,
		sizes:        nil,
		textSize:     0,
		frameCounter: 0,
		options:      options,
	}
}
