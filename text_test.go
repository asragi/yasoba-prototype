package main

import (
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func loadFont(size float64) font.Face {
	// フォントファイルを開く
	data, err := os.ReadFile("font/x12y12pxMaruMinya.ttf") // 日本語対応フォントを用意
	if err != nil {
		log.Fatal(err)
	}

	// パース
	ttf, err := opentype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	// 指定サイズでfont.Faceを生成
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return face
}
