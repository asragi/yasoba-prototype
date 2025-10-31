package widget

import (
	"reflect"
	"testing"

	"github.com/asragi/yasoba-prototype/drawing"
	"github.com/asragi/yasoba-prototype/frontend"
	"github.com/asragi/yasoba-prototype/util"
)

func TestText_ToChar(t *testing.T) {
	tests := []struct {
		text     string
		expected []Char
	}{
		{
			text: "こんにちは",
			expected: []Char{
				Char('こ'),
				Char('ん'),
				Char('に'),
				Char('ち'),
				Char('は'),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := ToChar(tt.text)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToChar(%s) = %v, want %v", tt.text, result, tt.expected)
			}
		})
	}
}

func TestText_Size(t *testing.T) {
	// ResourceManagerを作成
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		t.Fatalf("ResourceManagerの作成に失敗しました: %v", err)
	}

	// CreateNewText関数を作成
	newText := CreateNewText(resource)

	tests := []struct {
		name    string
		text    string
		options *TextOptionsNew
	}{
		{
			name: "単一行のテキスト",
			text: "Hello",
			options: &TextOptionsNew{
				Font:  frontend.MaruMinya,
				Scale: 1,
			},
		},
		{
			name: "複数行のテキスト",
			text: "Hello\nWorld",
			options: &TextOptionsNew{
				Font:  frontend.MaruMinya,
				Scale: 1,
			},
		},
		{
			name: "スケール2のテキスト",
			text: "Test",
			options: &TextOptionsNew{
				Font:  frontend.MaruMinya,
				Scale: 2,
			},
		},
		{
			name: "日本語テキスト",
			text: "こんにちは",
			options: &TextOptionsNew{
				Font:  frontend.MaruMinya,
				Scale: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Textオブジェクトを作成
			text := newText(tt.options)

			// テキストを設定
			text.SetText(tt.text, true)

			// Sizeメソッドをテスト
			result := text.Size()

			// 基本的な検証（正の値であることを確認）
			if result.X <= 0 {
				t.Errorf("X座標が正の値ではありません: %f", result.X)
			}
			if result.Y <= 0 {
				t.Errorf("Y座標が正の値ではありません: %f", result.Y)
			}

			// スケールが適用されているかを確認
			expectedScale := float64(tt.options.Scale)
			if tt.options.Scale > 1 {
				// スケールが1より大きい場合、サイズもそれに応じて大きくなることを確認
				// 同じテキストをスケール1で作成して比較
				compareOptions := *tt.options
				compareOptions.Scale = 1
				compareText := newText(&compareOptions)
				compareText.SetText(tt.text, true)
				compareSize := compareText.Size()

				// スケール倍された値の近似値かチェック
				expectedX := compareSize.X * expectedScale
				expectedY := compareSize.Y * expectedScale

				if util.Abs(result.X-expectedX) > 1.0 {
					t.Errorf("スケール適用後のX座標が期待値と異なります: 期待値 %f, 実際の値 %f", expectedX, result.X)
				}
				if util.Abs(result.Y-expectedY) > 1.0 {
					t.Errorf("スケール適用後のY座標が期待値と異なります: 期待値 %f, 実際の値 %f", expectedY, result.Y)
				}
			}

			t.Logf("テストケース: %s, サイズ: X=%f, Y=%f", tt.name, result.X, result.Y)
		})
	}
}

func TestText_getLineSpacing(t *testing.T) {
	// ResourceManagerを作成
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		t.Fatalf("ResourceManagerの作成に失敗しました: %v", err)
	}

	newText := CreateNewText(resource)

	tests := []struct {
		name     string
		options  *TextOptionsNew
		expected float64
	}{
		{
			name: "MaruMinyaフォントの行間隔",
			options: &TextOptionsNew{
				Font: frontend.MaruMinya,
			},
			expected: 12.0, // MaruMinyaフォントのサイズ
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Textオブジェクトを作成
			text := newText(tt.options)

			// getLineSpacingメソッドをテスト
			result := text.(*Text).getCharacterHeight()

			// 結果を検証（フォントサイズと一致するはず）
			if result != tt.expected {
				t.Errorf("行間隔が一致しません: 期待値 %f, 実際の値 %f", tt.expected, result)
			}
		})
	}
}

func TestText_DrawBeforeUpdateSkipsRendering(t *testing.T) {
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		t.Fatalf("failed to create resource manager: %v", err)
	}

	newText := CreateNewText(resource)
	text := newText(&TextOptionsNew{
		RelativePosition: drawing.VectorZero,
		Pivot:            drawing.PivotTopLeft,
		Font:             frontend.MaruMinya,
		Speed:            1,
		Depth:            drawing.DepthWindow,
	})
	text.SetText("dummy", true)

	drawCalled := false

	text.Draw(func(drawing.DrawArgFunc, drawing.Depth) {
		drawCalled = true
	})

	if drawCalled {
		t.Fatalf("expected Draw to skip rendering when parentPosition is nil")
	}
}

func TestText_DrawAfterUpdateRenders(t *testing.T) {
	resource, err := frontend.CreateResourceManager()
	if err != nil {
		t.Fatalf("failed to create resource manager: %v", err)
	}

	newText := CreateNewText(resource)
	textInterface := newText(&TextOptionsNew{
		RelativePosition: drawing.VectorZero,
		Pivot:            drawing.PivotTopLeft,
		Font:             frontend.MaruMinya,
		Speed:            1,
		Depth:            drawing.DepthWindow,
	})
	textInterface.SetText("dummy", true)

	position := &drawing.Vector{X: 10, Y: 20}
	textInterface.Update(position)

	drawCalled := false
	textInterface.Draw(func(drawing.DrawArgFunc, drawing.Depth) {
		drawCalled = true
	})

	if !drawCalled {
		t.Fatalf("expected Draw to render after Update sets parentPosition")
	}
}
