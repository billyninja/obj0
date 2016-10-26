package core

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

func Abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0
	}
	return x
}

func ThrotleValue(v float32, limitAbs float32) float32 {
	abs := Abs32(v)
	sign := Copysign32(1, v)
	if abs > limitAbs {
		return limitAbs * sign
	}
	return v
}

func Copysign32(x, y float32) float32 {
	const sign = 1 << 31
	return math.Float32frombits(math.Float32bits(x)&^sign | math.Float32bits(y)&sign)
}

func Center(r *sdl.Rect) *Vector2d {
	return &Vector2d{float32(r.X + (r.W / 2)), float32(r.Y + (r.H / 2))}
}

func ProjectHitBox(origin *Vector2d, orientation *Vector2d, length float32, distance *Vector2d, comp float32) *sdl.Rect {
	if distance == nil {
		distance = &Vector2d{0, 0}
	}

	lx := length + (length * Abs32(orientation.Y) * comp)
	ly := length + (length * Abs32(orientation.X) * comp)
	px := origin.X - (lx / 2) + (orientation.X * distance.X) + (orientation.X * lx)
	py := origin.Y - (ly / 2) + (orientation.Y * distance.Y) + (orientation.Y * ly)

	return &sdl.Rect{int32(px), int32(py), int32(lx), int32(ly)}
}

func CheckCol(r1 *sdl.Rect, r2 *sdl.Rect) bool {
	return (r1.X < (r2.X+r2.W) &&
		r1.X+r1.W > r2.X &&
		r1.Y < r2.Y+r2.H &&
		r1.Y+r1.H > r2.Y)
}

func FeetRect(pos *sdl.Rect) *sdl.Rect {
	third := pos.H / 3
	return &sdl.Rect{pos.X, pos.Y + third, pos.W, pos.H - third}
}
