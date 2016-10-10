package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
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
	KEY_SPACE_BAR             = 1073741824 // 32
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

type Event func(obj *Solid)

type InteractionHandlers struct {
	OnCollDmg   uint16
	OnCollPush  *Vector2d
	OnCollEvent Event
	OnPickUp    Event
	OnActDmg    uint16
	OnActPush   *Vector2d
	OnActEvent  Event
	DoorTo      *Scene
}

type Solid struct {
	Velocity  *Vector2d
	Position  *sdl.Rect
	Source    *sdl.Rect
	Facing    Facing
	Anim      *Animation
	Handlers  *InteractionHandlers
	Txt       *sdl.Texture
	Collision uint8
	// AI RELATED
	CPattern uint32
	MPattern []Movement
	Chase    *Solid
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
	Solid *Solid
	Buffs []*PowerUp

	Speed     int32
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

func ActHitBox(source *sdl.Rect, facing Vector2d) *sdl.Rect {
	println(">", facing.X, facing.Y)

	return &sdl.Rect{
		source.X + (facing.X * tSz),
		source.Y + (facing.Y * tSz),
		source.W,
		source.H,
	}
}

var (
	winTitle string = "Go-SDL2 Obj0"
	event    sdl.Event

	tilesetTxt     *sdl.Texture = nil
	spritesheetTxt *sdl.Texture = nil
	particlesTxt   *sdl.Texture = nil
	powerupsTxt    *sdl.Texture = nil
	glowTxt        *sdl.Texture = nil
	slimeTxt       *sdl.Texture = nil

	GRASS     *sdl.Rect = &sdl.Rect{0, 0, tSz, tSz}
	TREE                = &sdl.Rect{0, 32, tSz, tSz}
	DIRT                = &sdl.Rect{703, 0, tSz, tSz}
	WALL                = &sdl.Rect{0, 64, tSz, tSz}
	DOOR                = &sdl.Rect{256, 32, tSz, tSz}
	WOMAN               = &sdl.Rect{0, 128, tSz, tSz}
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

	MAN_UL_R *sdl.Rect = &sdl.Rect{96, 32, tSz, tSz}
	MAN_UL_N *sdl.Rect = &sdl.Rect{128, 32, tSz, tSz}
	MAN_UL_L *sdl.Rect = &sdl.Rect{160, 32, tSz, tSz}
	MAN_UR_R *sdl.Rect = &sdl.Rect{96, 96, tSz, tSz}
	MAN_UR_N *sdl.Rect = &sdl.Rect{128, 96, tSz, tSz}
	MAN_UR_L *sdl.Rect = &sdl.Rect{160, 96, tSz, tSz}
	MAN_DL_R *sdl.Rect = &sdl.Rect{96, 0, tSz, tSz}
	MAN_DL_N *sdl.Rect = &sdl.Rect{128, 0, tSz, tSz}
	MAN_DL_L *sdl.Rect = &sdl.Rect{160, 0, tSz, tSz}
	MAN_DR_R *sdl.Rect = &sdl.Rect{96, 64, tSz, tSz}
	MAN_DR_N *sdl.Rect = &sdl.Rect{128, 64, tSz, tSz}
	MAN_DR_L *sdl.Rect = &sdl.Rect{160, 64, tSz, tSz}

	MAN_WALK_DL [8]*sdl.Rect = [8]*sdl.Rect{MAN_DL_N, MAN_DL_L, MAN_DL_N, MAN_DL_R}
	MAN_WALK_DR [8]*sdl.Rect = [8]*sdl.Rect{MAN_DR_N, MAN_DR_L, MAN_DR_N, MAN_DR_R}
	MAN_WALK_UL [8]*sdl.Rect = [8]*sdl.Rect{MAN_UL_N, MAN_UL_L, MAN_UL_N, MAN_UL_R}
	MAN_WALK_UR [8]*sdl.Rect = [8]*sdl.Rect{MAN_UR_N, MAN_UR_L, MAN_UR_N, MAN_UR_R}

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

func actProc() {
	action_hit_box := ActHitBox(PC.Solid.Position, PC.Solid.Facing.Orientation)

	// Debug hint
	GUI = append(GUI, action_hit_box)

	for _, obj := range CullMap {
		if obj.Handlers != nil &&
			obj.Handlers.OnActEvent != nil &&
			checkCol(action_hit_box, obj.Position) {
			obj.Handlers.OnActEvent(obj)
			return
		}
	}
}

func onColHdk(o *Solid) {
	o.Destroy()
}

func (c *Char) peformHaduken() {
	r := ActHitBox(c.Solid.Position, c.Solid.Facing.Orientation)
	println(r.X, r.Y, "X", c.Solid.Position.X, c.Solid.Position.Y)
	h := &Solid{
		Position: r,
		Txt:      glowTxt,
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
	}
	Interactive = append(Interactive, h)
}

func handleKeyEvent(key sdl.Keycode) Vector2d {
	N := Vector2d{0, 0}
	switch key {
	case KEY_C:
		PC.peformHaduken()
		return N
	case KEY_SPACE_BAR:
		actProc()
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
		fr := feetRect(np)
		if checkCol(fr, obj.Position) && obj.Collision == 1 {
			return
		}
	}

	newScreenPos := worldToScreen(np, Cam)
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
	}

	return true
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
				},
			}
			Interactive = append(Interactive, sol)
			break
		case 4:
			mon := Char{
				Solid: &Solid{
					Position:  absolute_pos,
					Velocity:  &Vector2d{0, 0},
					Txt:       slimeTxt,
					Collision: 1,
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
				CurrentHP: 220,
				MaxHP:     250,
				CurrentST: 65,
				MaxST:     100,
			}
			Monsters = append(Monsters, &mon)
			break
		}
	}
}

