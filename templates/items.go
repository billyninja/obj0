package templates

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	GreenBlob *core.Item = &core.Item{
		Name:        "Green Blob",
		Description: "A chunck of slime.",
		Source:      &sdl.Rect{0, 0, 24, 24},
	}

	CrystalizedJelly *core.Item = &core.Item{
		Name:        "Crystalized Jelly",
		Description: "Some believe that the Slime's soul live within it",
		Source:      &sdl.Rect{24, 0, 24, 24},
		Weight:      2,
		BaseValue:   10,
	}
)

func BootstrapItems(itemsTxt *sdl.Texture) {
	GreenBlob.Txtr = itemsTxt
	CrystalizedJelly.Txtr = itemsTxt
	return
}
