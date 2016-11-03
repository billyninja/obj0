package core

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (t *TextEl) Bake(renderer *sdl.Renderer, limit int) (*sdl.Texture, int32, int32) {

	if t.Font == nil {
		return nil, 0, 0
	}

	if t.Content == t.BakedContent {
		return t.Txtr, t.TW, t.TH
	}

	surface, _ := t.Font.RenderUTF8_Blended_Wrapped(t.Content, t.Color, limit)
	defer surface.Free()
	t.Txtr, _ = renderer.CreateTextureFromSurface(surface)
	t.TW, t.TH = surface.W, surface.H

	return t.Txtr, t.TW, t.TH
}
