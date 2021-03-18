package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type n_tile struct {
	index        v2i
	adj          [8]*n_tile
	adjCount     int
	parent       *levelScean
	gfx          *n_aniSprite
	mine         bool
	water        bool
	flagged      bool
	flipped      bool
	lockedCount  int
	iced         bool
	icedCount    int
	shakeCounter int
}

func n_newTile(level *levelScean, index v2i, iced bool, locked bool, water bool) *n_tile {
	tile := &n_tile{
		parent: level,
		index:  index,
		iced:   iced,
	}

	if iced {
		tile.gfx = n_newAniSprite(
			iceImg[:],
			[]uint{6, 6, 6},
			false,
		)
	}

	if locked {
		tile.lockedCount = 15
	}

	if water {
		tile.gfx = n_newAniSprite(
			waterImg[:],
			[]uint{5, 5, 5, 5, 5, 5, 5},
			false,
		)
	}

	if !iced && !water {
		tile.gfx = n_newAniSprite(nil, []uint{5, 5, 5, 5, 5, 5, 5}, false)
		switch rand.Intn(5) {
		case 0:
			tile.gfx.sprs = blueGrassImg[:]
		case 1:
			tile.gfx.sprs = yellowGrassImg[:]
		case 2:
			tile.gfx.sprs = pinkGrassImg[:]
		default:
			tile.gfx.sprs = grassImg[:]
		}
	}

	return tile
}

func (t *n_tile) update(flipCount int) {
	t.gfx.update()
	t.shakeCounter--

	if t.flipped && flipCount > 0 && t.lockedCount > 0 {
		t.lockedCount -= flipCount
		t.shakeCounter = 10
	}

	if t.iced {
		var count int
		for _, adj := range t.adj {
			if adj == nil || adj.flipped || (adj.mine && adj.flagged) {
				count++
			}
		}
		if count > t.icedCount {
			t.icedCount = count
			t.shakeCounter = 10
			if count < 4 {
				t.gfx.next()
			}
		}
		if t.icedCount >= 4 {
			t.gfx.play()
			if t.gfx.done {
				t.gfx = n_newAniSprite(
					grassImg[:],
					[]uint{10, 5, 5, 5, 5, 5, 5},
					false,
				)
				t.flip()
				t.iced = false
			}
		}
	}
}

func (t *n_tile) shake() {
	t.shakeCounter = 30
}

func (t *n_tile) flip() int {
	if t.mine {
		t.flipped = true
		n_mineDog.play()
		return 1
	}

	if !t.flagged {
		t.gfx.play()
	}
	t.flipped = true
	count := 1

	// flip adjacent tiles if this is a zero
	if t.adjCount == 0 {
		for _, adj := range t.adj {
			if adj != nil {
				count += adj.flip()
			}
		}
	}

	return count
}

func (t *n_tile) flag() {
	if !t.flipped {
		t.flagged = !t.flagged
	}
}

func (t *n_tile) draw(screen *ebiten.Image) {
	// draw tile gfx
	coord := t.parent.boardXY
	coord.x += t.index.Float64().x * 17
	coord.y += t.index.Float64().y * 11
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(coord.x), float64(coord.y))
	if t.shakeCounter > 0 {
		op.GeoM.Translate(math.Sin(float64(tickCounter)), 0)
	}
	t.gfx.draw(screen, op)

	// you just lost
	if t.flipped && t.mine {
		op.GeoM.Translate(-7, -9)
		n_mineDog.draw(screen, op)
		return
	}

	if t.flagged {
		op.GeoM.Translate(-1, -9)
		n_markerFlag.draw(screen, op)
		op.GeoM.Translate(1, 9)
	}

	// draw adj numbers
	if t.flipped && t.adjCount > 0 {
		op.GeoM.Translate(6, 4)
		screen.DrawImage(numberSmall[t.adjCount], op)
		op.GeoM.Translate(-6, -4)
	}

	// draw a locked tile
	if t.flipped && t.lockedCount > 0 {
		screen.DrawImage(lock, op)
		op.GeoM.Translate(7, 5)
		if t.lockedCount > 9 {
			screen.DrawImage(numberSmall[t.lockedCount/10], op)
			op.GeoM.Translate(4, 0)
			screen.DrawImage(numberSmall[t.lockedCount%10], op)
			op.GeoM.Translate(-4, 0)
		} else {
			screen.DrawImage(numberSmall[t.lockedCount], op)
		}
		op.GeoM.Translate(-7, -5)
	}
}
