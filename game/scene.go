package game

import (
	"github.com/billyninja/obj0/assets"
	"github.com/billyninja/obj0/core"
	"github.com/billyninja/obj0/tmx"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"time"
)

func InitScene(mapname string, renderer *sdl.Renderer, pc *SceneEntity, winWidth, winHeight, charSize int32) *Scene {

	tmx, wld := tmx.LoadTMXFile(mapname, renderer)

	pc.Solid.Position = &sdl.Rect{500, 600, charSize, charSize}
	pc.Solid.Speed = pc.Char.BaseSpeed + pc.Char.SpeedMod
	cam := &Camera{
		P: core.Vector2d{
			float32(pc.Solid.Position.X - 600),
			float32(pc.Solid.Position.Y - 300),
		},
		DZx: 256,
		DZy: 144,
	}

	scn := &Scene{
		Codename:        mapname,
		Cam:             cam,
		Renderer:        renderer,
		TileSet:         tmx.Tilesets[0].Txtr,
		World:           wld,
		CellsX:          tmx.WidthTiles,
		CellsY:          tmx.HeightTiles,
		TileWidth:       tmx.TileW,
		TileHeight:      tmx.TileH,
		LimitW:          tmx.WidthTiles * tmx.TileW,
		LimitH:          tmx.HeightTiles * tmx.TileH,
		WinWidth:        winWidth,
		WinHeight:       winHeight,
		DBox:            &DBox{BGColor: sdl.Color{90, 90, 90, 255}},
		EventTick:       2,
		EventTickLength: 2,
		AiTick:          2,
		AiTickLength:    2,
	}

	scn.DBox.LoadText([]string{"Hello World!", "Again!"})

	var y int32 = 0
	for ; y < scn.CellsY; y++ {
		var x int32 = 0
		for ; x < scn.CellsX; x++ {
			scn.SolidFromTerrain(wld[x][y], x, y)
		}
	}
	scn.Populate(600, pc.Char)

	return scn
}

func (scn *Scene) SolidFromTerrain(terr *tmx.Terrain, cellX int32, cellY int32) {

	pos := &sdl.Rect{cellY * scn.TileWidth, cellX * scn.TileWidth, scn.TileWidth, scn.TileWidth}

	sE := &SceneEntity{
		Solid: &Solid{
			Position: pos,
		},
		Handlers: &SEventHandlers{},
	}

	for _, tt := range terr.TerrainTypes {

		if tt == nil {
			continue
		}

		switch tt.Name {
		case "COLL_BLOCK":
			sE.Solid.Collision = 1
			break
		case "DMG":
			sE.Handlers.OnCollDmg = 12
			sE.Handlers.OnCollPushBack = 12
			break
		case "DOOR":
			sE.Handlers.DoorTo = "cave.tmx"
			sE.Handlers.OnActEvent = OpenDoor
			break
		}

		if sE != nil {
			scn.Interactive = append(scn.Interactive, sE)
		}
	}
}

func (scn *Scene) Populate(population int, pc *Char) {

	for i := 0; i < population; i++ {
		cX, cY := rand.Int31n(scn.CellsX), rand.Int31n(scn.CellsY)

		absolute_pos := &sdl.Rect{cX * scn.TileWidth, cY * scn.TileWidth, 64, 64}

		switch rand.Int31n(9) {
		case 3:

			se := &SceneEntity{
				Solid: &Solid{
					Position:  absolute_pos,
					Txt:       assets.Textures.Sprites.Glow,
					Anim:      LIFE_ORB_ANIM,
					Collision: 0,
				},
				Handlers: &SEventHandlers{
					OnCollEvent: func(src, sbj *SceneEntity, scn *Scene) {
						char := src.Char

						src.Solid.SetAnimation(MAN_PU_ANIM)
						AddToInv(src, sbj, scn)

						if char != nil {
							char.CurrentHP += 10
						}
						sbj.Solid.Destroy()
					},
				},
			}

			scn.Interactive = append(scn.Interactive, se)
			break
		case 4:

			absolute_pos.W, absolute_pos.H = 128, 128
			for _, sp2 := range scn.Spawners {
				if core.CheckCol(absolute_pos, sp2.Position) {
					return
				}
			}

			rand.Seed(int64(time.Now().Nanosecond()))
			spw := &SpawnPoint{
				Position:  absolute_pos,
				Frequency: uint16(rand.Int31n(5)),
			}

			scn.Spawners = append(scn.Spawners, spw)
			break
		case 5:
			fnpc := SceneEntity{
				Char: &Char{
					ActionMap: pc.ActionMap,
					CurrentHP: 9999,
					MaxHP:     9999,
				},
				Handlers: &SEventHandlers{
					OnActEvent:   PlayDialog,
					DialogScript: []string{"...more", "npc", "chitchat"},
				},
				Solid: &Solid{
					Position:    absolute_pos,
					Velocity:    &core.Vector2d{0, 0},
					Orientation: &core.Vector2d{0, 1},
					Speed:       2,
					Txt:         assets.Textures.Sprites.MainChar,
					Collision:   1,
					CPattern:    0,
					MPattern: []Movement{
						Movement{F_UP, 90},
						Movement{F_RIGHT, 10},
						Movement{F_DOWN, 50},
						Movement{F_LEFT, 10},
					},
				},
			}
			fnpc.Solid.Anim = MAN_PB_ANIM
			fnpc.Solid.CharPtr = fnpc.Char

			scn.Monsters = append(scn.Monsters, &fnpc)
			break
		}

	}
}

