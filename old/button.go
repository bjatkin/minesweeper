package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type button struct {
	transform
	eid     uint64
	primary bool
	hovered bool
	clicked bool
	message string
}

func (b *button) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	if b.primary {
		switch {
		case b.hovered:
			screen.DrawImage(primaryButton[1], op)
			fmt.Println("Hovered")
		case b.clicked:
			screen.DrawImage(primaryButton[2], op)
			b.clicked = false
		default:
			screen.DrawImage(primaryButton[0], op)
		}
	} else {
		switch {
		case b.hovered:
			screen.DrawImage(secondButton[1], op)
			fmt.Println("Hovered")
		case b.clicked:
			screen.DrawImage(secondButton[2], op)
			b.clicked = false
		default:
			screen.DrawImage(secondButton[0], op)
		}
	}
}

func (b *button) update() error {
	return nil
}

func (b *button) setID(eid uint64) {
	b.eid = eid
}

func (b *button) getID() uint64 {
	return b.eid
}

func (b *button) hover(x, y int) {
	x1, x2 := b.x, b.x+56
	y1, y2 := b.y, b.y+16
	if float64(x) >= x1 && float64(x) <= x2 &&
		float64(y) >= y1 && float64(y) <= y2 {
		b.hovered = true
		return
	}
	b.hovered = false
}

func (b *button) click(x, y int) {
	x1, x2 := b.x, b.x+56
	y1, y2 := b.y, b.y+16
	if float64(x) >= x1 && float64(x) <= x2 &&
		float64(y) >= y1 && float64(y) <= y2 {
		b.clicked = true
	}
}
