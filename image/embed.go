package image

import _ "embed"

var (
	//go:embed cursor.png
	Cursor []byte
	//go:embed window.png
	Window []byte
	//go:embed face_lune_normal.png
	FaceLuneNormal []byte
	//go:embed face_sunny_normal.png
	FaceSunnyNormal []byte
)
