package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type gameBoard struct {
	transform
	eid        uint64
	rows, cols uint
	field      [][]*grassTile
	populated  bool
	mineCount  uint
}

func newGameBoard(rows, cols, mineCount uint) *gameBoard {
	ret := gameBoard{
		rows:      rows,
		cols:      cols,
		mineCount: mineCount,
	}
	for i := 0; i < int(cols); i++ {
		add := make([]*grassTile, rows)
		ret.field = append(ret.field, add)
	}

	for i := 0; i < int(cols); i++ {
		for ii := 0; ii < int(rows); ii++ {
			tile := newGrassTile(float64(i)*18, float64(ii)*11)
			registerEnt(tile, uint32(ii+1))
			ret.field[ii][i] = tile
		}
	}
	registerEnt(&ret, uint32(rows+2))

	return &ret
}

func (g *gameBoard) click(x, y int) (int, int) {
	for i := 0; i < int(g.cols); i++ {
		for ii := int(g.rows) - 1; ii >= 0; ii-- {
			tile := g.field[ii][i]
			if tile.colide(x, y) {
				return ii, i
			}
		}
	}
	return -1, -1
}

func (g *gameBoard) clearAround(x, y int) {
	if x < 0 || y < 0 || x > int(g.cols)-1 || y > int(g.rows)-1 {
		return
	}
	tile := g.field[x][y]
	if tile.adjMine == 0 || !tile.flipped {
		return
	}
	flags := 0
	if x > 0 && y > 0 && g.field[x-1][y-1].flagged {
		flags++
	}
	if x > 0 && y < int(g.rows)-1 && g.field[x-1][y+1].flagged {
		flags++
	}
	if x > 0 && g.field[x-1][y].flagged {
		flags++
	}
	if x < int(g.cols)-1 && y > 0 && g.field[x+1][y-1].flagged {
		flags++
	}
	if x < int(g.cols)-1 && y < int(g.rows)-1 && g.field[x+1][y+1].flagged {
		flags++
	}
	if x < int(g.cols)-1 && g.field[x+1][y].flagged {
		flags++
	}
	if y > 0 && g.field[x][y-1].flagged {
		flags++
	}
	if y < int(g.rows)-1 && g.field[x][y+1].flagged {
		flags++
	}

	if tile.adjMine == flags {
		// bomb it!
		if x > 0 && y > 0 && !g.field[x-1][y-1].flagged {
			g.field[x-1][y-1].flipped = true
			g.field[x-1][y-1].grass.play()
		}
		if x > 0 && y < int(g.rows)-1 && !g.field[x-1][y+1].flagged {
			g.field[x-1][y+1].flipped = true
			g.field[x-1][y+1].grass.play()
		}
		if x > 0 && !g.field[x-1][y].flagged {
			g.field[x-1][y].flipped = true
			g.field[x-1][y].grass.play()
		}
		if x < int(g.cols)-1 && y > 0 && !g.field[x+1][y-1].flagged {
			g.field[x+1][y-1].flipped = true
			g.field[x+1][y-1].grass.play()
		}
		if x < int(g.cols)-1 && y < int(g.rows)-1 && !g.field[x+1][y+1].flagged {
			g.field[x+1][y+1].flipped = true
			g.field[x+1][y+1].grass.play()
		}
		if x < int(g.cols)-1 && !g.field[x+1][y].flagged {
			g.field[x+1][y].flipped = true
			g.field[x+1][y].grass.play()
		}
		if y > 0 && !g.field[x][y-1].flagged {
			g.field[x][y-1].flipped = true
			g.field[x][y-1].grass.play()
		}
		if y < int(g.rows)-1 && !g.field[x][y+1].flagged {
			g.field[x][y+1].flipped = true
			g.field[x][y+1].grass.play()
		}
	}
}

func (g *gameBoard) flipTiles(x, y int) {
	if x < 0 || y < 0 || x > int(g.cols-1) || y > int(g.rows-1) {
		return
	}
	if !g.populated {
		g.addMines(x, y)
	}
	tile := g.field[x][y]
	if tile.flagged {
		return
	}
	if !tile.flipped {
		tile.flagged = false
		tile.flipped = true
		tile.grass.play()
		if !tile.isMine && tile.adjMine == 0 {
			g.flipTiles(x+1, y+1)
			g.flipTiles(x+1, y)
			g.flipTiles(x+1, y-1)
			g.flipTiles(x-1, y+1)
			g.flipTiles(x-1, y)
			g.flipTiles(x-1, y-1)
			g.flipTiles(x, y+1)
			g.flipTiles(x, y-1)
		}
	}
}

func (g *gameBoard) addMines(safex, safey int) {
	mines := g.mineCount
	for mines > 0 {
		x := rand.Intn(8)
		y := rand.Intn(8)
		if x == safex && y == safey {
			continue
		}
		if g.field[x][y].isMine {
			continue
		}
		g.field[x][y].isMine = true
		if x > 0 && y > 0 {
			g.field[x-1][y-1].adjMine++
		}
		if x > 0 && y < int(g.rows)-1 {
			g.field[x-1][y+1].adjMine++
		}
		if x > 0 {
			g.field[x-1][y].adjMine++
		}
		if x < int(g.cols)-1 && y > 0 {
			g.field[x+1][y-1].adjMine++
		}
		if x < int(g.cols)-1 && y < int(g.rows)-1 {
			g.field[x+1][y+1].adjMine++
		}
		if x < int(g.cols)-1 {
			g.field[x+1][y].adjMine++
		}
		if y > 0 {
			g.field[x][y-1].adjMine++
		}
		if y < int(g.rows)-1 {
			g.field[x][y+1].adjMine++
		}
		if g.field[safex][safey].adjMine > 0 {
			g.field[x][y].isMine = false
			if x > 0 && y > 0 {
				g.field[x-1][y-1].adjMine--
			}
			if x > 0 && y < int(g.rows)-1 {
				g.field[x-1][y+1].adjMine--
			}
			if x > 0 {
				g.field[x-1][y].adjMine--
			}
			if x < int(g.cols)-1 && y > 0 {
				g.field[x+1][y-1].adjMine--
			}
			if x < int(g.cols)-1 && y < int(g.rows)-1 {
				g.field[x+1][y+1].adjMine--
			}
			if x < int(g.cols)-1 {
				g.field[x+1][y].adjMine--
			}
			if y > 0 {
				g.field[x][y-1].adjMine--
			}
			if y < int(g.rows)-1 {
				g.field[x][y+1].adjMine--
			}
			continue
		}
		mines--
	}
	g.populated = true
}

func (g *gameBoard) translate(x, y float64) {
	g.x += x
	g.y += y
	for i := 0; i < int(g.rows); i++ {
		for ii := 0; ii < int(g.cols); ii++ {
			g.field[ii][i].x += x
			g.field[ii][i].y += y
		}
	}
}

func (g *gameBoard) draw(screen *ebiten.Image) {

}

func (g *gameBoard) update() error {

	return nil
}

func (g *gameBoard) setID(eid uint64) {
	g.eid = eid
}

func (g *gameBoard) getID() uint64 {
	return g.eid
}
