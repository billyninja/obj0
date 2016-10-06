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
	tileSize            int32 = 32
	cY                  int32 = winHeight / tileSize
	cX                  int32 = winWidth / tileSize
	WORLD_CELLS_X             = 500
	WORLD_CELLS_Y             = 200
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

type AnimatedObj struct {
	Position Vector2d
	Action   [8]*sdl.Rect
	Pose     uint8
	PoseTick uint32
}

type StillObj struct {
	Position  Vector2d
	Source    *sdl.Rect
	Collision uint8
	Anim      *AnimatedObj
}

var (
	winTitle string = "Go-SDL2 Obj0"
	event    sdl.Event
	GRASS    *sdl.Rect = &sdl.Rect{0, 0, tileSize, tileSize}
	TREE     *sdl.Rect = &sdl.Rect{0, 32, tileSize, tileSize}
	WALL     *sdl.Rect = &sdl.Rect{0, 64, tileSize, tileSize}
	DOOR     *sdl.Rect = &sdl.Rect{32, 288, tileSize, tileSize}
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
)

var (
	Cam = Camera{
		P:   Vector2d{300, 300},
		DZx: 30,
		DZy: 60,
	}
	PC = AnimatedObj{
		Position: Vector2d{Cam.P.X + 120, Cam.P.Y + 90},
		Action:   MAN_WALK_FRONT,
		Pose:     0,
		PoseTick: 16,
	}

	World   [WORLD_CELLS_X][WORLD_CELLS_Y]*sdl.Rect
	CullMap []*StillObj

	Obstacles  []*StillObj
	Explosions []*AnimatedObj
)

func checkCol(p1 Vector2d, p2 Vector2d) bool {
	return (p1.X < p2.X+tileSize &&
		p1.X+tileSize > p2.X &&
		p1.Y < p2.Y+tileSize &&
		p1.Y+tileSize > p2.Y)
}

func handleKeyEvent(key sdl.Keycode) {
	np := Vector2d{PC.Position.X, PC.Position.Y}

	switch key {
	case 1073741906:
		PC.Action = MAN_WALK_BACK
		np.Y -= 1
	case 1073741905:
		PC.Action = MAN_WALK_FRONT
		np.Y += 1
	case 1073741904:
		PC.Action = MAN_WALK_LEFT
		np.X -= 1
	case 1073741903:
		PC.Action = MAN_WALK_RIGHT
		np.X += 1
	}

	if np.X == PC.Position.X && np.Y == PC.Position.Y {
		PC.Pose = 0
	} else {

		if np.X <= 0 || np.Y <= 0 {
			return
		}

		for _, obj := range CullMap {
			if checkCol(np, obj.Position) {
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

		PC.Position = np

		PC.PoseTick -= 1
		if PC.PoseTick == 0 {
			PC.Pose = get_next_pose(PC.Action, PC.Pose)
			PC.PoseTick = 16
		}
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

func updateScene() {

	// update explosions poses
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

func worldToScreen(pos Vector2d, cam Camera) Vector2d {
	return Vector2d{
		X: pos.X - cam.P.X,
		Y: pos.Y - cam.P.Y,
	}
}

func inScreen(p Vector2d) bool {
	return (p.X > 0 && p.X < winWidth &&
		p.Y > 0 && p.Y < winHeight)
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

	CullMap = []*StillObj{}

	for _, obs := range Obstacles {
		scrPoint := worldToScreen(obs.Position, Cam)
		if inScreen(scrPoint) {
			CullMap = append(CullMap, obs)
		}
	}

	for _, exp := range Explosions {
		scrPoint := worldToScreen(exp.Position, Cam)
		if inScreen(scrPoint) {

			pos := sdl.Rect{scrPoint.X, scrPoint.Y, tileSize, tileSize}
			src := exp.Action[exp.Pose]
			renderer.Copy(ss, src, &pos)

			CullMap = append(CullMap, &StillObj{
				Position:  exp.Position,
				Source:    src,
				Collision: 1,
				Anim:      exp,
			})
		}
	}

	// Rendering the PC
	scrPoint := worldToScreen(PC.Position, Cam)
	pos := sdl.Rect{scrPoint.X, scrPoint.Y, tileSize, tileSize}
	renderer.Copy(ss, PC.Action[PC.Pose], &pos)

	renderer.Present()
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

	for i := 0; i < 2020; i++ {
		cX := rand.Int31n(WORLD_CELLS_X)
		cY := rand.Int31n(WORLD_CELLS_Y)
		World[cX][cY] = WALL
		Obstacles = append(Obstacles, &StillObj{
			Position:  Vector2d{cX * tileSize, cY * tileSize},
			Source:    WALL,
			Collision: 1,
		})
	}

	for i := 0; i < 800; i++ {
		cX := rand.Int31n(WORLD_CELLS_X)
		cY := rand.Int31n(WORLD_CELLS_Y)

		Explosions = append(Explosions, &AnimatedObj{
			Position: Vector2d{cX * tileSize, cY * tileSize},
			Action:   EXPLOSION_A,
			Pose:     0,
			PoseTick: 16,
		})
	}

	for running {
		running = catchEvents()
		then := time.Now()
		updateScene()
		renderScene(renderer, tilesetTxt, spritesheetTxt)
		println((time.Since(then)) / time.Microsecond)
		//sdl.Delay(33)
	}
}

func buildDummyWorld(cellsX int, cellsY int) {
	for i := 0; i < cellsX; i++ {
		for j := 0; j < cellsY; j++ {
			tile := GRASS

			if rand.Int31n(10) < 1 {
				tile = TREE
			}

			World[i][j] = tile
		}
	}
}
