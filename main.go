package main

import (
	"fmt"
	"github.com/billyninja/obj0/assets"
	"github.com/billyninja/obj0/core"
	"github.com/billyninja/obj0/game"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"runtime"
	"time"
)

const (
	CharSize            int32 = 64
	winWidth, winHeight int32 = 1280, 720
)

var (
	event        sdl.Event
	game_latency time.Duration
	Controls     *game.ControlState = &game.ControlState{}
	currScene    *game.Scene
	PC           = game.SceneEntity{
		Solid: &game.Solid{
			Velocity:    &core.Vector2d{0, 0},
			Orientation: &core.Vector2d{0, -1},
			Anim:        game.MAN_PU_ANIM,
		},
		Char: &game.Char{
			Lvl:       1,
			CurrentXP: 30,
			NextLvlXP: 100,
			SpeedMod:  0,
			BaseSpeed: 2,
			CurrentHP: 220,
			MaxHP:     250,
			CurrentST: 250,
			MaxST:     300,
			Inventory: []*game.ItemStack{{game.GreenBlob, 2}},
		},
	}
	PCptr *game.SceneEntity = &PC
)

func render(s *game.Scene, renderer *sdl.Renderer) {
	renderer.Clear()

	s.TerrainRender(renderer)
	s.SolidsRender(renderer)
	s.MonstersRender(renderer)
	s.PCRender(PCptr, renderer)
	s.VFXRender(renderer)
	s.GUIRender(PC.Char, renderer)

	debug_info(renderer)
	// FLUSH FRAME
	renderer.Present()
}

func update(scn *game.Scene) {
	scn.EventTick -= 1
	scn.AiTick -= 1

	now := time.Now().Unix()

	if PC.Char.Invinc > 0 && now > PC.Char.Invinc {
		PC.Char.Invinc = 0
	}

	if len(scn.DBox.Text) > 0 {
		scn.AiTick = scn.AiTickLength
	}

	scn.Recenter(PC.Solid.Position)

	for _, se := range scn.CullMap {
		sol := se.Solid
		if sol.Position == nil {
			continue
		}
		// Ttl kill
		if sol.Ttl > 0 && sol.Ttl < now {
			sol.Destroy()
			continue
		}

		if sol.CharPtr != nil && sol.CharPtr.CurrentHP <= 0 {
			if sol.CharPtr.Drop != nil {
				scn.PlaceDrop(sol.CharPtr.Drop, sol.Position)
			}
			scn.SpawnVFX(
				&sdl.Rect{sol.Position.X, sol.Position.Y, 92, 92},
				nil,
				game.Puff,
			)
			sol.Destroy()

			PC.Char.CurrentXP += uint16(sol.CharPtr.MaxHP / 10)
			if PC.Char.CurrentXP >= PC.Char.NextLvlXP {
				PC.Char.CurrentXP = 0
				PC.Char.Lvl++
				PC.Char.NextLvlXP = PC.Char.NextLvlXP * uint16(1+PC.Char.Lvl/2)
			}
			continue
		}

		if sol.Anim != nil {
			sol.PlayAnimation()
		}

		if scn.AiTick == 0 {

			if sol.Anim != nil && sol.Anim.PlayMode == 1 {
				continue
			}
			if sol.CharPtr != nil {
				sol.CharPtr.AtkCoolDownC -= sol.CharPtr.AtkSpeed
			}

			if sol.Chase != nil && sol.LoSCheck(PC.Solid) {
				scn.DoChase(se)
			} else {
				scn.PeformPattern(se)
			}
		}
	}

	if scn.EventTick == 0 {
		for _, spw := range scn.Spawners {
			if spw.Frequency <= 0 {
				spw.Frequency += 1
			}

			if uint16(rand.Int31n(1000)) < spw.Frequency {
				scn.SpawnMonster(spw, PCptr)
				spw.Frequency -= 1
			}
		}
		scn.EventTick = scn.EventTickLength
	}

	if scn.AiTick == 0 {
		scn.AiTick = scn.AiTickLength
	}

	for _, vi := range scn.Visual {
		if vi.Ttl > 0 && vi.Ttl < now {
			vi.Destroy()
			continue
		}
		vi.UpdateAnim()
	}
} // end update()

