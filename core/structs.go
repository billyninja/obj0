package core

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type Vector2d struct {
	X float32
	Y float32
}

type Event func(source *Solid, subject *Solid)

type TextEl struct {
	Font         *ttf.Font
	Color        sdl.Color
	Content      string
	BakedContent string
	Txtr         *sdl.Texture
	TW           int32
	TH           int32
}

type InteractionHandlers struct {
	OnCollDmg      float32
	OnCollPushBack int32
	OnCollEvent    Event
	OnPickUp       Event
	OnActDmg       uint16
	OnActPush      *Vector2d
	OnActEvent     Event
	DialogScript   []string
	DoorTo         string
}

type Solid struct {
	Velocity    *Vector2d
	Orientation *Vector2d
	Position    *sdl.Rect
	Source      *sdl.Rect
	Anim        *Animation
	Handlers    *InteractionHandlers
	Txt         *sdl.Texture
	Ttl         int64
	Collision   uint8
	Speed       float32

	// AI RELATED
	CPattern uint32
	MPattern []Movement
	Chase    *Solid
	CharPtr  *Char
	ItemPtr  *Item
	LoS      int32
}

type Char struct {
	Solid     *Solid
	ActionMap *ActionMap
	Inventory []*ItemStack
	//---
	SpeedMod     float32
	BaseSpeed    float32
	Lvl          uint8
	CurrentXP    uint16
	NextLvlXP    uint16
	CurrentHP    float32
	MaxHP        float32
	CurrentST    float32
	MaxST        float32
	Drop         *Item
	Invinc       int64
	AtkSpeed     float32
	AtkCoolDownC float32
	AtkCoolDown  float32
}

type Item struct {
	Name        string
	Description string
	Weight      float32
	BaseValue   uint16
	Txtr        *sdl.Texture
	Source      *sdl.Rect
}

type ItemInstance struct {
	ItemTpl *Item
	Solid   *Solid
}

type ItemStack struct {
	ItemTpl *Item
	Qty     int
}

type Animation struct {
	Action   [8]*sdl.Rect
	Pose     uint8
	PoseTick uint32
	PlayMode uint8
	After    Event
}

type Movement struct {
	Orientation Vector2d
	Ticks       uint8
}

type ActionMap struct {
	UP    *Animation
	DOWN  *Animation
	LEFT  *Animation
	RIGHT *Animation
	UL    *Animation
	UR    *Animation
	DL    *Animation
	DR    *Animation
	//---
	PB *Animation
	//---
	SA_1 *Animation
	SA_2 *Animation
	SA_3 *Animation
	SA_4 *Animation
	SA_5 *Animation
	SA_6 *Animation
	//---
}

type SpriteSheet struct {
	Txtr  *sdl.Texture
	StX   int32
	StY   int32
	StepW int32
	StepH int32
}

func (s *Solid) Destroy() {
	s.Position = nil
	s.Source = nil
	s.Anim = nil
	s.Handlers = nil
	s.Txt = nil
	s.Collision = 0
}

type VFX struct {
	Txtr         *sdl.Texture
	Strip        [8]*sdl.Rect
	DefaultSpeed uint8
}

type VFXInst struct {
	Vfx      *VFX
	Pos      *sdl.Rect
	Pose     uint8
	Tick     uint8
	Ttl      int64
	CurrTick uint8
	Loop     uint8
	Flip     Vector2d
	Text     *TextEl
}
