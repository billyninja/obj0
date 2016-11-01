package game

import (
	"github.com/billyninja/obj0/core"
)

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
}

func MeleeAttack(source, tgt *SceneEntity, scn *Scene) {
	source.Solid.SetAnimation(MAN_ATK1_ANIM)
	action_hit_box := core.ProjectHitBox(
		core.Center(source.Solid.Position), source.Solid.Orientation, 32, nil, 1)

	for _, se := range scn.CullMap {
		if se != source && se.Char != nil && core.CheckCol(action_hit_box, se.Solid.Position) {
			se.Char.DepletHP(12)
		}
	}
}

func OpenDoor(actor, door *SceneEntity, scn *Scene) {
	if door.Handlers != nil && door.Handlers.DoorTo != "" {
		println("NOT IMPLEMENTED! too much main dependency right now.", door.Handlers.DoorTo)
	} else {
		println("no DoorTo")
	}
}
