package core

type TextId string

type ServeTextData func(id TextId) *TextData

type TextData struct {
	Id   TextId
	Text string
}
