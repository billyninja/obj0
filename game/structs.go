package game

import (
	"github.com/billyninja/obj0/core"
	"github.com/billyninja/obj0/tmx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var (
	// FACING ORIENTATION
	F_LEFT  core.Vector2d = core.Vector2d{-1, 0}
	F_RIGHT               = core.Vector2d{1, 0}
	F_UP                  = core.Vector2d{0, -1}
	F_DOWN                = core.Vector2d{0, 1}
	F_DL                  = core.Vector2d{-1, 1}
	F_DR                  = core.Vector2d{1, 1}
	F_UL                  = core.Vector2d{-1, -1}
	F_UR                  = core.Vector2d{1, -1}
)

type SEvent func(source *SceneEntity, subject *SceneEntity, scene *Scene)

type SEventHandlers struct {
	OnCollDmg      float32
	OnCollPushBack int32
	OnCollEvent    SEvent
	OnPickUp       SEvent
	OnActDmg       uint16
	OnActPush      *core.Vector2d
	OnActEvent     SEvent

	DialogScript []string
	DoorTo       string
}

type Loot struct {
	Item *Item
	Perc float32
}

type SceneEntity struct {
	Char     *Char
	Solid    *Solid
	Handlers *SEventHandlers
	ItemPtr  *Item
}

type Solid struct {
	Velocity    *core.Vector2d
	Orientation *core.Vector2d
	Position    *sdl.Rect
	Source      *sdl.Rect
	Anim        *Animation
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
}

type Movement struct {
	Orientation core.Vector2d
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
	Flip     core.Vector2d
	Text     *TextEl
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

func (s *Solid) Destroy() {
	s.Position = nil
	s.Source = nil
	s.Anim = nil
	s.Txt = nil
	s.Collision = 0
}

type Camera struct {
	P   core.Vector2d
	DZx int32
	DZy int32
}

type Scene struct {
	codename   string
	StartPoint core.Vector2d
	TileSet    *sdl.Texture
	Cam        *Camera

	World           [][]*tmx.Terrain
	Visual          []*VFXInst
	Spawners        []*SpawnPoint
	Interactive     []*SceneEntity
	Monsters        []*SceneEntity
	CullMap         []*SceneEntity
	GUI             []*sdl.Rect
	DBox            *DBox
	WinWidth        int32
	WinHeight       int32
	CellsX          int32
	CellsY          int32
	TileWidth       int32
	TileHeight      int32
	LimitW          int32
	LimitH          int32
	EventTick       uint8
	AiTick          uint8
	EventTickLength uint8
	AiTickLength    uint8
}

type SpawnPoint struct {
	Position  *sdl.Rect
	Frequency uint16
	LvlMod    uint8
}
