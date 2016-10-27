package core

import (
	//	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

func ReleaseSpell(caster *Solid, tgt *Solid) {
	r := ProjectHitBox(Center(caster.Position), caster.Orientation, 48, &Vector2d{32, 32}, 0)
	ttl := time.Now().Add(3 * time.Second)

	spl := &Solid{
		Position: r,
		//Txt:      glowTxt,
		Ttl: ttl.Unix(),
		Anim: &Animation{
			Action:   BGLOW_A,
			Pose:     0,
			PoseTick: 16,
		},
		Handlers: &InteractionHandlers{
			OnCollDmg:   50,
			OnCollEvent: onColHdk,
		},
		CPattern: 0,
		MPattern: []Movement{
			Movement{*caster.Orientation, 255},
		},
		Collision: 1,
		Speed:     6,
	}

	println(spl)
}

func (ch *Char) PCMeleeAtk(cull []*Solid, visual []*VFXInst) {
	var stCost float32 = 5
	if stCost > ch.CurrentST {
		return
	}

	ch.CurrentST -= stCost
	ch.Solid.SetAnimation(MAN_ATK1_ANIM, nil)

	ch.MeleeAtk(cull, visual)
}

func (ch *Char) MeleeAtk(cull []*Solid, visual []*VFXInst) {

	r := ProjectHitBox(
		Center(ch.Solid.Position), ch.Solid.Orientation, 48, nil, 0)

	for _, cObj := range cull {
		if cObj.CharPtr != nil && CheckCol(r, cObj.Position) {
			var dmg float32 = 15.0
			cObj.CharPtr.DepletHP(dmg)
			r.W, r.H = 92, 92

			// MISS pTxt := PopText(font, r, fmt.Sprintf("%d", dmg), CL_WHITE)
			// visual = append(visual, pTxt)
			// visual = append(visual, Impact.Spawn(r, ch.Solid.Orientation))
		}
	}

	// MISS
	// visual = append(visual, Hit.Spawn(r, ch.Solid.Orientation))
}

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

func (ch *Char) ActProc(cull []*Solid) {
	action_hit_box := ProjectHitBox(
		Center(ch.Solid.Position), ch.Solid.Orientation, 32, nil, 1)

	for _, obj := range cull {
		if obj.Handlers != nil &&
			obj.Handlers.OnActEvent != nil &&
			CheckCol(action_hit_box, obj.Position) {
			obj.Handlers.OnActEvent(ch.Solid, obj)
			return
		}
	}
}

func (s *Solid) ProcMovement(cull []*Solid, limitW int32, limitH int32) {

	np := s.ApplyMovement()
	outbound := (np.X <= 0 || np.Y <= 0 || np.X > limitW || np.Y > limitH)

	if (np.X == s.Position.X && np.Y == s.Position.Y) || outbound {
		return
	}

	fr := FeetRect(np)
	for _, obj := range cull {
		if obj.Position == nil || obj == s {
			continue
		}
		if CheckCol(fr, obj.Position) && ResolveCol(s, obj) {
			return
		}
	}

	if s.CharPtr != nil {
		anim := s.CharPtr.CurrentFacing()
		if anim != nil && s.Anim != nil && s.Anim.PlayMode != 1 {
			s.Anim.PlayMode = anim.PlayMode
			s.Anim.Action = anim.Action
		}
	}

	s.Position = np
}

func (s *Solid) LoSCheck(tgt *Solid) bool {
	LoS := &sdl.Rect{
		s.Position.X - s.LoS,
		s.Position.Y - s.LoS,
		s.Position.W + (s.LoS * 2),
		s.Position.H + (s.LoS * 2),
	}

	if !CheckCol(tgt.Position, LoS) {
		return false
	}
	return true
}

func (s *Solid) DoChase(cull []*Solid, visual []*VFXInst) {

	s.Velocity.X = 0
	s.Velocity.Y = 0

	tgt := s.Chase

	diffX := Abs32(float32(s.Position.X - tgt.Position.X))
	diffY := Abs32(float32(s.Position.Y - tgt.Position.Y))

	if s.Position.X > tgt.Position.X {
		s.Velocity.X = -1
	}

	if s.Position.X < tgt.Position.X {
		s.Velocity.X = 1
	}
	if s.Position.Y > tgt.Position.Y {
		s.Velocity.Y = -1
	}

	if s.Position.Y < tgt.Position.Y {
		s.Velocity.Y = 1
	}

	if diffX < 80 && diffY < 80 && s.CharPtr != nil {
		chr := s.CharPtr
		if chr.AtkCoolDownC <= 0 {

			chr.MeleeAtk(cull, visual)
			chr.AtkCoolDownC += chr.AtkCoolDown

			r := ProjectHitBox(Center(s.Position), s.Orientation, 32, nil, 1)
			// MISS visual = append(visual, Hit.Spawn(r, s.Orientation))

			if CheckCol(r, tgt.Position) {
				tgt.CharPtr.DepletHP(s.Handlers.OnCollDmg)
				// MISS visual = append(visual, Impact.Spawn(r, s.Orientation))
				tgt.CharPtr.ApplyInvinc()
			}

		}

		return
	} else {
		s.ProcMovement(cull, 2000, 2000)
	}

	return
}

func (s *Solid) PeformPattern(sp float32, cull []*Solid) {

	anon := func(c uint32, mvs []Movement) *Movement {

		var sum uint32 = 0

		for _, mp := range mvs {
			sum += uint32(mp.Ticks)
			if sum > s.CPattern {
				return &mp
			}
		}
		s.CPattern = 0

		return nil
	}

	mov := anon(s.CPattern, s.MPattern)
	if mov != nil && s.Position != nil {
		s.Orientation = &mov.Orientation
		s.Velocity = &mov.Orientation

		s.ProcMovement(cull, 2000, 2000)

		s.CPattern += uint32(sp)
	}
}

func (ch *Char) ApplyInvinc() {
	if ch.Invinc == 0 {
		ch.Invinc = time.Now().Add(400 * time.Millisecond).Unix()
	}
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

func ResolveCol(ObjA *Solid, ObjB *Solid) bool {
	var halt bool
	if ObjB.Collision == 1 {
		halt = true
	}

	if ObjB.Handlers != nil && ObjB.Handlers.OnCollEvent != nil {
		ObjB.Handlers.OnCollEvent(ObjA, ObjB)
	}

	return halt
}

func PickUp(picker *Solid, item *Solid) {
	if item.Handlers != nil && item.Handlers.OnPickUp != nil {
		item.Handlers.OnPickUp(picker, item)
	}
	if picker.CharPtr != nil {
		picker.SetAnimation(MAN_PU_ANIM, nil)
	}
	item.Destroy()
}

func (ch *Char) PushBack(d int32, o *Vector2d) {
	ch.Solid.Position.X += int32(o.X) * d
	ch.Solid.Position.Y += int32(o.Y) * d
	// TODO WIRE PB ANIM into the ActionMap
	ch.Solid.SetAnimation(MAN_PB_ANIM, nil)
}

func (s *Solid) UpdateVelocity(dpad *Vector2d) {

	nv := &Vector2d{}
	*nv = *s.Velocity

	if dpad.X != 0 || s.Velocity.X != 0 {
		if dpad.X != 0 {
			nv.X = ThrotleValue(nv.X+dpad.X, 2)
		} else {
			nv.X = (Abs32(nv.X) - 1) * s.Orientation.X
		}
	}

	if dpad.Y != 0 || nv.Y != 0 {
		if dpad.Y != 0 {
			nv.Y = ThrotleValue(nv.Y+dpad.Y, 2)
		} else {
			nv.Y = (Abs32(nv.Y) - 1) * s.Orientation.Y
		}
	}

	*s.Velocity = *nv
}

func AddToInv(picker *Solid, item *Solid) {
	if picker.CharPtr != nil && item.ItemPtr != nil {
		for _, iStack := range picker.CharPtr.Inventory {
			if iStack.ItemTpl == item.ItemPtr {
				iStack.Qty += 1
				return
			}
		}
		picker.CharPtr.Inventory = append(picker.CharPtr.Inventory, &ItemStack{item.ItemPtr, 1})
	}
}

func onColHdk(tgt *Solid, hdk *Solid) {
	if tgt.CharPtr != nil {
		tgt.CharPtr.DepletHP(hdk.Handlers.OnCollDmg)
	}
	//Visual = append(Visual, impact.Spawn(tgt.Position, nil))
	hdk.Destroy()
}
