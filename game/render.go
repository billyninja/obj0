package game

import (
	"fmt"
	"github.com/billyninja/obj0/assets"
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

var (
	CL_WHITE          sdl.Color = sdl.Color{255, 255, 255, 255}
	CL_BLACK          sdl.Color = sdl.Color{0, 0, 0, 255}
	CL_LIGHT_PURPLE   sdl.Color = sdl.Color{151, 126, 254, 255}
	CL_AQUA_GREEN     sdl.Color = sdl.Color{131, 225, 191, 255}
	CL_LEATHER_BROWN  sdl.Color = sdl.Color{121, 59, 22, 255}
	CL_JUPITER_ORANGE sdl.Color = sdl.Color{255, 148, 0, 255}
	CL_LEAD_GRAY      sdl.Color = sdl.Color{70, 70, 70, 255}
)

func (s *Scene) TerrainRender(renderer *sdl.Renderer) {

	core.SetColor(renderer, CL_BLACK)

	var Source *sdl.Rect
	var init int32 = 0

	var offsetX, offsetY int32 = s.TileWidth, s.TileHeight
	var TSzi = s.TileWidth

	// Rendering the terrain
	for winY := init; winY < s.WinHeight; winY += offsetY {
		for winX := init; winX < s.WinWidth; winX += offsetX {

			offsetX = (TSzi - (int32(s.Cam.P.X)+winX)%TSzi)
			offsetY = (TSzi - (int32(s.Cam.P.Y)+winY)%TSzi)

			currCellX := (int32(s.Cam.P.X) + winX) / TSzi
			currCellY := (int32(s.Cam.P.Y) + winY) / TSzi
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
				renderer.Copy(s.TileSet, Source, &screenPos)
			}
		}
	}
}

func (s *Scene) SolidsRender(renderer *sdl.Renderer) {
	s.CullMap = []*SceneEntity{}

	for _, se := range s.Interactive {
		obj := se.Solid
		if obj.Position == nil {
			continue
		}

		scrPos := s.Cam.WorldToScreen(obj.Position)

		if s.InScreen(scrPos) {
			var src *sdl.Rect

			if obj.Anim != nil {
				src = obj.Anim.Action[obj.Anim.Pose]
			} else {
				src = obj.Source
			}

			if src != nil {
				renderer.Copy(obj.Txt, src, scrPos)
			}

			s.CullMap = append(s.CullMap, se)
		}
	}
}

func (s *Scene) MonstersRender(renderer *sdl.Renderer) {

	for _, mon := range s.Monsters {
		if mon.Solid == nil || mon.Solid.Anim == nil || mon.Solid.Position == nil {
			continue
		}
		scrPos := s.Cam.WorldToScreen(mon.Solid.Position)

		if s.InScreen(scrPos) {

			src := mon.Solid.Anim.Action[mon.Solid.Anim.Pose]
			renderer.Copy(assets.Textures.Sprites.MainChar, SHADOW, scrPos)
			scrPos.Y -= 12
			renderer.Copy(mon.Solid.Txt, src, scrPos)

			s.CullMap = append(s.CullMap, mon)

			// HP bar
			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.FillRect(&sdl.Rect{scrPos.X, scrPos.Y - 8, 32, 4})
			renderer.SetDrawColor(0, 255, 0, 255)
			renderer.FillRect(&sdl.Rect{scrPos.X, scrPos.Y - 8, int32(32 * CalcPerc(mon.Char.CurrentHP, mon.Char.MaxHP) / 100), 4})
		}
	}
}

