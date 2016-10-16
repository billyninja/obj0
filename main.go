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
	TSz                 float32 = 32
	TSzi                int32   = int32(TSz)
	winWidth, winHeight int32   = 640, 480
	cY                          = winHeight / TSzi
	cX                          = winWidth / TSzi
	WORLD_CELLS_X               = 500
	WORLD_CELLS_Y               = 200
	KEY_ARROW_UP                = 1073741906
	KEY_ARROW_DOWN              = 1073741905
	KEY_ARROW_LEFT              = 1073741904
	KEY_ARROW_RIGHT             = 1073741903
	KEY_LEFT_SHIT               = 1073742049
	KEY_SPACE_BAR               = 1073741824 // 32
	KEY_C                       = 99
	AI_TICK_LENGTH              = 2
	EVENT_TICK_LENGTH           = 3
)

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
	monstersTxt    *sdl.Texture = nil

	GRASS     *sdl.Rect = &sdl.Rect{0, 0, TSzi, TSzi}
	DIRT                = &sdl.Rect{703, 0, TSzi, TSzi}
	WALL                = &sdl.Rect{0, 64, TSzi, TSzi}
	DOOR                = &sdl.Rect{256, 32, TSzi, TSzi}
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
	MAN_PB_S1 *sdl.Rect = &sdl.Rect{96, 160, TSzi, TSzi}
	MAN_PB_S2 *sdl.Rect = &sdl.Rect{128, 160, TSzi, TSzi}
	MAN_PB_S3 *sdl.Rect = &sdl.Rect{160, 160, TSzi, TSzi}

	MAN_CS_S1     *sdl.Rect    = &sdl.Rect{192, 224, TSzi, TSzi}
	MAN_CS_S2     *sdl.Rect    = &sdl.Rect{224, 224, TSzi, TSzi}
	MAN_CS_S3     *sdl.Rect    = &sdl.Rect{256, 224, TSzi, TSzi}
	MAN_CS_S4     *sdl.Rect    = &sdl.Rect{224, 192, TSzi, TSzi}
	MAN_PUSH_BACK [8]*sdl.Rect = [8]*sdl.Rect{MAN_PB_S1, MAN_PB_S2}
	MAN_CAST      [8]*sdl.Rect = [8]*sdl.Rect{MAN_CS_S1, MAN_CS_S2, MAN_CS_S3, MAN_CS_S4}

	LAVA_S1 *sdl.Rect = &sdl.Rect{192, 0, TSzi, TSzi}
	LAVA_S2 *sdl.Rect = &sdl.Rect{224, 0, TSzi, TSzi}
	LAVA_S3 *sdl.Rect = &sdl.Rect{256, 0, TSzi, TSzi}

	LAVA_A [8]*sdl.Rect = [8]*sdl.Rect{LAVA_S1, LAVA_S2, LAVA_S3, LAVA_S3, LAVA_S2}

	YGLOW_S1 *sdl.Rect = &sdl.Rect{0, 0, TSzi * 2, TSzi * 2}
	YGLOW_S2 *sdl.Rect = &sdl.Rect{32, 0, TSzi * 2, TSzi * 2}
	YGLOW_S3 *sdl.Rect = &sdl.Rect{64, 0, TSzi * 2, TSzi * 2}

	YGLOW_A [8]*sdl.Rect = [8]*sdl.Rect{YGLOW_S1, YGLOW_S2, YGLOW_S3, YGLOW_S2}

	BGLOW_S1 *sdl.Rect = &sdl.Rect{224, 224, TSzi, TSzi}
	BGLOW_S2 *sdl.Rect = &sdl.Rect{256, 224, TSzi, TSzi}
	BGLOW_S3 *sdl.Rect = &sdl.Rect{288, 224, TSzi, TSzi}
	BGLOW_S4 *sdl.Rect = &sdl.Rect{320, 224, TSzi, TSzi}

	BGLOW_A [8]*sdl.Rect = [8]*sdl.Rect{BGLOW_S1, BGLOW_S2, BGLOW_S3, BGLOW_S4, BGLOW_S3, BGLOW_S2}

	LAVA_ANIM     = &Animation{Action: LAVA_A, PoseTick: 8}
	LIFE_ORB_ANIM = &Animation{Action: YGLOW_A, PoseTick: 8}
	MAN_PB_ANIM   = &Animation{Action: MAN_PUSH_BACK, PoseTick: 18, PlayMode: 1}
	MAN_CS_ANIM   = &Animation{Action: MAN_CAST, PoseTick: 18, PlayMode: 1}

	BatTPL MonsterTemplate = MonsterTemplate{}
	OrcTPL MonsterTemplate = MonsterTemplate{}

	LAVA_HANDLERS = &InteractionHandlers{
		OnCollDmg: 12,
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
			Velocity:    &Vector2d{0, 0},
			Orientation: &Vector2d{0, -1},
			Anim:        MAN_CS_ANIM,
		},
		Lvl:       1,
		CurrentXP: 0,
		NextLvlXP: 100,
		Buffs:     []*PowerUp{ATK_UP, DEF_UP},
		BaseSpeed: 1.5,
		Speed:     1.5,
		CurrentHP: 220,
		MaxHP:     250,
		CurrentST: 250,
		MaxST:     300,
		Inventory: []*ItemStack{{GreenBlob, 2}},
	}

	scene       *Scene
	World       [][]*sdl.Rect
	Interactive []*Solid
	Spawners    []*SpawnPoint
	Monsters    []*Char
	GUI         []*sdl.Rect
	CullMap     []*Solid
)

