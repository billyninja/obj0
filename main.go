package main

import (
	"fmt"
	"github.com/billyninja/obj0/core"
	"github.com/billyninja/obj0/templates"
	"github.com/billyninja/obj0/tmx"
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
	TSzi                int32   = int32(TSz)
	CharSize            int32   = 64
	CharThird                   = CharSize / 3
	winWidth, winHeight int32   = 1280, 720

	KEY_ARROW_UP    = 1073741906
	KEY_ARROW_DOWN  = 1073741905
	KEY_ARROW_LEFT  = 1073741904
	KEY_ARROW_RIGHT = 1073741903
	KEY_LEFT_SHIFT  = 1073742049
	KEY_SPACE_BAR   = 32 // 1073741824 | 32
	KEY_C           = 99
	KEY_X           = 120
	KEY_Z           = 80 // todo

	AI_TICK_LENGTH    = 2
	EVENT_TICK_LENGTH = 3
)

var (
	event        sdl.Event
	font         *ttf.Font
	game_latency time.Duration
	Controls     *ControlState = &ControlState{}

	CL_WHITE        sdl.Color    = sdl.Color{255, 255, 255, 255}
	spritesheetTxt  *sdl.Texture = nil
	particlesTxt    *sdl.Texture = nil
	powerupsTxt     *sdl.Texture = nil
	glowTxt         *sdl.Texture = nil
	monstersTxt     *sdl.Texture = nil
	transparencyTxt *sdl.Texture = nil
	puffTxt         *sdl.Texture = nil
	hitTxt          *sdl.Texture = nil

	SHADOW = &sdl.Rect{320, 224, TSzi, TSzi}

	Cam = Camera{
		P:   core.Vector2d{0, 0},
		DZx: 320,
		DZy: 256,
	}

	PC = core.Char{
		Solid: &core.Solid{
			Velocity:    &core.Vector2d{0, 0},
			Orientation: &core.Vector2d{0, -1},
			Anim:        core.MAN_PU_ANIM,
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
		Inventory: []*core.ItemStack{{templates.GreenBlob, 2}},
	}

	scene     *Scene
	Visual    []*core.VFXInst
	Spawners  []*SpawnPoint
	Monsters  []*core.Char
	GUI       []*sdl.Rect
	CullMap   []*core.Solid
	dbox      DBox  = DBox{BGColor: sdl.Color{90, 90, 90, 255}}
	EventTick uint8 = EVENT_TICK_LENGTH
	AiTick          = AI_TICK_LENGTH
)

type ControlState struct {
	DPAD        core.Vector2d
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
			cs.DPAD.Y -= 1
			break
		case KEY_ARROW_DOWN:
			cs.DPAD.Y += 1
			break
		case KEY_ARROW_LEFT:
			cs.DPAD.X -= 1
			break
		case KEY_ARROW_RIGHT:
			cs.DPAD.X += 1
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
			PC.Solid.Speed = PC.Solid.Speed + PC.SpeedMod
			break
		case KEY_ARROW_DOWN:
			cs.DPAD.Y = 0
			break
		case KEY_ARROW_LEFT:
			cs.DPAD.X = 0
			break
		case KEY_ARROW_RIGHT:
			cs.DPAD.X = 0
			break
		case KEY_Z:
			cs.ACTION_A = 0
			break
		case KEY_X:
			if PC.Solid.Anim.PlayMode != 1 {
				PC.PCMeleeAtk(CullMap, Visual)
			}
			cs.ACTION_B = 0
			break
		case KEY_C:
			if PC.Solid.Anim.PlayMode != 1 {
				PC.CurrentAction = templates.Fira.Cast(&PC, nil)
			}
			cs.ACTION_C = 0
			break
		case KEY_SPACE_BAR:
			if !dbox.NextText() {
				PC.ActProc(CullMap)
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

type Camera struct {
	P   core.Vector2d
	DZx int32
	DZy int32
}

type Scene struct {
	codename    string
	TileSet     *sdl.Texture
	World       [][]*tmx.Terrain
	Interactive []*core.Solid
	CellsX      int32
	CellsY      int32
	StartPoint  core.Vector2d
	CamPoint    core.Vector2d
	tileA       *sdl.Rect
	tileB       *sdl.Rect
}

type Event func(source *core.Solid, subject *core.Solid)

type SpawnPoint struct {
	Position  *sdl.Rect
	Frequency uint16
	LvlMod    uint8
}

type DBox struct {
	SPos     uint8
	CurrText uint8
	Text     []*core.TextEl
	BGColor  sdl.Color
	Char     *core.Char
}

func (db *DBox) LoadText(content []string) {
	db.Text = make([]*core.TextEl, len(content))
	for i, s := range content {
		db.Text[i] = &core.TextEl{
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
	txtr, w, h := ct.Bake(renderer, int(winWidth))
	br := &sdl.Rect{64, winHeight - 128, 512, 120}
	tr := &sdl.Rect{0, 0, w, h}
	bt := &sdl.Rect{64, winHeight - 128, w, h}

	renderer.Copy(transparencyTxt, &sdl.Rect{0, 0, 48, 48}, br)
	renderer.Copy(txtr, tr, bt)
}

func PlayDialog(listener *core.Solid, speaker *core.Solid) {
	if len(speaker.Handlers.DialogScript) > 0 {
		dbox.LoadText(speaker.Handlers.DialogScript)
	}
}

func BashDoor(actor *core.Solid, door *core.Solid) {
	if actor == PC.Solid {
		scene = load_scene(door.Handlers.DoorTo, nil)
	}
}

func (db *DBox) NextText() bool {
	if len(dbox.Text) == 0 {
		return false
	}
	dbox.CurrText += 1
	if int(dbox.CurrText+1) > len(dbox.Text) {
		dbox.Text = []*core.TextEl{}
		dbox.CurrText = 0
	}
	return true
}

func (s *Scene) SolidFromTerrain(terr *tmx.Terrain, cellX int32, cellY int32) {

	pos := &sdl.Rect{cellY * TSzi, cellX * TSzi, TSzi, TSzi}

	for _, tt := range terr.TerrainTypes {

		if tt == nil {
			continue
		}

		var sol *core.Solid

		switch tt.Name {
		case "COLL_BLOCK":
			sol = &core.Solid{
				Position:  pos,
				Collision: 1,
			}
			break
		case "DMG":
			sol = &core.Solid{
				Position:  pos,
				Collision: 0,
				Handlers: &core.InteractionHandlers{
					OnCollDmg: 12,
				},
			}
			break
		}

		if sol != nil {
			s.Interactive = append(s.Interactive, sol)
		}
	}
}

func load_scene(mapname string, renderer *sdl.Renderer) *Scene {

	tmx, wld := tmx.LoadTMXFile(mapname, renderer)

	PC.Solid.Position = V2R(core.Vector2d{2000, 2500}, CharSize, CharSize)
	PC.Solid.Speed = PC.BaseSpeed + PC.SpeedMod

	scn := &Scene{
		codename: mapname,
		TileSet:  tmx.Tilesets[0].Txtr,
		World:    wld,
		CellsX:   tmx.WidthTiles,
		CellsY:   tmx.HeightTiles,
	}

	// TODO func populate2
	var y int32 = 0
	for ; y < scn.CellsY; y++ {
		var x int32 = 0
		for ; x < scn.CellsX; x++ {
			scn.SolidFromTerrain(wld[x][y], x, y)
		}
	}
	//scn.populate(200)

	return scn
}

func (s *Scene) populate(population int) {

	for i := 0; i < population; i++ {
		cX, cY := rand.Int31n(s.CellsX), rand.Int31n(s.CellsY)

		absolute_pos := &sdl.Rect{cX * TSzi, cY * TSzi, CharSize, CharSize}
		sol := &core.Solid{}

		switch rand.Int31n(9) {
		case 3:

			absolute_pos.H, absolute_pos.W = 64, 64
			sol = &core.Solid{
				Position:  absolute_pos,
				Txt:       glowTxt,
				Anim:      core.LIFE_ORB_ANIM,
				Collision: 0,
				Handlers: &core.InteractionHandlers{
					OnCollEvent: core.PickUp,
					OnPickUp: func(healed *core.Solid, orb *core.Solid) {
						if healed.CharPtr != nil {
							healed.CharPtr.CurrentHP += 10
						}
					},
				},
			}

			s.Interactive = append(s.Interactive, sol)
			break
		case 4:

			absolute_pos.W, absolute_pos.H = 128, 128
			for _, sp2 := range Spawners {
				if core.CheckCol(absolute_pos, sp2.Position) {
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

			fnpc := core.Char{
				Solid: &core.Solid{
					Position:    absolute_pos,
					Velocity:    &core.Vector2d{0, 0},
					Orientation: &core.Vector2d{0, 1},
					Speed:       2,
					Txt:         spritesheetTxt,
					Collision:   1,
					Handlers: &core.InteractionHandlers{
						OnActEvent:   PlayDialog,
						DialogScript: []string{"more", "npc", "chitchat"},
					},
					CPattern: 0,
					MPattern: []core.Movement{
						core.Movement{core.F_UP, 90},
						core.Movement{core.F_RIGHT, 10},
						core.Movement{core.F_DOWN, 50},
						core.Movement{core.F_LEFT, 10},
					},
				},
				ActionMap: PC.ActionMap,
				CurrentHP: 9999,
				MaxHP:     9999,
			}
			fnpc.Solid.Anim = core.MAN_PB_ANIM
			fnpc.Solid.CharPtr = &fnpc

			Monsters = append(Monsters, &fnpc)
			break
		}

	}
}

func (s *Scene) _terrainRender2(renderer *sdl.Renderer) {
	var Source *sdl.Rect
	var init int32 = 0

	var offsetX, offsetY int32 = TSzi, TSzi
	// Rendering the terrain
	for winY := init; winY < winHeight; winY += offsetY {
		for winX := init; winX < winWidth; winX += offsetX {

			offsetX = (TSzi - (int32(Cam.P.X)+winX)%TSzi)
			offsetY = (TSzi - (int32(Cam.P.Y)+winY)%TSzi)

			currCellX := (int32(Cam.P.X) + winX) / TSzi
			currCellY := (int32(Cam.P.Y) + winY) / TSzi
			screenPos := sdl.Rect{winX, winY, offsetX, offsetY}

			if currCellX >= s.CellsX || currCellY >= s.CellsY || currCellX < 0 || currCellY < 0 {
				continue
			}
			cell := s.World[currCellY][currCellX]
			if cell.Source == nil {
				continue
			}
			gfx := cell.Source

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
	CullMap = []*core.Solid{}

	for _, obj := range s.Interactive {
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

			if src == nil {
				renderer.DrawRect(scrPos)
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

			if mon.Solid.Chase != nil {
				renderer.SetDrawColor(255, 0, 0, 255)
				renderer.FillRect(&sdl.Rect{scrPos.X, scrPos.Y - 8, 32, 4})
				renderer.SetDrawColor(0, 255, 0, 255)
				renderer.FillRect(&sdl.Rect{scrPos.X, scrPos.Y - 8, int32(32 * calcPerc(mon.CurrentHP, mon.MaxHP) / 100), 4})
			}
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
		counter := core.TextEl{
			Content: fmt.Sprintf("%d", stack.Qty),
			Font:    font,
			Color:   CL_WHITE,
		}

		counterTxtr, cW, cH := counter.Bake(renderer, int(winWidth))
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

	lvl_TextEl := core.TextEl{
		Font:    font,
		Content: fmt.Sprintf("Lvl. %d", PC.Lvl),
		Color:   CL_WHITE,
	}
	lvl_txtr, W, H := lvl_TextEl.Bake(renderer, int(winWidth))
	renderer.Copy(lvl_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{128, 68, W, H})
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
				txtr, w, h := vi.Text.Bake(renderer, int(winWidth))
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

	s._terrainRender2(renderer)
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

func (cam *Camera) RecenterCamera(target *sdl.Rect) {
	newScreenPos := worldToScreen(target, Cam)
	if (cam.DZx - newScreenPos.X) > 0 {
		cam.P.X -= float32(cam.DZx - newScreenPos.X)
	}

	if (winWidth - cam.DZx) < (newScreenPos.X + TSzi) {
		cam.P.X += float32((newScreenPos.X + TSzi) - (winWidth - cam.DZx))
	}

	if (cam.DZy - newScreenPos.Y) > 0 {
		cam.P.Y -= float32(cam.DZy - newScreenPos.Y)
	}

	if (winHeight - cam.DZy) < (newScreenPos.Y + TSzi) {
		cam.P.Y += float32((newScreenPos.Y + TSzi) - (winHeight - cam.DZy))
	}
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

	Cam.RecenterCamera(PC.Solid.Position)

	for _, obj := range CullMap {
		if obj.Position == nil {
			continue
		}
		// Ttl kill
		if obj.Ttl > 0 && obj.Ttl < now {
			obj.Destroy()
			continue
		}

		if obj.CharPtr != nil && obj.CharPtr.CurrentHP <= 0 {
			if obj.CharPtr.Drop != nil {
				PlaceDrop(obj.CharPtr.Drop, obj.Position)
			}
			Visual = append(Visual, templates.Puff.Spawn(&sdl.Rect{obj.Position.X, obj.Position.Y, 92, 92}, nil))
			obj.Destroy()

			PC.CurrentXP += uint16(obj.CharPtr.MaxHP / 10)
			if PC.CurrentXP >= PC.NextLvlXP {
				PC.CurrentXP = 0
				PC.Lvl++
				PC.NextLvlXP = PC.NextLvlXP * uint16(1+PC.Lvl/2)
			}
			continue
		}

		if obj.Anim != nil {
			obj.PlayAnimation()
		}

		if AiTick == 0 {

			if obj.Anim != nil && obj.Anim.PlayMode == 1 {
				continue
			}
			if obj.CharPtr != nil {
				obj.CharPtr.AtkCoolDownC -= obj.CharPtr.AtkSpeed
			}

			if obj.Chase != nil && obj.LoSCheck(PC.Solid) {
				obj.DoChase(CullMap, Visual)
			} else {
				var sp float32 = 2
				if obj.CharPtr != nil {
					sp = obj.Speed
				}
				obj.PeformPattern(sp, CullMap)
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
	dbg_content := fmt.Sprintf(
		"px %d py %d | Cx %.1f Cy %.1f |vx %.1f vy | %.1f (%.1f, %.1f) |"+
			" An:%d/%d/%d cull %d i %d cX %d cY %d L %dus ETick%d AiTick%d",
		PC.Solid.Position.X, PC.Solid.Position.Y, Controls.DPAD.X, Controls.DPAD.Y,
		PC.Solid.Velocity.X, PC.Solid.Velocity.Y, PC.Solid.Orientation.X,
		PC.Solid.Orientation.Y, PC.Solid.Anim.Pose, PC.Solid.Anim.PoseTick,
		PC.Solid.Anim.PlayMode, len(CullMap), len(scene.Interactive), Cam.P.X, Cam.P.Y,
		game_latency, EventTick, AiTick)

	dbg_TextEl := core.TextEl{
		Font:    font,
		Content: dbg_content,
		Color:   CL_WHITE,
	}
	dbg_txtr, W, H := dbg_TextEl.Bake(renderer, int(winWidth))
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

func isMoving(vel *core.Vector2d) bool {
	return (vel.X != 0 || vel.Y != 0)
}

func PlaceDrop(item *core.Item, origin *sdl.Rect) {
	instance := core.ItemInstance{
		ItemTpl: item,
		Solid: &core.Solid{
			ItemPtr: item,
			Txt:     item.Txtr,
			Source:  item.Source,
			Position: &sdl.Rect{
				origin.X,
				origin.Y,
				item.Source.W,
				item.Source.H,
			},
			Handlers: &core.InteractionHandlers{
				OnActEvent: core.PickUp,
				OnPickUp:   core.AddToInv,
			},
		},
	}

	scene.Interactive = append(scene.Interactive, instance.Solid)
}

func inScreen(r *sdl.Rect) bool {
	return (r.X > (r.W*-1) && r.X < winWidth && r.Y > (r.H*-1) && r.Y < winHeight)
}

func V2R(v core.Vector2d, w int32, h int32) *sdl.Rect {
	return &sdl.Rect{int32(v.X), int32(v.Y), w, h}
}

func (sp *SpawnPoint) Produce() {
	px := float32(rand.Int31n((sp.Position.X+sp.Position.W)-sp.Position.X) + sp.Position.X)
	py := float32(rand.Int31n((sp.Position.Y+sp.Position.H)-sp.Position.Y) + sp.Position.Y)
	mon := MonsterFactory(
		&templates.OrcTPL,
		sp.LvlMod,
		core.Vector2d{px, py},
	)

	Monsters = append(Monsters, mon)
}

func (scn *Scene) SpawnVFX(pos *sdl.Rect, o *core.Vector2d, vfx *core.VFX) {
	vi := vfx.Spawn(pos, o)
	Visual = append(Visual, vi)
}

func MonsterFactory(monsterTpl *templates.MonsterTemplate, lvlMod uint8, pos core.Vector2d) *core.Char {

	variance := uint8(math.Floor(float64(rand.Float32() * monsterTpl.LvlVariance * 100)))
	lvl := uint8((monsterTpl.Lvl + lvlMod) + variance)
	hp := monsterTpl.HP + float32(lvl*2)
	sizeMod := int32(float32(lvl-monsterTpl.Lvl) * monsterTpl.ScalingFactor)
	W, H := (monsterTpl.Size + sizeMod), (monsterTpl.Size + sizeMod)

	var DropItem *core.Item
	var sumP float32
	R := rand.Float32()
	for _, l := range monsterTpl.Loot {
		sumP += l.Perc
		if R < sumP {
			DropItem = l.Item
			break
		}
	}

	mon := core.Char{
		Lvl: lvl,
		Solid: &core.Solid{
			Position:    &sdl.Rect{int32(pos.X), int32(pos.Y), W, H},
			Velocity:    &core.Vector2d{0, 0},
			Orientation: &core.Vector2d{0, 0},
			Txt:         monsterTpl.Txtr,
			Speed:       1,
			Collision:   2,
			Handlers: &core.InteractionHandlers{
				OnCollDmg: 12,
			},
			CPattern: 0,
			LoS:      monsterTpl.LoS,
			MPattern: []core.Movement{
				core.Movement{core.F_DOWN, 50},
				core.Movement{core.F_UP, 90},
				core.Movement{core.F_RIGHT, 10},
				core.Movement{core.F_LEFT, 10},
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

	Visual = append(Visual, templates.Puff.Spawn(&sdl.Rect{int32(pos.X), int32(pos.Y), 92, 92}, nil))

	return &mon
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
	PC.Solid.UpdateVelocity(&Controls.DPAD)

	if Controls.DPAD.X != 0 {
		PC.Solid.Orientation.X = core.ThrotleValue(Controls.DPAD.X, 1)
	} else {
		if core.Abs32(Controls.DPAD.Y) > 1 {
			PC.Solid.Orientation.X -= 1
			if PC.Solid.Orientation.X < 0 {
				PC.Solid.Orientation.X = 0
			}
		}
	}
	if Controls.DPAD.Y != 0 {
		PC.Solid.Orientation.Y = core.ThrotleValue(Controls.DPAD.Y, 1)
	} else {
		if core.Abs32(Controls.DPAD.X) > 1 {
			PC.Solid.Orientation.Y -= 1
			if PC.Solid.Orientation.Y < 0 {
				PC.Solid.Orientation.Y = 0
			}
		}
	}

	if PC.Solid.Anim.PlayMode == 0 {
		PC.Solid.ProcMovement(CullMap, 3000, 3000)
	}

	if isMoving(PC.Solid.Velocity) && Controls.ACTION_MOD1 > 0 {
		// HACK - Play animation again when running
		PC.Solid.PlayAnimation()
		r := (rand.Int31n(128) * PC.Solid.Position.X * PC.Solid.Position.Y)
		if EventTick == 2 && r%3 == 0 {
			dust := core.FeetRect(PC.Solid.Position)
			dust.Y += int32(24 * PC.Solid.Orientation.Y)
			dust.X += int32(24 * PC.Solid.Orientation.X)
			Visual = append(Visual, templates.Puff.Spawn(dust, nil))
		}

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
	font, _ = ttf.OpenFont("data/assets/textures/PressStart2P.ttf", 12)

	spritesheetImg, _ := img.Load("data/assets/textures/main_char.png")
	powerupsImg, _ := img.Load("data/assets/textures/powerups_ts.png")
	glowImg, _ := img.Load("data/assets/textures/glowing_ts.png")
	monstersImg, _ := img.Load("data/assets/textures/monsters.png")
	transparencyImg, _ := img.Load("data/assets/textures/transparency.png")
	puffImg, _ := img.Load("data/assets/textures/puff.png")
	hitImg, _ := img.Load("data/assets/textures/hit.png")
	defer monstersImg.Free()
	defer spritesheetImg.Free()
	defer powerupsImg.Free()
	defer glowImg.Free()
	defer transparencyImg.Free()
	defer puffImg.Free()
	defer hitImg.Free()

	spritesheetTxt, _ = renderer.CreateTextureFromSurface(spritesheetImg)
	powerupsTxt, _ = renderer.CreateTextureFromSurface(powerupsImg)
	glowTxt, _ = renderer.CreateTextureFromSurface(glowImg)
	monstersTxt, _ = renderer.CreateTextureFromSurface(monstersImg)
	transparencyTxt, _ = renderer.CreateTextureFromSurface(transparencyImg)
	puffTxt, _ = renderer.CreateTextureFromSurface(puffImg)
	hitTxt, _ = renderer.CreateTextureFromSurface(hitImg)
	defer spritesheetTxt.Destroy()
	defer powerupsTxt.Destroy()
	defer glowTxt.Destroy()
	defer monstersTxt.Destroy()
	defer transparencyTxt.Destroy()
	defer puffTxt.Destroy()
	defer hitTxt.Destroy()

	templates.BootstrapMonsters(monstersTxt)
	templates.BootstrapItems(powerupsTxt)
	templates.BootstrapVfx(hitTxt, puffTxt)
	templates.BootstrapSpells(glowTxt, core.LIFE_ORB_ANIM)

	core.BootstrapResources(glowTxt, templates.Hit, templates.Impact)

	MainCharSS := &core.SpriteSheet{spritesheetTxt, 0, 0, 32, 32}
	MainCharActionMap := MainCharSS.BuildBasicActions(3, true)
	PC.ActionMap = MainCharActionMap
	PC.Solid.CharPtr = &PC

	renderer.SetDrawColor(0, 0, 255, 255)
	scene = load_scene("data/world.tmx", renderer)
	dbox.LoadText([]string{"Hello World!", "Again!"})

	var running bool = true
	for running {

		then := time.Now()

		running = catchEvents()
		scene.update()
		scene.ProcActions()
		scene.render(renderer)

		game_latency = (time.Since(then) / time.Microsecond)
		sdl.Delay(33 - uint32(game_latency/1000))
	}
}

func (scene *Scene) RunningActions() []core.ActionInterface {
	queue := []core.ActionInterface{}
	if PC.CurrentAction != nil {
		queue = append(queue, PC.CurrentAction)
	}

	return queue
}

func (scn *Scene) ProcActions() {

	for _, action := range scn.RunningActions() {

		actor := action.GetActor()

		switch action.GetState() {

		case -1:
			{
				continue
			}
			break

		case 0:
			{
				vfx := action.PreActionVFX()
				if vfx != nil {
					scn.SpawnVFX(actor.Position, actor.Orientation, vfx)
				}

				anim := action.PreActionAnim()
				if anim != nil {
					actor.SetAnimation(anim, action)
				} else {
					action.Step()
				}

				// important always should go to the intermediary step, at least!
				action.Step()
			}
			break

		case 1:
			{
				// TODO 1 - More explicit control like. action.IsWaiting()
				// TODO 2 - Timing Control!
				if !(actor.Anim != nil && actor.Anim.PlayMode == 1) {
					action.Step()
				}
			}
			break

		case 2:
			{
				vfx := action.ActionVFX()
				if vfx != nil {
					scn.SpawnVFX(actor.Position, actor.Orientation, vfx)
				}

				anim := action.ActionAnim()
				if anim != nil {
					actor.SetAnimation(anim, action)
				} else {
					action.Step()
				}

				output := action.Perform(actor, nil)
				scn.Interactive = append(scn.Interactive, output...)

				// important always should go to the intermediary step, at least!
				action.Step()
			}
			break

		case 3:
			{
				// Same observations as `Case 1`
				if !(actor.Anim != nil && actor.Anim.PlayMode == 1) {
					action.Step()
				}
			}
			break

		case 4:
			{
				vfx := action.PostActionVFX()
				if vfx != nil {
					scn.SpawnVFX(actor.Position, actor.Orientation, vfx)
				}

				anim := action.PostActionAnim()
				if anim != nil {
					actor.SetAnimation(anim, action)
				} else {
					action.Step()
				}

				action.SetFinished()
			}
			break

		}
	}
}
