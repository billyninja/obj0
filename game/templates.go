package game

import (
	"github.com/billyninja/obj0/assets"
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
)

type MonsterTemplate struct {
	Txtr          *sdl.Texture
	ActionMap     *ActionMap
	SpriteSheet   *SpriteSheet
	Lvl           uint8
	HP            float32
	Size          int32
	LvlVariance   float32
	ScalingFactor float32
	AtkCoolDown   float32
	AtkSpeed      float32
	LoS           int32
	Loot          [8]Loot
}

var (
	SHADOW = &sdl.Rect{320, 224, TSzi, TSzi}

	MAN_PB_S1 *sdl.Rect = &sdl.Rect{96, 160, TSzi, TSzi}
	MAN_PB_S2 *sdl.Rect = &sdl.Rect{128, 160, TSzi, TSzi}
	MAN_PB_S3 *sdl.Rect = &sdl.Rect{160, 160, TSzi, TSzi}
	MAN_PU_S1 *sdl.Rect = &sdl.Rect{0, 192, TSzi, TSzi}
	MAN_PU_S2 *sdl.Rect = &sdl.Rect{128, 192, TSzi, TSzi}
	MAN_PU_S3 *sdl.Rect = &sdl.Rect{160, 160, TSzi, TSzi}

	MAN_CS_S1     *sdl.Rect    = &sdl.Rect{192, 224, TSzi, TSzi}
	MAN_CS_S2     *sdl.Rect    = &sdl.Rect{224, 224, TSzi, TSzi}
	MAN_CS_S3     *sdl.Rect    = &sdl.Rect{256, 224, TSzi, TSzi}
	MAN_CS_S4     *sdl.Rect    = &sdl.Rect{224, 192, TSzi, TSzi}
	MAN_PUSH_BACK [8]*sdl.Rect = [8]*sdl.Rect{MAN_PB_S1, MAN_PB_S2}
	MAN_CAST      [8]*sdl.Rect = [8]*sdl.Rect{MAN_CS_S1, MAN_CS_S3, MAN_CS_S4}
	MAN_AT_1      [8]*sdl.Rect = [8]*sdl.Rect{MAN_CS_S4}
	MAN_PICK_UP   [8]*sdl.Rect = [8]*sdl.Rect{MAN_PU_S1, MAN_PU_S2, MAN_PU_S3}

	PUFF_S1 *sdl.Rect    = &sdl.Rect{0, 0, 64, 64}
	PUFF_S2 *sdl.Rect    = &sdl.Rect{64, 0, 64, 64}
	PUFF_S3 *sdl.Rect    = &sdl.Rect{128, 0, 64, 64}
	PUFF_S4 *sdl.Rect    = &sdl.Rect{192, 0, 64, 64}
	PUFF_S5 *sdl.Rect    = &sdl.Rect{0, 64, 64, 64}
	PUFF_S6 *sdl.Rect    = &sdl.Rect{64, 64, 64, 64}
	PUFF_A  [8]*sdl.Rect = [8]*sdl.Rect{PUFF_S1, PUFF_S2, PUFF_S3, PUFF_S4, PUFF_S5, PUFF_S6}

	BLANK  *sdl.Rect    = &sdl.Rect{2000, 200, 192, 192}
	HIT_S1 *sdl.Rect    = &sdl.Rect{576, 0, 192, 192}
	HIT_S2 *sdl.Rect    = &sdl.Rect{758, 0, 192, 192}
	HIT_S3 *sdl.Rect    = &sdl.Rect{192, 192, 192, 192}
	HIT_A  [8]*sdl.Rect = [8]*sdl.Rect{HIT_S1, HIT_S2, HIT_S3}

	HIT_S4 *sdl.Rect    = &sdl.Rect{0, 0, 192, 192}
	HIT_S5 *sdl.Rect    = &sdl.Rect{192, 0, 192, 192}
	HIT_S6 *sdl.Rect    = &sdl.Rect{384, 0, 192, 192}
	HIT_B  [8]*sdl.Rect = [8]*sdl.Rect{BLANK, BLANK, HIT_S4, HIT_S5, HIT_S6, HIT_S5, HIT_S4}

	LAVA_S1 *sdl.Rect    = &sdl.Rect{192, 0, TSzi, TSzi}
	LAVA_S2 *sdl.Rect    = &sdl.Rect{224, 0, TSzi, TSzi}
	LAVA_S3 *sdl.Rect    = &sdl.Rect{256, 0, TSzi, TSzi}
	LAVA_A  [8]*sdl.Rect = [8]*sdl.Rect{LAVA_S1, LAVA_S2, LAVA_S3, LAVA_S3, LAVA_S2}

	YGLOW_S1 *sdl.Rect    = &sdl.Rect{0, 0, TSzi * 2, TSzi * 2}
	YGLOW_S2 *sdl.Rect    = &sdl.Rect{32, 0, TSzi * 2, TSzi * 2}
	YGLOW_S3 *sdl.Rect    = &sdl.Rect{64, 0, TSzi * 2, TSzi * 2}
	YGLOW_A  [8]*sdl.Rect = [8]*sdl.Rect{YGLOW_S1, YGLOW_S2, YGLOW_S3, YGLOW_S2}

	BGLOW_S1 *sdl.Rect = &sdl.Rect{224, 224, TSzi, TSzi}
	BGLOW_S2 *sdl.Rect = &sdl.Rect{256, 224, TSzi, TSzi}
	BGLOW_S3 *sdl.Rect = &sdl.Rect{288, 224, TSzi, TSzi}
	BGLOW_S4 *sdl.Rect = &sdl.Rect{320, 224, TSzi, TSzi}

	BGLOW_A [8]*sdl.Rect = [8]*sdl.Rect{BGLOW_S1, BGLOW_S2, BGLOW_S3, BGLOW_S4, BGLOW_S3, BGLOW_S2}

	LAVA_ANIM     = &Animation{Action: LAVA_A, PoseTick: 8}
	LIFE_ORB_ANIM = &Animation{Action: YGLOW_A, PoseTick: 8}
	PUFF_ANIM     = &Animation{Action: PUFF_A, PoseTick: 8}
	MAN_PB_ANIM   = &Animation{Action: MAN_PUSH_BACK, PoseTick: 18, PlayMode: 1}
	MAN_CS_ANIM   = &Animation{Action: MAN_CAST, PoseTick: 18, PlayMode: 1}
	MAN_PU_ANIM   = &Animation{Action: MAN_PICK_UP, PoseTick: 18, PlayMode: 1}
	MAN_ATK1_ANIM = &Animation{Action: MAN_AT_1, PoseTick: 4, PlayMode: 1}

	// ITEMS
	GreenBlob *Item = &Item{
		Name:        "Green Blob",
		Description: "A chunck of slime.",
		Source:      &sdl.Rect{0, 0, 24, 24},
		Weight:      2,
		BaseValue:   10,
		Txtr:        assets.Textures.GUI.PowerUps,
	}

	CrystalizedJelly *Item = &Item{
		Name:        "Crystalized Jelly",
		Description: "Some believe that the Slime's soul live within it",
		Source:      &sdl.Rect{24, 0, 24, 24},
		Weight:      2,
		BaseValue:   10,
		Txtr:        assets.Textures.GUI.PowerUps,
	}

	// VFX
	Puff *VFX = &VFX{
		Strip:        PUFF_A,
		DefaultSpeed: 4,
	}
	Hit *VFX = &VFX{
		Strip:        HIT_A,
		DefaultSpeed: 4,
	}
	Impact *VFX = &VFX{
		Strip:        HIT_B,
		DefaultSpeed: 3,
	}

	// GFX
	FiraProjGFX *GFX = &GFX{
		Animation: LIFE_ORB_ANIM,
	}
	FiraExplGFX *GFX = &GFX{
		Animation: PUFF_ANIM,
	}

	// Monsters

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
		SpriteSheet:   &SpriteSheet{assets.Textures.Sprites.Monsters, 0, 0, 48, 48},
		Loot:          [8]Loot{{CrystalizedJelly, 0.5}, {GreenBlob, 0.5}},
		Lvl:           1,
		HP:            25,
		LoS:           90,
		Size:          32,
		LvlVariance:   0.3,
		ScalingFactor: 0.6,
	}

	OrcTPL = MonsterTemplate{
		SpriteSheet:   &SpriteSheet{assets.Textures.Sprites.Monsters, 288, 0, 48, 48},
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

func (tpl *MonsterTemplate) MonsterFactory(lvlMod uint8, pos core.Vector2d, target *Solid) *SceneEntity {

	variance := rand.Float32() * tpl.LvlVariance * 100.0
	lvl := uint8((tpl.Lvl + lvlMod) + uint8(variance))
	hp := tpl.HP + float32(lvl*2)
	sizeMod := int32(float32(lvl-tpl.Lvl) * tpl.ScalingFactor)
	W, H := (tpl.Size + sizeMod), (tpl.Size + sizeMod)

	var DropItem *Item
	var sumP float32
	R := rand.Float32()
	for _, l := range tpl.Loot {
		sumP += l.Perc
		if R < sumP {
			DropItem = l.Item
			break
		}
	}

	sol := &Solid{
		Position:    &sdl.Rect{int32(pos.X), int32(pos.Y), W, H},
		Velocity:    &core.Vector2d{0, 0},
		Orientation: &core.Vector2d{0, 0},
		Txt:         tpl.Txtr,
		Speed:       1,
		Collision:   2,
		CPattern:    0,
		LoS:         tpl.LoS,
		MPattern: []Movement{
			Movement{F_DOWN, 50},
			Movement{F_UP, 90},
			Movement{F_RIGHT, 10},
			Movement{F_LEFT, 10},
		},
		Chase: target,
	}

	mon := &Char{
		Lvl:         lvl,
		ActionMap:   tpl.ActionMap,
		AtkSpeed:    tpl.AtkSpeed,
		AtkCoolDown: tpl.AtkCoolDown,
		CurrentHP:   hp,
		MaxHP:       hp,
		Drop:        DropItem,
	}

	sol.SetAnimation(tpl.ActionMap.DOWN)
	sol.CharPtr = mon

	return &SceneEntity{
		Solid: sol,
		Char:  mon,
	}
}

func BootstrapItems() {
	GreenBlob.Txtr = assets.Textures.GUI.PowerUps
	CrystalizedJelly.Txtr = assets.Textures.GUI.PowerUps
}

func BootstrapVFX() {
	GreenBlob.Txtr = assets.Textures.GUI.PowerUps
	CrystalizedJelly.Txtr = assets.Textures.GUI.PowerUps

	Puff.Txtr = assets.Textures.Sprites.Puff
	Hit.Txtr = assets.Textures.Sprites.Hit
	Impact.Txtr = assets.Textures.Sprites.Hit

	FiraProjGFX.Txtr = assets.Textures.Sprites.Glow
	FiraExplGFX.Txtr = assets.Textures.Sprites.Puff
}

func BootstrapMonsters() {
	BatTPL.Txtr = assets.Textures.Sprites.Monsters
	OrcTPL.Txtr = assets.Textures.Sprites.Monsters
	BatTPL.ActionMap = BatTPL.SpriteSheet.BuildBasicActions(3, false)
	OrcTPL.ActionMap = OrcTPL.SpriteSheet.BuildBasicActions(3, false)
}

func BootstrapPC(PC *SceneEntity) {
	MainCharSS := &SpriteSheet{
		Txtr:  assets.Textures.Sprites.MainChar,
		StX:   0,
		StY:   0,
		StepW: 32,
		StepH: 32,
	}
	MainCharActionMap := MainCharSS.BuildBasicActions(3, true)
	PC.Char.ActionMap = MainCharActionMap
	PC.Solid.CharPtr = PC.Char
}