func (scn *Scene) SpawnMonster(spw *SpawnPoint, pc *SceneEntity) {
	mon, vi := spw.Produce(pc)
	scn.Monsters = append(scn.Monsters, mon)
	scn.VisualLower = append(scn.VisualLower, vi)
}

func (scn *Scene) ResolveCol(Moving *SceneEntity, Still *SceneEntity) bool {
	var halt bool
	if Still.Solid.Collision == 1 {
		halt = true
	}

	if Moving.Handlers != nil {
		if Still.Char != nil {
			if Moving.Handlers.OnCollDmg > 0 {
				Still.Char.DepletHP(Still.Handlers.OnCollDmg)
			}
			if Moving.Handlers.OnCollPushBack > 0 {
				Still.Solid.PushBack(Moving.Handlers.OnCollPushBack, Moving.Solid.Orientation)
			}
		}

		if Moving.Handlers.OnCollEvent != nil {
			Moving.Handlers.OnCollEvent(Moving, Still, scn)
		}
	}

	if Still != nil && Still.Handlers != nil {
		if Moving.Char != nil {
			if Still.Handlers.OnCollDmg > 0 {
				Moving.Char.DepletHP(Still.Handlers.OnCollDmg)
			}
			if Still.Handlers.OnCollPushBack > 0 {
				Moving.Solid.PushBack(Still.Handlers.OnCollPushBack, nil)
			}
		}

		if Still.Handlers.OnCollEvent != nil {
			Still.Handlers.OnCollEvent(Moving, Still, scn)
		}
	}

	if Moving != nil && Moving.Handlers != nil && Moving.Handlers.OnCollEvent != nil {
		Moving.Handlers.OnCollEvent(Moving, Still, scn)
	}

	return halt
}

func (scn *Scene) ActProc(source *SceneEntity) {
	sol := source.Solid
	action_hit_box := core.ProjectHitBox(
		core.Center(sol.Position), sol.Orientation, 32, nil, 1)

	for _, se := range scn.CullMap {
		if se.Handlers != nil &&
			se.Handlers.OnActEvent != nil &&
			core.CheckCol(action_hit_box, se.Solid.Position) {
			se.Handlers.OnActEvent(source, se, scn)
			return
		}
	}
}

func (scn *Scene) Travel(traveler *SceneEntity) {
	s := traveler.Solid
	np := s.ApplyMovement()

	outbound := (np.X <= 0 || np.Y <= 0 || np.X > scn.LimitW || np.Y > scn.LimitH)

	if (np.X == s.Position.X && np.Y == s.Position.Y) || outbound {
		return
	}

	fr := core.FeetRect(np)
	for _, se := range scn.CullMap {
		sol := se.Solid
		if sol == nil || sol == s || sol.Position == nil {
			continue
		}
		if core.CheckCol(fr, sol.Position) && scn.ResolveCol(traveler, se) {
			return
		}
	}

	if s.CharPtr != nil {
		anim := s.CharPtr.GetFacingAnim(s.Orientation)
		if anim != nil && s.Anim != nil && s.Anim.PlayMode != 1 {
			s.Anim.PlayMode = anim.PlayMode
			s.Anim.Action = anim.Action
		}
	}

	s.Position = np
}

