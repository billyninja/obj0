package main

import (
	//"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"math/rand"
	//"os"
)

const (
	winWidth, winHeight int32 = 640, 480
	tileSize            int32 = 32
	cY                  int32 = winHeight / tileSize
	cX                  int32 = winWidth / tileSize
)

var (
	winTitle string = "Go-SDL2 Obj0"
	event    sdl.Event
	GRASS    *sdl.Rect = &sdl.Rect{0, 0, tileSize, tileSize}
	TREE     *sdl.Rect = &sdl.Rect{0, 32, tileSize, tileSize}
	WALL     *sdl.Rect = &sdl.Rect{0, 64, tileSize, tileSize}
	WOMAN    *sdl.Rect = &sdl.Rect{0, 128, tileSize, tileSize}
)

var (
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
)

var (
	EXPLOSION_S1 *sdl.Rect = &sdl.Rect{128, 0, tileSize, tileSize}
	EXPLOSION_S2 *sdl.Rect = &sdl.Rect{128, 32, tileSize, tileSize}
	EXPLOSION_S3 *sdl.Rect = &sdl.Rect{128, 64, tileSize, tileSize}
	EXPLOSION_S4 *sdl.Rect = &sdl.Rect{128, 96, tileSize, tileSize}

	EXPLOSION_A [8]*sdl.Rect = [8]*sdl.Rect{EXPLOSION_S1, EXPLOSION_S2, EXPLOSION_S3, EXPLOSION_S4}
)

type Vector2d struct {
	X int32
	Y int32
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
}

var PC = AnimatedObj{
	Position: Vector2d{0, 0},
	Action:   MAN_WALK_FRONT,
	Pose:     0,
	PoseTick: 16,
}

var Obstacles []*StillObj
var Explosions []*AnimatedObj

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
		for _, obt := range Obstacles {
			if checkCol(np, obt.Position) {
				return
			}
		}
		for _, exp := range Explosions {
			if checkCol(np, exp.Position) {
				return
			}
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
	for _, exps := range Explosions {
		exps.PoseTick -= 1
		if exps.PoseTick == 0 {
			exps.Pose = get_next_pose(exps.Action, exps.Pose)
			exps.PoseTick = 16
		}
	}
}

func renderScene(renderer *sdl.Renderer, ts *sdl.Texture, ss *sdl.Texture) {
	renderer.Clear()

	var (
		ys int32
		xs int32
	)

	// Terrain
	for ys = 0; ys < cY; ys++ {
		for xs = 0; xs < cX; xs++ {
			pos := &sdl.Rect{xs * tileSize, ys * tileSize, tileSize, tileSize}
			renderer.Copy(ts, GRASS, pos)
		}
	}

	pos := &sdl.Rect{32, 0, tileSize, tileSize}
	renderer.Copy(ss, WOMAN, pos)

	for _, obt := range Obstacles {
		pos = &sdl.Rect{obt.Position.X, obt.Position.Y, tileSize, tileSize}
		renderer.Copy(ts, obt.Source, pos)
	}

	for _, exps := range Explosions {
		pos = &sdl.Rect{exps.Position.X, exps.Position.Y, tileSize, tileSize}
		renderer.Copy(ss, exps.Action[exps.Pose], pos)
	}

	pos = &sdl.Rect{PC.Position.X, PC.Position.Y, tileSize, tileSize}
	renderer.Copy(ss, PC.Action[PC.Pose], pos)

	renderer.Present()
}

func main() {
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
	var tick1 uint32 = sdl.GetTicks()
	var tick2 uint32

	renderer.SetDrawColor(0, 0, 255, 255)

	for i := 0; i < 20; i++ {
		cX := rand.Int31n(int32(winWidth / tileSize))
		cY := rand.Int31n(int32(winHeight / tileSize))

		Obstacles = append(Obstacles, &StillObj{
			Position:  Vector2d{cX * tileSize, cY * tileSize},
			Source:    WALL,
			Collision: 1,
		})
	}

	cellsX := int32(winWidth / tileSize)
	cellsY := int32(winHeight / tileSize)
	for i := 0; i < 8; i++ {

		cX := rand.Int31n(cellsX)
		cY := rand.Int31n(cellsY)

		Explosions = append(Explosions, &AnimatedObj{
			Position: Vector2d{cX * tileSize, cY * tileSize},
			Action:   EXPLOSION_A,
			Pose:     0,
			PoseTick: 16,
		})
	}

	for running {
		running = catchEvents()
		tick2 = sdl.GetTicks()
		dt := uint32(tick2 - tick1)
		println(dt)
		updateScene()
		renderScene(renderer, tilesetTxt, spritesheetTxt)
		tick1 = tick2
		sdl.Delay(33)
	}
}
