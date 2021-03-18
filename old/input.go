package main

import "github.com/hajimehoshi/ebiten/v2"

// We mostly only need this because ebiten v2 does not really support an 'on press' type event yet
// TODO: add this funciton to ebiten v2

var keyTracker map[ebiten.Key]uint

func trackKey(key ebiten.Key) {
	keyTracker[key] = 0
}

func updateKeys() {
	for key, _ := range keyTracker {
		if ebiten.IsKeyPressed(key) {
			keyTracker[key]++
		} else {
			keyTracker[key] = 0
		}
	}
}

func btnp(key ebiten.Key) bool {
	if keyTracker[key] == 1 {
		return true
	}
	return false
}

func btn(key ebiten.Key) bool {
	if keyTracker[key] > 0 {
		return true
	}
	return false
}

var leftMouseButtonCount uint
var rightMouseButtonCount uint

func updateMouse() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		leftMouseButtonCount++
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		rightMouseButtonCount++
	}
}

func mbtn(btn ebiten.MouseButton) bool {
	if btn == ebiten.MouseButtonLeft {
		return leftMouseButtonCount == 0
	}
	if btn == ebiten.MouseButtonRight {
		return rightMouseButtonCount == 0
	}
	return false
}

func mbtnp(btn ebiten.MouseButton) bool {
	if btn == ebiten.MouseButtonLeft {
		return leftMouseButtonCount > 0
	}
	if btn == ebiten.MouseButtonRight {
		return rightMouseButtonCount > 0
	}
	return false
}
