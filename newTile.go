package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type n_tile struct {
	index          v2i
	adj            [8]*n_tile
	adjCount       int
	parent         *levelScean
	gfx            *n_aniSprite
	mine           bool
	water          bool
	flagged        bool
	flaggedCount   int
	flipped        bool
	lockedCount    int
	iced           bool
	icedCount      int
	timeTile       bool
	timeTileHeight float64
	shakeCounter   int
	barkCounter    int
	bounce         bool
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
	t.barkCounter--
	t.flaggedCount--

	if t.parent.settings.fadeFlags && t.flaggedCount <= 0 {
		t.flagged = false
	}

	if t.flipped {
		t.timeTileHeight += 3
	}

	if t.flipped && flipCount > 0 && t.lockedCount > 0 {
		t.lockedCount -= flipCount
		t.shakeCounter = 10
	}

	if t.iced && t.parent.filled {
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
				t.iced = false
				for _, adj := range t.adj {
					if adj != nil && adj.adjCount == 0 && adj.flipped {
						t.flip()
						break
					}
				}
			}
		}
	}
}

func (t *n_tile) hovered() bool {
	m := mCoordsF()
	coord := t.parent.boardXY
	coord.x += t.parent.boardDXY.x
	coord.y += t.parent.boardDXY.y
	coord.x += t.index.Float64().x * 17
	coord.y += t.index.Float64().y * 11
	if m.x > coord.x && m.x < coord.x+16 &&
		m.y > coord.y && m.y < coord.y+12 {
		return true
	}
	return false
}

func (t *n_tile) shake() {
	t.shakeCounter = 30
}

func (t *n_tile) flip() int {
	if t.iced || t.flipped || t.flagged {
		return 0
	}

	if t.mine {
		t.flipped = true
		n_mineDog.play()
		return 1
	}
	t.gfx.play()
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
		t.flaggedCount = 600
	}
}

func (t *n_tile) draw(screen *ebiten.Image) {
	// draw tile gfx
	coord := t.parent.boardXY
	coord.x += t.parent.boardDXY.x
	coord.y += t.parent.boardDXY.y
	coord.x += t.index.Float64().x * 17
	coord.y += t.index.Float64().y * 11
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(coord.x), float64(coord.y))
	if t.shakeCounter > 0 {
		op.GeoM.Translate(math.Sin(float64(tickCounter)), 0)
	}
	if t.bounce {
		op.GeoM.Translate(0, -math.Abs(math.Sin(float64(tickCounter)/10)*3))
	}
	t.gfx.draw(screen, op)

	// you just lost
	if t.flipped && t.mine {
		op.GeoM.Translate(-7, -9)
		n_mineDog.draw(screen, op)
		return
	}

	// draw the marker flag
	if t.flagged {
		if t.parent.settings.fadeFlags {
			if t.flaggedCount > 60 {
				op.GeoM.Translate(-1, -9)
				n_markerFlag.draw(screen, op)
				op.GeoM.Translate(1, 9)
			}
			if t.flaggedCount > 0 &&
				t.flaggedCount <= 60 &&
				t.flaggedCount%2 == 0 {
				op.GeoM.Translate(-1, -9)
				n_markerFlag.draw(screen, op)
				op.GeoM.Translate(1, 9)

			}
		} else {
			op.GeoM.Translate(-1, -9)
			n_markerFlag.draw(screen, op)
			op.GeoM.Translate(1, 9)

		}
	}

	// draw the +time pow
	if t.timeTile && t.timeTileHeight < 240 {
		op.GeoM.Translate(0, -t.timeTileHeight)
		screen.DrawImage(timeToken, op)
		op.GeoM.Translate(0, t.timeTileHeight)
	}

	// draw adj numbers
	if t.flipped && t.adjCount > 0 {
		op.GeoM.Translate(6, 4)
		screen.DrawImage(numberSmallWhite[t.adjCount], op)
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

	// draw the barking sprite
	if t.barkCounter > 0 && !t.flipped && !t.flagged {
		op.GeoM.Translate(-8, -24)
		screen.DrawImage(n_dogBark, op)
	}
}
