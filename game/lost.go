package game

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

func (ch *Char) DepletHP(dmg float32) {
	if ch.Invinc > 0 {
		return
	}
	if dmg > ch.CurrentHP {
		ch.CurrentHP = 0
	} else {
		ch.CurrentHP -= dmg
	}
}

func (ch *Char) ApplyInvinc() {
	if ch.Invinc == 0 {
		ch.Invinc = time.Now().Add(400 * time.Millisecond).Unix()
	}
}

func (sol *Solid) LoSCheck(tgt *Solid) bool {
	LoS := &sdl.Rect{
		sol.Position.X - sol.LoS,
		sol.Position.Y - sol.LoS,
		sol.Position.W + (sol.LoS * 2),
		sol.Position.H + (sol.LoS * 2),
	}

	if !core.CheckCol(tgt.Position, LoS) {
		return false
	}
	return true
}

func (sol *Solid) ApplyMovement() *sdl.Rect {
	speed := sol.Speed

	if sol.Velocity.X == 0 && sol.Velocity.Y == 0 {
		return sol.Position
	}

	// Decreasing speed for diagonal movement
	if sol.Velocity.X != 0 && sol.Velocity.Y != 0 {
		speed -= 0.5
		if speed < 1 {
			speed = 1
		}
	}

	return &sdl.Rect{
		(sol.Position.X + int32(sol.Velocity.X*speed)),
		(sol.Position.Y + int32(sol.Velocity.Y*speed)),
		sol.Position.W,
		sol.Position.H,
	}
}

func (sol *Solid) UpdateVelocity(dpad *core.Vector2d) {

	nv := &core.Vector2d{}
	*nv = *sol.Velocity

	if dpad.X != 0 || sol.Velocity.X != 0 {
		if dpad.X != 0 {
			nv.X = core.ThrotleValue(nv.X+dpad.X, 2)
		} else {
			nv.X = (core.Abs32(nv.X) - 1) * sol.Orientation.X
		}
	}

	if dpad.Y != 0 || nv.Y != 0 {
		if dpad.Y != 0 {
			nv.Y = core.ThrotleValue(nv.Y+dpad.Y, 2)
		} else {
			nv.Y = (core.Abs32(nv.Y) - 1) * sol.Orientation.Y
		}
	}

	*sol.Velocity = *nv
}

func (sol *Solid) PatternStep() *Movement {
	var sum uint32 = 0

	for _, mp := range sol.MPattern {
		sum += uint32(mp.Ticks)
		if sum > sol.CPattern {
			return &mp
		}
	}
	sol.CPattern = 0

	return nil
}

func (sol *Solid) UpdatePCOrientation(ctrl *ControlState) {

	if ctrl.DPAD.X != 0 {
		sol.Orientation.X = core.ThrotleValue(ctrl.DPAD.X, 1)
	} else {
		if core.Abs32(ctrl.DPAD.Y) > 1 {
			sol.Orientation.X = 0
		}
	}

	if ctrl.DPAD.Y != 0 {
		sol.Orientation.Y = core.ThrotleValue(ctrl.DPAD.Y, 1)
	} else {
		if core.Abs32(ctrl.DPAD.X) > 1 {
			sol.Orientation.Y = 0
		}
	}
}

func (sol *Solid) PushBack(d int32, o *core.Vector2d) {

	if o == nil {
		o = &core.Vector2d{sol.Orientation.X * -1, sol.Orientation.Y * -1}
	}

	sol.Position.X += int32(o.X) * d
	sol.Position.Y += int32(o.Y) * d

	sol.SetAnimation(MAN_PB_ANIM)
}
