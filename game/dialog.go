package game

import (
	"github.com/billyninja/obj0/assets"
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
)

type DBox struct {
	SPos     uint8
	CurrText uint8
	Text     []*core.TextEl
	BGColor  sdl.Color
	Char     *Char
}

func (db *DBox) LoadText(content []string) {
	db.Text = make([]*core.TextEl, len(content))
	for i, s := range content {
		db.Text[i] = &core.TextEl{
			Font:    assets.Fonts.Default,
			Content: s,
			Color:   CL_WHITE,
		}
	}
}

func (db *DBox) Present(winWidth, winHeight int32, renderer *sdl.Renderer) {
	if len(db.Text) == 0 {
		return
	}

	ct := db.Text[db.CurrText]
	txtr, w, h := ct.Bake(renderer, int(winWidth))
	br := &sdl.Rect{64, winWidth - 128, 512, 120}
	tr := &sdl.Rect{0, 0, w, h}
	bt := &sdl.Rect{64, 720 - 128, w, h}

	renderer.Copy(assets.Textures.GUI.Transparency, &sdl.Rect{0, 0, 48, 48}, br)
	renderer.Copy(txtr, tr, bt)
}

func (db *DBox) NextText() bool {
	if len(db.Text) == 0 {
		return false
	}

	db.CurrText += 1
	if int(db.CurrText+1) > len(db.Text) {
		db.Text = []*core.TextEl{}
		db.CurrText = 0
	}

	return true
}
