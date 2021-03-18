package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type gameBoardScean struct {
	x, y                  float64
	mouseX, mouseY        int
	clickCount            uint
	rows, cols, mineCount uint
	board                 [][]*tile
	filled                bool
	won                   bool
	lost                  bool
	looseSel              int
	paused                bool
	pauseSel              int
	timer                 int64 // timer is an acumulator that will track all elapsed time from previous start/ stop cycles
	start                 int64 // start will be reset each time the game regains focus (e.g. on start/ after pausing)
	shuffel               bool
	shuffelHard           bool
	shuffelDone           bool
	shuffelTime           int
	maxTime               int
	flagCount             int
	levelID               int
	flipLockCount         int
	adjLockCount          int
	power                 int
}

func (gb *gameBoardScean) load() error {
	gb.x = 5
	gb.y = 25

	ss, err := getAsset("assets/sprites_test_2.png")
	if err != nil {
		return err
	}

	grassSprs = [7]*ebiten.Image{
		subImage(ss, 0, 0, 16, 16),
		subImage(ss, 16, 0, 16, 16),
		subImage(ss, 32, 0, 16, 16),
		subImage(ss, 48, 0, 16, 16),
		subImage(ss, 64, 0, 16, 16),
		subImage(ss, 80, 0, 16, 16),
		subImage(ss, 96, 0, 16, 16),
	}
	yellowGrassSprs = [7]*ebiten.Image{
		subImage(ss, 0, 16, 16, 16),
		subImage(ss, 16, 16, 16, 16),
		subImage(ss, 32, 16, 16, 16),
		subImage(ss, 48, 16, 16, 16),
		subImage(ss, 64, 16, 16, 16),
		subImage(ss, 80, 16, 16, 16),
		subImage(ss, 96, 16, 16, 16),
	}
	pinkGrassSprs = [7]*ebiten.Image{
		subImage(ss, 0, 32, 16, 16),
		subImage(ss, 16, 32, 16, 16),
		subImage(ss, 32, 32, 16, 16),
		subImage(ss, 48, 32, 16, 16),
		subImage(ss, 64, 32, 16, 16),
		subImage(ss, 80, 32, 16, 16),
		subImage(ss, 96, 32, 16, 16),
	}

	markerFlag = newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 0, 112, 16, 16),
			subImage(ss, 16, 112, 16, 16),
			subImage(ss, 32, 112, 16, 16),
			subImage(ss, 48, 112, 16, 16),
		},
		[]uint16{8, 8, 8, 8},
		true,
	)
	markerFlag.play()

	mineDog = newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 112, 0, 32, 16),
			subImage(ss, 144, 0, 32, 16),
			subImage(ss, 176, 0, 32, 16),
			subImage(ss, 208, 0, 32, 16),
		},
		[]uint16{4, 4, 4, 4},
		false,
	)

	dogBark = subImage(ss, 112, 16, 32, 32)

	addPU = subImage(ss, 64, 96, 16, 16)
	minusPU = subImage(ss, 80, 96, 16, 16)
	waterPU = subImage(ss, 96, 97, 16, 16)
	bombPU = subImage(ss, 112, 96, 16, 16)
	wistlePU = subImage(ss, 129, 96, 16, 16)

	for i := 0; i < int(gb.cols); i++ {
		add := make([]*tile, gb.rows)
		gb.board = append(gb.board, add)
	}

	for i := 0; i < int(gb.cols); i++ {
		for ii := 0; ii < int(gb.rows); ii++ {
			gb.board[i][ii] = newTile(gb, gb.x+float64(i)*18, gb.y+float64(ii)*11)
		}
	}

	// set up the adj tiles
	for i := 0; i < int(gb.cols); i++ {
		for ii := 0; ii < int(gb.rows); ii++ {
			tmp := gb.board[i][ii]
			if i > 0 && ii > 0 {
				tmp.adjTiles[0] = gb.board[i-1][ii-1]
			}
			if ii > 0 {
				tmp.adjTiles[1] = gb.board[i][ii-1]
			}
			if i < int(gb.cols)-1 && ii > 0 {
				tmp.adjTiles[2] = gb.board[i+1][ii-1]
			}
			if i > 0 {
				tmp.adjTiles[3] = gb.board[i-1][ii]
			}
			if i < int(gb.cols)-1 {
				tmp.adjTiles[4] = gb.board[i+1][ii]
			}
			if i > 0 && ii < int(gb.rows)-1 {
				tmp.adjTiles[5] = gb.board[i-1][ii+1]
			}
			if ii < int(gb.rows)-1 {
				tmp.adjTiles[6] = gb.board[i][ii+1]
			}
			if i < int(gb.cols)-1 && ii < int(gb.rows)-1 {
				tmp.adjTiles[7] = gb.board[i+1][ii+1]
			}
		}
	}

	return nil
}

