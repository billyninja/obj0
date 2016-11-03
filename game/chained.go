package game

import (
	"github.com/billyninja/obj0/assets"
	"github.com/billyninja/obj0/core"
	"time"
)

type ActionStep struct {
	Action         SEvent
	TimeToNextStep time.Duration
}

type TimedChainAction struct {
	Steps       []*ActionStep
	Actor       *SceneEntity
	Subject     *SceneEntity
	CurrStep    int
	NextTimeout int64
}

func (tca *TimedChainAction) Destroy() {
	tca.Steps = []*ActionStep{}
	tca.Actor = nil
	tca.Subject = nil
}

func CloneAndAssign(tca *TimedChainAction, source, tgt *SceneEntity) {
	clone := &TimedChainAction{}
	// DeepCopy ActionSteps?
	*clone = *tca
	clone.Actor = source
	clone.Subject = tgt
	source.ChainedAction = clone
}

func (se *SceneEntity) UpdateChainedAction(scn *Scene) {
	if se.ChainedAction == nil {
		return
	}

	now := time.Now()
	if now.Unix() >= se.ChainedAction.NextTimeout {
		ca := se.ChainedAction
		ca.CurrStep++

		if ca.CurrStep < len(ca.Steps) {
			newStep := ca.Steps[ca.CurrStep]
			ca.NextTimeout = now.Add(newStep.TimeToNextStep).Unix()
			newStep.Action(ca.Actor, ca.Subject, scn)
		} else {
			ca.Actor.ChainedAction = nil
		}
	}
}

func FiraPreCast(caster, tgt *SceneEntity, scn *Scene) {
	caster.Solid.SetAnimation(MAN_CS_ANIM)
}

func FiraProj(caster, tgt *SceneEntity, scn *Scene) {
	println("Fira spawn projectile!")

	proj := FactoryProjectile(
		caster.Solid.Position,
		caster.Solid.Orientation,
		assets.Textures.Sprites.Glow,
		LIFE_ORB_ANIM,
		time.Second*5,
		3,
		&core.Vector2d{32, 32},
	)

	scn.Projectiles = append(scn.Projectiles, proj)
}

func FiraAfterCast(caster, tgt *SceneEntity, scn *Scene) {
	caster.Solid.SetAnimation(MAN_PU_ANIM)
}

var (
	Fira *TimedChainAction = &TimedChainAction{
		Steps: []*ActionStep{
			&ActionStep{FiraPreCast, time.Millisecond * 1000},
			&ActionStep{FiraProj, time.Millisecond * 2000},
			&ActionStep{FiraAfterCast, time.Millisecond * 200},
		},
	}
)
