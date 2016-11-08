package game

import (
	"github.com/billyninja/obj0/core"
)

type SEvent func(source *SceneEntity, subject *SceneEntity, scene *Scene)

func PlayDialog(listener, speaker *SceneEntity, scn *Scene) {
	if len(speaker.Handlers.DialogScript) > 0 {
		scn.DBox.LoadText(speaker.Handlers.DialogScript)
	}
}

func AddToInv(picker, item *SceneEntity, scn *Scene) {
	char := picker.Char
	if picker != nil && char != nil && item != nil && item.ItemPtr != nil {
		for _, iStack := range char.Inventory {
			if iStack.ItemTpl == item.ItemPtr {
				iStack.Qty += 1
				return
			}
		}
		char.Inventory = append(char.Inventory, &ItemStack{item.ItemPtr, 1})
	}
}

func Pickup(picker, item *SceneEntity, scn *Scene) {
	picker.Solid.SetAnimation(MAN_PU_ANIM)
	AddToInv(picker, item, scn)
	item.Destroy()
}

func OpenDoor(actor, door *SceneEntity, scn *Scene) {
	if door.Handlers != nil && door.Handlers.DoorTo != "" {
		SceneTransition(door.Handlers.DoorTo, actor, scn)
	} else {
		println("no DoorTo?")
	}
}

func MeleeAttack(source, tgt *SceneEntity, scn *Scene) {
	sol := source.Solid
	sol.SetAnimation(MAN_ATK1_ANIM)
	action_hit_box := core.ProjectHitBox(
		core.Center(sol.Position), sol.Orientation, 32, nil, 1)

	scn.SpawnVFX(action_hit_box, sol.Orientation, Hit, 1)

	hitf := func(victim *SceneEntity) {
		scn.SpawnVFX(victim.Solid.Position, victim.Solid.Orientation, Impact, 1)
		victim.Char.DepletHP(12)
	}

	if tgt != nil {
		if source != tgt && core.CheckCol(action_hit_box, tgt.Solid.Position) {
			hitf(tgt)
		}
	} else {
		for _, se := range scn.CullMap {
			if se != source &&
				se.Char != nil &&
				se.Solid != nil &&
				core.CheckCol(action_hit_box, se.Solid.Position) {
				hitf(se)
			}
		}
	}
}

func CastSpell(source, tgt *SceneEntity, scn *Scene) {
	sol := source.Solid
	sol.SetAnimation(MAN_CS_ANIM)
	CloneAndAssign(Fira, source, tgt)
}