type Vector2d struct {
	X float32
	Y float32
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

type MonsterTemplate struct {
	Txtr        *sdl.Texture
	ActionMap   *ActionMap
	SpriteSheet *SpriteSheet

	Lvl           uint8
	HP            float32
	Size          int32
	LvlVariance   float32
	ScalingFactor float32
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
			println(
				o.X,
				o.Y,
				anim.Action[p].X,
				anim.Action[p].Y,
				anim.Action[p].W,
				anim.Action[p].H)
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

var (
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
)

type InteractionHandlers struct {
	OnCollDmg   float32
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
	Velocity    *Vector2d
	Orientation *Vector2d
	Position    *sdl.Rect
	Source      *sdl.Rect
	Anim        *Animation
	Handlers    *InteractionHandlers
	Txt         *sdl.Texture
	Ttl         int64
	Collision   uint8

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
	Ttl         int64
}

type Char struct {
	Solid     *Solid
	Buffs     []*PowerUp
	ActionMap *ActionMap
	Inventory []*ItemStack
	//---
	BaseSpeed float32
	Speed     float32
	Lvl       uint8
	CurrentXP uint16
	NextLvlXP uint16
	CurrentHP float32
	MaxHP     float32
	CurrentST float32
	MaxST     float32
	Drop      *Item
	Invinc    int64
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

func ActHitBox(source *sdl.Rect, orientation *Vector2d) *sdl.Rect {
	return &sdl.Rect{
		source.X + int32(orientation.X*TSz),
		source.Y + int32(orientation.Y*TSz),
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

func (t *TextEl) Bake(renderer *sdl.Renderer) (*sdl.Texture, int32, int32) {

	surface, _ := t.Font.RenderUTF8_Blended_Wrapped(t.Content, t.Color, int(winWidth))
	defer surface.Free()
	txtr, _ := renderer.CreateTextureFromSurface(surface)

	return txtr, surface.W, surface.H
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
	sol := &Solid{
		Position: &sdl.Rect{hdk.Position.X, hdk.Position.Y, TSzi, TSzi},
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

	var stCost float32 = 15

	if stCost > c.CurrentST {
		return
	}
	c.CurrentST -= stCost

	r := ActHitBox(c.Solid.Position, c.Solid.Orientation)
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
			OnCollDmg:   30,
			OnCollEvent: onColHdk,
		},
		CPattern: 0,
		MPattern: []Movement{
			Movement{*c.Solid.Orientation, 255},
		},
		Collision: 1,
	}
	Interactive = append(Interactive, h)
	PC.Solid.SetAnimation(MAN_CS_ANIM)
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
		PC.Speed = (PC.BaseSpeed * 2)
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
		if PC.Solid.Anim.PlayMode != 1 {
			PC.peformHaduken()
		}
	case KEY_LEFT_SHIT:
		PC.Speed = PC.BaseSpeed
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

func (s *Solid) SetAnimation(an *Animation) {
	var nA *Animation = &Animation{}
	*nA = *an

	s.Anim = nA
	s.Anim.Pose = 0
	s.Anim.PoseTick = 16
}

func (s *Solid) PlayAnimation() {
	s.Anim.PoseTick -= 1

	if s.Anim.PoseTick == 0 {
		s.Anim.PoseTick = 5
		if s.CharPtr != nil {
			anim := s.CharPtr.CurrentFacing()
			pvrPose := s.Anim.Pose
			s.Anim.Pose = getNextPose(s.Anim.Action, s.Anim.Pose)
			if anim != nil && s.Anim.Pose < pvrPose && s.Anim.PlayMode == 1 {
				s.SetAnimation(anim)
			}
		}
	}
}

func (s *Solid) procMovement(speed float32) {
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
	for _, obj := range CullMap {
		if obj == s || obj.Position == nil {
			continue
		}
		fr := feetRect(np)
		if checkCol(fr, obj.Position) {
			if obj.Handlers != nil && obj.Handlers.OnCollEvent != nil {
				obj.Handlers.OnCollEvent(s, obj)
			}
			if obj.Collision == 1 {
				return
			}
		}
	}

	if s.CharPtr != nil {
		anim := s.CharPtr.CurrentFacing()
		if anim != nil && s.Anim != nil && s.Anim.PlayMode != 1 {
			s.Anim.PlayMode = anim.PlayMode
			s.Anim.Action = anim.Action
		}
	}

	*s.Orientation = *s.Velocity
	s.Position = np
}

func catchEvents() bool {
	var c bool

	if PC.Solid.Anim.PlayMode == 1 {
		PC.Solid.PlayAnimation()
	}

	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyDownEvent:
			v := handleKeyEvent(t.Keysym.Sym)
			if v.X != 0 {
				PC.Solid.Velocity.X = v.X
				PC.Solid.Orientation.X = v.X
				c = true
			}
			if v.Y != 0 {
				PC.Solid.Velocity.Y = v.Y
				PC.Solid.Orientation.Y = v.Y
				c = true
			}
		case *sdl.KeyUpEvent:
			handleKeyUpEvent(t.Keysym.Sym)
		}
	}

	if c && PC.Solid.Anim.PlayMode == 0 {
		PC.Solid.PlayAnimation()
		PC.Solid.procMovement(PC.Speed)

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
	}

	if isMoving(PC.Solid.Velocity) && PC.Speed > PC.BaseSpeed {
		dpl := (PC.MaxST * 0.0009)

		if PC.CurrentST <= dpl {
			PC.CurrentST = 0
			PC.Speed = PC.BaseSpeed
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

		cX := rand.Int31n(s.CellsX)
		cY := rand.Int31n(s.CellsY)

		absolute_pos := &sdl.Rect{cX * TSzi, cY * TSzi, TSzi, TSzi}
		sol := &Solid{}

		switch rand.Int31n(9) {
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
					DoorTo:     doorTo,
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

			absolute_pos.H = 128
			absolute_pos.W = 128

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
					Position:  absolute_pos,
					Velocity:  &Vector2d{0, 0},
					Txt:       spritesheetTxt,
					Collision: 1,
					Anim:      PC.ActionMap.DOWN,
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
	EventTick uint8 = EVENT_TICK_LENGTH
	AiTick          = AI_TICK_LENGTH
	dbox      DBox  = DBox{BGColor: sdl.Color{90, 90, 90, 255}}
)

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
				OnActEvent: pickUp,
				OnPickUp:   addToInv,
			},
		},
	}

	Interactive = append(Interactive, instance.Solid)
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

	for _, cObj := range CullMap {
		// Ttl kill
		if cObj.Ttl > 0 && cObj.Ttl < now {
			cObj.Destroy()
			continue
		}

		// Kill logic
		if cObj.CharPtr != nil && cObj.CharPtr.CurrentHP <= 0 {

			if cObj.CharPtr.Drop != nil {
				PlaceDrop(cObj.CharPtr.Drop, cObj.Position)
			}

			cObj.Destroy()

			PC.CurrentXP += uint16(cObj.CharPtr.MaxHP / 10)
			if PC.CurrentXP >= PC.NextLvlXP {
				PC.CurrentXP = 0
				PC.Lvl++
				PC.NextLvlXP = PC.NextLvlXP * uint16(1+PC.Lvl/2)
			}
		}

		if cObj.Anim != nil {
			cObj.PlayAnimation()
		}

		if cObj.Handlers != nil && EventTick == 0 {
			fr := feetRect(PC.Solid.Position)
			if checkCol(fr, cObj.Position) {
				if cObj.Handlers.OnCollDmg != 0 {
					PC.depletHP(cObj.Handlers.OnCollDmg)
					if PC.Invinc == 0 {
						PC.Invinc = time.Now().Add(2 * time.Second).Unix()
					}
				}
				if cObj.Handlers.OnCollEvent != nil {
					cObj.Handlers.OnCollEvent(PC.Solid, cObj)
				}
			}
		}

		if AiTick == 0 {
			if cObj.Chase != nil && cObj.LoSCheck() {
				cObj.chase()
			} else {
				cObj.peformPattern(1)
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

} // end update()

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

	if diffX > 24 && s.Position.X > s.Chase.Position.X {
		s.Velocity.X = -1
	}

	if diffX > 24 && s.Position.X < s.Chase.Position.X {
		s.Velocity.X = 1
	}
	if diffY > 18 && s.Position.Y > s.Chase.Position.Y {
		s.Velocity.Y = -1
	}

	if diffY > 18 && s.Position.Y < s.Chase.Position.Y {
		s.Velocity.Y = 1
	}

	s.procMovement(s.CharPtr.Speed)

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
		s.procMovement(sp)
		s.CPattern += uint32(sp)
	}
}

func pickUp(picker *Solid, item *Solid) {
	if item.Handlers != nil && item.Handlers.OnPickUp != nil {
		item.Handlers.OnPickUp(picker, item)
	}
	item.Destroy()
}

func addToInv(picker *Solid, item *Solid) {
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
		return ch.ActionMap.DR
	}

	if ch.Solid.Orientation.X == 1 && ch.Solid.Orientation.Y == -1 {
		return ch.ActionMap.UR
	}

	if ch.Solid.Orientation.X == -1 && ch.Solid.Orientation.Y == 1 {
		return ch.ActionMap.DL
	}

	if ch.Solid.Orientation.X == -1 && ch.Solid.Orientation.Y == -1 {
		return ch.ActionMap.UL
	}

	return ch.ActionMap.DOWN
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

func (s *Solid) Destroy() {
	s.Position = nil
	s.Source = nil
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

		renderer.SetDrawColor(255, 255, 0, 255)
		renderer.DrawRect(scrPos)

		if inScreen(scrPos) {

			src := mon.Solid.Anim.Action[mon.Solid.Anim.Pose]

			renderer.Copy(mon.Solid.Txt, src, scrPos)
			renderer.DrawRect(scrPos)
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
	renderer.SetDrawColor(190, 190, 0, 255)
	renderer.FillRect(&sdl.Rect{10, 38, int32(calcPerc(float32(PC.CurrentXP), float32(PC.NextLvlXP))), 4})

	renderer.SetDrawColor(90, 90, 90, 255)
	renderer.FillRect(&sdl.Rect{0, 60, 240, 30})

	for i, stack := range PC.Inventory {
		counter := TextEl{
			Content: fmt.Sprintf("%d", stack.Qty),
			Font:    font,
			Color:   sdl.Color{255, 255, 255, 255},
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
		Color:   sdl.Color{255, 255, 255, 255},
	}
	lvl_txtr, W, H := lvl_TextEl.Bake(renderer)
	renderer.Copy(lvl_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{128, 60, W, H})

	dbg_content := fmt.Sprintf(
		"px %d py %d|vx %.1f vy %.1f (%.1f, %.1f) An:%d/%d/%d cull %d i %d cX %.1f cY %.1f L %dus ETick%d AiTick%d",
		PC.Solid.Position.X, PC.Solid.Position.Y, PC.Solid.Velocity.X, PC.Solid.Velocity.Y, PC.Solid.Orientation.X,
		PC.Solid.Orientation.Y, PC.Solid.Anim.Pose, PC.Solid.Anim.PoseTick, PC.Solid.Anim.PlayMode, len(CullMap),
		len(Interactive), Cam.P.X, Cam.P.Y, game_latency, EventTick, AiTick,
	)
	dbg_TextEl := TextEl{
		Font:    font,
		Content: dbg_content,
		Color:   sdl.Color{255, 255, 255, 255},
	}
	dbg_txtr, W, H := dbg_TextEl.Bake(renderer)
	renderer.Copy(dbg_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{0, winHeight - H, W, H})

	dbox.Present(renderer)

	for _, spw := range Spawners {
		renderer.DrawRect(worldToScreen(spw.Position, Cam))
	}
}

func (s *Scene) render(renderer *sdl.Renderer) {
	renderer.Clear()

	s._terrainRender(renderer)
	s._solidsRender(renderer)
	s._monstersRender(renderer)
	// Rendering the PC
	if !(PC.Invinc > 0 && EventTick == 2) {
		scrPos := worldToScreen(PC.Solid.Position, Cam)
		renderer.Copy(spritesheetTxt, PC.Solid.Anim.Action[PC.Solid.Anim.Pose], scrPos)
	}
	s._GUIRender(renderer)

	// FLUSH FRAME
	renderer.Present()
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

func inScreen(r *sdl.Rect) bool {
	return (r.X > (r.W*-1) && r.X < winWidth && r.Y > (r.H*-1) && r.Y < winHeight)
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
	ch.PushBack(12)
}

func (c *Char) PushBack(d float32) {
	f := c.Solid.Orientation
	c.Solid.Position.X -= int32(f.X * d * 4)
	c.Solid.Position.Y -= int32(f.Y * d * 4)
	// TODO WIRE PB ANIM
	c.Solid.SetAnimation(MAN_PB_ANIM)
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
	monstersImg, _ := img.Load("assets/textures/monsters.png")
	defer monstersImg.Free()
	defer tilesetImg.Free()
	defer spritesheetImg.Free()
	defer powerupsImg.Free()
	defer glowImg.Free()

	tilesetTxt, _ = renderer.CreateTextureFromSurface(tilesetImg)
	spritesheetTxt, _ = renderer.CreateTextureFromSurface(spritesheetImg)
	powerupsTxt, _ = renderer.CreateTextureFromSurface(powerupsImg)
	glowTxt, _ = renderer.CreateTextureFromSurface(glowImg)
	monstersTxt, _ = renderer.CreateTextureFromSurface(monstersImg)
	defer tilesetTxt.Destroy()
	defer spritesheetTxt.Destroy()
	defer powerupsTxt.Destroy()
	defer glowTxt.Destroy()
	defer monstersTxt.Destroy()

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
		Loot:          [8]Loot{},
		Lvl:           1,
		HP:            70,
		LoS:           120,
		Size:          32,
		LvlVariance:   0.5,
		ScalingFactor: 0.1,
	}

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

		sdl.Delay(24)
	}
}

func V2R(v Vector2d, w int32, h int32) *sdl.Rect {
	return &sdl.Rect{int32(v.X), int32(v.Y), w, h}
}

func change_scene(new_scene *Scene, staring_pos *Vector2d) {
	new_scene.build()
	Interactive = []*Solid{}
	new_scene.populate(200)
	scene = new_scene

	PC.Solid.Position = V2R(new_scene.StartPoint, TSzi, TSzi)
	Cam.P = new_scene.CamPoint
}

func feetRect(pos *sdl.Rect) *sdl.Rect {
	return &sdl.Rect{pos.X, pos.Y + 16, pos.W, pos.H - 16}
}

type SpawnPoint struct {
	Position  *sdl.Rect
	Frequency uint16
	LvlMod    uint8
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
		ActionMap: monsterTpl.ActionMap,
		Speed:     1,
		BaseSpeed: 1,
		CurrentHP: hp,
		MaxHP:     hp,
		Drop:      DropItem,
	}
	mon.Solid.SetAnimation(monsterTpl.ActionMap.DOWN)
	mon.Solid.CharPtr = &mon

	return &mon
}
