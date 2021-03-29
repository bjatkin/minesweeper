package main

import (
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type powUnlock struct {
	scrollIn    float64
	scrollOut   float64
	powerUpType int
	slot        bool
	done        bool
	closed      bool
}

var exclaim *ebiten.Image

func (pow *powUnlock) update() {
	if pow.done {
		pow.scrollOut += (1.0 / 15.0)
	} else {
		pow.scrollIn += (1.0 / 15.0)
	}

	if btnp(ebiten.KeyEnter) || btnp(ebiten.KeyEscape) || mbtnp(ebiten.MouseButtonLeft) {
		pow.done = true
	}

	if pow.scrollOut > 1.5 {
		pow.closed = true
	}
}

func (pow *powUnlock) draw(screen *ebiten.Image) {
	icons := []*ebiten.Image{
		locked[0],
		addMine[0],
		minusMine[0],
		tidalWave[0],
		scaredyCat[0],
		dogWistle[0],
		shuffel[0],
		dogABone[0],
	}

	desc := []string{
		"use this slot to bring a powerup into a level",
		"select any {flipped}next to one or more {grass}\nif possible, a {dog}will be added under \na random {grass}",
		"select any {flipped}next to one or more {grass}\nif possible, a {dog}will be removed from \na random {grass}",
		"convert up to 4 random {grass}into {water}\n {water}will never contain a {dog}",
		"select any {grass}. That tile and its neighbors are flipped\nAny {dog}will be removed before being flipping",
		"a random {grass}will be hilighted, \ninidicating that {grass}is hiding a {dog}",
		"shuffle all unflagged {dog}to different {grass}",
		"prevent a \"game over\" once per level",
	}

	// Power Up Unlock Menu
	powOp := &ebiten.DrawImageOptions{}
	if pow.done {
		powOp.GeoM.Translate(lerp(-1, -300, pow.scrollOut), 25)
	} else {
		powOp.GeoM.Translate(lerp(-250, -1, pow.scrollIn), 25)
	}
	screen.DrawImage(powerUpUnlock, powOp)

	powOp.GeoM.Translate(75, 18)
	if pow.scrollIn > 1 && pow.scrollOut == 0 {
		if pow.slot {
			screen.DrawImage(powerUpUnlockSlot, powOp)
			drawDesc(screen, desc[0])
		} else {
			screen.DrawImage(powerUpUnlockPow, powOp)
			powOp.GeoM.Translate(7, 3)
			screen.DrawImage(icons[pow.powerUpType], powOp)
			powOp.GeoM.Translate(-7, -3)
			drawDesc(screen, desc[pow.powerUpType])
		}
	}

	var shakeX float64
	if pow.scrollIn < 3 {
		shakeX = math.Sin(float64(tickCounter * 2))
	}
	powOp.GeoM.Translate(-36+shakeX, -14)
	screen.DrawImage(exclaim, powOp)
}

func drawDesc(screen *ebiten.Image, desc string) {
	lines := strings.Split(desc, "\n")
	for i, line := range lines {
		flipped := []v2i{}
		grass := []v2i{}
		water := []v2i{}
		dog := []v2i{}
		var done bool
		for !done {
			done = true
			f := strings.Index(line, "{flipped}")
			if f > 0 {
				flipped = append(flipped, v2i{pxlLen(mainFont, line[:f]) - 6, 73 + i*16})
				done = false
			}
			line = strings.Replace(line, "{flipped}", "  ", 1)

			g := strings.Index(line, "{grass}")
			if g > 0 {
				grass = append(grass, v2i{pxlLen(mainFont, line[:g]) - 6, 73 + i*16})
				done = false
			}
			line = strings.Replace(line, "{grass}", "  ", 1)

			w := strings.Index(line, "{water}")
			if w > 0 {
				water = append(water, v2i{pxlLen(mainFont, line[:w]) - 6, 73 + i*16})
			}
			line = strings.Replace(line, "{water}", "  ", 1)

			d := strings.Index(line, "{dog}")
			if d > 0 {
				dog = append(dog, v2i{pxlLen(mainFont, line[:d]) - 6, 73 + i*16})
			}
			line = strings.Replace(line, "{dog}", "  ", 1)
		}

		x := (240 - float64(pxlLen(mainFont, line))) / 2
		tileOp := &ebiten.DrawImageOptions{}
		for _, flip := range flipped {
			tileOp.GeoM.Reset()
			tileOp.GeoM.Translate(flip.Float64().x+x, flip.Float64().y)
			screen.DrawImage(flipIcon, tileOp)
		}
		for _, g := range grass {
			tileOp.GeoM.Reset()
			tileOp.GeoM.Translate(g.Float64().x+x, g.Float64().y)
			screen.DrawImage(grassIcon, tileOp)
		}
		for _, w := range water {
			tileOp.GeoM.Reset()
			tileOp.GeoM.Translate(w.Float64().x+x, w.Float64().y)
			screen.DrawImage(waterIcon, tileOp)
		}
		for _, d := range dog {
			tileOp.GeoM.Reset()
			tileOp.GeoM.Translate(d.Float64().x+x, d.Float64().y)
			screen.DrawImage(dogIcon, tileOp)
		}

		pxlPrint(screen, mainFont, x, 75+float64(i*16), line)
	}
}
