package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type uiButton struct {
	coord        v2f
	size         v2i
	main         *ebiten.Image
	hover        *ebiten.Image
	click        *ebiten.Image
	hovered      bool
	clicked      bool
	clickCounter int
}

func newUIButton(coord v2f, states [3]*ebiten.Image) *uiButton {
	w := states[0].Bounds().Max.X - states[0].Bounds().Min.X
	h := states[0].Bounds().Max.Y - states[0].Bounds().Min.Y
	ret := &uiButton{
		main:  states[0],
		hover: states[1],
		click: states[2],
		size:  v2i{w, h},
		coord: coord,
	}
	return ret
}

func (b *uiButton) wasClicked() bool {
	return mbtnp(ebiten.MouseButtonLeft) && b.clicked
}

func (btn *uiButton) update() {
	mouse := mCoordsF()
	btn.hovered = false
	if mouse.x > btn.coord.x &&
		mouse.x < btn.coord.x+btn.size.Float64().x &&
		mouse.y > btn.coord.y &&
		mouse.y < btn.coord.y+btn.size.Float64().y {
		btn.hovered = true
	}

	if btn.hovered && mbtn(ebiten.MouseButtonLeft) {
		btn.clicked = true
	}

	if btn.clicked {
		btn.clickCounter++
	}

	if btn.clickCounter > 15 {
		btn.clicked = false
		btn.clickCounter = 0
	}
}

func (btn *uiButton) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(btn.coord.x, btn.coord.y)

	switch {
	case btn.clicked:
		if btn.clickCounter < 1 {
			screen.DrawImage(btn.hover, op)
		} else if btn.clickCounter < 4 {
			screen.DrawImage(btn.main, op)
		} else {
			screen.DrawImage(btn.click, op)
		}
	case btn.hovered:
		screen.DrawImage(btn.hover, op)
	default:
		screen.DrawImage(btn.main, op)
	}
}
