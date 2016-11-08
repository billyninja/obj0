package game

import (
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
