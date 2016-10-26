package templates

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	Puff core.VFX = core.VFX{
		Strip:        core.PUFF_A,
		DefaultSpeed: 4,
	}
	Hit core.VFX = core.VFX{
		Strip:        core.HIT_A,
		DefaultSpeed: 4,
	}
	Impact core.VFX = core.VFX{
		Strip:        core.HIT_B,
		DefaultSpeed: 3,
	}
)

func BootstrapVfx(hitTxt, puffTxt *sdl.Texture) {
	Hit.Txtr = hitTxt
	Impact.Txtr = hitTxt
	Puff.Txtr = puffTxt
}
