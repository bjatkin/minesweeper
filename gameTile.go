package main

import (
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

// these are initalized by the game board loading
var grassSprs [7]*ebiten.Image
var yellowGrassSprs [7]*ebiten.Image
var pinkGrassSprs [7]*ebiten.Image

var markerFlag aniSprite
var mineDog aniSprite
var dogBark *ebiten.Image

// Powerups
var addPU *ebiten.Image
var minusPU *ebiten.Image
var waterPU *ebiten.Image
var bombPU *ebiten.Image
var wistlePU *ebiten.Image

type tile struct {
	op          *ebiten.DrawImageOptions
	grass       aniSprite
	safe        bool
	flagged     bool
	flipped     bool
	mine        bool
	blown       bool
	adjCount    int
	adjLock     int
	flipLock    int
	adjTiles    [8]*tile
	water       bool
	doBark      int
	timeBack    int
	parentBoard *gameBoardScean
}

func newTile(parent *gameBoardScean, x, y float64) *tile {
	var grass aniSprite
	timing := []uint16{3, 3, 3, 3, 3, 3, 3}
	switch rand.Intn(5) {
	case 0, 1, 2:
		grass = newAniSprite(grassSprs[:], timing, false)
	case 3:
		grass = newAniSprite(yellowGrassSprs[:], timing, false)
	case 4:
		grass = newAniSprite(pinkGrassSprs[:], timing, false)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	return &tile{
		parentBoard: parent,
		op:          op,
		grass:       grass,
	}
}

func (t *tile) draw(screen *ebiten.Image) {
	t.grass.op.GeoM.Reset()
	t.grass.op.GeoM.Translate(getX(t.op), getY(t.op))
	t.grass.draw(screen)

	if t.timeBack > 0 && !t.flipped {
		pxlPrint(screen, mainFont, getX(t.op)+7, getY(t.op)+1, "+"+strconv.Itoa(t.timeBack)+"s")
	}

	if t.adjLock > 0 && !t.flipped {
		pxlPrint(screen, mainFont, getX(t.op)+7, getY(t.op)+1, "L")
	}

	if t.flipped && t.flipLock > 0 {
		pxlPrint(screen, mainFont, getX(t.op)+7, getY(t.op)+1, "?")
	}

	if t.flipped && t.flipLock <= 0 && t.adjCount > 0 {
		pxlPrint(screen, mainFont, getX(t.op)+7, getY(t.op)+1, strconv.Itoa(t.adjCount))
	}

	if !t.flipped && t.water {
		pxlPrint(screen, mainFont, getX(t.op)+7, getY(t.op)+1, "W")
	}

	if t.flagged {
		markerFlag.op.GeoM.Reset()
		markerFlag.op.GeoM.Translate(getX(t.op)-1, getY(t.op)-9)
		markerFlag.draw(screen)
	}

	if t.doBark > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(getX(t.op)-10, getY(t.op)-24)
		screen.DrawImage(dogBark, op)
	}
}

func (t *tile) flag() {
	if !t.flipped {
		t.flagged = !t.flagged
	}
}

func (t *tile) flip() {
	if t.flagged || t.flipped {
		return
	}

	var adjLockCount int
	for _, adj := range t.adjTiles {
		if adj == nil {
			adjLockCount++
		}
		if adj != nil && adj.flipped {
			adjLockCount++
		}
		if adj != nil && adj.mine && adj.flagged {
			adjLockCount++
		}
	}

	if t.adjLock > 0 && adjLockCount < t.adjLock {
		return
	}

	if t.mine {
		t.blown = true
		mineDog.op.GeoM.Reset()
		mineDog.op.GeoM.Translate(getX(t.op)-7, getY(t.op)-10)
		mineDog.play()
		return
	}

	t.flipped = true
	t.grass.play()
	if t.parentBoard.maxTime > 0 {
		t.parentBoard.maxTime += t.timeBack
	}

	if t.adjCount == 0 {
		for _, tile := range t.adjTiles {
			if tile != nil {
				tile.flip()
			}
		}
	}
}
