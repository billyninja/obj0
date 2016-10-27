package templates

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	Fira *core.SpellTemplate = &core.SpellTemplate{
		Name:        "Fira",
		Description: "Basic fire spell",
		BaseSPCost:  12.0,
		Elemental:   0,
		BaseDmg:     20,

		ProjectileTtl:   500,
		ProjectileSpeed: 4,
		ProjectileSize:  &core.Vector2d{24, 24},

		PreCastVfx:   Hit,
		PreCastAnim:  core.MAN_CS_ANIM,
		CastVfx:      Hit,
		PostCastVfx:  Impact,
		PostCastAnim: core.MAN_PU_ANIM,
	}
)

func BootstrapSpells(spellTxtr *sdl.Texture, life_orb *core.Animation) {
	Fira.ProjectileTxtr = spellTxtr
	Fira.ProjectileAnim = life_orb
}
