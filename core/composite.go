package core

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	glowTxt *sdl.Texture = nil
	Hit     *VFX         = nil
	Impact  *VFX         = nil

	// FACING ORIENTATION
	F_LEFT  Vector2d = Vector2d{-1, 0}
	F_RIGHT          = Vector2d{1, 0}
	F_UP             = Vector2d{0, -1}
	F_DOWN           = Vector2d{0, 1}
	F_DL             = Vector2d{-1, 1}
	F_DR             = Vector2d{1, 1}
	F_UL             = Vector2d{-1, -1}
	F_UR             = Vector2d{1, -1}

	CL_WHITE sdl.Color = sdl.Color{255, 255, 255, 255}
	SHADOW             = &sdl.Rect{320, 224, TSzi, TSzi}

	MAN_PB_S1 *sdl.Rect = &sdl.Rect{96, 160, TSzi, TSzi}
	MAN_PB_S2 *sdl.Rect = &sdl.Rect{128, 160, TSzi, TSzi}
	MAN_PB_S3 *sdl.Rect = &sdl.Rect{160, 160, TSzi, TSzi}
	MAN_PU_S1 *sdl.Rect = &sdl.Rect{0, 192, TSzi, TSzi}
	MAN_PU_S2 *sdl.Rect = &sdl.Rect{128, 192, TSzi, TSzi}
	MAN_PU_S3 *sdl.Rect = &sdl.Rect{160, 160, TSzi, TSzi}

	MAN_CS_S1     *sdl.Rect    = &sdl.Rect{192, 224, TSzi, TSzi}
	MAN_CS_S2     *sdl.Rect    = &sdl.Rect{224, 224, TSzi, TSzi}
	MAN_CS_S3     *sdl.Rect    = &sdl.Rect{256, 224, TSzi, TSzi}
	MAN_CS_S4     *sdl.Rect    = &sdl.Rect{224, 192, TSzi, TSzi}
	MAN_PUSH_BACK [8]*sdl.Rect = [8]*sdl.Rect{MAN_PB_S1, MAN_PB_S2}
	MAN_CAST      [8]*sdl.Rect = [8]*sdl.Rect{MAN_CS_S1, MAN_CS_S3, MAN_CS_S4}
	MAN_AT_1      [8]*sdl.Rect = [8]*sdl.Rect{MAN_CS_S4}
	MAN_PICK_UP   [8]*sdl.Rect = [8]*sdl.Rect{MAN_PU_S1, MAN_PU_S2, MAN_PU_S3}

	PUFF_S1 *sdl.Rect    = &sdl.Rect{0, 0, 64, 64}
	PUFF_S2 *sdl.Rect    = &sdl.Rect{64, 0, 64, 64}
	PUFF_S3 *sdl.Rect    = &sdl.Rect{128, 0, 64, 64}
	PUFF_S4 *sdl.Rect    = &sdl.Rect{192, 0, 64, 64}
	PUFF_S5 *sdl.Rect    = &sdl.Rect{0, 64, 64, 64}
	PUFF_S6 *sdl.Rect    = &sdl.Rect{64, 64, 64, 64}
	PUFF_A  [8]*sdl.Rect = [8]*sdl.Rect{PUFF_S1, PUFF_S2, PUFF_S3, PUFF_S4, PUFF_S5, PUFF_S6}

	BLANK  *sdl.Rect    = &sdl.Rect{2000, 200, 192, 192}
	HIT_S1 *sdl.Rect    = &sdl.Rect{576, 0, 192, 192}
	HIT_S2 *sdl.Rect    = &sdl.Rect{758, 0, 192, 192}
	HIT_S3 *sdl.Rect    = &sdl.Rect{192, 192, 192, 192}
	HIT_A  [8]*sdl.Rect = [8]*sdl.Rect{HIT_S1, HIT_S2, HIT_S3}

	HIT_S4 *sdl.Rect    = &sdl.Rect{0, 0, 192, 192}
	HIT_S5 *sdl.Rect    = &sdl.Rect{192, 0, 192, 192}
	HIT_S6 *sdl.Rect    = &sdl.Rect{384, 0, 192, 192}
	HIT_B  [8]*sdl.Rect = [8]*sdl.Rect{BLANK, BLANK, HIT_S4, HIT_S5, HIT_S6, HIT_S5, HIT_S4}

	LAVA_S1 *sdl.Rect    = &sdl.Rect{192, 0, TSzi, TSzi}
	LAVA_S2 *sdl.Rect    = &sdl.Rect{224, 0, TSzi, TSzi}
	LAVA_S3 *sdl.Rect    = &sdl.Rect{256, 0, TSzi, TSzi}
	LAVA_A  [8]*sdl.Rect = [8]*sdl.Rect{LAVA_S1, LAVA_S2, LAVA_S3, LAVA_S3, LAVA_S2}

	YGLOW_S1 *sdl.Rect    = &sdl.Rect{0, 0, TSzi * 2, TSzi * 2}
	YGLOW_S2 *sdl.Rect    = &sdl.Rect{32, 0, TSzi * 2, TSzi * 2}
	YGLOW_S3 *sdl.Rect    = &sdl.Rect{64, 0, TSzi * 2, TSzi * 2}
	YGLOW_A  [8]*sdl.Rect = [8]*sdl.Rect{YGLOW_S1, YGLOW_S2, YGLOW_S3, YGLOW_S2}

	BGLOW_S1 *sdl.Rect = &sdl.Rect{224, 224, TSzi, TSzi}
	BGLOW_S2 *sdl.Rect = &sdl.Rect{256, 224, TSzi, TSzi}
	BGLOW_S3 *sdl.Rect = &sdl.Rect{288, 224, TSzi, TSzi}
	BGLOW_S4 *sdl.Rect = &sdl.Rect{320, 224, TSzi, TSzi}

	BGLOW_A [8]*sdl.Rect = [8]*sdl.Rect{BGLOW_S1, BGLOW_S2, BGLOW_S3, BGLOW_S4, BGLOW_S3, BGLOW_S2}

	LAVA_ANIM     = &Animation{Action: LAVA_A, PoseTick: 8}
	LIFE_ORB_ANIM = &Animation{Action: YGLOW_A, PoseTick: 8}
	MAN_PB_ANIM   = &Animation{Action: MAN_PUSH_BACK, PoseTick: 18, PlayMode: 1}
	MAN_CS_ANIM   = &Animation{Action: MAN_CAST, PoseTick: 18, PlayMode: 1}
	MAN_PU_ANIM   = &Animation{Action: MAN_PICK_UP, PoseTick: 18, PlayMode: 1}
	MAN_ATK1_ANIM = &Animation{Action: MAN_AT_1, PoseTick: 4, PlayMode: 1}
)
