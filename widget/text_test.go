package widget

import (
	"reflect"
	"testing"

	"github.com/asragi/yasoba-prototype/frontend"
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
		name     string
		text     string
		options  *TextOptionsNew
		expected *frontend.Vector
	}{
		{
			name: "単一行のテキスト",
			text: "Hello",
			options: &TextOptionsNew{
				Font:  frontend.MaruMinya,
				Scale: 1,
			},
			expected: &frontend.Vector{
				X: 30.0, // text.Measureで計算された幅
				Y: 12.0, // text.Measureで計算された高さ
			},
		},
		{
			name: "複数行のテキスト",
			text: "Hello\nWorld",
			options: &TextOptionsNew{
				Font:  frontend.MaruMinya,
				Scale: 1,
			},
			expected: &frontend.Vector{
				X: 30.0, // text.Measureで計算された幅（最初の行の長さ）
				Y: 28.0, // text.Measureで計算された高さ（2行分）
			},
		},
		{
			name: "スケール2のテキスト",
			text: "Test",
			options: &TextOptionsNew{
				Font:  frontend.MaruMinya,
				Scale: 2,
			},
			expected: &frontend.Vector{
				X: 48.0, // text.Measureで計算された幅 * 2
				Y: 24.0, // text.Measureで計算された高さ * 2
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

			eps := 0.1

			// 結果を検証（実際の値に近い範囲でチェック）
			if result.X < tt.expected.X-eps || result.X > tt.expected.X+eps {
				t.Errorf("X座標が期待範囲外です: 期待値 %f, 実際の値 %f", tt.expected.X, result.X)
				t.Logf("テストケース: %s, 実際の値: X=%f, Y=%f", tt.name, result.X, result.Y)
			}
			if result.Y < tt.expected.Y-eps || result.Y > tt.expected.Y+eps {
				t.Errorf("Y座標が期待範囲外です: 期待値 %f, 実際の値 %f", tt.expected.Y, result.Y)
				t.Logf("テストケース: %s, 実際の値: X=%f, Y=%f", tt.name, result.X, result.Y)
			}
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
			result := text.(*Text).getLineSpacing()

			// 結果を検証（フォントサイズと一致するはず）
			if result != tt.expected {
				t.Errorf("行間隔が一致しません: 期待値 %f, 実際の値 %f", tt.expected, result)
			}
		})
	}
}
