package main

import (
	//"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"math/rand"
	"runtime"
	"time"
	//"os"
)

const (
	winWidth, winHeight int32 = 640, 480
	tileSize                  = 32
	cY                        = winHeight / tileSize
	cX                        = winWidth / tileSize
	WORLD_CELLS_X             = 500
	WORLD_CELLS_Y             = 200
	KEY_ARROW_UP              = 1073741906
	KEY_ARROW_DOWN            = 1073741905
	KEY_ARROW_LEFT            = 1073741904
	KEY_ARROW_RIGHT           = 1073741903
	KEY_LEFT_SHIT             = 1073742049
	KEY_SPACE_BAR             = 32
)

type Vector2d struct {
	X int32
	Y int32
}

type Camera struct {
	P   Vector2d
	DZx int32
	DZy int32
}

type Event func(uint16)

type InteractionHandlers struct {
	OnCollDmg   uint16
	OnCollPush  *Vector2d
	OnCollEvent Event
	OnActDmg    uint16
	OnActPush   *Vector2d
	OnActEvent  Event
}

type Solid struct {
	Position  Vector2d
	Source    *sdl.Rect
	Anim      *Animation
	Handlers  *InteractionHandlers
	Txt       *sdl.Texture
	Collision uint8
}

type Animation struct {
	Action   [8]*sdl.Rect
	Pose     uint8
	PoseTick uint32
}

type Char struct {
	Solid *Solid

	Speed     int32
	CurrentHP uint16
	MaxHP     uint16

	CurrentST uint16
	MaxST     uint16
}

func Facing() Vector2d {
	off := Vector2d{0, 0}
	switch PC.Solid.Anim.Action {
	case MAN_WALK_BACK:
		off = Vector2d{0, -1}
	case MAN_WALK_FRONT:
		off = Vector2d{0, 1}
	case MAN_WALK_LEFT:
		off = Vector2d{-1, 0}
	case MAN_WALK_RIGHT:
		off = Vector2d{1, 0}
	}
	return off
}

func ActHitBox(source, facing Vector2d) *sdl.Rect {
	return &sdl.Rect{
		source.X + (facing.X * tileSize),
		source.Y + (facing.Y * tileSize),
		tileSize,
		tileSize,
	}
}

var (
	winTitle string = "Go-SDL2 Obj0"
	event    sdl.Event
	GRASS    *sdl.Rect = &sdl.Rect{0, 0, tileSize, tileSize}
	TREE     *sdl.Rect = &sdl.Rect{0, 32, tileSize, tileSize}
	DIRT     *sdl.Rect = &sdl.Rect{703, 0, tileSize, tileSize}
	WALL     *sdl.Rect = &sdl.Rect{0, 64, tileSize, tileSize}
	DOOR     *sdl.Rect = &sdl.Rect{256, 32, tileSize, tileSize}
	WOMAN    *sdl.Rect = &sdl.Rect{0, 128, tileSize, tileSize}

	// MAIN CHAR POSES AND ANIMATIONS
	MAN_FRONT_R *sdl.Rect = &sdl.Rect{0, 0, tileSize, tileSize}
	MAN_FRONT_N *sdl.Rect = &sdl.Rect{32, 0, tileSize, tileSize}
	MAN_FRONT_L *sdl.Rect = &sdl.Rect{64, 0, tileSize, tileSize}
	MAN_LEFT_R  *sdl.Rect = &sdl.Rect{0, 32, tileSize, tileSize}
	MAN_LEFT_N  *sdl.Rect = &sdl.Rect{32, 32, tileSize, tileSize}
	MAN_LEFT_L  *sdl.Rect = &sdl.Rect{64, 32, tileSize, tileSize}
	MAN_RIGHT_R *sdl.Rect = &sdl.Rect{0, 64, tileSize, tileSize}
	MAN_RIGHT_N *sdl.Rect = &sdl.Rect{32, 64, tileSize, tileSize}
	MAN_RIGHT_L *sdl.Rect = &sdl.Rect{64, 64, tileSize, tileSize}
	MAN_BACK_R  *sdl.Rect = &sdl.Rect{0, 96, tileSize, tileSize}
	MAN_BACK_N  *sdl.Rect = &sdl.Rect{32, 96, tileSize, tileSize}
	MAN_BACK_L  *sdl.Rect = &sdl.Rect{64, 96, tileSize, tileSize}

	MAN_WALK_FRONT [8]*sdl.Rect = [8]*sdl.Rect{MAN_FRONT_N, MAN_FRONT_R, MAN_FRONT_N, MAN_FRONT_L}
	MAN_WALK_LEFT  [8]*sdl.Rect = [8]*sdl.Rect{MAN_LEFT_N, MAN_LEFT_R, MAN_LEFT_N, MAN_LEFT_L}
	MAN_WALK_RIGHT [8]*sdl.Rect = [8]*sdl.Rect{MAN_RIGHT_N, MAN_RIGHT_R, MAN_RIGHT_N, MAN_RIGHT_L}
	MAN_WALK_BACK  [8]*sdl.Rect = [8]*sdl.Rect{MAN_BACK_N, MAN_BACK_R, MAN_BACK_N, MAN_BACK_L}

	EXPLOSION_S1 *sdl.Rect = &sdl.Rect{128, 0, tileSize, tileSize}
	EXPLOSION_S2 *sdl.Rect = &sdl.Rect{128, 32, tileSize, tileSize}
	EXPLOSION_S3 *sdl.Rect = &sdl.Rect{128, 64, tileSize, tileSize}
	EXPLOSION_S4 *sdl.Rect = &sdl.Rect{128, 96, tileSize, tileSize}

	EXPLOSION_A [8]*sdl.Rect = [8]*sdl.Rect{EXPLOSION_S1, EXPLOSION_S2, EXPLOSION_S3, EXPLOSION_S4}

	LAVA_S1 *sdl.Rect = &sdl.Rect{192, 0, tileSize, tileSize}
	LAVA_S2 *sdl.Rect = &sdl.Rect{224, 0, tileSize, tileSize}
	LAVA_S3 *sdl.Rect = &sdl.Rect{256, 0, tileSize, tileSize}

	LAVA_A [8]*sdl.Rect = [8]*sdl.Rect{LAVA_S1, LAVA_S2, LAVA_S3, LAVA_S3, LAVA_S2}

	LAVA_ANIMATION = &Animation{
		Action:   LAVA_A,
		Pose:     0,
		PoseTick: 16,
	}

	LAVA_HANDLERS = &InteractionHandlers{
		OnCollDmg: 3,
	}
)

