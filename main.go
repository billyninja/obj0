package main

import (
	//"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
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
	TREE     *sdl.Rect = &sdl.Rect{32, 32, tileSize, tileSize}
	WOMAN    *sdl.Rect = &sdl.Rect{128, 0, tileSize, tileSize}
)

var (
	MAN_FRONT_R    *sdl.Rect    = &sdl.Rect{0, 0, tileSize, tileSize}
	MAN_FRONT_N    *sdl.Rect    = &sdl.Rect{32, 0, tileSize, tileSize}
	MAN_FRONT_L    *sdl.Rect    = &sdl.Rect{64, 0, tileSize, tileSize}
	MAN_LEFT_R     *sdl.Rect    = &sdl.Rect{0, 32, tileSize, tileSize}
	MAN_LEFT_N     *sdl.Rect    = &sdl.Rect{32, 32, tileSize, tileSize}
	MAN_LEFT_L     *sdl.Rect    = &sdl.Rect{64, 32, tileSize, tileSize}
	MAN_RIGHT_R    *sdl.Rect    = &sdl.Rect{0, 64, tileSize, tileSize}
	MAN_RIGHT_N    *sdl.Rect    = &sdl.Rect{32, 64, tileSize, tileSize}
	MAN_RIGHT_L    *sdl.Rect    = &sdl.Rect{64, 64, tileSize, tileSize}
	MAN_BACK_R     *sdl.Rect    = &sdl.Rect{0, 96, tileSize, tileSize}
	MAN_BACK_N     *sdl.Rect    = &sdl.Rect{32, 96, tileSize, tileSize}
	MAN_BACK_L     *sdl.Rect    = &sdl.Rect{64, 96, tileSize, tileSize}
	MAN_WALK_FRONT [8]*sdl.Rect = [8]*sdl.Rect{MAN_FRONT_N, MAN_FRONT_R, MAN_FRONT_N, MAN_FRONT_L}
	MAN_WALK_LEFT  [8]*sdl.Rect = [8]*sdl.Rect{MAN_LEFT_N, MAN_LEFT_R, MAN_LEFT_N, MAN_LEFT_L}
	MAN_WALK_RIGHT [8]*sdl.Rect = [8]*sdl.Rect{MAN_RIGHT_N, MAN_RIGHT_R, MAN_RIGHT_N, MAN_RIGHT_L}
	MAN_WALK_BACK  [8]*sdl.Rect = [8]*sdl.Rect{MAN_BACK_N, MAN_BACK_R, MAN_BACK_N, MAN_BACK_L}
)

type Vector2d struct {
	X int32
	Y int32
}

type CharGFX struct {
	Position Vector2d
	Action   [8]*sdl.Rect
	Pose     uint8
	PoseTick uint32
}

var PC = CharGFX{
	Position: Vector2d{0, 0},
	Action:   MAN_WALK_FRONT,
	Pose:     0,
	PoseTick: 16,
}

func handleKeyEvent(key sdl.Keycode) {
	var (
		vel   int32 = 1
		still bool  = true
	)

	switch key {
	case 1073741906:
		PC.Action = MAN_WALK_BACK
		PC.Position.Y -= vel
		still = false
	case 1073741905:
		PC.Action = MAN_WALK_FRONT
		PC.Position.Y += vel
		still = false
	case 1073741904:
		PC.Action = MAN_WALK_LEFT
		PC.Position.X -= vel
		still = false
	case 1073741903:
		PC.Action = MAN_WALK_RIGHT
		PC.Position.X += vel
		still = false
	}

	if still {
		PC.Pose = 0
	} else {
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

func renderScene(renderer *sdl.Renderer, ts *sdl.Texture, ss *sdl.Texture) {
	renderer.Clear()

	var (
		ys int32
		xs int32
	)

	for ys = 0; ys < cY; ys++ {
		for xs = 0; xs < cX; xs++ {
			pos := &sdl.Rect{xs * tileSize, ys * tileSize, tileSize, tileSize}
			renderer.Copy(ts, GRASS, pos)
		}
	}

	pos := &sdl.Rect{PC.Position.X, PC.Position.Y, tileSize, tileSize}
	renderer.Copy(ss, PC.Action[PC.Pose], pos)

	pos = &sdl.Rect{32, 0, tileSize, tileSize}
	renderer.Copy(ss, WOMAN, pos)

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

	var running bool = true
	var tick1 uint32 = sdl.GetTicks()
	var tick2 uint32

	renderer.SetDrawColor(0, 0, 255, 255)
	for running {
		running = catchEvents()
		tick2 = sdl.GetTicks()
		dt := uint32(tick2 - tick1)
		println(dt)
		renderScene(renderer, tilesetTxt, spritesheetTxt)
		tick1 = tick2
		sdl.Delay(33)
	}
}
