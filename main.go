package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

var (
	winTitle string = "Go-SDL2 Obj0"
	event    sdl.Event
)

const (
	winWidth, winHeight int32 = 640, 480
	tileSize            int32 = 32
)

func handleKeyEvent(key sdl.Keycode) {
	return
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

func main() {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(winWidth), int(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	var running bool = true
	for running {
		running = catchEvents()
		sdl.Delay(2)
	}
}
