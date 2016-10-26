package templates

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	SlimeTPL MonsterTemplate = MonsterTemplate{
		Lvl:           1,
		HP:            25,
		LoS:           90,
		Size:          32,
		LvlVariance:   0.3,
		ScalingFactor: 0.6,
		Loot: [8]Loot{
			{CrystalizedJelly, 0.5},
			{GreenBlob, 0.5},
		},
	}

	BatTPL = MonsterTemplate{
		Loot:          [8]Loot{{CrystalizedJelly, 0.5}, {GreenBlob, 0.5}},
		Lvl:           1,
		HP:            25,
		LoS:           90,
		Size:          32,
		LvlVariance:   0.3,
		ScalingFactor: 0.6,
	}

	OrcTPL = MonsterTemplate{
		Loot:          [8]Loot{{CrystalizedJelly, 0.5}, {GreenBlob, 0.5}},
		Lvl:           5,
		HP:            70,
		LoS:           100,
		Size:          64,
		LvlVariance:   0.5,
		ScalingFactor: 0.1,
		AtkCoolDown:   60.0,
		AtkSpeed:      1,
	}
)

func BootstrapMonsters(monstersTxt *sdl.Texture) {

	BatSS := &core.SpriteSheet{monstersTxt, 0, 0, 48, 48}
	BatActionMap := BatSS.BuildBasicActions(3, false)
	BatTPL.Txtr = monstersTxt
	BatTPL.ActionMap = BatActionMap

	OrcSS := &core.SpriteSheet{monstersTxt, 288, 0, 48, 48}
	OrcActionMap := OrcSS.BuildBasicActions(3, false)
	OrcTPL.Txtr = monstersTxt
	OrcTPL.ActionMap = OrcActionMap

	return
}