func (gb *gameBoardScean) unload() error {
	// we probably want to keep the most of this stuff in
	// memory since we're going to be using them alot for the main game loop
	return nil
}

func (gb *gameBoardScean) update() error {
	var mouseDX, mouseDY int
	leftClick := false
	maxClickLen := uint(12)

	if mbtn(ebiten.MouseButtonLeft) {
		gb.clickCount++
	}

	if !mbtn(ebiten.MouseButtonLeft) {
		if gb.clickCount > 0 && gb.clickCount <= maxClickLen {
			leftClick = true
		}
		gb.clickCount = 0
	}

	x, y := ebiten.CursorPosition()
	cursorHold = false
	if gb.clickCount > maxClickLen {
		mouseDX, mouseDY = x-gb.mouseX, y-gb.mouseY
		cursorHold = true
	}
	gb.mouseX, gb.mouseY = x, y

	if gb.lost {

		gb.x += float64(mouseDX)
		gb.y += float64(mouseDY)

		for i := 0; i < int(gb.cols); i++ {
			for ii := 0; ii < int(gb.rows); ii++ {
				gb.board[i][ii].op.GeoM.Translate(float64(mouseDX), float64(mouseDY))
			}
		}

		mineDog.op.GeoM.Translate(float64(mouseDX), float64(mouseDY))
		mineDog.update()
		markerFlag.update()

		if btnp(ebiten.KeyDown) || btnp(ebiten.KeyS) {
			gb.looseSel++
		}
		if btnp(ebiten.KeyUp) || btnp(ebiten.KeyW) {
			gb.looseSel--
		}
		if gb.looseSel < 0 {
			gb.looseSel = 1
		}
		if gb.looseSel > 1 {
			gb.looseSel = 0
		}

		if btnp(ebiten.KeyEnter) && gb.looseSel == 0 {
			currentScean = &gameBoardScean{
				rows:          gb.rows,
				cols:          gb.cols,
				mineCount:     gb.mineCount,
				maxTime:       gb.maxTime,
				shuffel:       gb.shuffel,
				shuffelHard:   gb.shuffelHard,
				shuffelTime:   gb.shuffelTime,
				adjLockCount:  gb.adjLockCount,
				flipLockCount: gb.flipLockCount,
				levelID:       gb.levelID,
			}
			var err error
			err = currentScean.load()
			if err != nil {
				return err
			}

			err = gb.unload()
			if err != nil {
				return err
			}
		}

		if btnp(ebiten.KeyEscape) || (btnp(ebiten.KeyEnter) && gb.looseSel == 1) {
			currentScean = &levelSelectScean{}
			var err error
			err = currentScean.load()
			if err != nil {
				return err
			}

			err = gb.unload()
			if err != nil {
				return err
			}
		}

		return nil

	}

	// The game was won
	if gb.won {
		for i := 0; i < int(gb.cols); i++ {
			for ii := 0; ii < int(gb.rows); ii++ {
				gb.board[i][ii].grass.update()
				if !gb.board[i][ii].flagged {
					gb.board[i][ii].flip()
				}
			}
		}

		gb.x += float64(mouseDX)
		gb.y += float64(mouseDY)

		for i := 0; i < int(gb.cols); i++ {
			for ii := 0; ii < int(gb.rows); ii++ {
				gb.board[i][ii].op.GeoM.Translate(float64(mouseDX), float64(mouseDY))
			}
		}

		markerFlag.update()

		if btnp(ebiten.KeyEscape) || btnp(ebiten.KeyEnter) {
			currentScean = &levelSelectScean{}
			var err error
			err = currentScean.load()
			if err != nil {
				return err
			}

			err = gb.unload()
			if err != nil {
				return err
			}
		}

		return nil
	}

	if btnp(ebiten.KeyEscape) {
		gb.paused = !gb.paused
		if gb.paused && gb.filled {
			gb.timer += time.Now().Unix() - gb.start
			gb.start = 0
		}
		if !gb.paused && gb.filled {
			gb.start = time.Now().Unix()
		}
	}

	if gb.paused {
		if btnp(ebiten.KeyEnter) {
			if gb.pauseSel == 0 {
				gb.paused = false
				if gb.filled {
					gb.start = time.Now().Unix()
				}
			}
			if gb.pauseSel == 1 {
				currentScean = &gameBoardScean{
					rows:          gb.rows,
					cols:          gb.cols,
					mineCount:     gb.mineCount,
					maxTime:       gb.maxTime,
					shuffel:       gb.shuffel,
					shuffelHard:   gb.shuffelHard,
					shuffelTime:   gb.shuffelTime,
					adjLockCount:  gb.adjLockCount,
					flipLockCount: gb.flipLockCount,
					levelID:       gb.levelID,
				}
				var err error
				err = currentScean.load()
				if err != nil {
					return nil
				}

				err = gb.unload()
				if err != nil {
					return err
				}
			}
			if gb.pauseSel == 2 {
				currentScean = &levelSelectScean{}
				var err error
				err = currentScean.load()
				if err != nil {
					return err
				}

				err = gb.unload()
				if err != nil {
					return err
				}
			}
		}
		if btnp(ebiten.KeyW) || btnp(ebiten.KeyUp) {
			gb.pauseSel--
		}
		if btnp(ebiten.KeyS) || btnp(ebiten.KeyDown) {
			gb.pauseSel++
		}
		if gb.pauseSel < 0 {
			gb.pauseSel = 2
		}
		if gb.pauseSel > 2 {
			gb.pauseSel = 0
		}
		// Don't do any board updates if were paused
		return nil
	}

	if btnp(ebiten.Key1) {
		if gb.power == 1 {
			gb.power = 0
		} else {
			gb.power = 1
		}
	}
	if btnp(ebiten.Key2) {
		if gb.power == 2 {
			gb.power = 0
		} else {
			gb.power = 2
		}
	}
	if btnp(ebiten.Key3) {
		if gb.power == 3 {
			gb.power = 0
		} else {
			gb.power = 3
		}
	}
	if btnp(ebiten.Key4) {
		if gb.power == 4 {
			gb.power = 0
		} else {
			gb.power = 4
		}
	}
	if btnp(ebiten.Key5) {
		if gb.power == 5 {
			gb.power = 0
		} else {
			gb.power = 5
		}
	}

	var overTile *tile
	cx, cy := ebiten.CursorPosition()
	for i := 0; i < int(gb.cols); i++ {
		for ii := int(gb.rows) - 1; ii >= 0; ii-- {
			tmp := gb.board[i][ii]
			tx := int(getX(tmp.op))
			ty := int(getY(tmp.op))
			x1, x2 := tx, tx+16
			y1, y2 := ty, ty+13
			if cx >= x1 && cx <= x2 &&
				cy >= y1 && cy <= y2 {
				overTile = tmp
				break
			}
		}
		if overTile != nil {
			break
		}
	}

	if overTile != nil && overTile.flipped && mbtn(ebiten.MouseButtonLeft) && mbtn(ebiten.MouseButtonRight) {
		// do the double clicky thingy (give some feed back even if it fails)
		flags := 0
		for _, adj := range overTile.adjTiles {
			if adj != nil && adj.flagged {
				flags++
			}
		}
		if overTile.flipLock > 0 {
			flags = -1
		}
		if flags == overTile.adjCount {
			flipped := false
			for _, adj := range overTile.adjTiles {
				if adj != nil {
					before := adj.flipped
					adj.flip()
					if !before && adj.flipped {
						flipped = true
					}
				}
			}

			if flipped && gb.flipLockCount > 0 {
				for i := 0; i < int(gb.cols); i++ {
					for ii := 0; ii < int(gb.rows); ii++ {
						if gb.board[i][ii].flipped {
							gb.board[i][ii].flipLock--
						}
					}
				}
			}
		}
	}

	if overTile != nil && leftClick && !mbtnp(ebiten.MouseButtonRight) {
		// check for the tile I'm clicking on and flip it
		if gb.power == 0 {
			if !gb.filled {
				gb.filled = true
				overTile.safe = true
				gb.start = time.Now().Unix()
				gb.addMines()
			}
			before := overTile.flipped
			overTile.flip()

			if !before && overTile.flipped {
				if gb.flipLockCount > 0 {
					for i := 0; i < int(gb.cols); i++ {
						for ii := 0; ii < int(gb.rows); ii++ {
							if gb.board[i][ii].flipped {
								gb.board[i][ii].flipLock--
							}
						}
					}
				}
			}
		}

		// add one mine to the surounding tiles
		if gb.power == 1 && overTile.flipped && overTile.adjCount > 0 {
			var openMines []*tile
			for _, adj := range overTile.adjTiles {
				if adj != nil && !adj.mine && !adj.flipped && !adj.flagged && !adj.water {
					openMines = append(openMines, adj)
				}
			}
			if len(openMines) > 0 {
				newMine := openMines[rand.Intn(len(openMines))]
				newMine.mine = true
				for _, adj := range newMine.adjTiles {
					if adj != nil {
						adj.adjCount++
					}
				}
				gb.mineCount++
				gb.power = 0
			}
		}

		// minus one mine from the surounding tiles
		if gb.power == 2 && overTile.flipped && overTile.adjCount > 0 {
			var openMines []*tile
			for _, adj := range overTile.adjTiles {
				if adj != nil && adj.mine && !adj.flagged {
					openMines = append(openMines, adj)
				}
			}
			if len(openMines) > 0 {
				oldMine := openMines[rand.Intn(len(openMines))]
				oldMine.mine = false
				gb.mineCount--
				for _, adj := range oldMine.adjTiles {
					if adj != nil {
						adj.adjCount--
						if adj.adjCount == 0 && !adj.mine {
							adj.flipped = false
							adj.flip()
						}
					}
				}
				gb.power = 0
			}
		}

		// water gun power up, soak a bunch of tiles
		if gb.power == 3 {
			var waterTiles []*tile
			for i := 0; i < int(gb.cols); i++ {
				for ii := 0; ii < int(gb.rows); ii++ {
					tile := gb.board[i][ii]
					if !tile.flagged && !tile.flipped && !tile.mine {
						waterTiles = append(waterTiles, tile)
					}
				}
			}
			if len(waterTiles) > 0 {
				for i := 0; i < 4; i++ {
					waterTiles[rand.Intn(len(waterTiles))].water = true
				}
				gb.power = 0
			}
		}

		// mini bomb power up, clear a 3x3 area
		if gb.power == 4 && !overTile.flipped {
			fmt.Println("HERE")
			if overTile.mine {
				gb.mineCount--
				for _, adj := range overTile.adjTiles {
					adj.adjCount--
				}
				overTile.mine = false
				overTile.flip()
			}

			for _, adj := range overTile.adjTiles {
				if adj != nil {
					if adj.mine {
						gb.mineCount--
						for _, metaAdj := range adj.adjTiles {
							metaAdj.adjCount--
						}
						adj.mine = false
					}
					adj.flip()
				}
			}
			gb.power = 0
		}

		// dog wistle powerup, reveal one random mine
		if gb.power == 5 {
			var mines []*tile
			for i := 0; i < int(gb.cols); i++ {
				for ii := 0; ii < int(gb.rows); ii++ {
					if gb.board[i][ii].mine && !gb.board[i][ii].flagged {
						mines = append(mines, gb.board[i][ii])
					}
				}
			}
			if len(mines) > 0 {
				mines[rand.Intn(len(mines))].doBark = 600
				gb.power = 0
			}
		}
	}

	if overTile != nil && !mbtnp(ebiten.MouseButtonLeft) && mbtnp(ebiten.MouseButtonRight) {
		// check for the tile I'm clicking on and flag/ unflag it
		overTile.flag()
	}

	var tx, ty float64
	speed := 1.2
	if btn(ebiten.KeyA) || btn(ebiten.KeyLeft) {
		tx = -speed
	}
	if btn(ebiten.KeyD) || btn(ebiten.KeyRight) {
		tx = speed
	}
	if btn(ebiten.KeyW) || btn(ebiten.KeyUp) {
		ty = -speed
	}
	if btn(ebiten.KeyS) || btn(ebiten.KeyDown) {
		ty = speed
	}

	gb.x += tx + float64(mouseDX)
	gb.y += ty + float64(mouseDY)

	// do the sprite updates
	markerFlag.update()
	mineDog.update()

	var unflipped uint
	var found uint
	gb.flagCount = 0
	for i := 0; i < int(gb.cols); i++ {
		for ii := 0; ii < int(gb.rows); ii++ {
			gb.board[i][ii].op.GeoM.Translate(tx+float64(mouseDX), ty+float64(mouseDY))
			gb.board[i][ii].grass.update()
			gb.board[i][ii].doBark--

			if gb.board[i][ii].blown {
				gb.lost = true
			}

			if gb.board[i][ii].flagged {
				gb.flagCount++
			}

			if gb.board[i][ii].mine && gb.board[i][ii].flagged {
				found++
			}

			if !gb.board[i][ii].flipped {
				unflipped++
			}
		}
	}

	if found == gb.mineCount && uint(gb.flagCount) == gb.mineCount {
		gb.won = true
		levelCleared[gb.levelID] = true
	}

	if unflipped == gb.mineCount {
		gb.won = true
		levelCleared[gb.levelID] = true
	}

	if gb.shuffel {
		total := gb.timer
		if gb.start > 1 {
			total += (time.Now().Unix() - gb.start)
		}
		if total > 0 && total%int64(gb.shuffelTime) == 0 && !gb.shuffelDone {
			if gb.shuffelHard {
				gb.shuffleMinesHard()
			} else {
				gb.shuffleMines()
			}

			for i := 0; i < int(gb.cols); i++ {
				for ii := 0; ii < int(gb.rows); ii++ {
					if gb.board[i][ii].flipped && gb.board[i][ii].adjCount == 0 {
						gb.board[i][ii].flipped = false
						gb.board[i][ii].flip()
					}
				}
			}

			gb.shuffelDone = true
		}
		if (total-1)%int64(gb.shuffelTime) == 0 {
			gb.shuffelDone = false
		}
	}

	if gb.maxTime > 0 {
		total := gb.timer
		if gb.start > 1 {
			total += (time.Now().Unix() - gb.start)
		}
		if total > int64(gb.maxTime) {
			gb.lost = true
		}
	}

	return nil
}