func debug_info(renderer *sdl.Renderer) {
	dbg_content := fmt.Sprintf(
		"px %d py %d | Cx %.1f Cy %.1f |vx %.1f vy | %.1f (%.1f, %.1f) |"+
			" An:%d/%d/%d cull %d i %d cX %d cY %d L %dus ETick%d AiTick%d",
		PC.Solid.Position.X,
		PC.Solid.Position.Y,
		Controls.DPAD.X,
		Controls.DPAD.Y,
		PC.Solid.Velocity.X,
		PC.Solid.Velocity.Y,
		PC.Solid.Orientation.X,
		PC.Solid.Orientation.Y,
		PC.Solid.Anim.Pose,
		PC.Solid.Anim.PoseTick,
		PC.Solid.Anim.PlayMode,
		len(currScene.CullMap),
		len(currScene.Interactive),
		currScene.Cam.P.X,
		currScene.Cam.P.Y,
		game_latency,
		currScene.EventTick,
		currScene.AiTick)

	dbg_TextEl := core.TextEl{
		Font:    assets.Fonts.Default,
		Content: dbg_content,
		Color:   sdl.Color{255, 255, 255, 255},
	}
	dbg_txtr, W, H := dbg_TextEl.Bake(renderer, int(winWidth))
	renderer.Copy(dbg_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{0, winHeight - H, W, H})
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

	Controls.Update(currScene, PCptr, KDs, KUs)

	PC.Solid.UpdateVelocity(&Controls.DPAD)
	PC.Solid.UpdatePCOrientation(Controls)

	if PC.Solid.Anim.PlayMode == 0 {
		currScene.Travel(PCptr)
	}

	if core.IsMoving(PC.Solid.Velocity) && Controls.ACTION_MOD1 > 0 {
		// HACK - Play animation again when running
		PC.Solid.PlayAnimation()
		r := (rand.Int31n(128) * PC.Solid.Position.X * PC.Solid.Position.Y)
		if currScene.EventTick == 2 && r%3 == 0 {
			dust := core.FeetRect(PC.Solid.Position)
			dust.Y += int32(24 * PC.Solid.Orientation.Y)
			dust.X += int32(24 * PC.Solid.Orientation.X)
			currScene.SpawnVFX(dust, PC.Solid.Orientation, game.Puff)
		}

		dpl := (PC.Char.MaxST * 0.0009)
		if PC.Char.CurrentST <= dpl {
			PC.Char.CurrentST = 0
			PC.Solid.Speed = PC.Char.BaseSpeed + PC.Char.SpeedMod
		} else {
			PC.Char.CurrentST -= dpl
		}
	} else {
		if !core.IsMoving(PC.Solid.Velocity) && PC.Char.CurrentST < PC.Char.MaxST {
			PC.Char.CurrentST += (PC.Char.MaxST * 0.001)
		}
	}

	return true
}

func main() {

	runtime.GOMAXPROCS(1)

	var window *sdl.Window
	var renderer *sdl.Renderer

	window, _ = sdl.CreateWindow(
		"Go-SDL2 Obj0",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int(winWidth),
		int(winHeight),
		sdl.WINDOW_SHOWN,
	)

	renderer, _ = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer window.Destroy()
	defer renderer.Destroy()

	assets.Load(renderer)
	game.BootstrapMonsters()
	game.BootstrapItems()
	game.BootstrapVFX()
	game.BootstrapPC(PCptr)

	renderer.SetDrawColor(0, 0, 255, 255)
	currScene = game.InitScene("data/world.tmx", renderer, &PC, winWidth, winHeight, CharSize)

	var running bool = true
	for running {

		then := time.Now()
		running = catchEvents()

		update(currScene)
		render(currScene, renderer)

		game_latency = (time.Since(then) / time.Microsecond)
		sdl.Delay(33 - uint32(game_latency/1000))
	}
}
