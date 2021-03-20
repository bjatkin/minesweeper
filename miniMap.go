package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// pico 8 colors
var (
	green     = color.RGBA{0, 228, 54, 255}
	darkGreen = color.RGBA{0, 135, 81, 255}
	gray      = color.RGBA{95, 87, 79, 255}
	darkGray  = color.RGBA{29, 43, 83, 255}
)

type miniMap struct {
	coord            v2f
	scale            int // 1x is a 14 x 14 window, 2x is a 7x7 window, 3x could be a 4x4 or 5x5 window
	tiles            *[]n_tile
	width, height    int
	flippedTiles     [32][32]uint8
	grassTiles       [32][32]uint8
	parent           *levelScean
	mines            int
	flippedTileCount int
	flaggedTileCount int
	tileCount        int
}

func newMiniMap(level *levelScean, coord v2f, tiles *[]n_tile, mines int) *miniMap {
	maxX, minX := 0, 9999
	maxY, minY := 0, 9999
	for _, tile := range *tiles {
		if tile.index.x > maxX {
			maxX = tile.index.x
		}
		if tile.index.x < minX {
			minX = tile.index.x
		}
		if tile.index.y > maxY {
			maxY = tile.index.y
		}
		if tile.index.y < minY {
			minY = tile.index.y
		}
	}
	width, height := maxX-minX, maxY-minY
	scale := 1
	if width > 32 || height > 32 {
		scale = 2
	}
	if width > 64 || height > 64 {
		scale = 3
	}

	return &miniMap{
		coord:     coord,
		tiles:     tiles,
		mines:     mines,
		scale:     scale,
		parent:    level,
		width:     width,
		height:    height,
		tileCount: len(*tiles),
	}
}

func (mm *miniMap) update() {
	mm.flippedTileCount = 0
	mm.flaggedTileCount = 0
	for x := 0; x < 32; x++ {
		for y := 0; y < 32; y++ {
			mm.flippedTiles[x][y] = 0
			mm.grassTiles[x][y] = 0
		}
	}

	for _, tile := range *mm.tiles {
		if tile.flipped {
			mm.flippedTiles[tile.index.x/mm.scale][tile.index.y/mm.scale]++
			mm.flippedTileCount++
		} else {
			mm.grassTiles[tile.index.x/mm.scale][tile.index.y/mm.scale]++
		}
		if tile.flagged {
			mm.flaggedTileCount++
		}
	}
}

func (mm *miniMap) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(mm.coord.x, mm.coord.y)
	screen.DrawImage(miniMapHud, op)

	miniMap := ebiten.NewImage(32, 32)
	for x := 0; x < 32; x++ {
		for y := 0; y < 32; y++ {
			grass := mm.grassTiles[x][y]
			flip := mm.flippedTiles[x][y]
			total := grass + flip
			if int(total) > (mm.scale*mm.scale)/2 {
				if grass > flip {
					// put a green pixel here
					miniMap.Set(x, y, green)
				} else {
					// put a flipped pixel here
					miniMap.Set(x, y, gray)
				}
				continue
			}

			if y > 0 && miniMap.At(x, y-1) == green {
				miniMap.Set(x, y, darkGreen)
			}
			if y > 0 && miniMap.At(x, y-1) == gray {
				miniMap.Set(x, y, darkGray)
			}
		}
	}

	offset := v2f{float64((32 - mm.width) / 2), float64((32 - mm.height) / 2)}
	op.GeoM.Translate(offset.x, offset.y+4)
	screen.DrawImage(miniMap, op)

	// total mines
	cop := &ebiten.DrawImageOptions{}
	cop.GeoM.Translate(mm.coord.x+67, mm.coord.y+9)
	if mm.mines > 99 {
		screen.DrawImage(numberSmallGray[mm.mines/100], cop)
		cop.GeoM.Translate(4, 0)
		screen.DrawImage(numberSmallGray[(mm.mines%100)/10], cop)
		cop.GeoM.Translate(4, 0)
		screen.DrawImage(numberSmallGray[(mm.mines%100)%10], cop)
	} else if mm.mines > 9 {
		screen.DrawImage(numberSmallGray[mm.mines/10], cop)
		cop.GeoM.Translate(4, 0)
		screen.DrawImage(numberSmallGray[mm.mines%10], cop)
	} else {
		screen.DrawImage(numberSmallGray[mm.mines], cop)
	}

	// flagged tiles
	cop.GeoM.Reset()
	cop.GeoM.Translate(mm.coord.x+50, mm.coord.y+9)
	if mm.flaggedTileCount > 99 {
		screen.DrawImage(numberBigBlue[mm.flaggedTileCount/100], cop)
		cop.GeoM.Translate(6, 0)
		screen.DrawImage(numberBigBlue[(mm.flaggedTileCount%100)/10], cop)
		cop.GeoM.Translate(6, 0)
		screen.DrawImage(numberBigBlue[(mm.flaggedTileCount%100)%10], cop)
	} else if mm.mines > 9 {
		screen.DrawImage(numberBigBlue[mm.flaggedTileCount/10], cop)
		cop.GeoM.Translate(6, 0)
		screen.DrawImage(numberBigBlue[mm.flaggedTileCount%10], cop)
	} else {
		screen.DrawImage(numberBigBlue[mm.flaggedTileCount], cop)
	}

	// total tiles - mines
	cop.GeoM.Reset()
	cop.GeoM.Translate(mm.coord.x+71, mm.coord.y+24)
	totalTiles := mm.tileCount - mm.mines
	if totalTiles > 99 {
		screen.DrawImage(numberSmallGray[totalTiles/100], cop)
		cop.GeoM.Translate(4, 0)
		screen.DrawImage(numberSmallGray[(totalTiles%100)/10], cop)
		cop.GeoM.Translate(4, 0)
		screen.DrawImage(numberSmallGray[(totalTiles%100)%10], cop)
	} else if mm.mines > 9 {
		screen.DrawImage(numberSmallGray[totalTiles/10], cop)
		cop.GeoM.Translate(4, 0)
		screen.DrawImage(numberSmallGray[totalTiles%10], cop)
	} else {
		screen.DrawImage(numberSmallGray[totalTiles], cop)
	}

	// flipped tiles
	cop.GeoM.Reset()
	cop.GeoM.Translate(mm.coord.x+53, mm.coord.y+24)
	if mm.flippedTileCount > 99 {
		screen.DrawImage(numberBigBlue[mm.flippedTileCount/100], cop)
		cop.GeoM.Translate(6, 0)
		screen.DrawImage(numberBigBlue[(mm.flaggedTileCount%100)/10], cop)
		cop.GeoM.Translate(6, 0)
		screen.DrawImage(numberBigBlue[(mm.flippedTileCount%100)%10], cop)
	} else if mm.mines > 9 {
		screen.DrawImage(numberBigBlue[mm.flippedTileCount/10], cop)
		cop.GeoM.Translate(6, 0)
		screen.DrawImage(numberBigBlue[mm.flippedTileCount%10], cop)
	} else {
		screen.DrawImage(numberBigBlue[mm.flippedTileCount], cop)
	}
}