func (scn *Scene) DoChase(chaser *SceneEntity) {
	sol := chaser.Solid

	sol.Velocity.X = 0
	sol.Velocity.Y = 0

	tgt := sol.Chase

	diffX := core.Abs32(float32(sol.Position.X - tgt.Position.X))
	diffY := core.Abs32(float32(sol.Position.Y - tgt.Position.Y))

	if sol.Position.X > tgt.Position.X {
		sol.Velocity.X = -1
	}

	if sol.Position.X < tgt.Position.X {
		sol.Velocity.X = 1
	}
	if sol.Position.Y > tgt.Position.Y {
		sol.Velocity.Y = -1
	}

	if sol.Position.Y < tgt.Position.Y {
		sol.Velocity.Y = 1
	}
	chr := chaser.Char
	if diffX < 80 && diffY < 80 && chr != nil {
		if chr.AtkCoolDownC <= 0 {
			chr.AtkCoolDownC += chr.AtkCoolDown
			MeleeAttack(chaser, nil, scn)
		}
		return
	}

	scn.Travel(chaser)
	return
}

func (scn *Scene) PeformPattern(se *SceneEntity) {
	sol := se.Solid
	mov := sol.PatternStep()

	if mov != nil && sol.Position != nil {
		sol.Orientation = &mov.Orientation
		sol.Velocity = &mov.Orientation
		scn.Travel(se)
		sol.CPattern += uint32(sol.Speed)
	}
}

func (scn *Scene) SpawnVFX(pos *sdl.Rect, o *core.Vector2d, vfx *VFX, layer uint8) {
	vi := vfx.Spawn(pos, o)
	if layer == 0 {
		scn.VisualLower = append(scn.VisualLower, vi)
	} else {
		scn.VisualUpper = append(scn.VisualUpper, vi)
	}
}

func (scn *Scene) UpdateVFX(now int64) {
	upd := func(vi *VFXInst, now int64) {
		if vi.Ttl > 0 && vi.Ttl < now {
			vi.Destroy()
			return
		}
		vi.UpdateAnim()
	}
	for _, vi := range scn.VisualLower {
		upd(vi, now)
	}
	for _, vi := range scn.VisualUpper {
		upd(vi, now)
	}
}

func (scn *Scene) UpdateProjectiles(now int64) {
	for _, prj := range scn.Projectiles {
		sol := prj.Solid
		if sol == nil || sol.Position == nil {
			continue
		}
		if sol.Ttl > 0 && sol.Ttl < now {
			prj.Destroy()
			continue
		}
		sol.PlayAnimation()
		sol.Position = sol.ApplyMovement()
		if scn.EventTick == 0 {
			for _, mon := range scn.CullMap {
				if (mon.Solid.Collision == 1 && mon.Char != nil) && core.CheckCol(sol.Position, mon.Solid.Position) {
					scn.ResolveCol(prj, mon)
					break
				}
			}
		}
	}
}

func (scn *Scene) PlaceDrop(item *Item, origin *sdl.Rect) {
	instance := &SceneEntity{
		ItemPtr: item,
		Solid: &Solid{
			Txt:    item.Txtr,
			Source: item.Source,
			Position: &sdl.Rect{
				origin.X,
				origin.Y,
				item.Source.W,
				item.Source.H,
			},
		},
		Handlers: &SEventHandlers{
			OnActEvent: Pickup,
		},
	}

	scn.Interactive = append(scn.Interactive, instance)
}

func (sp *SpawnPoint) Produce(pc *SceneEntity) (*SceneEntity, *VFXInst) {
	px := float32(rand.Int31n((sp.Position.X+sp.Position.W)-sp.Position.X) + sp.Position.X)
	py := float32(rand.Int31n((sp.Position.Y+sp.Position.H)-sp.Position.Y) + sp.Position.Y)
	mon := OrcTPL.MonsterFactory(sp.LvlMod, core.Vector2d{px, py}, pc.Solid)
	vi := Puff.Spawn(mon.Solid.Position, mon.Solid.Orientation)

	return mon, vi
}
