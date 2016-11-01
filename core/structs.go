package core

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Vector2d struct {
	X float32
	Y float32
}

type TextEl struct {
	Font         *ttf.Font
	Color        sdl.Color
	Content      string
	BakedContent string
	Txtr         *sdl.Texture
	TW           int32
	TH           int32
}
