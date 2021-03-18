package main

import (
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

// var grassAniA []*ebiten.Image
// var grassAniB []*ebiten.Image
// var grassAniC []*ebiten.Image
// var flagAni []*ebiten.Image

type grassTile struct {
	transform
	grass           aniSprite
	eid             uint64
	frames          uint
	isMine          bool
	isBlown         bool
	adjMine         int
	flagged         bool
	flaggedCoolDown int
	flipped         bool
}

func newGrassTile(x, y float64) *grassTile {
	ret := grassTile{
		frames: uint(rand.Intn(30)),
	}
	switch rand.Intn(5) {
	case 0, 1, 2:
		ret.grass = newAniSprite(grassSpr[:], []uint16{3, 3, 3, 3, 3, 3, 3}, false)
	case 3:
		ret.grass = newAniSprite(yellowGrassSpr[:], []uint16{3, 3, 3, 3, 3, 3, 3}, false)
	case 4:
		ret.grass = newAniSprite(pinkGrassSpr[:], []uint16{3, 3, 3, 3, 3, 3, 3}, false)
	}

	ret.x, ret.y = x, y

	return &ret
}

func (g *grassTile) reset() {
	g.grass.reset()
	g.grass.pause()
	g.frames = uint(rand.Intn(15))
	g.isMine = false
	g.isBlown = false
	g.adjMine = 0
	g.flagged = false
	g.flaggedCoolDown = 0
	g.flipped = false
}

func (g *grassTile) update() error {
	g.frames++
	g.flaggedCoolDown--
	g.grass.update()
	return nil
}

func (g *grassTile) draw(screen *ebiten.Image) {
	if g.grass.isPlaying() && g.isMine {
		// don't flip becuase we're actually a mine
		g.grass.reset()
		g.grass.pause()
		dogAni.play()
		g.isBlown = true
	}

	grassY := 6.0
	if g.frames <= 30.0 {
		grassY = lerp(200, 6, float64(g.frames)/30.0)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, g.y)
	screen.DrawImage(grassSpot, op)

	g.grass.draw(screen, g.x, g.y-grassY)

	if g.grass.frame >= 6 && g.adjMine > 0 {
		pxlPrint(screen, mainFont, g.x+6, g.y-5, strconv.Itoa(g.adjMine))
	}

	if g.flagged {
		flagAni.draw(screen, g.x-1, g.y-14)
	}

	if g.isBlown {
		dogAni.draw(screen, g.x-7, g.y-18)
	}
}

func (g *grassTile) setID(eid uint64) {
	g.eid = eid
}

func (g *grassTile) getID() uint64 {
	return g.eid
}

func (g *grassTile) colide(x, y int) bool {
	x1, x2 := int(g.x), int(g.x+16)
	y1, y2 := int(g.y-3), int(g.y+13)
	if x >= x1 && x <= x2 &&
		y >= y1 && y <= y2 {
		return true
	}
	return false
}

// func flipGrassTile(x, y int) {
// 	if x < 0 || y < 0 || x > 7 || y > 7 {
// 		return
// 	}
// 	tile := board[x][y]
// 	if tile.flagged {
// 		return
// 	}
// 	if firstFlip {
// 		firstFlip = false
// 		setupBoard(x, y)
// 	}
// 	if !tile.flipped {
// 		tile.flagged = false
// 		tile.flipped = true
// 		tile.grass.play()
// 		if !tile.isMine && tile.adjMine == 0 {
// 			flipGrassTile(x+1, y+1)
// 			flipGrassTile(x+1, y)
// 			flipGrassTile(x+1, y-1)
// 			flipGrassTile(x-1, y+1)
// 			flipGrassTile(x-1, y)
// 			flipGrassTile(x-1, y-1)
// 			flipGrassTile(x, y+1)
// 			flipGrassTile(x, y-1)
// 		}
// 	}
// }

// func clearAround(x, y int) {
// 	if x < 0 || y < 0 || x > 7 || y > 7 {
// 		return
// 	}
// 	tile := board[x][y]
// 	if tile.adjMine == 0 || !tile.flipped {
// 		return
// 	}
// 	flags := 0
// 	if x > 0 && y > 0 && board[x-1][y-1].flagged {
// 		flags++
// 	}
// 	if x > 0 && y < 7 && board[x-1][y+1].flagged {
// 		flags++
// 	}
// 	if x > 0 && board[x-1][y].flagged {
// 		flags++
// 	}
// 	if x < 7 && y > 0 && board[x+1][y-1].flagged {
// 		flags++
// 	}
// 	if x < 7 && y < 7 && board[x+1][y+1].flagged {
// 		flags++
// 	}
// 	if x < 7 && board[x+1][y].flagged {
// 		flags++
// 	}
// 	if y > 0 && board[x][y-1].flagged {
// 		flags++
// 	}
// 	if y < 7 && board[x][y+1].flagged {
// 		flags++
// 	}

// 	if tile.adjMine == flags {
// 		// bomb it!
// 		if x > 0 && y > 0 && !board[x-1][y-1].flagged {
// 			board[x-1][y-1].flipped = true
// 			board[x-1][y-1].grass.play()
// 		}
// 		if x > 0 && y < 7 && !board[x-1][y+1].flagged {
// 			board[x-1][y+1].flipped = true
// 			board[x-1][y+1].grass.play()
// 		}
// 		if x > 0 && !board[x-1][y].flagged {
// 			board[x-1][y].flipped = true
// 			board[x-1][y].grass.play()
// 		}
// 		if x < 7 && y > 0 && !board[x+1][y-1].flagged {
// 			board[x+1][y-1].flipped = true
// 			board[x+1][y-1].grass.play()
// 		}
// 		if x < 7 && y < 7 && !board[x+1][y+1].flagged {
// 			board[x+1][y+1].flipped = true
// 			board[x+1][y+1].grass.play()
// 		}
// 		if x < 7 && !board[x+1][y].flagged {
// 			board[x+1][y].flipped = true
// 			board[x+1][y].grass.play()
// 		}
// 		if y > 0 && !board[x][y-1].flagged {
// 			board[x][y-1].flipped = true
// 			board[x][y-1].grass.play()
// 		}
// 		if y < 7 && !board[x][y+1].flagged {
// 			board[x][y+1].flipped = true
// 			board[x][y+1].grass.play()
// 		}
// 	}
// }