func (s *Scene) GUIRender(pc *Char, renderer *sdl.Renderer) {

	// Gray overlay
	renderer.Copy(
		assets.Textures.GUI.Transparency,
		&sdl.Rect{0, 0, 48, 48},
		&sdl.Rect{0, 0, 120, 48},
	)
	renderer.Copy(
		assets.Textures.GUI.Transparency,
		&sdl.Rect{0, 0, 48, 48},
		&sdl.Rect{0, 48, 210, 28},
	)

	// HEALTH BAR
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.FillRect(&sdl.Rect{10, 10, 100, 4})
	renderer.SetDrawColor(0, 255, 0, 255)
	renderer.FillRect(&sdl.Rect{10, 10, int32(CalcPerc(pc.CurrentHP, pc.MaxHP)), 4})

	// MANA BAR
	renderer.SetDrawColor(190, 0, 120, 255)
	renderer.FillRect(&sdl.Rect{10, 24, 100, 4})
	renderer.SetDrawColor(0, 0, 255, 255)
	renderer.FillRect(&sdl.Rect{10, 24, int32(CalcPerc(pc.CurrentST, pc.MaxST)), 4})

	// XP BAR
	renderer.SetDrawColor(90, 90, 0, 255)
	renderer.SetDrawColor(190, 190, 0, 255)
	renderer.FillRect(&sdl.Rect{10, 38, int32(CalcPerc(float32(pc.CurrentXP), float32(pc.NextLvlXP))), 4})

	for i, stack := range pc.Inventory {
		counter := core.TextEl{
			Content: fmt.Sprintf("%d", stack.Qty),
			Font:    assets.Fonts.Default,
			Color:   CL_WHITE,
		}

		counterTxtr, cW, cH := counter.Bake(renderer, int(s.WinWidth))
		pos := sdl.Rect{8 + (int32(i) * 32), 48, 24, 24}
		renderer.Copy(stack.ItemTpl.Txtr, stack.ItemTpl.Source, &pos)
		pos.Y += 16
		pos.X += 16
		pos.W = cW
		pos.H = cH
		renderer.Copy(counterTxtr, &sdl.Rect{0, 0, cW, cH}, &pos)
	}

	for _, el := range s.GUI {
		scrPos := s.Cam.WorldToScreen(el)
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.DrawRect(scrPos)
	}

	lvl_TextEl := core.TextEl{
		Font:    assets.Fonts.Default,
		Content: fmt.Sprintf("Lvl. %d", pc.Lvl),
		Color:   CL_WHITE,
	}
	lvl_txtr, W, H := lvl_TextEl.Bake(renderer, int(s.WinWidth))
	renderer.Copy(lvl_txtr, &sdl.Rect{0, 0, W, H}, &sdl.Rect{128, 64, W, H})

	s.DBox.Present(s.WinWidth, s.WinHeight, renderer)

	for _, spw := range s.Spawners {
		renderer.DrawRect(s.Cam.WorldToScreen(spw.Position))
	}
	FW := &sdl.Rect{0, 0, 1280, 720}
	UpdateInventoryMenu(ROOT_MENU.SubMenus[0], pc.Inventory)
	ROOT_MENU.Render(renderer, FW)
}

func (s *Scene) VFXRender(Els []*VFXInst, renderer *sdl.Renderer) {
	for _, vi := range Els {
		if vi.Pos == nil {
			continue
		}
		scrp := s.Cam.WorldToScreen(vi.Pos)
		if s.InScreen(scrp) {
			if vi.Text != nil {
				txtr, w, h := vi.Text.Bake(renderer, int(s.WinWidth))
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

func (s *Scene) PCRender(pc *SceneEntity, renderer *sdl.Renderer) {
	scrPos := s.Cam.WorldToScreen(pc.Solid.Position)

	if !(pc.Char.Invinc > 0 && time.Now().Unix()%2 == 3) {
		renderer.Copy(assets.Textures.Sprites.MainChar, pc.Solid.Anim.Action[pc.Solid.Anim.Pose], scrPos)
		scrPos.Y += 12
		renderer.Copy(assets.Textures.Sprites.MainChar, SHADOW, scrPos)
	}
}

func (scn *Scene) ProjectilesRender(renderer *sdl.Renderer) {
	for _, prj := range scn.Projectiles {
		sol := prj.Solid
		if sol == nil || sol.Position == nil {
			continue
		}

		scrp := scn.Cam.WorldToScreen(sol.Position)
		if scn.InScreen(scrp) {
			frame := sol.Anim.Action[sol.Anim.Pose]
			renderer.Copy(sol.Txt, frame, scrp)
		}
	}
}

func (scn *Scene) PostEffectRender(renderer *sdl.Renderer) {
	renderer.Copy(assets.Textures.GUI.Transparency, &sdl.Rect{0, 0, 48, 48}, &sdl.Rect{0, 0, scn.WinWidth, scn.WinHeight})
}
