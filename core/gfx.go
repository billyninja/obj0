package core

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"time"
)

const (
	TSz  float32 = 32
	TSzi int32   = int32(TSz)
)

// SPRITESHEET RELATED
func (ss *SpriteSheet) BuildBasicActions(actLength uint8, hasDiagonals bool) *ActionMap {
	O := []Vector2d{F_UP, F_DOWN, F_LEFT, F_RIGHT}
	if hasDiagonals {
		O = append(O, []Vector2d{F_DL, F_DR, F_UL, F_UR}...)
	}
	var AM = &ActionMap{}
	for _, o := range O {
		anim := &Animation{}

		for p := 0; p < int(actLength); p++ {
			anim.Action[p] = ss.GetPose(o, uint8(p))
		}

		switch o {
		case F_UP:
			AM.UP = anim
			break
		case F_DOWN:
			AM.DOWN = anim
			break
		case F_LEFT:
			AM.LEFT = anim
			break
		case F_RIGHT:
			AM.RIGHT = anim
			break
		case F_DL:
			AM.DL = anim
			break
		case F_DR:
			AM.DR = anim
			break
		case F_UL:
			AM.UL = anim
			break
		case F_UR:
			AM.UR = anim
			break
		} // end switch
	} // end for
	return AM
}

func (ss *SpriteSheet) GetPose(o Vector2d, p uint8) *sdl.Rect {

	var (
		poseY int32 = 0
		poseX int32 = 0
	)

	if o.X == -1 && o.Y == 0 {
		poseY = ss.StepW
	}
	if o.X == 1 && o.Y == 0 {
		poseY = ss.StepW * 2
	}
	if o.Y == -1 && o.X == 0 {
		poseY = ss.StepW * 3
	}
	if o.Y == 1 && o.X == 0 {
		poseY = 0
	}
	switch o {
	case F_UP:
		poseY = ss.StepH * 3
		break
	case F_DOWN:
		break
	case F_LEFT:
		poseY = ss.StepH * 1
		break
	case F_RIGHT:
		poseY = ss.StepH * 2
		break
	case F_DL:
		poseX = ss.StepW * 3
		break
	case F_DR:
		poseY = ss.StepH * 2
		poseX = ss.StepW * 3
		break
	case F_UL:
		poseY = ss.StepH * 1
		poseX = ss.StepW * 3
		break
	case F_UR:
		poseY = ss.StepH * 3
		poseX = ss.StepW * 3
		break
	}

	poseX += int32(p) * ss.StepW

	return &sdl.Rect{ss.StX + poseX, ss.StY + poseY, ss.StepW, ss.StepH}
}

// SOLID-ANIMATION RELATED

func (s *Solid) SetAnimation(an *Animation, evt Event) {
	var nA *Animation = &Animation{}
	*nA = *an

	s.Anim = nA
	s.Anim.Pose = 0
	s.Anim.PoseTick = 16
	s.Anim.After = evt
}

func (s *Solid) PlayAnimation() {

	s.Anim.PoseTick -= 1

	if s.Anim.PoseTick <= 0 {
		s.Anim.PoseTick = 12
		if s.CharPtr != nil {
			anim := s.CharPtr.CurrentFacing()
			prvPose := s.Anim.Pose
			s.Anim.Pose = getNextPose(s.Anim.Action, s.Anim.Pose)
			if anim != nil && s.Anim.Pose <= prvPose && s.Anim.PlayMode == 1 {
				if s.Anim.After != nil {
					s.Anim.After(s, nil)
				}
				s.SetAnimation(anim, nil)
			}
		} else {
			s.Anim.Pose = getNextPose(s.Anim.Action, s.Anim.Pose)
		}
	}
}