func (gb *gameBoardScean) draw(screen *ebiten.Image) {
	if gb.paused {
		pxlPrint(screen, mainFont, 15, 5, "Resume")
		pxlPrint(screen, mainFont, 15, 15, "Restart Level")
		pxlPrint(screen, mainFont, 15, 25, "Exit Level")
		pxlPrint(screen, mainFont, 5, float64(5+10*gb.pauseSel), ">")
		return
	}

	for i := 0; i < int(gb.cols); i++ {
		for ii := 0; ii < int(gb.rows); ii++ {
			gb.board[i][ii].draw(screen)
		}
	}

	if !gb.won && !gb.lost {
		total := gb.timer
		if gb.start > 0 {
			total += (time.Now().Unix() - gb.start)
		}
		if gb.maxTime > 0 {
			total = int64(gb.maxTime) - total
		}

		min := strconv.Itoa(int(total / 60))
		if len(min) == 1 {
			min = "0" + min
		}
		sec := strconv.Itoa(int(total % 60))
		if len(sec) == 1 {
			sec = "0" + sec
		}
		pxlPrint(screen, mainFont, 5, 3, min+":"+sec)

		markerFlag.op.GeoM.Reset()
		markerFlag.op.GeoM.Translate(185, -4)
		markerFlag.draw(screen)
		pxlPrint(screen, mainFont, 200, 3, "x"+strconv.Itoa(gb.flagCount)+"/"+strconv.Itoa(int(gb.mineCount)))

		power := ""
		switch gb.power {
		case 1:
			power = "power: add mine"
		case 2:
			power = "power: minus mine"
		case 3:
			power = "power: watter gun"
		case 4:
			power = "power: mini bomb"
		case 5:
			power = "power: dog wistle"
		}

		pxlPrint(screen, mainFont, 80, 3, power)
	}

	if gb.lost {

		mineDog.draw(screen)

		bg := ebiten.NewImage(110, 38)
		bg.Fill(color.RGBA{255, 204, 170, 255}) // pico-8 light tan
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(60, 0)
		screen.DrawImage(bg, op)

		pxlPrint(screen, mainFont, 80, 3, "A Looser is You!")
		pxlPrint(screen, mainFont, 100, 13, "Retry")
		pxlPrint(screen, mainFont, 100, 23, "Exit Level")
		pxlPrint(screen, mainFont, 95, float64(13+10*gb.looseSel), ">")

		return
	}

	if gb.won {
		bg := ebiten.NewImage(110, 15)
		bg.Fill(color.RGBA{255, 204, 170, 255}) // pico-8 light tan
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(60, 0)
		screen.DrawImage(bg, op)

		pxlPrint(screen, mainFont, 80, 3, "A Winner is You!")

		return
	}
}