var (
	Cam = Camera{
		P:   Vector2d{300, 300},
		DZx: 30,
		DZy: 60,
	}
	PC = Char{
		Solid: &Solid{
			Position: Vector2d{Cam.P.X + 120, Cam.P.Y + 90},
			Anim: &Animation{
				Action:   MAN_WALK_FRONT,
				Pose:     0,
				PoseTick: 16,
			},
		},
		Speed:     1,
		CurrentHP: 220,
		MaxHP:     250,
		CurrentST: 65,
		MaxST:     100,
	}

	World       [WORLD_CELLS_X][WORLD_CELLS_Y]*sdl.Rect
	Interactive []*Solid
	GUI         []*sdl.Rect
	CullMap     []*Solid
)

func checkCol(p1 Vector2d, p2 Vector2d) bool {
	return (p1.X < p2.X+tileSize &&
		p1.X+tileSize > p2.X &&
		p1.Y < p2.Y+tileSize &&
		p1.Y+tileSize > p2.Y)
}

func ActProc() {
	action_hit_box := ActHitBox(PC.Solid.Position, Facing())
	action_origin := Vector2d{action_hit_box.X, action_hit_box.Y}

	// Debug hint
	GUI = append(GUI, action_hit_box)

	for _, obj := range CullMap {
		if obj.Handlers != nil &&
			obj.Handlers.OnActEvent != nil &&
			checkCol(action_origin, obj.Position) {
			obj.Handlers.OnActEvent(1)
			return
		}
	}
}

func handleKeyEvent(key sdl.Keycode) {
	np := Vector2d{PC.Solid.Position.X, PC.Solid.Position.Y}

	switch key {
	case KEY_SPACE_BAR:
		ActProc()
		return
	case KEY_LEFT_SHIT:
		PC.Speed = 3
		np.Y -= PC.Speed
	case KEY_ARROW_UP:
		PC.Solid.Anim.Action = MAN_WALK_BACK
		np.Y -= PC.Speed
	case KEY_ARROW_DOWN:
		PC.Solid.Anim.Action = MAN_WALK_FRONT
		np.Y += PC.Speed
	case KEY_ARROW_LEFT:
		PC.Solid.Anim.Action = MAN_WALK_LEFT
		np.X -= PC.Speed
	case KEY_ARROW_RIGHT:
		PC.Solid.Anim.Action = MAN_WALK_RIGHT
		np.X += PC.Speed
	}

	var outbound bool = (np.X <= 0 || np.Y <= 0)
	if np.X == PC.Solid.Position.X && np.Y == PC.Solid.Position.Y || outbound {
		PC.Solid.Anim.Pose = 0
		return
	}

	for _, obj := range CullMap {
		if checkCol(np, obj.Position) && obj.Collision == 1 {
			return
		}
	}

	newScreenPos := worldToScreen(np, Cam)
	if (Cam.DZx - newScreenPos.X) > 0 {
		Cam.P.X -= (Cam.DZx - newScreenPos.X)
	}

	if (winWidth - Cam.DZx) < (newScreenPos.X + tileSize) {
		Cam.P.X += (newScreenPos.X + tileSize) - (winWidth - Cam.DZx)
	}

	if (Cam.DZy - newScreenPos.Y) > 0 {
		Cam.P.Y -= (Cam.DZy - newScreenPos.Y)
	}

	if (winHeight - Cam.DZy) < (newScreenPos.Y + tileSize) {
		Cam.P.Y += (newScreenPos.Y + tileSize) - (winHeight - Cam.DZy)
	}

	PC.Solid.Position = np

	PC.Solid.Anim.PoseTick -= 1
	if PC.Solid.Anim.PoseTick == 0 {
		PC.Solid.Anim.Pose = get_next_pose(PC.Solid.Anim.Action, PC.Solid.Anim.Pose)
		PC.Solid.Anim.PoseTick = 16
	}
}

