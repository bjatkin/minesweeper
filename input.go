package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var keyTracker = map[ebiten.Key]uint{
	ebiten.KeyEnter:  0,
	ebiten.KeyLeft:   0,
	ebiten.KeyRight:  0,
	ebiten.KeyUp:     0,
	ebiten.KeyDown:   0,
	ebiten.KeyW:      0,
	ebiten.KeyA:      0,
	ebiten.KeyS:      0,
	ebiten.KeyD:      0,
	ebiten.KeyEscape: 0,
	ebiten.Key0:      0,
	ebiten.Key1:      0,
	ebiten.Key2:      0,
	ebiten.Key3:      0,
	ebiten.Key4:      0,
	ebiten.Key5:      0,
}

func updateKeys() {
	for key := range keyTracker {
		if ebiten.IsKeyPressed(key) {
			keyTracker[key]++
		} else {
			keyTracker[key] = 0
		}
	}
}

func btnp(key ebiten.Key) bool {
	count, found := keyTracker[key]
	if !found {
		log.Fatalf("Untracked Key: %v", key)
		return false
	}
	if count == 1 {
		return true
	}
	return false
}

func btn(key ebiten.Key) bool {
	count, found := keyTracker[key]
	if !found {
		log.Fatalf("Untracked Key: %v", key)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

var leftMouseButtonCount uint
var rightMouseButtonCount uint

func updateMouse() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		leftMouseButtonCount++
	} else {
		leftMouseButtonCount = 0
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		rightMouseButtonCount++
	} else {
		rightMouseButtonCount = 0
	}
}

func mbtnp(btn ebiten.MouseButton) bool {
	if btn == ebiten.MouseButtonLeft {
		return leftMouseButtonCount == 1
	}
	if btn == ebiten.MouseButtonRight {
		return rightMouseButtonCount == 1
	}
	log.Fatalf("Untracked Mouse Button: %v", btn)
	return false
}

func mbtn(btn ebiten.MouseButton) bool {
	if btn == ebiten.MouseButtonLeft {
		return leftMouseButtonCount > 0
	}
	if btn == ebiten.MouseButtonRight {
		return rightMouseButtonCount > 0
	}
	log.Fatalf("Untracked Mouse Button: %v", btn)
	return false
}

func mCoords() v2i {
	x, y := ebiten.CursorPosition()
	return v2i{x: x, y: y}
}

func mCoordsF() v2f {
	return mCoords().Float64()
}