var (
	EventTick uint8 = 16
	AiTick          = 16
	blinkTick       = 32
)

func (s *Scene) update() {
	EventTick -= 1
	AiTick -= 1

	for _, cObj := range CullMap {

		if cObj.Anim != nil {
			// update Interactives poses
			animObj := cObj.Anim
			animObj.PoseTick -= 1
			if animObj.PoseTick == 0 {
				animObj.Pose = getNextPose(animObj.Action, animObj.Pose)
				animObj.PoseTick = 16
			}
		}

		if cObj.Handlers != nil && EventTick == 0 {
			fr := feetRect(PC.Solid.Position)
			if checkCol(fr, cObj.Position) {
				if cObj.Handlers.OnCollDmg != 0 {
					depletHP(cObj.Handlers.OnCollDmg)
				}
				if cObj.Handlers.OnCollEvent != nil {
					cObj.Handlers.OnCollEvent(cObj)
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
	// TODO - LINE OF SIGHT CHECK
	new_pos := &sdl.Rect{s.Position.X, s.Position.Y, s.Position.W, s.Position.H}

	if s.Chase.Position.X > s.Position.X {
		new_pos.X += 1
	}

	if s.Chase.Position.X < s.Position.X {
		new_pos.X -= 1
	}

	if s.Chase.Position.Y < s.Position.Y {
		new_pos.Y -= 1
	}

	if s.Chase.Position.Y > s.Position.Y {
		new_pos.Y += 1
	}

	f := feetRect(new_pos)
	for _, obj := range CullMap {
		if obj.Position != nil && s != obj && checkCol(f, obj.Position) {
			if obj.Handlers.OnCollEvent != nil {
				obj.Handlers.OnCollEvent(obj)
			}
			if obj.Collision == 1 {
				return
			}
		}
	}
	s.Position = new_pos

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

func pickUp(s *Solid) {
	if s.Handlers.OnPickUp != nil {
		s.Handlers.OnPickUp(s)
	}
	s.Destroy()
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
			// renderer.SetDrawColor(0, 0, 255, 255)
			// renderer.DrawRect(&screenPos)
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
	s.Collision = 4
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

			//renderer.SetDrawColor(0, 255, 0, 255)
			//renderer.DrawRect(scrPos)
			CullMap = append(CullMap, obj)
		}
	}
}

func (s *Scene) _monstersRender(renderer *sdl.Renderer) {

	renderer.SetDrawColor(0, 255, 0, 255)

	for _, mon := range Monsters {
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

	// MANA BAR BG
	renderer.SetDrawColor(190, 0, 120, 255)
	renderer.FillRect(&sdl.Rect{10, 24, 100, 4})
	renderer.SetDrawColor(0, 0, 255, 255)
	renderer.FillRect(&sdl.Rect{10, 24, int32(calcPerc(PC.CurrentST, PC.MaxST)), 4})

	for i, b := range PC.Buffs {
		pos := sdl.Rect{8 + (int32(i) * 32), 32, 24, 24}
		renderer.Copy(powerupsTxt, b.Ico, &pos)
	}
	for _, el := range GUI {
		scrPos := worldToScreen(el, Cam)
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(scrPos)
	}
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
	return &sdl.Rect{
		pos.X - cam.P.X,
		pos.Y - cam.P.Y,
		pos.W,
		pos.H,
	}
}

func inScreen(r *sdl.Rect) bool {
	return (r.X > (r.W*-1) && r.X < winWidth &&
		r.Y > (r.H*-1) && r.Y < winHeight)
}

func depletHP(dmg uint16) {
	if dmg > PC.CurrentHP {
		PC.CurrentHP = 0
	} else {
		PC.CurrentHP -= dmg
		PC.PushBack(24)
	}
}

func (c *Char) PushBack(d int32) {
	f := c.Solid.Facing.Orientation
	c.Solid.Position.X -= f.X * d
	c.Solid.Position.Y -= f.Y * d
}

func BashDoor(obj *Solid) {
	change_scene(obj.Handlers.DoorTo, nil)
}

func main() {

	runtime.GOMAXPROCS(1)

	var window *sdl.Window
	var renderer *sdl.Renderer

	window, _ = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(winWidth), int(winHeight), sdl.WINDOW_SHOWN)
	defer window.Destroy()

	renderer, _ = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	tilesetImg, _ := img.Load("assets/textures/ts1.bmp")
	defer tilesetImg.Free()

	tilesetTxt, _ = renderer.CreateTextureFromSurface(tilesetImg)
	defer tilesetTxt.Destroy()

	spritesheetImg, _ := img.Load("assets/textures/main_char.png")
	defer spritesheetImg.Free()

	spritesheetTxt, _ = renderer.CreateTextureFromSurface(spritesheetImg)
	defer spritesheetTxt.Destroy()

	powerupsImg, _ := img.Load("assets/textures/powerups_ts.png")
	defer powerupsImg.Free()

	powerupsTxt, _ = renderer.CreateTextureFromSurface(powerupsImg)
	defer powerupsTxt.Destroy()

	glowImg, _ := img.Load("assets/textures/glowing_ts.png")
	defer glowImg.Free()

	glowTxt, _ = renderer.CreateTextureFromSurface(glowImg)
	defer glowTxt.Destroy()

	slimeImg, _ := img.Load("assets/textures/slimes.png")
	defer slimeImg.Free()

	slimeTxt, _ = renderer.CreateTextureFromSurface(slimeImg)
	defer slimeTxt.Destroy()

	var running bool = true

	for _, scn := range SCENES {
		scn.TileSet = tilesetTxt
	}

	renderer.SetDrawColor(0, 0, 255, 255)
	scene = SCENES[0]
	change_scene(scene, nil)

	for running {
		then := time.Now()

		scene.update()
		scene.render(renderer)

		println((time.Since(then)) / time.Microsecond)
		running = catchEvents()
		sdl.Delay(16)
	}
}

func V2R(v Vector2d, w int32, h int32) *sdl.Rect {
	return &sdl.Rect{v.X, v.Y, w, h}
}

func R2Vo(r *sdl.Rect) Vector2d {
	return Vector2d{r.X, r.Y}
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
