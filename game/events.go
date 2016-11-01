package game

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
)

func PlayDialog(listener *SceneEntity, speaker *SceneEntity, scn *Scene) {
	if len(speaker.Handlers.DialogScript) > 0 {
		scn.DBox.LoadText(speaker.Handlers.DialogScript)
	}
}

func AddToInv(picker *SceneEntity, item *SceneEntity, scn *Scene) {
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

func BashDoor(actor *Solid, door *Solid) {
	//currScene = load_scene(door.Handlers.DoorTo, nil)
}

func (scn *Scene) PlaceDrop(item *Item, origin *sdl.Rect) {
	instance := &SceneEntity{
		ItemPtr: item,
		Solid: &Solid{
			Txt:    item.Txtr,
			Source: item.Source,
			Position: &sdl.Rect{
				origin.X,
				origin.Y,
				item.Source.W,
				item.Source.H,
			},
		},
		Handlers: &SEventHandlers{
			OnPickUp: AddToInv,
		},
	}

	scn.Interactive = append(scn.Interactive, instance)
}