func (gb *gameBoardScean) addMines() {
	mines := gb.mineCount
	for mines > 0 {
		i := rand.Intn(int(gb.cols))
		ii := rand.Intn(int(gb.rows))
		tile := gb.board[i][ii]
		valid := !tile.mine && !tile.safe

		for _, adj := range tile.adjTiles {
			if adj != nil {
				if adj.safe {
					valid = false
					break
				}
			}
		}

		if !valid {
			continue
		}

		tile.mine = true
		for _, adj := range tile.adjTiles {
			if adj != nil {
				adj.adjCount++
			}
		}
		mines--
	}

	adjLock := gb.adjLockCount
	fmt.Println("adjLock:", gb.adjLockCount)
	for adjLock > 0 {
		fmt.Println("adding adj lock", adjLock)
		i := rand.Intn(int(gb.cols))
		ii := rand.Intn(int(gb.rows))
		tile := gb.board[i][ii]
		valid := tile.adjCount > 0 && !tile.mine

		if valid {
			tile.adjLock = 4
			adjLock--
		}
	}

	flippLock := gb.flipLockCount
	fmt.Println("flipLock:", gb.flipLockCount)
	for flippLock > 0 {
		fmt.Println("adding flip lock", flippLock)
		i := rand.Intn(int(gb.cols))
		ii := rand.Intn(int(gb.rows))
		tile := gb.board[i][ii]
		valid := tile.adjCount > 0 && !tile.mine

		if valid {
			tile.flipLock = 3
			flippLock--
		}
	}
}