func (ch *Char) CurrentFacing() *Animation {

	if ch.Solid.Orientation.X == 0 && ch.Solid.Orientation.Y == 1 {
		return ch.ActionMap.DOWN
	}

	if ch.Solid.Orientation.X == 0 && ch.Solid.Orientation.Y == -1 {
		return ch.ActionMap.UP
	}

	if ch.Solid.Orientation.X == -1 && ch.Solid.Orientation.Y == 0 {
		return ch.ActionMap.LEFT
	}

	if ch.Solid.Orientation.X == 1 && ch.Solid.Orientation.Y == 0 {
		return ch.ActionMap.RIGHT
	}

	if ch.Solid.Orientation.X == 1 && ch.Solid.Orientation.Y == 1 {
		if ch.ActionMap.DR != nil {
			return ch.ActionMap.DR
		} else {
			return ch.ActionMap.RIGHT
		}
	}

	if ch.Solid.Orientation.X == 1 && ch.Solid.Orientation.Y == -1 {
		if ch.ActionMap.UR != nil {
			return ch.ActionMap.UR
		} else {
			return ch.ActionMap.RIGHT
		}
	}

	if ch.Solid.Orientation.X == -1 && ch.Solid.Orientation.Y == 1 {
		if ch.ActionMap.DL != nil {
			return ch.ActionMap.DL
		} else {
			return ch.ActionMap.LEFT
		}
	}

	if ch.Solid.Orientation.X == -1 && ch.Solid.Orientation.Y == -1 {
		if ch.ActionMap.UL != nil {
			return ch.ActionMap.UL
		} else {
			return ch.ActionMap.LEFT
		}
	}

	return ch.ActionMap.DOWN
}

func getNextPose(action [8]*sdl.Rect, currPose uint8) uint8 {
	if action[currPose+1] == nil {
		return 0
	} else {
		return currPose + 1
	}
}

// TEXT RELATED

func (t *TextEl) Bake(renderer *sdl.Renderer, limit int) (*sdl.Texture, int32, int32) {
	if t.Content == t.BakedContent {
		return t.Txtr, t.TW, t.TH
	}

	surface, _ := t.Font.RenderUTF8_Blended_Wrapped(t.Content, t.Color, limit)
	defer surface.Free()
	t.Txtr, _ = renderer.CreateTextureFromSurface(surface)
	t.TW, t.TH = surface.W, surface.H

	return t.Txtr, t.TW, t.TH
}

// VFX RELATED

func PopText(font *ttf.Font, pos *sdl.Rect, content string, color sdl.Color) *VFXInst {

	tEl := &TextEl{
		Font:    font,
		Content: content,
		Color:   color,
	}
	tPos := &sdl.Rect{pos.X, pos.Y - 30, 20, 20}

	return &VFXInst{
		Text: tEl,
		Pos:  tPos,
		Ttl:  time.Now().Add(400 * time.Millisecond).Unix(),
	}
}

func (vi *VFXInst) UpdateAnim() {
	if vi.Text != nil {
		vi.Pos.Y -= 1
		vi.Pos.X -= 1
		return
	}
	if vi.Vfx == nil {
		return
	}

	vi.CurrTick -= 1
	if vi.CurrTick == 0 {
		vi.Pose += 1
		if vi.Vfx.Strip[vi.Pose] == nil {
			if vi.Loop <= 0 {
				vi.Destroy()
				return
			} else {
				vi.Pose = 0
			}
		}
		vi.CurrTick = vi.Tick
	}
}

func (vi *VFXInst) Destroy() {
	vi.Vfx, vi.Pos = nil, nil
}

func (vi *VFXInst) CurrentFrame() *sdl.Rect {
	return vi.Vfx.Strip[vi.Pose]
}

func (v *VFX) Spawn(Position *sdl.Rect, flip *Vector2d) *VFXInst {
	var ttl int64

	i := &VFXInst{
		Vfx:      v,
		Pos:      Position,
		Ttl:      ttl,
		Tick:     v.DefaultSpeed,
		CurrTick: v.DefaultSpeed,
	}

	if flip != nil {
		i.Flip = *flip
	}

	return i
}
