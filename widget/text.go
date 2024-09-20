package widget

import (
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
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

func (t *Text) Size() *frontend.Vector {
	scale := float64(t.options.Scale)
	// TODO: Size should be calculated from font Size
	tmp := &frontend.Vector{
		X: float64(len(t.characterSet[0])*t.options.XSpacing) - 1,
		Y: float64(len(t.characterSet))*16 - 4,
	}
	return tmp.Multiply(scale)
}

func (t *Text) drawText(
	characters []string,
	currentIndex int,
	parentPosition *frontend.Vector,
	line int,
	drawFunc frontend.DrawFunc,
) {
	// TODO: characterSizeX should be calculated from font Size
	const lineHeight = 16
	scale := float64(t.options.Scale)
	diffSet := []*frontend.Vector{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
	pivotDiff := t.options.Pivot.ApplyToSize(t.Size())
	characterPosition := func() []*frontend.Vector {
		result := make([]*frontend.Vector, t.textSize)
		for i := 0; i < t.textSize; i++ {
			tmp := frontend.Vector{
				X: t.options.RelativePosition.X + float64(i*t.options.XSpacing)*scale,
				Y: t.options.RelativePosition.Y,
			}
			result[i] = tmp.Sub(pivotDiff)
		}
		return result
	}()
	for i := 0; i < currentIndex; i++ {
		op := &text.DrawOptions{}
		x := characterPosition[i].X + parentPosition.X
		y := characterPosition[i].Y + parentPosition.Y + float64(line*lineHeight)*scale
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x, y)
		targetCharacter := characters[i]
		drawFunc(
			func(screen *ebiten.Image) {
				if t.options.EnableOutline {
					outlineOp := &text.DrawOptions{}
					*outlineOp = *op
					outlineOp.ColorScale.ScaleWithColor(t.options.OutlineColor)
					for j := 0; j < len(diffSet); j++ {
						v := diffSet[j].Multiply(scale)
						outlineOp.GeoM.Translate(v.X, v.Y)
						text.Draw(screen, targetCharacter, t.options.TextFace, outlineOp)
						outlineOp.GeoM.Translate(-v.X, -v.Y)
					}
				}
				op.ColorScale.ScaleWithColor(t.options.Color)
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
	Color            color.Color
	OutlineColor     color.Color
	EnableOutline    bool
	Scale            int
	XSpacing         int
}

func NewText(options *TextOptions) *Text {
	const characterSizeX = 13
	if options.Color == nil {
		options.Color = color.White
	}
	if options.OutlineColor == nil {
		options.OutlineColor = color.Black
	}
	if options.Scale == 0 {
		options.Scale = 1
	}
	if options.XSpacing == 0 {
		options.XSpacing = characterSizeX
	}
	return &Text{
		currentIndex: 0,
		characterSet: nil,
		sizes:        nil,
		textSize:     0,
		frameCounter: 0,
		options:      options,
	}
}
