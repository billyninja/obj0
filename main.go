package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"math/rand"
	"runtime"
	"time"
)

const (
	winWidth, winHeight int32 = 640, 480
	tSz                       = 32
	cY                        = winHeight / tSz
	cX                        = winWidth / tSz
	WORLD_CELLS_X             = 500
	WORLD_CELLS_Y             = 200
	KEY_ARROW_UP              = 1073741906
	KEY_ARROW_DOWN            = 1073741905
	KEY_ARROW_LEFT            = 1073741904
	KEY_ARROW_RIGHT           = 1073741903
	KEY_LEFT_SHIT             = 1073742049
	KEY_SPACE_BAR             = 1073741824 //32
	KEY_C                     = 99
)

type Vector2d struct {
	X int32
	Y int32
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

type Camera struct {
	P   Vector2d
	DZx int32
	DZy int32
}

type Event func(source *Solid, subject *Solid)

type InteractionHandlers struct {
	OnCollDmg   uint16
	OnCollPush  *Vector2d
	OnCollEvent Event
	OnPickUp    Event
	OnActDmg    uint16

	OnActPush    *Vector2d
	OnActEvent   Event
	DialogScript []string
	DoorTo       *Scene
}

type Solid struct {
	Velocity  *Vector2d
	Position  *sdl.Rect
	Source    *sdl.Rect
	Facing    Facing
	Anim      *Animation
	Handlers  *InteractionHandlers
	Txt       *sdl.Texture
	Ttl       int64
	Collision uint8

	// AI RELATED
	CPattern uint32
	MPattern []Movement
	Chase    *Solid
	CharPtr  *Char
}

type Animation struct {
	Action   [8]*sdl.Rect
	Pose     uint8
	PoseTick uint32
}

type Movement struct {
	Orientation Vector2d
	Ticks       uint8
}

type PowerUp struct {
	Code        uint8
	Name        string
	Description string
	Ico         *sdl.Rect
	IcoTxt      *sdl.Texture
}

type Char struct {
	Solid     *Solid
	Buffs     []*PowerUp
	Speed     int32
	Lvl       uint8
	CurrentXP uint16
	NextLvlXP uint16
	CurrentHP uint16
	MaxHP     uint16
	CurrentST uint16
	MaxST     uint16
}

type Facing struct {
	Orientation Vector2d
	Up          [8]*sdl.Rect
	Down        [8]*sdl.Rect
	Left        [8]*sdl.Rect
	Right       [8]*sdl.Rect
	DownLeft    [8]*sdl.Rect
	DownRight   [8]*sdl.Rect
	UpLeft      [8]*sdl.Rect
	UpRight     [8]*sdl.Rect
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

func (db *DBox) LoadText(content []string) {
	db.Text = make([]*TextEl, len(content))
	for i, s := range content {
		db.Text[i] = &TextEl{
			Font:    font,
			Content: s,
			Color:   sdl.Color{255, 255, 255, 255},
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

	renderer.SetDrawColor(db.BGColor.R, db.BGColor.G, db.BGColor.B, db.BGColor.A)
	renderer.FillRect(br)
	renderer.Copy(txtr, tr, bt)
}

func ActHitBox(source *sdl.Rect, facing Vector2d) *sdl.Rect {
	return &sdl.Rect{
		source.X + (facing.X * tSz),
		source.Y + (facing.Y * tSz),
		source.W,
		source.H,
	}
}

var (
	winTitle     string = "Go-SDL2 Obj0"
	event        sdl.Event
	font         *ttf.Font
	game_latency time.Duration

	tilesetTxt     *sdl.Texture = nil
	spritesheetTxt *sdl.Texture = nil
	particlesTxt   *sdl.Texture = nil
	powerupsTxt    *sdl.Texture = nil
	glowTxt        *sdl.Texture = nil
	slimeTxt       *sdl.Texture = nil

	GRASS     *sdl.Rect = &sdl.Rect{0, 0, tSz, tSz}
	DIRT                = &sdl.Rect{703, 0, tSz, tSz}
	WALL                = &sdl.Rect{0, 64, tSz, tSz}
	DOOR                = &sdl.Rect{256, 32, tSz, tSz}
	BF_ATK_UP           = &sdl.Rect{72, 24, 24, 24}
	BF_DEF_UP           = &sdl.Rect{96, 24, 24, 24}

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
	MAN_FRONT_R *sdl.Rect = &sdl.Rect{0, 0, tSz, tSz}
	MAN_FRONT_N *sdl.Rect = &sdl.Rect{32, 0, tSz, tSz}
	MAN_FRONT_L *sdl.Rect = &sdl.Rect{64, 0, tSz, tSz}
	MAN_LEFT_R  *sdl.Rect = &sdl.Rect{0, 32, tSz, tSz}
	MAN_LEFT_N  *sdl.Rect = &sdl.Rect{32, 32, tSz, tSz}
	MAN_LEFT_L  *sdl.Rect = &sdl.Rect{64, 32, tSz, tSz}
	MAN_RIGHT_R *sdl.Rect = &sdl.Rect{0, 64, tSz, tSz}
	MAN_RIGHT_N *sdl.Rect = &sdl.Rect{32, 64, tSz, tSz}
	MAN_RIGHT_L *sdl.Rect = &sdl.Rect{64, 64, tSz, tSz}
	MAN_BACK_R  *sdl.Rect = &sdl.Rect{0, 96, tSz, tSz}
	MAN_BACK_N  *sdl.Rect = &sdl.Rect{32, 96, tSz, tSz}
	MAN_BACK_L  *sdl.Rect = &sdl.Rect{64, 96, tSz, tSz}
	MAN_UL_R    *sdl.Rect = &sdl.Rect{96, 32, tSz, tSz}
	MAN_UL_N    *sdl.Rect = &sdl.Rect{128, 32, tSz, tSz}
	MAN_UL_L    *sdl.Rect = &sdl.Rect{160, 32, tSz, tSz}
	MAN_UR_R    *sdl.Rect = &sdl.Rect{96, 96, tSz, tSz}
	MAN_UR_N    *sdl.Rect = &sdl.Rect{128, 96, tSz, tSz}
	MAN_UR_L    *sdl.Rect = &sdl.Rect{160, 96, tSz, tSz}
	MAN_DL_R    *sdl.Rect = &sdl.Rect{96, 0, tSz, tSz}
	MAN_DL_N    *sdl.Rect = &sdl.Rect{128, 0, tSz, tSz}
	MAN_DL_L    *sdl.Rect = &sdl.Rect{160, 0, tSz, tSz}
	MAN_DR_R    *sdl.Rect = &sdl.Rect{96, 64, tSz, tSz}
	MAN_DR_N    *sdl.Rect = &sdl.Rect{128, 64, tSz, tSz}
	MAN_DR_L    *sdl.Rect = &sdl.Rect{160, 64, tSz, tSz}

	MAN_WALK_DL    [8]*sdl.Rect = [8]*sdl.Rect{MAN_DL_N, MAN_DL_L, MAN_DL_N, MAN_DL_R}
	MAN_WALK_DR    [8]*sdl.Rect = [8]*sdl.Rect{MAN_DR_N, MAN_DR_L, MAN_DR_N, MAN_DR_R}
	MAN_WALK_UL    [8]*sdl.Rect = [8]*sdl.Rect{MAN_UL_N, MAN_UL_L, MAN_UL_N, MAN_UL_R}
	MAN_WALK_UR    [8]*sdl.Rect = [8]*sdl.Rect{MAN_UR_N, MAN_UR_L, MAN_UR_N, MAN_UR_R}
	MAN_WALK_FRONT [8]*sdl.Rect = [8]*sdl.Rect{MAN_FRONT_N, MAN_FRONT_R, MAN_FRONT_N, MAN_FRONT_L}
	MAN_WALK_LEFT  [8]*sdl.Rect = [8]*sdl.Rect{MAN_LEFT_N, MAN_LEFT_R, MAN_LEFT_N, MAN_LEFT_L}
	MAN_WALK_RIGHT [8]*sdl.Rect = [8]*sdl.Rect{MAN_RIGHT_N, MAN_RIGHT_R, MAN_RIGHT_N, MAN_RIGHT_L}
	MAN_WALK_BACK  [8]*sdl.Rect = [8]*sdl.Rect{MAN_BACK_N, MAN_BACK_R, MAN_BACK_N, MAN_BACK_L}

	LAVA_S1 *sdl.Rect = &sdl.Rect{192, 0, tSz, tSz}
	LAVA_S2 *sdl.Rect = &sdl.Rect{224, 0, tSz, tSz}
	LAVA_S3 *sdl.Rect = &sdl.Rect{256, 0, tSz, tSz}

	LAVA_A [8]*sdl.Rect = [8]*sdl.Rect{LAVA_S1, LAVA_S2, LAVA_S3, LAVA_S3, LAVA_S2}

	YGLOW_S1 *sdl.Rect = &sdl.Rect{0, 0, tSz * 2, tSz * 2}
	YGLOW_S2 *sdl.Rect = &sdl.Rect{32, 0, tSz * 2, tSz * 2}
	YGLOW_S3 *sdl.Rect = &sdl.Rect{64, 0, tSz * 2, tSz * 2}

	YGLOW_A [8]*sdl.Rect = [8]*sdl.Rect{YGLOW_S1, YGLOW_S2, YGLOW_S3, YGLOW_S2}

	BGLOW_S1 *sdl.Rect = &sdl.Rect{224, 224, tSz, tSz}
	BGLOW_S2 *sdl.Rect = &sdl.Rect{256, 224, tSz, tSz}
	BGLOW_S3 *sdl.Rect = &sdl.Rect{288, 224, tSz, tSz}
	BGLOW_S4 *sdl.Rect = &sdl.Rect{320, 224, tSz, tSz}

	BGLOW_A [8]*sdl.Rect = [8]*sdl.Rect{BGLOW_S1, BGLOW_S2, BGLOW_S3, BGLOW_S4, BGLOW_S3, BGLOW_S2}

	LAVA_ANIM = &Animation{
		Action:   LAVA_A,
		Pose:     0,
		PoseTick: 16,
	}

	LIFE_ORB_ANIM = &Animation{
		Action:   YGLOW_A,
		Pose:     0,
		PoseTick: 16,
	}

	DEFAULT_FACING = Facing{
		Up:        MAN_WALK_BACK,
		Down:      MAN_WALK_FRONT,
		Left:      MAN_WALK_LEFT,
		Right:     MAN_WALK_RIGHT,
		DownLeft:  MAN_WALK_DL,
		DownRight: MAN_WALK_DR,
		UpLeft:    MAN_WALK_UL,
		UpRight:   MAN_WALK_UR,
	}

	LAVA_HANDLERS = &InteractionHandlers{
		OnCollDmg: 3,
	}

	ATK_UP *PowerUp = &PowerUp{
		Name:        "Attack up!",
		Description: "+20% dgm. dealed",
		Ico:         BF_ATK_UP,
	}
	DEF_UP *PowerUp = &PowerUp{
		Name:        "Defense up!",
		Description: "+20% dgm. reduction",
		Ico:         BF_DEF_UP,
	}

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

	SCENES []*Scene = []*Scene{
		SCN_PLAINS,
		SCN_CAVE,
	}

	Cam = Camera{
		DZx: 30,
		DZy: 60,
	}

	PC = Char{
		Solid: &Solid{
			Velocity: &Vector2d{0, 0},
			Facing:   DEFAULT_FACING,
			Anim: &Animation{
				Action:   MAN_WALK_FRONT,
				Pose:     0,
				PoseTick: 24,
			},
		},

		Lvl:       1,
		CurrentXP: 0,
		NextLvlXP: 100,

		Buffs:     []*PowerUp{ATK_UP, DEF_UP},
		Speed:     1,
		CurrentHP: 220,
		MaxHP:     250,
		CurrentST: 65,
		MaxST:     100,
	}

	scene       *Scene
	World       [][]*sdl.Rect
	Interactive []*Solid
	Monsters    []*Char
	GUI         []*sdl.Rect
	CullMap     []*Solid
)

func checkCol(r1 *sdl.Rect, r2 *sdl.Rect) bool {
	return (r1.X < (r2.X+r2.W) &&
		r1.X+r1.W > r2.X &&
		r1.Y < r2.Y+r2.H &&
		r1.Y+r1.H > r2.Y)
}

func (t *TextEl) Bake(renderer *sdl.Renderer) (*sdl.Texture, int32, int32) {

	surface, _ := t.Font.RenderUTF8_Blended_Wrapped(t.Content, t.Color, int(winWidth))
	defer surface.Free()
	txtr, _ := renderer.CreateTextureFromSurface(surface)

	return txtr, surface.W, surface.H
}

func actProc() {
	action_hit_box := ActHitBox(PC.Solid.Position, PC.Solid.Facing.Orientation)

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

func onColHdk(hdk *Solid, tgt *Solid) {
	if tgt.CharPtr != nil {
		tgt.CharPtr.depletHP(hdk.Handlers.OnCollDmg)
	}
	sol := &Solid{
		Position: &sdl.Rect{hdk.Position.X, hdk.Position.Y, tSz, tSz},
		Txt:      glowTxt,
		Anim: &Animation{
			Action:   YGLOW_A,
			Pose:     0,
			PoseTick: 16,
		},
		Ttl:       time.Now().Add(100 * time.Millisecond).Unix(),
		Collision: 0,
	}

	Interactive = append(Interactive, sol)

	hdk.Destroy()
}

func (c *Char) peformHaduken() {
	var stCost uint16 = 12

	if stCost > c.CurrentST {
		return
	}
	c.CurrentST -= stCost

	r := ActHitBox(c.Solid.Position, c.Solid.Facing.Orientation)
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
			OnCollDmg:   12,
			OnCollEvent: onColHdk,
		},
		CPattern: 0,
		MPattern: []Movement{
			Movement{c.Solid.Facing.Orientation, 255},
		},
		Collision: 1,
	}
	Interactive = append(Interactive, h)
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
func handleKeyEvent(key sdl.Keycode) Vector2d {
	N := Vector2d{0, 0}
	switch key {
	case KEY_SPACE_BAR:
		if !dbox.NextText() {
			actProc()
		}
		return N
	case KEY_LEFT_SHIT:
		PC.Speed = 2
	case KEY_ARROW_UP:
		return F_UP
	case KEY_ARROW_DOWN:
		return F_DOWN
	case KEY_ARROW_LEFT:
		return F_LEFT
	case KEY_ARROW_RIGHT:
		return F_RIGHT
	}

	return N
}

func handleKeyUpEvent(key sdl.Keycode) {
	switch key {
	case KEY_C:
		PC.peformHaduken()
	case KEY_LEFT_SHIT:
		PC.Speed = 1
	case KEY_ARROW_UP:
		PC.Solid.Velocity.Y = 0
	case KEY_ARROW_DOWN:
		PC.Solid.Velocity.Y = 0
	case KEY_ARROW_LEFT:
		PC.Solid.Velocity.X = 0
	case KEY_ARROW_RIGHT:
		PC.Solid.Velocity.X = 0
	}
}

func (s *Solid) procMovement(speed int32) {
	np := &sdl.Rect{
		(s.Position.X + (s.Velocity.X * speed)),
		(s.Position.Y + (s.Velocity.Y * speed)),
		s.Position.W,
		s.Position.H,
	}
	var outbound bool = (np.X <= 0 ||
		np.Y <= 0 ||
		np.X > int32(scene.CellsX*tSz) ||
		np.Y > int32(scene.CellsY*tSz))

	if (np.X == s.Position.X && np.Y == s.Position.Y) || outbound {
		return
	}

	for _, obj := range CullMap {
		if obj == s || obj.Position == nil {
			continue
		}
		fr := feetRect(np)
		if checkCol(fr, obj.Position) && obj.Collision == 1 {
			if obj.Handlers != nil && obj.Handlers.OnCollEvent != nil {
				obj.Handlers.OnCollEvent(obj, s)
			}
			return
		}
	}

	s.Facing.Orientation = *s.Velocity
	_, act := GetFacing(&s.Facing, *s.Velocity)
	s.Anim.Action = act

	s.Position = np
}

func catchEvents() bool {
	var c bool
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyDownEvent:
			v := handleKeyEvent(t.Keysym.Sym)
			if v.X != 0 {
				PC.Solid.Velocity.X = v.X
				c = true
			}
			if v.Y != 0 {
				PC.Solid.Velocity.Y = v.Y
				c = true
			}
		case *sdl.KeyUpEvent:
			handleKeyUpEvent(t.Keysym.Sym)
		}
	}

	if c {
		PC.Solid.procMovement(PC.Speed)
		PC.Solid.Anim.PoseTick -= 1
		if PC.Solid.Anim.PoseTick == 0 {
			PC.Solid.Anim.Pose = getNextPose(PC.Solid.Anim.Action, PC.Solid.Anim.Pose)
			PC.Solid.Anim.PoseTick = 8
		}

		newScreenPos := worldToScreen(PC.Solid.Position, Cam)
		if (Cam.DZx - newScreenPos.X) > 0 {
			Cam.P.X -= (Cam.DZx - newScreenPos.X)
		}

		if (winWidth - Cam.DZx) < (newScreenPos.X + tSz) {
			Cam.P.X += (newScreenPos.X + tSz) - (winWidth - Cam.DZx)
		}

		if (Cam.DZy - newScreenPos.Y) > 0 {
			Cam.P.Y -= (Cam.DZy - newScreenPos.Y)
		}

		if (winHeight - Cam.DZy) < (newScreenPos.Y + tSz) {
			Cam.P.Y += (newScreenPos.Y + tSz) - (winHeight - Cam.DZy)
		}

	}
	if time.Now().Nanosecond()%2 == 1 {
		if isMoving(PC.Solid.Velocity) && PC.Speed > 1 {
			if PC.CurrentST <= uint16(PC.Speed) {
				PC.CurrentST = 0
				PC.Speed = 1
			} else {
				PC.CurrentST -= 1
			}
		} else {
			if !isMoving(PC.Solid.Velocity) && PC.CurrentST < PC.MaxST {
				PC.CurrentST += 1
			}
		}
	}

	return true
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

func (s *Scene) build() {
	ni := int(s.CellsX) + 1
	nj := int(s.CellsY) + 1

	println("start build: ", s.codename)

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
	dt := SCN_PLAINS
	if s.codename == "plains" {
		dt = SCN_CAVE
	}
	//population = 0
	for i := 0; i < population; i++ {

		cX := rand.Int31n(s.CellsX)
		cY := rand.Int31n(s.CellsY)

		absolute_pos := &sdl.Rect{cX * tSz, cY * tSz, tSz, tSz}
		obj_type := rand.Int31n(10)
		sol := &Solid{}

		switch obj_type {
		case 1:
			sol = &Solid{
				Position: absolute_pos,
				Txt:      s.TileSet,
				Anim: &Animation{
					Action:   LAVA_A,
					Pose:     0,
					PoseTick: 32,
				},
				Handlers:  LAVA_HANDLERS,
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
					DoorTo:     dt,
				},
			}
			Interactive = append(Interactive, sol)
			break
		case 3:
			absolute_pos.H = 64
			absolute_pos.W = 64
			sol = &Solid{
				Position:  absolute_pos,
				Txt:       glowTxt,
				Anim:      LIFE_ORB_ANIM,
				Collision: 0,
				Handlers: &InteractionHandlers{
					OnCollEvent: pickUp,
					OnPickUp: func(healed *Solid, healer *Solid) {
						if healed.CharPtr != nil {
							healed.CharPtr.CurrentHP += 10
						}
					},
				},
			}
			Interactive = append(Interactive, sol)
			break
		case 4:
			mon := Char{
				Lvl: 10,
				Solid: &Solid{
					Position:  absolute_pos,
					Velocity:  &Vector2d{0, 0},
					Txt:       slimeTxt,
					Collision: 2,
					Facing:    DEFAULT_FACING,
					Anim: &Animation{
						Action:   MAN_WALK_FRONT,
						Pose:     0,
						PoseTick: 16,
					},
					Handlers: &InteractionHandlers{
						OnCollDmg: 12,
					},
					CPattern: 0,
					MPattern: []Movement{
						Movement{F_DOWN, 50},
						Movement{F_UP, 90},
						Movement{F_RIGHT, 10},
						Movement{F_LEFT, 10},
					},
					Chase: PC.Solid,
				},
				Speed:     1,
				CurrentHP: 50,
				MaxHP:     50,
				CurrentST: 50,
				MaxST:     50,
			}
			mon.Solid.CharPtr = &mon
			Monsters = append(Monsters, &mon)
			break
		case 5:
			fnpc := Char{
				Solid: &Solid{
					Position:  absolute_pos,
					Velocity:  &Vector2d{0, 0},
					Txt:       spritesheetTxt,
					Collision: 1,
					Facing:    DEFAULT_FACING,
					Anim: &Animation{
						Action:   MAN_WALK_FRONT,
						Pose:     0,
						PoseTick: 16,
					},
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
				Speed:     1,
				CurrentHP: 9999,
				MaxHP:     9999,
			}
			fnpc.Solid.CharPtr = &fnpc
			Monsters = append(Monsters, &fnpc)
			break
		}

	}
}

var (
	EventTick uint8 = 16
	AiTick          = 16
	dbox      DBox  = DBox{BGColor: sdl.Color{90, 90, 90, 255}}
)

func (s *Scene) update() {
	EventTick -= 1
	AiTick -= 1

	now := time.Now().Unix()
	for _, cObj := range CullMap {
		// Ttl kill
		if cObj.Ttl > 0 && cObj.Ttl < now {
			cObj.Destroy()
		}

		// Kill logic
		if cObj.CharPtr != nil && cObj.CharPtr.CurrentHP <= 0 {
			drop := &Solid{
				Position: cObj.Position,
				Txt:      powerupsTxt,
				Source:   BF_ATK_UP,
				Handlers: &InteractionHandlers{
					OnActEvent: pickUp,
					OnPickUp: func(picker *Solid, picked *Solid) {
						if picker.CharPtr != nil {
							picker.CharPtr.CurrentST = picker.CharPtr.MaxST
						}
					},
				},
			}
			drop.Position.W = 24
			drop.Position.H = 24
			Interactive = append(Interactive, drop)

			cObj.Destroy()

			PC.CurrentXP += uint16(cObj.CharPtr.MaxHP / 10)
			if PC.CurrentXP >= PC.NextLvlXP {
				PC.CurrentXP = 0
				PC.Lvl++
				PC.NextLvlXP = PC.NextLvlXP * uint16(1+PC.Lvl/2)
			}
		}

		if cObj.Anim != nil {
			// update Interactives poses
			animObj := cObj.Anim
			animObj.PoseTick -= 1
			if animObj.PoseTick == 0 {
				animObj.Pose = getNextPose(animObj.Action, animObj.Pose)
				animObj.PoseTick = 16
			}
		}

		if len(dbox.Text) > 0 {
			AiTick = 16
			return
		}

		if cObj.Handlers != nil && EventTick == 0 {
			fr := feetRect(PC.Solid.Position)
			if checkCol(fr, cObj.Position) {
				if cObj.Handlers.OnCollDmg != 0 {
					PC.depletHP(cObj.Handlers.OnCollDmg)
				}
				if cObj.Handlers.OnCollEvent != nil {
					cObj.Handlers.OnCollEvent(PC.Solid, cObj)
				}
			}
		}

		if AiTick == 0 {
			if cObj.Chase != nil && cObj.LoSCheck(32) {
				cObj.chase()
			} else {
				cObj.peformPattern(1)
			}
		}
	}
	if EventTick == 0 {
		EventTick = 16
	}
	if AiTick == 0 {
		AiTick = 3
	}
}

func (s *Solid) LoSCheck(int32) bool {
	LoS := &sdl.Rect{s.Position.X - 128, s.Position.Y - 128, s.Position.W + 256, s.Position.H + 256}
	if !checkCol(PC.Solid.Position, LoS) {
		return false
	}
	return true
}

func (s *Solid) chase() {
	s.Velocity.X = 0
	s.Velocity.Y = 0

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

	s.procMovement(s.CharPtr.Speed)

	return
}

func (s *Solid) peformPattern(sp int32) {
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
		applyMov(s.Position, mov.Orientation, sp)
		s.CPattern += uint32(sp)

		if s.Facing.Up[0] != nil {
			_, na := GetFacing(&s.Facing, mov.Orientation)
			if na != s.Anim.Action {
				s.Anim.Action = na
			}
		}
	}
}

func pickUp(picker *Solid, item *Solid) {
	if item.Handlers != nil {
		item.Handlers.OnPickUp(picker, item)
	}
	item.Destroy()
}

func PlayDialog(listener *Solid, speaker *Solid) {
	if len(speaker.Handlers.DialogScript) > 0 {
		dbox.LoadText(speaker.Handlers.DialogScript)
	}
}

func GetFacing(f *Facing, o Vector2d) (Vector2d, [8]*sdl.Rect) {
	if o.X == 0 && o.Y == -1 {
		return F_UP, f.Up
	}

	if o.X == -1 && o.Y == 0 {
		return F_LEFT, f.Left
	}

	if o.X == 1 && o.Y == 0 {
		return F_RIGHT, f.Right
	}

	if o.X == 1 && o.Y == 1 {
		return F_DR, f.DownRight
	}

	if o.X == 1 && o.Y == -1 {
		return F_UR, f.UpRight
	}

	if o.X == -1 && o.Y == 1 {
		return F_DL, f.DownLeft
	}

	if o.X == -1 && o.Y == -1 {
		return F_UL, f.UpLeft
	}

	return F_DOWN, f.Down
}

func applyMov(p *sdl.Rect, o Vector2d, s int32) {
	p.X += o.X * s
	p.Y += o.Y * s
}

func (s *Scene) _terrainRender(renderer *sdl.Renderer) {
	var Source *sdl.Rect
	var init int32 = 0

	if Cam.P.X < 0 {
		Cam.P.X = 0
	}

	if Cam.P.Y < 0 {
		Cam.P.Y = 0
	}

	var offsetX, offsetY int32 = tSz, tSz
	// Rendering the terrain
	for winY := init; winY < winHeight; winY += offsetY {
		for winX := init; winX < winWidth; winX += offsetX {

			offsetX = (tSz - ((Cam.P.X + winX) % tSz))
			offsetY = (tSz - ((Cam.P.Y + winY) % tSz))

			worldCellX := uint16((Cam.P.X + winX) / tSz)
			worldCellY := uint16((Cam.P.Y + winY) / tSz)
			screenPos := sdl.Rect{winX, winY, offsetX, offsetY}

			if worldCellX > uint16(s.CellsX) || worldCellY > uint16(s.CellsY) || worldCellX < 0 || worldCellY < 0 {
				continue
			}

			gfx := World[worldCellX][worldCellY]

			if offsetX != int32(tSz) || offsetY != int32(tSz) {
				Source = &sdl.Rect{gfx.X + (tSz - offsetX), gfx.Y + (tSz - offsetY), offsetX, offsetY}
			} else {
				Source = gfx
			}
			if Source != nil && &screenPos != nil {
			}
			renderer.Copy(s.TileSet, Source, &screenPos)
		}
	}
}

func (s *Solid) Destroy() {
	s.Position = nil
	s.Source = nil
	s.Facing = Facing{}
	s.Anim = nil
	s.Handlers = nil
	s.Txt = nil
	s.Collision = 0
}

func (s *Scene) _solidsRender(renderer *sdl.Renderer) {
	CullMap = []*Solid{}

	for _, obj := range Interactive {
		if obj.Position == nil {
			continue
		}

		scrPos := worldToScreen(obj.Position, Cam)

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
			renderer.DrawRect(scrPos)
			CullMap = append(CullMap, mon.Solid)
		}
	}
}

func (s *Scene) _GUIRender(renderer *sdl.Renderer) {

	// Gray overlay
	renderer.SetDrawColor(60, 60, 60, 10)
	renderer.FillRect(&sdl.Rect{0, 0, 120, 60})

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
	renderer.FillRect(&sdl.Rect{10, 38, 100, 4})
	renderer.SetDrawColor(190, 190, 0, 255)
	renderer.FillRect(&sdl.Rect{10, 38, int32(calcPerc(PC.CurrentXP, PC.NextLvlXP)), 4})

	for i, b := range PC.Buffs {
		pos := sdl.Rect{8 + (int32(i) * 32), 48, 24, 24}
		renderer.Copy(powerupsTxt, b.Ico, &pos)
	}
	for _, el := range GUI {
		scrPos := worldToScreen(el, Cam)
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(scrPos)
	}

	lvl_TextEl := TextEl{
		Font:    font,
		Content: fmt.Sprintf("Lvl. %d", PC.Lvl),
		Color:   sdl.Color{255, 255, 255, 255},
	}
	lvl_txtr, W, H := lvl_TextEl.Bake(renderer)
	renderer.Copy(lvl_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{128, 60, W, H})

	dbg_content := fmt.Sprintf(
		"px %d py %d|vx %d vy %d cull %d i %d cX %d cY %d L %dus ETick%d AiTick%d",
		PC.Solid.Position.X,
		PC.Solid.Position.Y,
		PC.Solid.Velocity.X,
		PC.Solid.Velocity.Y,
		len(CullMap),
		len(Interactive),
		Cam.P.X,
		Cam.P.Y,
		game_latency,
		EventTick,
		AiTick,
	)

	dbox.Present(renderer)

	dbg_TextEl := TextEl{
		Font:    font,
		Content: dbg_content,
		Color:   sdl.Color{255, 255, 255, 255},
	}
	dbg_txtr, W, H := dbg_TextEl.Bake(renderer)
	renderer.Copy(dbg_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{0, winHeight - H, W, H})
}

func (s *Scene) render(renderer *sdl.Renderer) {
	renderer.Clear()

	s._terrainRender(renderer)
	s._solidsRender(renderer)
	s._monstersRender(renderer)
	// Rendering the PC
	scrPos := worldToScreen(PC.Solid.Position, Cam)
	renderer.Copy(spritesheetTxt, PC.Solid.Anim.Action[PC.Solid.Anim.Pose], scrPos)
	s._GUIRender(renderer)

	// FLUSH FRAME
	renderer.Present()
}

func calcPerc(v1 uint16, v2 uint16) float32 {
	return (float32(v1) / float32(v2) * 100)
}

func worldToScreen(pos *sdl.Rect, cam Camera) *sdl.Rect {
	return &sdl.Rect{(pos.X - cam.P.X), (pos.Y - cam.P.Y), pos.W, pos.H}
}

func inScreen(r *sdl.Rect) bool {
	return (r.X > (r.W*-1) && r.X < winWidth && r.Y > (r.H*-1) && r.Y < winHeight)
}

func (ch *Char) depletHP(dmg uint16) {
	if dmg > ch.CurrentHP {
		ch.CurrentHP = 0
	} else {
		ch.CurrentHP -= dmg
	}
	ch.PushBack(8)
}

func (c *Char) PushBack(d int32) {
	f := c.Solid.Facing.Orientation
	c.Solid.Position.X -= f.X * d
	c.Solid.Position.Y -= f.Y * d
}

func BashDoor(actor *Solid, door *Solid) {
	change_scene(door.Handlers.DoorTo, nil)
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

	tilesetImg, _ := img.Load("assets/textures/ts1.bmp")
	spritesheetImg, _ := img.Load("assets/textures/main_char.png")
	powerupsImg, _ := img.Load("assets/textures/powerups_ts.png")
	glowImg, _ := img.Load("assets/textures/glowing_ts.png")
	slimeImg, _ := img.Load("assets/textures/slimes.png")
	defer slimeImg.Free()
	defer tilesetImg.Free()
	defer spritesheetImg.Free()
	defer powerupsImg.Free()
	defer glowImg.Free()

	tilesetTxt, _ = renderer.CreateTextureFromSurface(tilesetImg)
	spritesheetTxt, _ = renderer.CreateTextureFromSurface(spritesheetImg)
	powerupsTxt, _ = renderer.CreateTextureFromSurface(powerupsImg)
	glowTxt, _ = renderer.CreateTextureFromSurface(glowImg)
	slimeTxt, _ = renderer.CreateTextureFromSurface(slimeImg)
	defer tilesetTxt.Destroy()
	defer spritesheetTxt.Destroy()
	defer powerupsTxt.Destroy()
	defer glowTxt.Destroy()
	defer slimeTxt.Destroy()

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

		game_latency = (time.Since(then) / time.Microsecond)

		running = catchEvents()
		sdl.Delay(22)
	}
}

func V2R(v Vector2d, w int32, h int32) *sdl.Rect {
	return &sdl.Rect{v.X, v.Y, w, h}
}

func change_scene(new_scene *Scene, staring_pos *Vector2d) {

	new_scene.build()
	Interactive = []*Solid{}
	new_scene.populate(200)
	scene = new_scene

	PC.Solid.Position = V2R(new_scene.StartPoint, tSz, tSz)
	Cam.P = new_scene.CamPoint
}

func feetRect(pos *sdl.Rect) *sdl.Rect {
	return &sdl.Rect{pos.X, pos.Y + 16, pos.W, pos.H - 16}
}