func catchEvents() bool {
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyDownEvent:
			handleKeyEvent(t.Keysym.Sym)
		}
	}
	return true
}

func get_next_pose(action [8]*sdl.Rect, currPose uint8) uint8 {
	if action[currPose+1] == nil {
		return 0
	} else {
		return currPose + 1
	}
}

var EventTick uint8 = 16

func updateScene() {
	EventTick -= 1
	if EventTick == 0 {
		for _, obj := range CullMap {
			if checkCol(PC.Solid.Position, obj.Position) {
				if obj.Handlers != nil && obj.Handlers.OnCollDmg != 0 {
					depletHP(obj.Handlers.OnCollDmg)
				}
			}
		}
		//PC.Speed = 1
		EventTick = 16
	}

	// update Interactives poses
	for _, cObj := range CullMap {
		if cObj.Anim == nil {
			continue
		}
		animObj := cObj.Anim
		animObj.PoseTick -= 1
		if animObj.PoseTick == 0 {
			animObj.Pose = get_next_pose(animObj.Action, animObj.Pose)
			animObj.PoseTick = 16
		}
	}
}

func calcPerc(v1 uint16, v2 uint16) float32 {
	return (float32(v1) / float32(v2) * 100)
}

func worldToScreen(pos Vector2d, cam Camera) Vector2d {
	return Vector2d{
		X: pos.X - cam.P.X,
		Y: pos.Y - cam.P.Y,
	}
}

func inScreen(p Vector2d) bool {
	return (p.X > (tileSize*-1) && p.X < winWidth &&
		p.Y > (tileSize*-1) && p.Y < winHeight)
}

func renderScene(renderer *sdl.Renderer, ts *sdl.Texture, ss *sdl.Texture) {
	var init int32 = 0
	var Source *sdl.Rect

	var offsetX, offsetY int32 = tileSize, tileSize

	renderer.Clear()
	renderer.SetDrawColor(0, 0, 0, 255)

	// Rendering the terrain
	for winY := init; winY < winHeight; winY += offsetY {
		for winX := init; winX < winWidth; winX += offsetX {

			offsetX = (tileSize - ((Cam.P.X + winX) % tileSize))
			offsetY = (tileSize - ((Cam.P.Y + winY) % tileSize))

			worldCellX := uint16((Cam.P.X + winX) / tileSize)
			worldCellY := uint16((Cam.P.Y + winY) / tileSize)

			if worldCellX <= 0 ||
				worldCellX > WORLD_CELLS_X ||
				worldCellY <= 0 ||
				worldCellY > WORLD_CELLS_Y {
				continue
			}

			gfx := World[worldCellX][worldCellY]

			if offsetX != int32(tileSize) || offsetY != int32(tileSize) {
				Source = &sdl.Rect{gfx.X + (tileSize - offsetX), gfx.Y + (tileSize - offsetY), offsetX, offsetY}
			} else {
				Source = gfx
			}

			screenPos := sdl.Rect{winX, winY, offsetX, offsetY}
			renderer.Copy(ts, Source, &screenPos)
		}
	}

	CullMap = []*Solid{}

	for _, obj := range Interactive {
		scrPoint := worldToScreen(obj.Position, Cam)
		if inScreen(scrPoint) {

			pos := sdl.Rect{scrPoint.X, scrPoint.Y, tileSize, tileSize}
			var src *sdl.Rect
			if obj.Anim != nil {
				src = obj.Anim.Action[obj.Anim.Pose]
			} else {
				src = obj.Source
			}
			renderer.Copy(obj.Txt, src, &pos)

			CullMap = append(CullMap, &Solid{
				Position:  obj.Position,
				Source:    src,
				Collision: obj.Collision,
				Anim:      obj.Anim,
				Handlers:  obj.Handlers,
			})
		}
	}

	// Rendering the PC
	scrPoint := worldToScreen(PC.Solid.Position, Cam)
	pos := sdl.Rect{scrPoint.X, scrPoint.Y, tileSize, tileSize}
	renderer.Copy(ss, PC.Solid.Anim.Action[PC.Solid.Anim.Pose], &pos)

	// HEALTH BAR BG
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.FillRect(&sdl.Rect{20, 20, 100, 12})

	// HEALTH BAR FG
	renderer.SetDrawColor(0, 255, 0, 255)
	renderer.FillRect(&sdl.Rect{20, 20, int32(calcPerc(PC.CurrentHP, PC.MaxHP)), 12})

	// MANA BAR BG
	renderer.SetDrawColor(60, 60, 60, 255)
	renderer.FillRect(&sdl.Rect{20, 40, 100, 12})

	// MANA BAR FG
	renderer.SetDrawColor(0, 0, 255, 255)
	x := int32(calcPerc(PC.CurrentST, PC.MaxST))
	renderer.FillRect(&sdl.Rect{20, 40, x, 12})

	for _, el := range GUI {
		scrPoint := worldToScreen(Vector2d{el.X, el.Y}, Cam)
		rect := &sdl.Rect{scrPoint.X, scrPoint.Y, el.W, el.H}
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(rect)
	}

	// FLUSH FRAME
	renderer.Present()

}