func (gb *gameBoardScean) shuffleMines() {
	fmt.Println("DO SHUFFEL")
	for i := 0; i < int(gb.cols); i++ {
		for ii := 0; ii < int(gb.rows); ii++ {
			gb.board[i][ii].adjCount = 0
		}
	}

	mines := gb.mineCount

	for i := 0; i < int(gb.cols); i++ {
		for ii := 0; ii < int(gb.rows); ii++ {
			tile := gb.board[i][ii]

			// leave flagged mines alone
			// if tile.flagged && tile.mine {
			// 	mines--

			// 	for _, adj := range tile.adjTiles {
			// 		if adj != nil {
			// 			adj.adjCount++
			// 		}
			// 	}
			// 	continue
			// }
			tile.mine = false
			tile.flagged = false
		}
	}

	for mines > 0 {
		i := rand.Intn(int(gb.cols))
		ii := rand.Intn(int(gb.rows))
		tile := gb.board[i][ii]
		valid := !tile.mine && !tile.flipped
		if !valid {
			continue
		}

		tile.mine = true
		for _, adj := range tile.adjTiles {
			if adj != nil {
				adj.adjCount++
			}
		}
		mines--
	}
}

func (gb *gameBoardScean) shuffleMinesHard() {
	fmt.Println("DO HARD SHUFFEL")
	for i := 0; i < int(gb.cols); i++ {
		for ii := 0; ii < int(gb.rows); ii++ {
			tile := gb.board[i][ii]
			tile.adjCount = 0
			tile.mine = false
			tile.flagged = false
			if tile.flipped && rand.Intn(25) == 1 {
				tile.flipped = false
				tile.grass.reset()
				tile.grass.pause()
			}
		}
	}

	mines := gb.mineCount
	for mines > 0 {
		i := rand.Intn(int(gb.cols))
		ii := rand.Intn(int(gb.rows))
		tile := gb.board[i][ii]
		valid := !tile.mine && !tile.flipped
		if !valid {
			continue
		}

		tile.mine = true
		for _, adj := range tile.adjTiles {
			if adj != nil {
				adj.adjCount++
			}
		}
		mines--
	}
}

func getX(op *ebiten.DrawImageOptions) float64 {
	return op.GeoM.Element(0, 2)
}

func getY(op *ebiten.DrawImageOptions) float64 {
	return op.GeoM.Element(1, 2)
}
