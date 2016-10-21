package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"math"
	"math/rand"
	"runtime"
	"time"
)

const (
	winTitle            string  = "Go-SDL2 Obj0"
	TSz                 float32 = 32
	CharSize            int32   = 64
	TSzi                int32   = int32(TSz)
	winWidth, winHeight int32   = 1280, 720
	cY                          = winHeight / TSzi
	cX                          = winWidth / TSzi
	CenterX                     = winWidth / 2
	CenterY                     = winHeight / 2
	WORLD_CELLS_X               = 500
	WORLD_CELLS_Y               = 200

	KEY_ARROW_UP    = 1073741906
	KEY_ARROW_DOWN  = 1073741905
	KEY_ARROW_LEFT  = 1073741904
	KEY_ARROW_RIGHT = 1073741903
	KEY_LEFT_SHIFT  = 1073742049
	KEY_SPACE_BAR   = 32 // 1073741824
	KEY_C           = 99
	KEY_X           = 120
	KEY_Z           = 80 // todo

	AI_TICK_LENGTH    = 2
	EVENT_TICK_LENGTH = 2
)

var (
	event        sdl.Event
	font         *ttf.Font
	game_latency time.Duration
	Controls     *ControlState = &ControlState{}

	CL_WHITE        sdl.Color    = sdl.Color{255, 255, 255, 255}
	tilesetTxt      *sdl.Texture = nil
	spritesheetTxt  *sdl.Texture = nil
	particlesTxt    *sdl.Texture = nil
	powerupsTxt     *sdl.Texture = nil
	glowTxt         *sdl.Texture = nil
	monstersTxt     *sdl.Texture = nil
	transparencyTxt *sdl.Texture = nil
	puffTxt         *sdl.Texture = nil
	hitTxt          *sdl.Texture = nil

	GRASS     *sdl.Rect = &sdl.Rect{0, 0, TSzi, TSzi}
	DIRT                = &sdl.Rect{703, 0, TSzi, TSzi}
	WALL                = &sdl.Rect{0, 64, TSzi, TSzi}
	DOOR                = &sdl.Rect{256, 32, TSzi, TSzi}
	BF_ATK_UP           = &sdl.Rect{72, 24, 24, 24}
	BF_DEF_UP           = &sdl.Rect{96, 24, 24, 24}
	SHADOW              = &sdl.Rect{320, 224, TSzi, TSzi}

	// FACING ORIENTATION
	F_LEFT  Vector2d = Vector2d{-1, 0}
	F_RIGHT          = Vector2d{1, 0}
	F_UP             = Vector2d{0, -1}
	F_DOWN           = Vector2d{0, 1}
	F_DL             = Vector2d{-1, 1}
	F_DR             = Vector2d{1, 1}
	F_UL             = Vector2d{-1, -1}
	F_UR             = Vector2d{1, -1}

	// MAIN CHAR POSES AND ANIMATIONS
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
	MAN_PB_ANIM   = &Animation{Action: MAN_PUSH_BACK, PoseTick: 18, PlayMode: 1}
	MAN_CS_ANIM   = &Animation{Action: MAN_CAST, PoseTick: 18, PlayMode: 1}
	MAN_PU_ANIM   = &Animation{Action: MAN_PICK_UP, PoseTick: 18, PlayMode: 1}
	MAN_ATK1_ANIM = &Animation{Action: MAN_AT_1, PoseTick: 4, PlayMode: 1}

	BatTPL MonsterTemplate = MonsterTemplate{}
	OrcTPL MonsterTemplate = MonsterTemplate{}
	puff   VFX             = VFX{}
	hit    VFX             = VFX{}
	impact VFX             = VFX{}

	SCN_PLAINS *Scene = &Scene{
		codename:   "plains",
		CellsX:     50,
		CellsY:     30,
		CamPoint:   Vector2d{120, 120},
		StartPoint: Vector2d{300, 200},
		tileA:      GRASS,
		tileB:      DIRT,
	}
	SCN_CAVE = &Scene{
		codename:   "cave",
		CellsX:     15,
		CellsY:     20,
		StartPoint: Vector2d{100, 100},
		CamPoint:   Vector2d{0, 0},
		tileA:      DIRT,
		tileB:      LAVA_S1,
	}

	SCENES []*Scene = []*Scene{SCN_PLAINS, SCN_CAVE}

	Cam = Camera{Vector2d{0, 0}, 320, 256}

	PC = Char{
		Solid: &Solid{
			Velocity:    &Vector2d{0, 0},
			Orientation: &Vector2d{0, -1},
			Anim:        MAN_PU_ANIM,
		},
		Lvl:       1,
		CurrentXP: 0,
		NextLvlXP: 100,
		SpeedMod:  0,
		BaseSpeed: 2,
		CurrentHP: 220,
		MaxHP:     250,
		CurrentST: 250,
		MaxST:     300,
		Inventory: []*ItemStack{{GreenBlob, 2}},
	}

	GreenBlob *Item = &Item{
		Name:        "Green Blob",
		Description: "A chunck of slime.",
		Txtr:        powerupsTxt,
		Source:      &sdl.Rect{0, 0, 24, 24},
	}
	CrystalizedJelly *Item = &Item{
		Name:        "Crystalized Jelly",
		Description: "Some believe that the Slime's soul live within it",
		Txtr:        powerupsTxt,
		Source:      &sdl.Rect{24, 0, 24, 24},
		Weight:      2,
		BaseValue:   10,
	}

	SlimeTPL MonsterTemplate = MonsterTemplate{
		Txtr:          monstersTxt,
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

	scene       *Scene
	World       [][]*sdl.Rect
	Interactive []*Solid
	Visual      []*VFXInst
	Spawners    []*SpawnPoint
	Monsters    []*Char
	GUI         []*sdl.Rect
	CullMap     []*Solid
	dbox        DBox  = DBox{BGColor: sdl.Color{90, 90, 90, 255}}
	EventTick   uint8 = EVENT_TICK_LENGTH
	AiTick            = AI_TICK_LENGTH
)

type Vector2d struct {
	X float32
	Y float32
}

type ControlState struct {
	DPAD        Vector2d
	ACTION_A    int32
	ACTION_B    int32
	ACTION_C    int32
	ACTION_MAIN int32
	ACTION_MOD1 int32
}

func (cs *ControlState) Update(keydown []sdl.Keycode, keyup []sdl.Keycode) {
	for _, key := range keydown {
		switch key {
		case KEY_ARROW_UP:
			cs.DPAD.Y = -1
			PC.Solid.Orientation.Y = -1
			break
		case KEY_ARROW_DOWN:
			cs.DPAD.Y = 1
			PC.Solid.Orientation.Y = 1
			break
		case KEY_ARROW_LEFT:
			cs.DPAD.X = -1
			PC.Solid.Orientation.X = -1
			break
		case KEY_ARROW_RIGHT:
			cs.DPAD.X = 1
			PC.Solid.Orientation.X = 1
			break
		case KEY_Z:
			cs.ACTION_A += 1
			break
		case KEY_X:
			cs.ACTION_B += 1
			break
		case KEY_C:
			cs.ACTION_C += 1
			break
		case KEY_SPACE_BAR:
			cs.ACTION_MAIN += 1
			break
		case KEY_LEFT_SHIFT:
			cs.ACTION_MOD1 += 1
			PC.Solid.Speed = (PC.Solid.Speed + PC.SpeedMod) * 1.6
			break
		}
	}
	for _, key := range keyup {
		switch key {
		case KEY_ARROW_UP:
			cs.DPAD.Y = 0
			PC.Solid.Orientation.Y = 0
			PC.Solid.Speed = PC.Solid.Speed + PC.SpeedMod
			break
		case KEY_ARROW_DOWN:
			cs.DPAD.Y = 0
			PC.Solid.Orientation.Y = 0
			break
		case KEY_ARROW_LEFT:
			cs.DPAD.X = 0
			PC.Solid.Orientation.X = 0
			break
		case KEY_ARROW_RIGHT:
			cs.DPAD.X = 0
			PC.Solid.Orientation.X = 0
			break
		case KEY_Z:
			cs.ACTION_A = 0
			break
		case KEY_X:
			if PC.Solid.Anim.PlayMode != 1 {
				PC.MeleeAtk()
			}
			cs.ACTION_B = 0
			break
		case KEY_C:
			if PC.Solid.Anim.PlayMode != 1 {
				PC.CastSpell()
			}
			cs.ACTION_C = 0
			break
		case KEY_SPACE_BAR:
			if !dbox.NextText() {
				actProc()
			}
			cs.ACTION_MAIN = 0
			break
		case KEY_LEFT_SHIFT:
			PC.Solid.Speed = PC.BaseSpeed + PC.SpeedMod
			cs.ACTION_MOD1 = 0
			break
		}
	}
}

func ThrotleValue(v float32, limitAbs float64) float32 {
	v64 := float64(v)
	abs := math.Abs(v64)
	sign := math.Copysign(1, v64)
	if abs > limitAbs {
		return float32(limitAbs * sign)
	}
	return v
}

func ThrotleCeil(v float32, limitAbs float64) float32 {
	v64 := float64(v)
	abs := math.Abs(v64)
	sign := math.Copysign(1, v64)
	if abs > limitAbs {
		return float32(limitAbs * sign)
	}
	return v
}

func ThrotleFloor(v float32, limitAbs float64) float32 {
	if v < 0 {
		return 0
	}
	return v
}

func (s *Solid) UpdateOrientation(inc float32, dcr float32) {
	o = s.Orientation

	if inc.X != 0 {
		o.X += incX
		o.X = ThrotleCeil(o.X, 3)
		if incr.Y == 0 && dcr.Y == 0 && o.X < 3 {
			o.Y -= 1
		}
	}

	if inc.Y != 0 {
		o.Y += inc.Y
		o.Y = ThrotleCeil(o.Y, 3)
		if inc.X == 0 && dcr.X == 0 && o.Y < 3 {
			o.X -= 1
		}
	}

	if dcr.X != 0 {
		o.X -= dcr.X
		o.X = ThrotleFloor(o.X, 0)
	}

	if dcr.Y != 0 {
		o.Y -= dcr.Y
		o.Y = ThrotleFloor(o.Y, 0)
	}
}

func (s *Solid) UpdateVelocity(cs *ControlState) {

	nv := &Vector2d{}
	*nv = *s.Velocity
	if cs.DPAD.X != 0 || s.Velocity.X != 0 {
		if cs.DPAD.X != 0 {
			nv.X += cs.DPAD.X
			nv.X = ThrotleValue(nv.X, 2)
		} else {
			nv.X = float32(math.Abs(float64(nv.X))-1) * s.Orientation.X
		}
	}

	if cs.DPAD.Y != 0 || nv.Y != 0 {
		if cs.DPAD.Y != 0 {
			nv.Y += cs.DPAD.Y
			nv.Y = ThrotleValue(nv.Y, 2)
		} else {
			nv.Y = float32(math.Abs(float64(nv.Y))-1) * s.Orientation.Y
		}
	}
	*s.Velocity = *nv
}

type Camera struct {
	P   Vector2d
	DZx int32
	DZy int32
}

type Scene struct {
	codename   string
	TileSet    *sdl.Texture
	CellsX     int32
	CellsY     int32
	StartPoint Vector2d
	CamPoint   Vector2d
	tileA      *sdl.Rect
	tileB      *sdl.Rect
}

type Event func(source *Solid, subject *Solid)

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

type Loot struct {
	Item *Item
	Perc float32
}

type SpawnPoint struct {
	Position  *sdl.Rect
	Frequency uint16
	LvlMod    uint8
}

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

type InteractionHandlers struct {
	OnCollDmg      float32
	OnCollPushBack int32
	OnCollEvent    Event
	OnPickUp       Event
	OnActDmg       uint16
	OnActPush      *Vector2d
	OnActEvent     Event
	DialogScript   []string
	DoorTo         *Scene
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

type TextEl struct {
	Font         *ttf.Font
	Content      string
	Color        sdl.Color
	BakedContent string
}

type DBox struct {
	SPos     uint8
	CurrText uint8
	Text     []*TextEl
	BGColor  sdl.Color
	Char     *Char
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

func (ss *SpriteSheet) BuildBasicActions(actLength uint8, hasDiagonals bool) *ActionMap {
	O := []Vector2d{F_UP, F_DOWN, F_LEFT, F_RIGHT}
	if hasDiagonals {
		O = append(O, []Vector2d{F_DL, F_DR, F_UL, F_UR}...)
	}
	var AM = &ActionMap{}
	for _, o := range O {
		anim := &Animation{}

		for p := 0; p < int(actLength); p++ {
			anim.Action[p] = ss.GetPose(o, uint8(p))
		}

		switch o {
		case F_UP:
			AM.UP = anim
			break
		case F_DOWN:
			AM.DOWN = anim
			break
		case F_LEFT:
			AM.LEFT = anim
			break
		case F_RIGHT:
			AM.RIGHT = anim
			break
		case F_DL:
			AM.DL = anim
			break
		case F_DR:
			AM.DR = anim
			break
		case F_UL:
			AM.UL = anim
			break
		case F_UR:
			AM.UR = anim
			break
		} // end switch
	} // end for
	return AM
}

func (ss *SpriteSheet) GetPose(o Vector2d, p uint8) *sdl.Rect {

	var (
		poseY int32 = 0
		poseX int32 = 0
	)

	if o.X == -1 && o.Y == 0 {
		poseY = ss.StepW
	}
	if o.X == 1 && o.Y == 0 {
		poseY = ss.StepW * 2
	}
	if o.Y == -1 && o.X == 0 {
		poseY = ss.StepW * 3
	}
	if o.Y == 1 && o.X == 0 {
		poseY = 0
	}
	switch o {
	case F_UP:
		poseY = ss.StepH * 3
		break
	case F_DOWN:
		break
	case F_LEFT:
		poseY = ss.StepH * 1
		break
	case F_RIGHT:
		poseY = ss.StepH * 2
		break
	case F_DL:
		poseX = ss.StepW * 3
		break
	case F_DR:
		poseY = ss.StepH * 2
		poseX = ss.StepW * 3
		break
	case F_UL:
		poseY = ss.StepH * 1
		poseX = ss.StepW * 3
		break
	case F_UR:
		poseY = ss.StepH * 3
		poseX = ss.StepW * 3
		break
	}

	poseX += int32(p) * ss.StepW

	return &sdl.Rect{ss.StX + poseX, ss.StY + poseY, ss.StepW, ss.StepH}
}

func (db *DBox) LoadText(content []string) {
	db.Text = make([]*TextEl, len(content))
	for i, s := range content {
		db.Text[i] = &TextEl{
			Font:    font,
			Content: s,
			Color:   CL_WHITE,
		}
	}
}

func (db *DBox) Present(renderer *sdl.Renderer) {
	if len(db.Text) == 0 {
		return
	}

	ct := db.Text[db.CurrText]
	txtr, w, h := ct.Bake(renderer)
	br := &sdl.Rect{64, winHeight - 128, 512, 120}
	tr := &sdl.Rect{0, 0, w, h}
	bt := &sdl.Rect{64, winHeight - 128, w, h}

	renderer.Copy(transparencyTxt, &sdl.Rect{0, 0, 48, 48}, br)
	renderer.Copy(txtr, tr, bt)
}

func (t *TextEl) Bake(renderer *sdl.Renderer) (*sdl.Texture, int32, int32) {
	surface, _ := t.Font.RenderUTF8_Blended_Wrapped(t.Content, t.Color, int(winWidth))
	defer surface.Free()
	txtr, _ := renderer.CreateTextureFromSurface(surface)

	return txtr, surface.W, surface.H
}

func ActHitBox(source *sdl.Rect, orientation *Vector2d) *sdl.Rect {
	return &sdl.Rect{
		source.X + (int32(orientation.X) * source.W),
		source.Y + (int32(orientation.Y) * source.H),
		source.W,
		source.H,
	}
}

func checkCol(r1 *sdl.Rect, r2 *sdl.Rect) bool {
	return (r1.X < (r2.X+r2.W) &&
		r1.X+r1.W > r2.X &&
		r1.Y < r2.Y+r2.H &&
		r1.Y+r1.H > r2.Y)
}

func (ch *Char) ApplyInvinc() {
	if PC.Invinc == 0 {
		PC.Invinc = time.Now().Add(400 * time.Millisecond).Unix()
	}
}

func ResolveCol(ObjA *Solid, ObjB *Solid) {
	if ObjB.Handlers != nil && ObjB.Handlers.OnCollEvent != nil {
		ObjB.Handlers.OnCollEvent(ObjA, ObjB)
	}
	if ObjB.Collision == 1 {
		ObjA.Velocity = &Vector2d{0, 0}
	}
}

func actProc() {
	action_hit_box := ActHitBox(PC.Solid.Position, PC.Solid.Orientation)

	// Debug hint
	GUI = append(GUI, action_hit_box)
	for _, obj := range CullMap {
		if obj.Handlers != nil &&
			obj.Handlers.OnActEvent != nil &&
			checkCol(action_hit_box, obj.Position) {
			obj.Handlers.OnActEvent(PC.Solid, obj)
			return
		}
	}
}

func onColHdk(tgt *Solid, hdk *Solid) {
	if tgt.CharPtr != nil {
		tgt.CharPtr.depletHP(hdk.Handlers.OnCollDmg)
	}
	Visual = append(Visual, impact.Spawn(tgt.Position, nil))
	hdk.Destroy()
}

func ReleaseSpell(caster *Solid, tgt *Solid) {

	r := ActHitBox(caster.Position, caster.Orientation)
	ttl := time.Now().Add(3 * time.Second)

	h := &Solid{
		Position: r,
		Txt:      glowTxt,
		Ttl:      ttl.Unix(),
		Anim: &Animation{
			Action:   BGLOW_A,
			Pose:     0,
			PoseTick: 16,
		},
		Handlers: &InteractionHandlers{
			OnCollDmg:   50,
			OnCollEvent: onColHdk,
		},
		CPattern: 0,
		MPattern: []Movement{
			Movement{*caster.Orientation, 255},
		},
		Collision: 1,
		Speed:     6,
	}
	Interactive = append(Interactive, h)
}

func PickUp(picker *Solid, item *Solid) {
	if item.Handlers != nil && item.Handlers.OnPickUp != nil {
		item.Handlers.OnPickUp(picker, item)
	}
	if picker.CharPtr != nil {
		picker.SetAnimation(MAN_PU_ANIM, nil)
	}
	item.Destroy()
}

func AddToInv(picker *Solid, item *Solid) {
	if picker.CharPtr != nil && item.ItemPtr != nil {
		for _, iStack := range picker.CharPtr.Inventory {
			if iStack.ItemTpl == item.ItemPtr {
				iStack.Qty += 1
				return
			}
		}
		picker.CharPtr.Inventory = append(picker.CharPtr.Inventory, &ItemStack{item.ItemPtr, 1})
	}
}

func PlayDialog(listener *Solid, speaker *Solid) {
	if len(speaker.Handlers.DialogScript) > 0 {
		dbox.LoadText(speaker.Handlers.DialogScript)
	}
}

func BashDoor(actor *Solid, door *Solid) {
	if actor == PC.Solid {
		change_scene(door.Handlers.DoorTo, nil)
	}
}

func (ch *Char) MeleeAtk() {
	var stCost float32 = 5
	if stCost > ch.CurrentST {
		return
	}
	ch.CurrentST -= stCost
	ch.Solid.SetAnimation(MAN_ATK1_ANIM, nil)
	r := ActHitBox(ch.Solid.Position, ch.Solid.Orientation)
	for _, cObj := range CullMap {
		if cObj.CharPtr != nil && checkCol(r, cObj.Position) {
			cObj.CharPtr.depletHP(15)
			r.W, r.H = 92, 92
			Visual = append(Visual, impact.Spawn(r, ch.Solid.Orientation))
		}
	}
	Visual = append(Visual, hit.Spawn(r, ch.Solid.Orientation))
}

func (ch *Char) CastSpell() {
	var stCost float32 = 20
	if stCost > ch.CurrentST {
		return
	}
	ch.CurrentST -= stCost
	ch.Solid.SetAnimation(MAN_CS_ANIM, ReleaseSpell)
}

func (ch *Char) CurrentFacing() *Animation {

	if ch.Solid.Orientation.X == 0 && ch.Solid.Orientation.Y == 1 {
		return ch.ActionMap.DOWN
	}

	if ch.Solid.Orientation.X == 0 && ch.Solid.Orientation.Y == -1 {
		return ch.ActionMap.UP
	}

	if ch.Solid.Orientation.X == -1 && ch.Solid.Orientation.Y == 0 {
		return ch.ActionMap.LEFT
	}

	if ch.Solid.Orientation.X == 1 && ch.Solid.Orientation.Y == 0 {
		return ch.ActionMap.RIGHT
	}

	if ch.Solid.Orientation.X == 1 && ch.Solid.Orientation.Y == 1 {
		if ch.ActionMap.DR != nil {
			return ch.ActionMap.DR
		} else {
			return ch.ActionMap.RIGHT
		}
	}

	if ch.Solid.Orientation.X == 1 && ch.Solid.Orientation.Y == -1 {
		if ch.ActionMap.UR != nil {
			return ch.ActionMap.UR
		} else {
			return ch.ActionMap.RIGHT
		}
	}

	if ch.Solid.Orientation.X == -1 && ch.Solid.Orientation.Y == 1 {
		if ch.ActionMap.DL != nil {
			return ch.ActionMap.DL
		} else {
			return ch.ActionMap.LEFT
		}
	}

	if ch.Solid.Orientation.X == -1 && ch.Solid.Orientation.Y == -1 {
		if ch.ActionMap.UL != nil {
			return ch.ActionMap.UL
		} else {
			return ch.ActionMap.LEFT
		}
	}

	return ch.ActionMap.DOWN
}

func (ch *Char) depletHP(dmg float32) {
	if ch.Invinc > 0 {
		return
	}
	if dmg > ch.CurrentHP {
		ch.CurrentHP = 0
	} else {
		ch.CurrentHP -= dmg
	}

	PopText(ch.Solid.Position, fmt.Sprintf("%.0f", dmg), CL_WHITE)
}

func PopText(pos *sdl.Rect, content string, color sdl.Color) {
	tEl := &TextEl{Font: font, Content: content, Color: color}
	tPos := &sdl.Rect{pos.X, pos.Y - 30, 20, 20}
	vfi := &VFXInst{Text: tEl, Pos: tPos, Ttl: time.Now().Add(400 * time.Millisecond).Unix()}
	Visual = append(Visual, vfi)
}

func (ch *Char) PushBack(d int32, o *Vector2d) {
	ch.Solid.Position.X += int32(o.X) * d
	ch.Solid.Position.Y += int32(o.Y) * d
	// TODO WIRE PB ANIM into the ActionMap
	ch.Solid.SetAnimation(MAN_PB_ANIM, nil)
}

func (db *DBox) NextText() bool {
	if len(dbox.Text) == 0 {
		return false
	}
	dbox.CurrText += 1
	if int(dbox.CurrText+1) > len(dbox.Text) {
		dbox.Text = []*TextEl{}
		dbox.CurrText = 0
	}
	return true
}

func (s *Solid) SetAnimation(an *Animation, evt Event) {
	var nA *Animation = &Animation{}
	*nA = *an

	s.Anim = nA
	s.Anim.Pose = 0
	s.Anim.PoseTick = 16
	s.Anim.After = evt
}

func (s *Solid) PlayAnimation() {

	s.Anim.PoseTick -= 1

	if s.Anim.PoseTick <= 0 {
		s.Anim.PoseTick = 12
		if s.CharPtr != nil {
			anim := s.CharPtr.CurrentFacing()
			prvPose := s.Anim.Pose
			s.Anim.Pose = getNextPose(s.Anim.Action, s.Anim.Pose)
			if anim != nil && s.Anim.Pose <= prvPose && s.Anim.PlayMode == 1 {
				if s.Anim.After != nil {
					s.Anim.After(s, nil)
				}
				s.SetAnimation(anim, nil)
			}
		} else {
			s.Anim.Pose = getNextPose(s.Anim.Action, s.Anim.Pose)
		}
	}
}

func (s *Solid) procMovement() {
	speed := s.Speed
	if s.Velocity.X != 0 && s.Velocity.Y != 0 {
		speed -= 0.5
		if speed < 1 {
			speed = 1
		}
	}
	np := &sdl.Rect{
		(s.Position.X + int32(s.Velocity.X*speed)),
		(s.Position.Y + int32(s.Velocity.Y*speed)),
		s.Position.W,
		s.Position.H,
	}
	var outbound bool = (np.X <= 0 ||
		np.Y <= 0 ||
		np.X > int32(scene.CellsX*TSzi) ||
		np.Y > int32(scene.CellsY*TSzi))

	if (np.X == s.Position.X && np.Y == s.Position.Y) || outbound {
		return
	}

	if s.CharPtr != nil {
		anim := s.CharPtr.CurrentFacing()
		if anim != nil && s.Anim != nil && s.Anim.PlayMode != 1 {
			s.Anim.PlayMode = anim.PlayMode
			s.Anim.Action = anim.Action
		}
	}

	s.Position = np
}

func (s *Solid) LoSCheck() bool {
	LoS := &sdl.Rect{
		s.Position.X - s.LoS,
		s.Position.Y - s.LoS,
		s.Position.W + (s.LoS * 2),
		s.Position.H + (s.LoS * 2),
	}

	if !checkCol(PC.Solid.Position, LoS) {
		return false
	}
	return true
}

func (s *Solid) chase() {

	s.Velocity.X = 0
	s.Velocity.Y = 0

	diffX := math.Abs(float64(s.Position.X - s.Chase.Position.X))
	diffY := math.Abs(float64(s.Position.Y - s.Chase.Position.Y))

	if s.Position.X > s.Chase.Position.X {
		s.Velocity.X = -1
	}

	if s.Position.X < s.Chase.Position.X {
		s.Velocity.X = 1
	}
	if s.Position.Y > s.Chase.Position.Y {
		s.Velocity.Y = -1
	}

	if s.Position.Y < s.Chase.Position.Y {
		s.Velocity.Y = 1
	}

	if int32(diffX) < CharSize+10 && int32(diffY) < CharSize+10 && s.CharPtr != nil {

		if s.CharPtr.AtkCoolDownC <= 0 {
			r := ActHitBox(s.Position, s.Orientation)
			Visual = append(Visual, hit.Spawn(r, s.Orientation))
			if checkCol(r, PC.Solid.Position) {
				PC.depletHP(s.Handlers.OnCollDmg)
				Visual = append(Visual, impact.Spawn(r, s.Orientation))
				PC.ApplyInvinc()
			}
			s.CharPtr.AtkCoolDownC += s.CharPtr.AtkCoolDown
		}

		return
	} else {
		s.procMovement()
	}

	return
}

func (s *Solid) peformPattern(sp float32) {
	anon := func(c uint32, mvs []Movement) *Movement {
		var sum uint32 = 0
		for _, mp := range mvs {
			sum += uint32(mp.Ticks)
			if sum > s.CPattern {
				return &mp
			}
		}
		s.CPattern = 0
		return nil
	}

	mov := anon(s.CPattern, s.MPattern)
	if mov != nil && s.Position != nil {
		s.Orientation = &mov.Orientation
		s.Velocity = &mov.Orientation
		s.procMovement()
		s.CPattern += uint32(sp)
	}
}

func (s *Solid) Destroy() {
	s.Position = nil
	s.Source = nil
	s.Anim = nil
	s.Handlers = nil
	s.Txt = nil
	s.Collision = 0
}

func change_scene(new_scene *Scene, staring_pos *Vector2d) {
	new_scene.build()
	Interactive = []*Solid{}
	new_scene.populate(200)
	scene = new_scene

	PC.Solid.Position = V2R(new_scene.StartPoint, CharSize, CharSize)
	PC.Solid.Speed = PC.BaseSpeed + PC.SpeedMod
}

func (s *Scene) build() {
	ni, nj := int(s.CellsX)+1, int(s.CellsY)+1

	println("start build: ", s.codename)
	rand.Seed(int64(time.Now().Nanosecond()))
	World = make([][]*sdl.Rect, ni)
	println("Allocation rows", ni)

	for i := 0; i < ni; i++ {
		World[i] = make([]*sdl.Rect, nj)
		println("Allocation columns for", i, nj)

		for j := 0; j < nj; j++ {
			println("c", i, j)
			tile := s.tileA
			if rand.Int31n(100) < 10 {
				tile = s.tileB
			}

			World[i][j] = tile
		}
	}

	println("==finish build: ", s.codename)
}

func (s *Scene) populate(population int) {
	doorTo := SCN_PLAINS
	if s.codename == "plains" {
		doorTo = SCN_CAVE
	}

	for i := 0; i < population; i++ {
		cX, cY := rand.Int31n(s.CellsX), rand.Int31n(s.CellsY)

		absolute_pos := &sdl.Rect{cX * TSzi, cY * TSzi, CharSize, CharSize}
		sol := &Solid{}

		switch rand.Int31n(9) {
		case 1:
			sol = &Solid{
				Position: absolute_pos,
				Txt:      s.TileSet,
				Anim:     LAVA_ANIM,
				Handlers: &InteractionHandlers{
					OnCollDmg:      12,
					OnCollPushBack: 16,
				},
				Collision: 0,
			}
			Interactive = append(Interactive, sol)
			break
		case 2:
			sol = &Solid{
				Position:  absolute_pos,
				Txt:       s.TileSet,
				Source:    DOOR,
				Anim:      nil,
				Collision: 1,
				Handlers: &InteractionHandlers{
					OnActEvent: BashDoor,
					DoorTo:     doorTo,
				},
			}
			Interactive = append(Interactive, sol)
			break
		case 3:
			absolute_pos.H, absolute_pos.W = 64, 64
			sol = &Solid{
				Position:  absolute_pos,
				Txt:       glowTxt,
				Anim:      LIFE_ORB_ANIM,
				Collision: 0,
				Handlers: &InteractionHandlers{
					OnCollEvent: PickUp,
					OnPickUp: func(healed *Solid, orb *Solid) {
						if healed.CharPtr != nil {
							healed.CharPtr.CurrentHP += 10
						}
					},
				},
			}
			Interactive = append(Interactive, sol)
			break
		case 4:
			absolute_pos.W, absolute_pos.H = 128, 128

			for _, sp2 := range Spawners {
				if checkCol(absolute_pos, sp2.Position) {
					return
				}
			}

			rand.Seed(int64(time.Now().Nanosecond()))
			spw := &SpawnPoint{
				Position:  absolute_pos,
				Frequency: uint16(rand.Int31n(5)),
			}
			Spawners = append(Spawners, spw)
			break
		case 5:
			fnpc := Char{
				Solid: &Solid{
					Position:    absolute_pos,
					Velocity:    &Vector2d{0, 0},
					Orientation: &Vector2d{0, 1},
					Speed:       2,
					Txt:         spritesheetTxt,
					Collision:   1,
					Handlers: &InteractionHandlers{
						OnActEvent:   PlayDialog,
						DialogScript: []string{"more", "npc", "chitchat"},
					},
					CPattern: 0,
					MPattern: []Movement{
						Movement{F_UP, 90},
						Movement{F_RIGHT, 10},
						Movement{F_DOWN, 50},
						Movement{F_LEFT, 10},
					},
				},
				ActionMap: PC.ActionMap,
				CurrentHP: 9999,
				MaxHP:     9999,
			}
			fnpc.Solid.Anim = MAN_PB_ANIM

			fnpc.Solid.CharPtr = &fnpc
			Monsters = append(Monsters, &fnpc)
			break
		}

	}
}

func (s *Scene) _terrainRender(renderer *sdl.Renderer) {
	var Source *sdl.Rect
	var init int32 = 0

	var offsetX, offsetY int32 = TSzi, TSzi
	// Rendering the terrain
	for winY := init; winY < winHeight; winY += offsetY {
		for winX := init; winX < winWidth; winX += offsetX {

			offsetX = (TSzi - (int32(Cam.P.X)+winX)%TSzi)
			offsetY = (TSzi - (int32(Cam.P.Y)+winY)%TSzi)

			worldCellX := uint16((int32(Cam.P.X) + winX) / TSzi)
			worldCellY := uint16((int32(Cam.P.Y) + winY) / TSzi)
			screenPos := sdl.Rect{winX, winY, offsetX, offsetY}

			if worldCellX > uint16(s.CellsX) || worldCellY > uint16(s.CellsY) || worldCellX < 0 || worldCellY < 0 {
				continue
			}

			gfx := World[worldCellX][worldCellY]

			if offsetX != TSzi || offsetY != TSzi {
				Source = &sdl.Rect{gfx.X + (TSzi - offsetX), gfx.Y + (TSzi - offsetY), offsetX, offsetY}
			} else {
				Source = gfx
			}
			if Source != nil && &screenPos != nil {
			}
			renderer.Copy(s.TileSet, Source, &screenPos)
		}
	}
}

func (s *Scene) _solidsRender(renderer *sdl.Renderer) {
	CullMap = []*Solid{}

	for _, obj := range Interactive {
		if obj.Position == nil {
			continue
		}

		scrPos := worldToScreen(obj.Position, Cam)

		renderer.SetDrawColor(0, 255, 0, 255)
		renderer.DrawRect(scrPos)

		if inScreen(scrPos) {
			var src *sdl.Rect
			if obj.Anim != nil {
				src = obj.Anim.Action[obj.Anim.Pose]
			} else {
				src = obj.Source
			}
			renderer.Copy(obj.Txt, src, scrPos)
			CullMap = append(CullMap, obj)
		}
	}
}

func (s *Scene) _monstersRender(renderer *sdl.Renderer) {

	for _, mon := range Monsters {
		if mon.Solid.Anim == nil || mon.Solid == nil || mon.Solid.Position == nil {
			continue
		}
		scrPos := worldToScreen(mon.Solid.Position, Cam)

		if inScreen(scrPos) {

			src := mon.Solid.Anim.Action[mon.Solid.Anim.Pose]
			renderer.Copy(mon.Solid.Txt, src, scrPos)
			scrPos.Y += mon.Solid.Position.H / 8
			renderer.Copy(spritesheetTxt, SHADOW, scrPos)

			CullMap = append(CullMap, mon.Solid)

			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.FillRect(&sdl.Rect{scrPos.X, scrPos.Y - 8, 32, 4})
			renderer.SetDrawColor(0, 255, 0, 255)
			renderer.FillRect(&sdl.Rect{scrPos.X, scrPos.Y - 8, int32(32 * calcPerc(mon.CurrentHP, mon.MaxHP) / 100), 4})
		}
	}
}

func (s *Scene) _GUIRender(renderer *sdl.Renderer) {

	// Gray overlay
	renderer.Copy(transparencyTxt, &sdl.Rect{0, 0, 48, 48}, &sdl.Rect{0, 0, 120, 60})

	// HEALTH BAR
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.FillRect(&sdl.Rect{10, 10, 100, 4})
	renderer.SetDrawColor(0, 255, 0, 255)
	renderer.FillRect(&sdl.Rect{10, 10, int32(calcPerc(PC.CurrentHP, PC.MaxHP)), 4})

	// MANA BAR
	renderer.SetDrawColor(190, 0, 120, 255)
	renderer.FillRect(&sdl.Rect{10, 24, 100, 4})
	renderer.SetDrawColor(0, 0, 255, 255)
	renderer.FillRect(&sdl.Rect{10, 24, int32(calcPerc(PC.CurrentST, PC.MaxST)), 4})

	// XP BAR
	renderer.SetDrawColor(90, 90, 0, 255)
	renderer.SetDrawColor(190, 190, 0, 255)
	renderer.FillRect(&sdl.Rect{10, 38, int32(calcPerc(float32(PC.CurrentXP), float32(PC.NextLvlXP))), 4})

	renderer.Copy(transparencyTxt, &sdl.Rect{0, 0, 48, 30}, &sdl.Rect{0, 60, 240, 30})

	for i, stack := range PC.Inventory {
		counter := TextEl{
			Content: fmt.Sprintf("%d", stack.Qty),
			Font:    font,
			Color:   CL_WHITE,
		}

		counterTxtr, cW, cH := counter.Bake(renderer)
		pos := sdl.Rect{8 + (int32(i) * 32), 60, 24, 24}
		renderer.Copy(stack.ItemTpl.Txtr, stack.ItemTpl.Source, &pos)
		pos.Y += 16
		pos.X += 16
		pos.W = cW
		pos.H = cH
		renderer.Copy(counterTxtr, &sdl.Rect{0, 0, cW, cH}, &pos)
	}

	for _, el := range GUI {
		scrPos := worldToScreen(el, Cam)
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(scrPos)
	}

	lvl_TextEl := TextEl{
		Font:    font,
		Content: fmt.Sprintf("Lvl. %d", PC.Lvl),
		Color:   CL_WHITE,
	}
	lvl_txtr, W, H := lvl_TextEl.Bake(renderer)
	renderer.Copy(lvl_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{128, 60, W, H})
	debug_info(renderer)
	dbox.Present(renderer)

	for _, spw := range Spawners {
		renderer.DrawRect(worldToScreen(spw.Position, Cam))
	}
}

func (s *Scene) _VFXRender(renderer *sdl.Renderer) {
	for _, vi := range Visual {
		if vi.Pos == nil {
			continue
		}
		scrp := worldToScreen(vi.Pos, Cam)
		if inScreen(scrp) {
			if vi.Text != nil {
				txtr, w, h := vi.Text.Bake(renderer)
				renderer.Copy(txtr, &sdl.Rect{0, 0, w, h}, scrp)
			} else {
				if vi.Flip.X == -1 {
					renderer.CopyEx(vi.Vfx.Txtr, vi.CurrentFrame(), scrp, 0, nil, sdl.FLIP_HORIZONTAL)
				} else {
					renderer.Copy(vi.Vfx.Txtr, vi.CurrentFrame(), scrp)
				}
			}
		}
	}
}

func (s *Scene) render(renderer *sdl.Renderer) {
	renderer.Clear()
	scrPos := worldToScreen(PC.Solid.Position, Cam)

	s._terrainRender(renderer)
	s._solidsRender(renderer)
	s._monstersRender(renderer)
	// Rendering the PC
	if !(PC.Invinc > 0 && EventTick == 2) {
		renderer.Copy(spritesheetTxt, PC.Solid.Anim.Action[PC.Solid.Anim.Pose], scrPos)
		scrPos.Y += 12
		renderer.Copy(spritesheetTxt, SHADOW, scrPos)
	}
	s._VFXRender(renderer)
	s._GUIRender(renderer)

	// FLUSH FRAME
	renderer.Present()
}

func (s *Scene) update() {
	EventTick -= 1
	AiTick -= 1

	now := time.Now().Unix()

	if PC.Invinc > 0 && now > PC.Invinc {
		println("pc invincibility ended!")
		PC.Invinc = 0
	}

	if len(dbox.Text) > 0 {
		AiTick = AI_TICK_LENGTH
	}

	newScreenPos := worldToScreen(PC.Solid.Position, Cam)
	if (Cam.DZx - newScreenPos.X) > 0 {
		Cam.P.X -= float32(Cam.DZx - newScreenPos.X)
	}

	if (winWidth - Cam.DZx) < (newScreenPos.X + TSzi) {
		Cam.P.X += float32((newScreenPos.X + TSzi) - (winWidth - Cam.DZx))
	}

	if (Cam.DZy - newScreenPos.Y) > 0 {
		Cam.P.Y -= float32(Cam.DZy - newScreenPos.Y)
	}

	if (winHeight - Cam.DZy) < (newScreenPos.Y + TSzi) {
		Cam.P.Y += float32((newScreenPos.Y + TSzi) - (winHeight - Cam.DZy))
	}

	for _, obj := range CullMap {
		if obj.Position == nil {
			continue
		}
		fr := feetRect(PC.Solid.Position)
		if checkCol(fr, obj.Position) {
			ResolveCol(PC.Solid, obj)
		}
	}

	for _, i := range CullMap {
		if i.Position == nil {
			continue
		}
		// Ttl kill
		if i.Ttl > 0 && i.Ttl < now {
			i.Destroy()
			continue
		}

		if i.CharPtr != nil && i.CharPtr.CurrentHP <= 0 {
			if i.CharPtr.Drop != nil {
				PlaceDrop(i.CharPtr.Drop, i.Position)
			}
			Visual = append(Visual, puff.Spawn(&sdl.Rect{i.Position.X, i.Position.Y, 92, 92}, nil))
			i.Destroy()

			PC.CurrentXP += uint16(i.CharPtr.MaxHP / 10)
			if PC.CurrentXP >= PC.NextLvlXP {
				PC.CurrentXP = 0
				PC.Lvl++
				PC.NextLvlXP = PC.NextLvlXP * uint16(1+PC.Lvl/2)
			}
			continue
		}

		fr := feetRect(PC.Solid.Position)
		if checkCol(fr, i.Position) {
			ResolveCol(PC.Solid, i)
		}

		if i.Anim != nil {
			i.PlayAnimation()
		}

		if AiTick == 0 {

			if i.Anim != nil && i.Anim.PlayMode == 1 {
				continue
			}

			if i.Chase != nil && i.LoSCheck() {
				i.chase()
			} else {
				var sp float32 = 2
				if i.CharPtr != nil {
					sp = i.Speed
				}
				i.peformPattern(sp)
			}
		}

		for _, j := range CullMap {
			if j == i || j.Position == nil {
				continue
			}
			fr := feetRect(i.Position)
			if checkCol(fr, j.Position) {
				ResolveCol(i, j)
			}
		}
	}

	if EventTick == 0 {
		for _, spw := range Spawners {
			if spw.Frequency <= 0 {
				spw.Frequency += 1
			}

			if uint16(rand.Int31n(1000)) < spw.Frequency {
				spw.Produce()
				spw.Frequency -= 1
			}
		}

		EventTick = EVENT_TICK_LENGTH
	}

	if AiTick == 0 {
		AiTick = AI_TICK_LENGTH
	}

	for _, vi := range Visual {
		if vi.Ttl > 0 && vi.Ttl < now {
			vi.Destroy()
			continue
		}
		vi.UpdateAnim()
	}
} // end update()

func debug_info(renderer *sdl.Renderer) {
	dbg_content := fmt.Sprintf("px %d py %d|Cx %.1f Cy %.1f | vx %.1f vy %.1f (%.1f, %.1f) An:%d/%d/%d cull %d i %d cX %d cY %d L %dus ETick%d AiTick%d",
		PC.Solid.Position.X, PC.Solid.Position.Y, Controls.DPAD.X, Controls.DPAD.Y, PC.Solid.Velocity.X, PC.Solid.Velocity.Y, PC.Solid.Orientation.X,
		PC.Solid.Orientation.Y, PC.Solid.Anim.Pose, PC.Solid.Anim.PoseTick, PC.Solid.Anim.PlayMode, len(CullMap),
		len(Interactive), Cam.P.X, Cam.P.Y, game_latency, EventTick, AiTick)

	dbg_TextEl := TextEl{
		Font:    font,
		Content: dbg_content,
		Color:   CL_WHITE,
	}
	dbg_txtr, W, H := dbg_TextEl.Bake(renderer)
	renderer.Copy(dbg_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{0, winHeight - H, W, H})
}

func calcPerc(v1 float32, v2 float32) float32 {
	return (float32(v1) / float32(v2) * 100)
}

func worldToScreen(pos *sdl.Rect, cam Camera) *sdl.Rect {
	return &sdl.Rect{
		(pos.X - int32(cam.P.X)),
		(pos.Y - int32(cam.P.Y)),
		pos.W, pos.H,
	}
}

func isMoving(vel *Vector2d) bool {
	return (vel.X != 0 || vel.Y != 0)
}

func getNextPose(action [8]*sdl.Rect, currPose uint8) uint8 {
	if action[currPose+1] == nil {
		return 0
	} else {
		return currPose + 1
	}
}

func PlaceDrop(item *Item, origin *sdl.Rect) {
	instance := ItemInstance{
		ItemTpl: item,
		Solid: &Solid{
			ItemPtr: item,
			Txt:     item.Txtr,
			Source:  item.Source,
			Position: &sdl.Rect{
				origin.X,
				origin.Y,
				item.Source.W,
				item.Source.H,
			},
			Handlers: &InteractionHandlers{
				OnActEvent: PickUp,
				OnPickUp:   AddToInv,
			},
		},
	}

	Interactive = append(Interactive, instance.Solid)
}

func inScreen(r *sdl.Rect) bool {
	return (r.X > (r.W*-1) && r.X < winWidth && r.Y > (r.H*-1) && r.Y < winHeight)
}

func V2R(v Vector2d, w int32, h int32) *sdl.Rect {
	return &sdl.Rect{int32(v.X), int32(v.Y), w, h}
}

func feetRect(pos *sdl.Rect) *sdl.Rect {
	third := pos.H / 3
	return &sdl.Rect{pos.X, pos.Y + third, pos.W, pos.H - third}
}

func (sp *SpawnPoint) Produce() {
	px := float32(rand.Int31n((sp.Position.X+sp.Position.W)-sp.Position.X) + sp.Position.X)
	py := float32(rand.Int31n((sp.Position.Y+sp.Position.H)-sp.Position.Y) + sp.Position.Y)
	mon := MonsterFactory(&OrcTPL, sp.LvlMod, Vector2d{px, py})

	Monsters = append(Monsters, mon)
}

func MonsterFactory(monsterTpl *MonsterTemplate, lvlMod uint8, pos Vector2d) *Char {

	variance := uint8(math.Floor(float64(rand.Float32() * monsterTpl.LvlVariance * 100)))
	lvl := uint8((monsterTpl.Lvl + lvlMod) + variance)
	hp := monsterTpl.HP + float32(lvl*2)
	sizeMod := int32(float32(lvl-monsterTpl.Lvl) * monsterTpl.ScalingFactor)
	W := (monsterTpl.Size + sizeMod)
	H := (monsterTpl.Size + sizeMod)

	var DropItem *Item
	var sumP float32
	R := rand.Float32()
	for _, l := range monsterTpl.Loot {
		sumP += l.Perc
		if R < sumP {
			DropItem = l.Item
			break
		}
	}

	mon := Char{
		Lvl: lvl,
		Solid: &Solid{
			Position:    &sdl.Rect{int32(pos.X), int32(pos.Y), W, H},
			Velocity:    &Vector2d{0, 0},
			Orientation: &Vector2d{0, 0},
			Txt:         monsterTpl.Txtr,
			Speed:       1,
			Collision:   2,
			Handlers: &InteractionHandlers{
				OnCollDmg: 12,
			},
			CPattern: 0,
			LoS:      monsterTpl.LoS,
			MPattern: []Movement{
				Movement{F_DOWN, 50},
				Movement{F_UP, 90},
				Movement{F_RIGHT, 10},
				Movement{F_LEFT, 10},
			},
			Chase: PC.Solid,
		},
		ActionMap:   monsterTpl.ActionMap,
		AtkSpeed:    monsterTpl.AtkSpeed,
		AtkCoolDown: monsterTpl.AtkCoolDown,
		CurrentHP:   hp,
		MaxHP:       hp,
		Drop:        DropItem,
	}
	mon.Solid.SetAnimation(monsterTpl.ActionMap.DOWN, nil)
	mon.Solid.CharPtr = &mon

	Visual = append(Visual, puff.Spawn(&sdl.Rect{int32(pos.X), int32(pos.Y), 92, 92}, nil))

	return &mon
}

func (vi *VFXInst) UpdateAnim() {
	if vi.Text != nil {
		vi.Pos.Y -= 1
		vi.Pos.X -= 1
		return
	}
	if vi.Vfx == nil {
		return
	}

	vi.CurrTick -= 1
	if vi.CurrTick == 0 {
		vi.Pose += 1
		if vi.Vfx.Strip[vi.Pose] == nil {
			if vi.Loop <= 0 {
				vi.Destroy()
				return
			} else {
				vi.Pose = 0
			}
		}
		vi.CurrTick = vi.Tick
	}
}

func (vi *VFXInst) Destroy() {
	vi.Vfx, vi.Pos = nil, nil
}

func (vi *VFXInst) CurrentFrame() *sdl.Rect {
	return vi.Vfx.Strip[vi.Pose]
}

func (v *VFX) Spawn(Position *sdl.Rect, flip *Vector2d) *VFXInst {
	var ttl int64

	i := &VFXInst{
		Vfx:      v,
		Pos:      Position,
		Ttl:      ttl,
		Tick:     v.DefaultSpeed,
		CurrTick: v.DefaultSpeed,
	}

	if flip != nil {
		i.Flip = *flip
	}

	return i
}

func catchEvents() bool {

	PC.Solid.PlayAnimation()

	KUs := []sdl.Keycode{}
	KDs := []sdl.Keycode{}

	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyDownEvent:
			KDs = append(KDs, t.Keysym.Sym)
		case *sdl.KeyUpEvent:
			KUs = append(KUs, t.Keysym.Sym)
		}
	}

	Controls.Update(KDs, KUs)
	PC.Solid.UpdateVelocity(Controls)

	if PC.Solid.Anim.PlayMode == 0 {
		PC.Solid.procMovement()
	}

	if isMoving(PC.Solid.Velocity) && PC.Solid.Speed > PC.BaseSpeed+PC.SpeedMod {
		// Play animation again when running
		PC.Solid.PlayAnimation()

		dpl := (PC.MaxST * 0.0009)

		if PC.CurrentST <= dpl {
			PC.CurrentST = 0
			PC.Solid.Speed = PC.BaseSpeed + PC.SpeedMod
		} else {
			PC.CurrentST -= dpl
		}
	} else {
		if !isMoving(PC.Solid.Velocity) && PC.CurrentST < PC.MaxST {
			PC.CurrentST += (PC.MaxST * 0.001)
		}
	}

	return true
}

func main() {

	runtime.GOMAXPROCS(1)

	var window *sdl.Window
	var renderer *sdl.Renderer

	window, _ = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(winWidth), int(winHeight), sdl.WINDOW_SHOWN)
	renderer, _ = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer window.Destroy()
	defer renderer.Destroy()

	ttf.Init()
	font, _ = ttf.OpenFont("assets/textures/PressStart2P.ttf", 12)

	tilesetImg, _ := img.Load("assets/textures/ts1.png")
	spritesheetImg, _ := img.Load("assets/textures/main_char.png")
	powerupsImg, _ := img.Load("assets/textures/powerups_ts.png")
	glowImg, _ := img.Load("assets/textures/glowing_ts.png")
	monstersImg, _ := img.Load("assets/textures/monsters.png")
	transparencyImg, _ := img.Load("assets/textures/transparency.png")
	puffImg, _ := img.Load("assets/textures/puff.png")
	hitImg, _ := img.Load("assets/textures/hit.png")
	defer monstersImg.Free()
	defer tilesetImg.Free()
	defer spritesheetImg.Free()
	defer powerupsImg.Free()
	defer glowImg.Free()
	defer transparencyImg.Free()
	defer puffImg.Free()
	defer hitImg.Free()

	tilesetTxt, _ = renderer.CreateTextureFromSurface(tilesetImg)
	spritesheetTxt, _ = renderer.CreateTextureFromSurface(spritesheetImg)
	powerupsTxt, _ = renderer.CreateTextureFromSurface(powerupsImg)
	glowTxt, _ = renderer.CreateTextureFromSurface(glowImg)
	monstersTxt, _ = renderer.CreateTextureFromSurface(monstersImg)
	transparencyTxt, _ = renderer.CreateTextureFromSurface(transparencyImg)
	puffTxt, _ = renderer.CreateTextureFromSurface(puffImg)
	hitTxt, _ = renderer.CreateTextureFromSurface(hitImg)
	defer tilesetTxt.Destroy()
	defer spritesheetTxt.Destroy()
	defer powerupsTxt.Destroy()
	defer glowTxt.Destroy()
	defer monstersTxt.Destroy()
	defer transparencyTxt.Destroy()
	defer puffTxt.Destroy()
	defer hitTxt.Destroy()

	GreenBlob.Txtr = powerupsTxt
	CrystalizedJelly.Txtr = powerupsTxt
	SlimeTPL.Txtr = monstersTxt

	BatSS := &SpriteSheet{monstersTxt, 0, 0, 48, 48}
	BatActionMap := BatSS.BuildBasicActions(3, false)
	OrcSS := &SpriteSheet{monstersTxt, 288, 0, 48, 48}
	OrcActionMap := OrcSS.BuildBasicActions(3, false)

	MainCharSS := &SpriteSheet{spritesheetTxt, 0, 0, 32, 32}
	MainCharActionMap := MainCharSS.BuildBasicActions(3, true)
	PC.ActionMap = MainCharActionMap

	BatTPL = MonsterTemplate{
		Txtr:          monstersTxt,
		ActionMap:     BatActionMap,
		Loot:          [8]Loot{{CrystalizedJelly, 0.5}, {GreenBlob, 0.5}},
		Lvl:           1,
		HP:            25,
		LoS:           90,
		Size:          32,
		LvlVariance:   0.3,
		ScalingFactor: 0.6,
	}

	OrcTPL = MonsterTemplate{
		Txtr:          monstersTxt,
		ActionMap:     OrcActionMap,
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

	puff = VFX{Txtr: puffTxt, Strip: PUFF_A, DefaultSpeed: 4}
	hit = VFX{Txtr: hitTxt, Strip: HIT_A, DefaultSpeed: 4}
	impact = VFX{Txtr: hitTxt, Strip: HIT_B, DefaultSpeed: 3}

	var running bool = true

	for _, scn := range SCENES {
		scn.TileSet = tilesetTxt
	}

	PC.Solid.CharPtr = &PC
	renderer.SetDrawColor(0, 0, 255, 255)
	scene = SCENES[0]
	dbox.LoadText([]string{"Hello World!", "Again!"})
	change_scene(scene, nil)
	for running {
		then := time.Now()

		scene.update()
		scene.render(renderer)
		running = catchEvents()

		game_latency = (time.Since(then) / time.Microsecond)

		sdl.Delay(33 - uint32(game_latency/1000))
	}
}
