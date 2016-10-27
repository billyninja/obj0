package core

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

/***********************************
* 	ACTION STATES
* ---------------------------------
*  0 - PreAction
*  1 - From PreAction to Action
*  2 - Action itself
*  3 - From Action to Post Action
*  4 - Post Action
**********************************/

type ActionInterface interface {
	GetState() int
	GetActor() *Solid
	Step()

	PreActionVFX() *VFX
	PreActionAnim() *Animation

	ActionAnim() *Animation
	ActionVFX() *VFX
	Perform(*Solid, *Solid) []*Solid

	PostActionVFX() *VFX
	PostActionAnim() *Animation

	SetFinished()
}

type SpellTemplate struct {

	// Basic
	Name        string
	Description string
	BaseSPCost  float32
	Elemental   int32 // TODO - SEE HOW TO REPRESENT THIS
	BaseDmg     float32

	// Projectile Related (see if should decople)
	ProjectileTxtr  *sdl.Texture
	ProjectileAnim  *Animation
	ProjectileTtl   time.Duration
	ProjectileSpeed float32
	ProjectileSize  *Vector2d
	PreCastVfx      *VFX
	PreCastAnim     *Animation
	CastVfx         *VFX
	CastAnim        *Animation
	PostCastVfx     *VFX
	PostCastAnim    *Animation
}

func FactoryProjectile(p *sdl.Rect, o *Vector2d, t *sdl.Texture, a *Animation, ttl time.Duration, spd float32, sz *Vector2d) *Solid {

	pos := &sdl.Rect{
		p.X + p.W*int32(o.X),
		p.Y + p.H*int32(o.Y),
		p.W,
		p.H,
	}

	return &Solid{
		Position:  pos,
		Velocity:  &Vector2d{o.X, o.Y},
		Txt:       t,
		Anim:      a,
		Speed:     spd,
		Collision: 1,
	}
}

func (spt *SpellTemplate) perform(caster, subject *Solid) []*Solid {
	output := []*Solid{}

	// TODO - A Single spell might output multiple projectiles!
	if spt.ProjectileTxtr != nil {
		proj := FactoryProjectile(
			caster.Position,
			caster.Orientation,
			spt.ProjectileTxtr,
			spt.ProjectileAnim,
			spt.ProjectileTtl,
			spt.ProjectileSpeed,
			spt.ProjectileSize,
		)
		output = append(output, proj)
	}

	return output
}

func (spt *SpellTemplate) Cast(caster *Char, subject *Solid) *SpellCasting {
	return &SpellCasting{
		Caster:      caster,
		Template:    spt,
		currentStep: 0,
		Dmg:         spt.BaseDmg,
		SPCost:      spt.BaseSPCost,
	}
}

type SpellCasting struct {
	Template *SpellTemplate

	Dmg    float32
	SPCost float32
	Caster *Char

	// ActionInterface Related
	currentStep int
}

// Implementing Action Interface bellow

func (spl *SpellCasting) GetState() int {
	return spl.currentStep
}

func (spl *SpellCasting) GetActor() *Solid {
	return spl.Caster.Solid
}

func (spl *SpellCasting) Step() {
	println("step from", spl.currentStep, spl.currentStep+1)
	spl.currentStep += 1
}

func (spl *SpellCasting) PreActionVFX() *VFX {
	return spl.Template.PreCastVfx
}

func (spl *SpellCasting) PreActionAnim() *Animation {
	return spl.Template.PreCastAnim
}

func (spl *SpellCasting) ActionAnim() *Animation {
	return spl.Template.CastAnim
}

func (spl *SpellCasting) ActionVFX() *VFX {
	return spl.Template.CastVfx
}

func (spl *SpellCasting) Perform(caster *Solid, subject *Solid) []*Solid {
	var stCost float32 = 20
	if stCost > caster.CharPtr.CurrentST {
		return []*Solid{}
	}
	caster.CharPtr.CurrentST -= stCost

	return spl.Template.perform(caster, nil)
}

func (spl *SpellCasting) PostActionVFX() *VFX {
	return spl.Template.PostCastVfx
}

func (spl *SpellCasting) PostActionAnim() *Animation {
	return spl.Template.PostCastAnim
}

func (spl *SpellCasting) SetFinished() {
	spl.currentStep = -1
	spl.Caster.CurrentAction = nil
}
