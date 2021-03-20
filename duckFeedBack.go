package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	duckSurprised = iota
	duckNormal
	duckCool
	duckDead
)

type duckFeedBack struct {
	coord    v2f
	state    int
	powOne   *powerUp
	powTwo   *powerUp
	powThree *powerUp
}

func newDuckFeedBack(coord v2f, powOne, powTwo, powThree *powerUp) *duckFeedBack {
	powOne.coord = coord
	powOne.coord.x += 6
	powOne.coord.y += 21
	powTwo.coord = coord
	powTwo.coord.x += 27
	powTwo.coord.y += 21
	powThree.coord = coord
	powThree.coord.x += 48
	powThree.coord.y += 21
	return &duckFeedBack{
		coord:    coord,
		powOne:   powOne,
		powTwo:   powTwo,
		powThree: powThree,
		state:    duckNormal,
	}
}

func (d *duckFeedBack) update() {
	// TODO: were gonna need to add some stuf here
}

func (d *duckFeedBack) draw(screen *ebiten.Image) {
	d.powOne.draw(screen)
	d.powTwo.draw(screen)
	d.powThree.draw(screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(d.coord.x+64, d.coord.y)
	screen.DrawImage(powerUpHudBG, op)
	screen.DrawImage(sideDuck[d.state], op)
	op.GeoM.Translate(-64, 0)
	screen.DrawImage(powerUpHud, op)

	// should this be in an update function?
	if d.state == duckSurprised {
		d.state = duckNormal
	}
}