func depletHP(dmg uint16) {
	hp := PC.CurrentHP
	if hp <= 0 {
		return
	}
	if hp-dmg < 0 {
		hp = 0
	} else {
		hp -= dmg
	}

	PC.CurrentHP = hp
}

func BashDoor(a uint16) {
	PC.Solid.Position.X = 200
	PC.Solid.Position.Y = 100
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

	tilesetImg, _ := img.Load("assets/textures/ts1.png")
	defer tilesetImg.Free()

	tilesetTxt, _ := renderer.CreateTextureFromSurface(tilesetImg)
	defer tilesetTxt.Destroy()

	spritesheetImg, _ := img.Load("assets/textures/actor3.png")
	defer spritesheetImg.Free()

	spritesheetTxt, _ := renderer.CreateTextureFromSurface(spritesheetImg)
	defer spritesheetTxt.Destroy()

	particlesImg, _ := img.Load("assets/textures/actor3.png")
	defer particlesImg.Free()

	particlesTxt, _ := renderer.CreateTextureFromSurface(particlesImg)
	defer particlesTxt.Destroy()

	var running bool = true

	renderer.SetDrawColor(0, 0, 255, 255)

	buildDummyWorld(WORLD_CELLS_X, WORLD_CELLS_Y)

	rand.Seed(int64(300 * WORLD_CELLS_X * WORLD_CELLS_Y))
	for i := 0; i < 5000; i++ {

		cX := rand.Int31n(WORLD_CELLS_X)
		cY := rand.Int31n(WORLD_CELLS_Y)
		obj_type := rand.Int31n(10)
		pos := Vector2d{cX * tileSize, cY * tileSize}
		sol := &Solid{}

		switch obj_type {
		case 1:
			sol = &Solid{
				Position:  pos,
				Txt:       tilesetTxt,
				Anim:      LAVA_ANIMATION,
				Handlers:  LAVA_HANDLERS,
				Collision: 0,
			}
			break
		case 2:
			sol = &Solid{
				Position:  pos,
				Txt:       tilesetTxt,
				Source:    DOOR,
				Anim:      nil,
				Collision: 1,
				Handlers: &InteractionHandlers{
					OnActEvent: BashDoor,
				},
			}
			break
		}

		Interactive = append(Interactive, sol)

	}

	for running {
		running = catchEvents()
		then := time.Now()
		updateScene()
		renderScene(renderer, tilesetTxt, spritesheetTxt)
		println((time.Since(then)) / time.Microsecond)
		sdl.Delay(24)
	}
}

func buildDummyWorld(cellsX int, cellsY int) {

	rand.Seed(int64(cellsX * cellsY))

	for i := 0; i < cellsX; i++ {
		rand.Seed(int64(cellsX * i))
		for j := 0; j < cellsY; j++ {

			seed := int64(i * j)
			if seed%2 == 1 {
				rand.Seed(seed)
			}

			tile := GRASS
			r := rand.Int31n(100)
			if r < 10 {
				tile = DIRT
			}

			World[i][j] = tile
		}
	}
}
