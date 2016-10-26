package templates

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
)

type Loot struct {
	Item *core.Item
	Perc float32
}

type MonsterTemplate struct {
	Txtr          *sdl.Texture
	ActionMap     *core.ActionMap
	SpriteSheet   *core.SpriteSheet
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
