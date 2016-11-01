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

func (s *Solid) LoSCheck(tgt *Solid) bool {
	LoS := &sdl.Rect{
		s.Position.X - s.LoS,
		s.Position.Y - s.LoS,
		s.Position.W + (s.LoS * 2),
		s.Position.H + (s.LoS * 2),
	}

	if !core.CheckCol(tgt.Position, LoS) {
		return false
	}
	return true
}

func (s *Solid) ApplyMovement() *sdl.Rect {
	speed := s.Speed

	if s.Velocity.X == 0 && s.Velocity.Y == 0 {
		return s.Position
	}

	// Decreasing speed for diagonal movement
	if s.Velocity.X != 0 && s.Velocity.Y != 0 {
		speed -= 0.5
		if speed < 1 {
			speed = 1
		}
	}

	return &sdl.Rect{
		(s.Position.X + int32(s.Velocity.X*speed)),
		(s.Position.Y + int32(s.Velocity.Y*speed)),
		s.Position.W,
		s.Position.H,
	}
}

func (s *Solid) UpdateVelocity(dpad *core.Vector2d) {

	nv := &core.Vector2d{}
	*nv = *s.Velocity

	if dpad.X != 0 || s.Velocity.X != 0 {
		if dpad.X != 0 {
			nv.X = core.ThrotleValue(nv.X+dpad.X, 2)
		} else {
			nv.X = (core.Abs32(nv.X) - 1) * s.Orientation.X
		}
	}

	if dpad.Y != 0 || nv.Y != 0 {
		if dpad.Y != 0 {
			nv.Y = core.ThrotleValue(nv.Y+dpad.Y, 2)
		} else {
			nv.Y = (core.Abs32(nv.Y) - 1) * s.Orientation.Y
		}
	}

	*s.Velocity = *nv
}

func (s *Solid) PatternStep() *Movement {
	var sum uint32 = 0

	for _, mp := range s.MPattern {
		sum += uint32(mp.Ticks)
		if sum > s.CPattern {
			return &mp
		}
	}
	s.CPattern = 0

	return nil
}

func (s *Solid) UpdatePCOrientation(ctrl *ControlState) {

	if ctrl.DPAD.X != 0 {
		s.Orientation.X = core.ThrotleValue(ctrl.DPAD.X, 1)
	} else {
		if core.Abs32(ctrl.DPAD.Y) > 1 {
			s.Orientation.X = 0
		}
	}

	if ctrl.DPAD.Y != 0 {
		s.Orientation.Y = core.ThrotleValue(ctrl.DPAD.Y, 1)
	} else {
		if core.Abs32(ctrl.DPAD.X) > 1 {
			s.Orientation.Y = 0
		}
	}
}

func (sol *Solid) PushBack(d int32, o *core.Vector2d) {
	sol.Position.X += int32(o.X) * d
	sol.Position.Y += int32(o.Y) * d

	sol.SetAnimation(MAN_PB_ANIM)
}
