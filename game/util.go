package game

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

func (cam *Camera) WorldToScreen(pos *sdl.Rect) *sdl.Rect {
	return &sdl.Rect{
		(pos.X - int32(cam.P.X)),
		(pos.Y - int32(cam.P.Y)),
		pos.W, pos.H,
	}
}

func (scn *Scene) InScreen(r *sdl.Rect) bool {
	return (r.X > (r.W*-1) && r.X < scn.WinWidth && r.Y > (r.H*-1) && r.Y < scn.WinHeight)
}

func CalcPerc(v1 float32, v2 float32) float32 {
	return ((v1 / v2) * 100)
}

func (scn *Scene) Recenter(target *sdl.Rect) {
	newScreenPos := scn.Cam.WorldToScreen(target)

	if (scn.Cam.DZx - newScreenPos.X) > 0 {
		scn.Cam.P.X -= float32(scn.Cam.DZx - newScreenPos.X)
	}

	if (scn.WinWidth - scn.Cam.DZx) < (newScreenPos.X + scn.TileWidth) {
		scn.Cam.P.X += float32((newScreenPos.X + scn.TileWidth) - (scn.WinWidth - scn.Cam.DZx))
	}

	if (scn.Cam.DZy - newScreenPos.Y) > 0 {
		scn.Cam.P.Y -= float32(scn.Cam.DZy - newScreenPos.Y)
	}

	if (scn.WinHeight - scn.Cam.DZy) < (newScreenPos.Y + scn.TileWidth) {
		scn.Cam.P.Y += float32((newScreenPos.Y + scn.TileWidth) - (scn.WinHeight - scn.Cam.DZy))
	}
}

func FactoryProjectileSolid(p *sdl.Rect, o *core.Vector2d, t *sdl.Texture, a *Animation, ttl time.Duration, spd float32, sz *core.Vector2d) *Solid {

	pos := &sdl.Rect{
		p.X + p.W*int32(o.X),
		p.Y + p.H*int32(o.Y),
		p.W,
		p.H,
	}

	sol := &Solid{
		Position:    pos,
		Orientation: &core.Vector2d{o.X, o.Y},
		Velocity:    &core.Vector2d{o.X, o.Y},
		Txt:         t,
		Anim:        a,
		Speed:       spd,
		Collision:   1,
	}

	if ttl > 0 {
		sol.Ttl = time.Now().Add(ttl).Unix()
	}

	return sol
}

func FactoryProjectileEntity(sol *Solid, collDmg float32, collPushBack int32, explRadius, explDmg, explPushBack float32) *SceneEntity {
	han := &SEventHandlers{
		OnCollDmg:      collDmg,
		OnCollPushBack: collPushBack,
	}
	if explDmg > 0 || explPushBack > 0 {
		han.OnCollEvent = Explode
	}

	return &SceneEntity{
		Solid:    sol,
		Handlers: han,
	}
}

func (scn *Scene) Destroy() {
	scn.GameState = nil
	scn.TileSet = nil
	scn.Cam = nil
	scn.Renderer = nil

	scn.Interactive = []*SceneEntity{}
	scn.Monsters = []*SceneEntity{}
	scn.CullMap = []*SceneEntity{}
	scn.Projectiles = []*SceneEntity{}
	scn.GUI = []*sdl.Rect{}
	scn.DBox = nil
}

func SceneTransition(doorTo string, pc *SceneEntity, scn *Scene) {
	gs := scn.GameState
	rdr := scn.Renderer
	w := scn.WinWidth
	h := scn.WinHeight
	scn.Destroy()
	newScene := InitScene(
		"data/"+doorTo,
		rdr,
		pc,
		w,
		h,
		64)
	gs.CurrentScene = newScene
	gs.CurrentScene.GameState = gs
}
