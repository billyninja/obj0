package game

import (
	"github.com/billyninja/obj0/core"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	KEY_ARROW_UP    = 1073741906
	KEY_ARROW_DOWN  = 1073741905
	KEY_ARROW_LEFT  = 1073741904
	KEY_ARROW_RIGHT = 1073741903
	KEY_LEFT_SHIFT  = 1073742049
	KEY_SPACE_BAR   = 32 // 1073741824 | 32
	KEY_C           = 99
	KEY_X           = 120
	KEY_Z           = 80 // todo
)

type ControlState struct {
	DPAD        core.Vector2d
	ACTION_A    int32
	ACTION_B    int32
	ACTION_C    int32
	ACTION_MAIN int32
	ACTION_MOD1 int32
}

func (cs *ControlState) Update(scn *Scene, PC *SceneEntity, keydown []sdl.Keycode, keyup []sdl.Keycode) {
	for _, key := range keydown {
		switch key {
		case KEY_ARROW_UP:
			cs.DPAD.Y -= 1
			break
		case KEY_ARROW_DOWN:
			cs.DPAD.Y += 1
			break
		case KEY_ARROW_LEFT:
			cs.DPAD.X -= 1
			break
		case KEY_ARROW_RIGHT:
			cs.DPAD.X += 1
			break
		case KEY_Z:
			cs.ACTION_A += 1
			break
		case KEY_X:
			cs.ACTION_B += 1
			break
		case KEY_C:
			cs.ACTION_C += 1
			break
		case KEY_SPACE_BAR:
			cs.ACTION_MAIN += 1
			break
		case KEY_LEFT_SHIFT:
			cs.ACTION_MOD1 += 1
			PC.Solid.Speed = (PC.Solid.Speed + PC.Char.SpeedMod) * 1.6
			break
		}
	}
	for _, key := range keyup {
		switch key {
		case KEY_ARROW_UP:
			cs.DPAD.Y = 0
			PC.Solid.Speed = PC.Solid.Speed + PC.Char.SpeedMod
			break
		case KEY_ARROW_DOWN:
			cs.DPAD.Y = 0
			break
		case KEY_ARROW_LEFT:
			cs.DPAD.X = 0
			break
		case KEY_ARROW_RIGHT:
			cs.DPAD.X = 0
			break
		case KEY_Z:
			cs.ACTION_A = 0
			break
		case KEY_X:
			if PC.Solid.Anim.PlayMode != 1 {
				MeleeAttack(PC, nil, scn)
			}
			cs.ACTION_B = 0
			break
		case KEY_C:
			if PC.Solid.Anim.PlayMode != 1 {
				println("REIMPLEMENT!!!")
			}
			cs.ACTION_C = 0
			break
		case KEY_SPACE_BAR:
			if !scn.DBox.NextText() {
				scn.ActProc(PC)
			}
			cs.ACTION_MAIN = 0
			break
		case KEY_LEFT_SHIFT:
			PC.Solid.Speed = PC.Char.BaseSpeed + PC.Char.SpeedMod
			cs.ACTION_MOD1 = 0
			break
		}
	}
}
