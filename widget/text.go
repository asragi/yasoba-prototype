package widget

import (
	"image/color"
	"strings"

	"github.com/asragi/yasoba-prototype/adapter/ebiten/frontend"
	"github.com/asragi/yasoba-prototype/toolkit/drawing"
	"github.com/asragi/yasoba-prototype/util"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	lineMargin = 4
	marginX    = 1
)

type TextInterface interface {
	PositionUpdater
	Drawer
	ForceComplete()
	SetText(text string, displayAll bool)
	Size() *drawing.Vector
	CheckIsEnd() bool
	SetTextColor(color color.Color)
}

type Char string

func (c Char) String() string {
	return string(c)
}

func (c Char) Size(face *text.GoTextFace) *drawing.Vector {
	width, height := text.Measure(string(c), face, 1)
	return &drawing.Vector{
		X: width,
		Y: height,
	}
}

func ToChar(s string) []Char {
	runes := []rune(s)
	result := make([]Char, len(runes))
	for i, c := range runes {
		result[i] = Char(c)
	}
	return result
}

// charactersSetType は文字セットを表す2次元文字列スライス
// 1次元目: 改行で分割された各行
// 2次元目: 各行の文字を1文字ずつに分割した配列
//
// 例: "Hello\nWorld" の場合
// [
//
//	["H", "e", "l", "l", "o"],  // 1行目
//	["W", "o", "r", "l", "d"]   // 2行目
//
// ]
type charactersSetType [][]Char

// 文字を一つずつ表示するなどの機能を持たせたテキストレンダリング
type Text struct {
	currentIndex   int
	characterSet   charactersSetType
	fullText       string
	sizes          []int
	textSize       int
	frameCounter   int
	options        *TextOptionsNew
	textFace       *text.GoTextFace
	parentPosition *drawing.Vector
	drawTextFunc   drawing.DrawTextFunc
}

type FontProvider interface {
	GetFont(frontend.FontId) *text.GoTextFace
}

func (t *Text) ForceComplete() {
	t.currentIndex = t.textSize
}

func (t *Text) CheckIsEnd() bool {
	return t.currentIndex >= t.textSize
}

func (t *Text) SetText(textString string, displayAll bool) {
	t.fullText = textString
	charactersSet, sizes, textSize := func(textString string) ([][]Char, []int, int) {
		texts := strings.Split(textString, "\n")
		textSize := 0
		characters := make(charactersSetType, len(texts))
		sizes := make([]int, len(texts))
		for i, t := range texts {
			characters[i] = ToChar(t)
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

func (t *Text) Update(parentPosition *drawing.Vector) {
	t.frameCounter++
	t.currentIndex = util.ClampInt(t.frameCounter/t.options.Speed, t.currentIndex, t.textSize)
	t.parentPosition = parentPosition
}

func (t *Text) Draw(drawFunc drawing.DrawFunc) {
	if t.parentPosition == nil {
		return
	}
	for i := 0; i < len(t.characterSet); i++ {
		tmpCurrentIndex := t.currentIndex
		for j := 0; j < i; j++ {
			tmpCurrentIndex -= t.sizes[j]
		}
		tmpCurrentIndex = util.ClampInt(tmpCurrentIndex, 0, t.sizes[i])
		t.drawText(t.characterSet[i], tmpCurrentIndex, t.parentPosition, i, drawFunc)
	}
}

func (t *Text) Size() *drawing.Vector {
	scale := float64(t.options.Scale)
	characterHeight := t.getCharacterHeight()
	maxWidth := 0.0
	for _, line := range t.characterSet {
		width := 0.0
		for _, character := range line {
			characterWidth, _ := text.Measure(character.String(), t.textFace, 1)
			width += characterWidth + marginX
		}
		if width > maxWidth {
			maxWidth = width
		}
	}
	// 最後の文字のmarginXを引く
	maxWidth -= marginX
	height := characterHeight*float64(len(t.characterSet)) + lineMargin*float64(len(t.characterSet)-1)
	return &drawing.Vector{
		X: maxWidth * scale,
		Y: height * scale,
	}
}

func (t *Text) SetTextColor(color color.Color) {
	t.options.Color = color
}

func (t *Text) getCharacterHeight() float64 {
	_, height := text.Measure("あ", t.textFace, 1)
	return height
}

// 与えられた1行の文字列を描画する
func (t *Text) drawText(
	characters []Char,
	currentIndex int,
	parentPosition *drawing.Vector,
	line int,
	drawFunc drawing.DrawFunc,
) {
	length := len(characters)
	// TODO: characterSizeX should be calculated from font Size
	const lineHeight = 16
	scale := float64(t.options.Scale)
	diffSet := []*drawing.Vector{
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 0},
		{X: -1, Y: 0},
	}
	pivotDiff := t.options.Pivot.ApplyToSize(t.Size())
	characterPosition := func() []*drawing.Vector {
		result := make([]*drawing.Vector, length)
		xPosition := 0.0
		for i := 0; i < length; i++ {
			targetCharacter := characters[i]
			tmp := drawing.Vector{
				X: t.options.RelativePosition.X + xPosition*scale,
				Y: t.options.RelativePosition.Y,
			}
			xPosition += targetCharacter.Size(t.textFace).X + marginX
			result[i] = tmp.Sub(pivotDiff)
		}
		return result
	}()
	for i := 0; i < currentIndex; i++ {
		op := drawing.NewTextDrawOptions()
		x := characterPosition[i].X + parentPosition.X
		y := characterPosition[i].Y + parentPosition.Y + float64(line*lineHeight)*scale
		op.SetScale(scale, scale)
		op.Translate(x, y)
		targetCharacter := characters[i]
		drawFunc(
			func(screen drawing.Image) {
				if t.options.EnableOutline {
					outlineOp := drawing.NewTextDrawOptions()
					outlineOp.SetScale(scale, scale)
					outlineOp.Translate(x, y)
					outlineOp.SetColorScale(t.options.OutlineColor)
					for j := 0; j < len(diffSet); j++ {
						v := diffSet[j].Multiply(scale)
						outlineOp.Translate(v.X, v.Y)
						t.drawTextFunc(screen, targetCharacter.String(), t.textFace, outlineOp)
						outlineOp.Translate(-v.X, -v.Y)
					}
				}
				op.SetColorScale(t.options.Color)
				t.drawTextFunc(screen, targetCharacter.String(), t.textFace, op)
			}, t.options.Depth,
		)
	}
}

type TextOptionsNew struct {
	RelativePosition *drawing.Vector
	Pivot            *drawing.Pivot
	Font             frontend.FontId
	Speed            int
	Depth            drawing.Depth
	Color            color.Color
	OutlineColor     color.Color
	EnableOutline    bool
	Scale            int
}

type NewTextFunc func(*TextOptionsNew) TextInterface

func CreateNewText(
	resource FontProvider,
	drawText drawing.DrawTextFunc,
) NewTextFunc {
	if drawText == nil {
		panic("widget: drawText func is required")
	}
	return func(options *TextOptionsNew) TextInterface {
		if options.Color == nil {
			options.Color = color.White
		}
		if options.OutlineColor == nil {
			options.OutlineColor = color.Black
		}
		if options.Scale == 0 {
			options.Scale = 1
		}
		if options.Speed == 0 {
			options.Speed = 1
		}
		return &Text{
			currentIndex:   0,
			characterSet:   nil,
			sizes:          nil,
			textSize:       0,
			frameCounter:   0,
			options:        options,
			parentPosition: nil,
			textFace:       resource.GetFont(options.Font),
			drawTextFunc:   drawText,
		}
	}
}
