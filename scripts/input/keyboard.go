package input

import "github.com/hajimehoshi/ebiten/v2"

// =========== KeyBinding ==========

type KeyBinding struct {
	Binding
	keys       [2]ebiten.Key
	prevActive bool
}

func (b *KeyBinding) Update(dt float64) {
	b.prevActive = b.active
	for i := range b.keys {
		if ebiten.IsKeyPressed(b.keys[i]) {
			b.active = true
			return
		}
	}
	b.active = false
}

func (b *KeyBinding) IsActive() bool {
	return b.active
}

func (b *KeyBinding) JustReleased() bool {
	return !b.active && b.prevActive
}

func (b *KeyBinding) JustPressed() bool {
	return b.active && !b.prevActive
}

// =========== KeyPressBinding ==========

type KeyPressPhase uint8

const (
	KeyPressPhaseNone KeyPressPhase = iota
	KeyPressPhaseDown
	KeyPressPhaseHold
	KeyPressPhaseUp
)

// =========== KeyHoldBinding ===========

type KeyHoldBinding struct {
	KeyBinding
	holdTime   float64 // current hold duration
	holdThresh float64 // threshold for hold time
}

func NewKeyHoldBinding(action InputAction, holdThresh float64, keys [2]ebiten.Key) *KeyHoldBinding {
	return &KeyHoldBinding{
		KeyBinding: KeyBinding{
			Binding: Binding{action: action},
			keys:    keys,
		},
		holdThresh: holdThresh,
	}
}

func (b *KeyHoldBinding) Update(dt float64) {
	b.KeyBinding.Update(dt)
	if b.active {
		b.holdTime += dt
		if b.holdTime >= b.holdThresh {
			b.active = true
		} else {
			b.active = false
		}
	} else {
		b.holdTime = 0
		b.active = false
	}
}

func (b *KeyHoldBinding) IsActive() bool {
	return b.active
}

// =========== KeyPressDurationBinding ==========

type KeyPressDurationBinding struct {
	KeyBinding
	holdMax  float64 // duration the key can be held before deactivating
	holdTime float64 // current hold time
	pressure float64 // normalized pressure (0.0 to 1.0)
}

func NewKeyPressDurationBinding(action InputAction, holdTime float64, keys [2]ebiten.Key) *KeyPressDurationBinding {
	return &KeyPressDurationBinding{
		KeyBinding: KeyBinding{
			Binding: Binding{action: action},
			keys:    keys,
		},
		holdMax: holdTime,
	}
}

func (b *KeyPressDurationBinding) Update(dt float64) {
	b.KeyBinding.Update(dt)
	if b.active {
		b.holdTime += dt
		if b.holdTime >= b.holdMax {
			b.active = false
		}
		b.pressure = b.holdTime / b.holdMax
	} else {
		b.holdTime = 0
	}
}

func (b *KeyPressDurationBinding) Pressure() float64 {
	return b.pressure
}
