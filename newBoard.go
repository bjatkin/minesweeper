package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type levelScean struct {
	boardXY          v2f
	boardDXY         v2f
	mouseAnchor      v2f
	clickCount       int
	panning          bool
	mouseXY          v2i
	board            *[]n_tile
	filled           bool
	win, loose       bool
	paused           bool
	bestTime         *timer
	quit             *uiButton
	restart          *uiButton
	timerAccumulator int64
	start            int64
	flagCount        int
	settings         *n_levelData
	usingPowerUp     bool
	usingPowerUpID   int
	powerUps         [3]*powerUp
	powerUpTypes     [3]int

	miniMap          *miniMap
	duckCharacterUI  *duckFeedBack
	levelTimer       *boardTimer
	levelStarCounter *starCounter
}

func newLevelScean(data *n_levelData, powerUpTypes [3]int) *levelScean {
	ret := &levelScean{
		settings:     data,
		powerUpTypes: powerUpTypes,
	}

	minX, maxX := 999999, 0
	minY, maxY := 999999, 0
	tiles := make([]n_tile, 0, len(data.layout))
	for _, tile := range data.layout {
		tiles = append(tiles, *n_newTile(ret, tile.index, false, false, false))
		if tile.index.x < minX {
			minX = tile.index.x
		}
		if tile.index.x > maxX {
			maxX = tile.index.x
		}
		if tile.index.y < minY {
			minY = tile.index.y
		}
		if tile.index.y > maxY {
			maxY = tile.index.y
		}
	}
	width, height := maxX-minX, maxY-minY
	x := float64(width)/2*17 + 8
	y := float64(height)/2*11 + 5
	ret.boardXY = v2f{120 - x, 70 - y}
	// get the middle tile to the middle of the screen

	for i := 0; i < len(tiles); i++ {
		for ii, adj := range data.layout[i].adj {
			if adj > -1 {
				tiles[i].adj[ii] = &tiles[adj]
			}
		}
	}

	ret.board = &tiles

	return ret
}

// Level Assets, all of these should be loaded from the main sprite sheet
var (
	grassImg       [7]*ebiten.Image
	yellowGrassImg [7]*ebiten.Image
	pinkGrassImg   [7]*ebiten.Image
	blueGrassImg   [7]*ebiten.Image
	waterImg       [7]*ebiten.Image
	iceImg         [3]*ebiten.Image
	lock           *ebiten.Image
	n_mineDog      *n_aniSprite
	n_dogBark      *ebiten.Image
	// TODO: dead duck needs an update to make his beak a little lighter
	// cool duck could probably also use a little more dark brown in his hair
	sideDuck        [4]*ebiten.Image
	activePowerUp   [7]*ebiten.Image
	inactivePowerUp [7]*ebiten.Image
	n_markerFlag    *n_aniSprite
	powerUpHud      *ebiten.Image
	powerUpHudBG    *ebiten.Image
	numberBig       [12]*ebiten.Image
	numberSmall     [11]*ebiten.Image
	numberBigBlue   [12]*ebiten.Image
	numberSmallBlue [11]*ebiten.Image
	numberBigGray   [12]*ebiten.Image
	numberSmallGray [11]*ebiten.Image
	miniMapHud      *ebiten.Image
	starCountHud    *ebiten.Image
	starCountHudBG  *ebiten.Image
	star            *ebiten.Image
	pauseMenu       *ebiten.Image
	restartBtn      [3]*ebiten.Image
	quitBtn         [3]*ebiten.Image
	timerBG         *ebiten.Image
	timerPlayBtn    [3]*ebiten.Image
	timerPauseBtn   [3]*ebiten.Image
)

func (l *levelScean) load() error {
	ss, err := getAsset("assets/sprite_sheet.png")
	if err != nil {
		return err
	}

	grassImg = [7]*ebiten.Image{
		subImage(ss, 0, 192, 16, 16),
		subImage(ss, 16, 192, 16, 16),
		subImage(ss, 32, 192, 16, 16),
		subImage(ss, 48, 192, 16, 16),
		subImage(ss, 64, 192, 16, 16),
		subImage(ss, 80, 192, 16, 16),
		subImage(ss, 96, 192, 16, 16),
	}

	yellowGrassImg = [7]*ebiten.Image{
		subImage(ss, 0, 208, 16, 16),
		subImage(ss, 16, 208, 16, 16),
		subImage(ss, 32, 208, 16, 16),
		subImage(ss, 48, 208, 16, 16),
		subImage(ss, 64, 208, 16, 16),
		subImage(ss, 80, 208, 16, 16),
		subImage(ss, 96, 208, 16, 16),
	}

	pinkGrassImg = [7]*ebiten.Image{
		subImage(ss, 0, 224, 16, 16),
		subImage(ss, 16, 224, 16, 16),
		subImage(ss, 32, 224, 16, 16),
		subImage(ss, 48, 224, 16, 16),
		subImage(ss, 64, 224, 16, 16),
		subImage(ss, 80, 224, 16, 16),
		subImage(ss, 96, 224, 16, 16),
	}

	blueGrassImg = [7]*ebiten.Image{
		subImage(ss, 0, 240, 16, 16),
		subImage(ss, 16, 240, 16, 16),
		subImage(ss, 32, 240, 16, 16),
		subImage(ss, 48, 240, 16, 16),
		subImage(ss, 64, 240, 16, 16),
		subImage(ss, 80, 240, 16, 16),
		subImage(ss, 96, 240, 16, 16),
	}

	waterImg = [7]*ebiten.Image{
		subImage(ss, 0, 256, 16, 16),
		subImage(ss, 16, 256, 16, 16),
		subImage(ss, 32, 256, 16, 16),
		subImage(ss, 48, 256, 16, 16),
		subImage(ss, 64, 256, 16, 16),
		subImage(ss, 80, 256, 16, 16),
		subImage(ss, 96, 256, 16, 16),
	}

	iceImg = [3]*ebiten.Image{
		subImage(ss, 112, 208, 16, 16),
		subImage(ss, 128, 208, 16, 16),
		subImage(ss, 112, 224, 16, 16),
	}

	lock = subImage(ss, 112, 240, 16, 16)

	n_mineDog = n_newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 0, 272, 32, 16),
			subImage(ss, 32, 272, 32, 16),
			subImage(ss, 64, 272, 32, 16),
			subImage(ss, 96, 272, 32, 16),
		},
		[]uint{3, 3, 3, 3},
		false,
	)

	n_dogBark = subImage(ss, 128, 256, 32, 32)

	sideDuck = [4]*ebiten.Image{
		subImage(ss, 0, 288, 40, 40),
		subImage(ss, 40, 288, 40, 40),
		subImage(ss, 80, 288, 40, 40),
		subImage(ss, 120, 288, 40, 40),
	}

	activePowerUp = [7]*ebiten.Image{
		subImage(ss, 0, 0, 16, 16),
		subImage(ss, 16, 0, 16, 16),
		subImage(ss, 32, 0, 16, 16),
		subImage(ss, 48, 0, 16, 16),
		subImage(ss, 64, 0, 16, 16),
		subImage(ss, 80, 0, 16, 16),
		subImage(ss, 96, 0, 16, 16),
	}

	inactivePowerUp = [7]*ebiten.Image{
		subImage(ss, 112, 0, 16, 16),
		subImage(ss, 128, 0, 16, 16),
		subImage(ss, 144, 0, 16, 16),
		subImage(ss, 160, 0, 16, 16),
		subImage(ss, 176, 0, 16, 16),
		subImage(ss, 192, 0, 16, 16),
		subImage(ss, 208, 0, 16, 16),
	}

	n_markerFlag = n_newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 0, 328, 16, 16),
			subImage(ss, 16, 328, 16, 16),
			subImage(ss, 32, 328, 16, 16),
			subImage(ss, 48, 328, 16, 16),
		},
		[]uint{8, 8, 8, 8},
		true,
	)
	n_markerFlag.play()

	powerUpHud = subImage(ss, 0, 152, 104, 40)
	powerUpHudBG = subImage(ss, 104, 152, 40, 40)

	numberBig = [12]*ebiten.Image{
		subImage(ss, 80, 16, 8, 8),
		subImage(ss, 88, 16, 8, 8),
		subImage(ss, 96, 16, 8, 8),
		subImage(ss, 104, 16, 8, 8),
		subImage(ss, 112, 16, 8, 8),
		subImage(ss, 120, 16, 8, 8),
		subImage(ss, 128, 16, 8, 8),
		subImage(ss, 136, 16, 8, 8),
		subImage(ss, 144, 16, 8, 8),
		subImage(ss, 152, 16, 8, 8),
		subImage(ss, 160, 16, 8, 8),
		subImage(ss, 168, 16, 8, 8),
	}

	numberSmall = [11]*ebiten.Image{
		subImage(ss, 80, 24, 8, 8),
		subImage(ss, 88, 24, 8, 8),
		subImage(ss, 96, 24, 8, 8),
		subImage(ss, 104, 24, 8, 8),
		subImage(ss, 112, 24, 8, 8),
		subImage(ss, 120, 24, 8, 8),
		subImage(ss, 128, 24, 8, 8),
		subImage(ss, 136, 24, 8, 8),
		subImage(ss, 144, 24, 8, 8),
		subImage(ss, 152, 24, 8, 8),
		subImage(ss, 160, 24, 8, 8),
	}

	numberBigBlue = [12]*ebiten.Image{
		subImage(ss, 272, 32, 8, 8),
		subImage(ss, 280, 32, 8, 8),
		subImage(ss, 288, 32, 8, 8),
		subImage(ss, 296, 32, 8, 8),
		subImage(ss, 304, 32, 8, 8),
		subImage(ss, 312, 32, 8, 8),
		subImage(ss, 320, 32, 8, 8),
		subImage(ss, 328, 32, 8, 8),
		subImage(ss, 336, 32, 8, 8),
		subImage(ss, 344, 32, 8, 8),
		subImage(ss, 352, 32, 8, 8),
		subImage(ss, 352, 32, 8, 8),
	}
	numberSmallBlue = [11]*ebiten.Image{
		subImage(ss, 272, 40, 8, 8),
		subImage(ss, 280, 40, 8, 8),
		subImage(ss, 288, 40, 8, 8),
		subImage(ss, 296, 40, 8, 8),
		subImage(ss, 304, 40, 8, 8),
		subImage(ss, 312, 40, 8, 8),
		subImage(ss, 320, 40, 8, 8),
		subImage(ss, 328, 40, 8, 8),
		subImage(ss, 336, 40, 8, 8),
		subImage(ss, 344, 40, 8, 8),
		subImage(ss, 352, 40, 8, 8),
	}

	numberBigGray = [12]*ebiten.Image{
		subImage(ss, 272, 48, 8, 8),
		subImage(ss, 280, 48, 8, 8),
		subImage(ss, 288, 48, 8, 8),
		subImage(ss, 296, 48, 8, 8),
		subImage(ss, 304, 48, 8, 8),
		subImage(ss, 312, 48, 8, 8),
		subImage(ss, 320, 48, 8, 8),
		subImage(ss, 328, 48, 8, 8),
		subImage(ss, 336, 48, 8, 8),
		subImage(ss, 344, 48, 8, 8),
		subImage(ss, 352, 48, 8, 8),
		subImage(ss, 352, 48, 8, 8),
	}

	numberSmallGray = [11]*ebiten.Image{
		subImage(ss, 272, 56, 8, 8),
		subImage(ss, 280, 56, 8, 8),
		subImage(ss, 288, 56, 8, 8),
		subImage(ss, 296, 56, 8, 8),
		subImage(ss, 304, 56, 8, 8),
		subImage(ss, 312, 56, 8, 8),
		subImage(ss, 320, 56, 8, 8),
		subImage(ss, 328, 56, 8, 8),
		subImage(ss, 336, 56, 8, 8),
		subImage(ss, 344, 56, 8, 8),
		subImage(ss, 352, 56, 8, 8),
	}

	addMine = [2]*ebiten.Image{
		subImage(ss, 0, 0, 16, 16),
		subImage(ss, 112, 0, 16, 16),
	}
	scaredyCat = [2]*ebiten.Image{
		subImage(ss, 16, 0, 16, 16),
		subImage(ss, 128, 0, 16, 16),
	}
	tidalWave = [2]*ebiten.Image{
		subImage(ss, 32, 0, 16, 16),
		subImage(ss, 144, 0, 16, 16),
	}
	minusMine = [2]*ebiten.Image{
		subImage(ss, 48, 0, 16, 16),
		subImage(ss, 160, 0, 16, 16),
	}
	dogWistle = [2]*ebiten.Image{
		subImage(ss, 64, 0, 16, 16),
		subImage(ss, 176, 0, 16, 16),
	}
	shuffel = [2]*ebiten.Image{
		subImage(ss, 80, 0, 16, 16),
		subImage(ss, 192, 0, 16, 16),
	}
	dogABone = [2]*ebiten.Image{
		subImage(ss, 96, 0, 16, 16),
		subImage(ss, 208, 0, 16, 16),
	}

	// TODO: add a set of white numbers

	// TODO: the mini map hud needs to be changed to support 3 digit flags and tile counts
	// so we don't have GUI issues on larger maps
	miniMapHud = subImage(ss, 184, 32, 88, 40)
	starCountHud = subImage(ss, 80, 48, 40, 11)
	starCountHudBG = subImage(ss, 80, 59, 40, 11)
	star = subImage(ss, 136, 48, 16, 16)
	pauseMenu = subImage(ss, 184, 72, 64, 72)

	restartBtn = [3]*ebiten.Image{
		subImage(ss, 248, 104, 56, 16), // normal
		subImage(ss, 248, 88, 56, 16),  // hover
		subImage(ss, 248, 72, 56, 16),  // clicked
	}

	quitBtn = [3]*ebiten.Image{
		subImage(ss, 248, 136, 56, 16), // normal
		subImage(ss, 248, 120, 56, 16), // hover
		subImage(ss, 248, 72, 56, 16),  // clicked
	}

	timerBG = subImage(ss, 96, 72, 49, 16)

	timerPlayBtn = [3]*ebiten.Image{
		subImage(ss, 80, 72, 16, 16),  // normal
		subImage(ss, 80, 88, 16, 16),  // hover
		subImage(ss, 128, 88, 16, 16), // clicked
	}

	timerPauseBtn = [3]*ebiten.Image{
		subImage(ss, 112, 88, 16, 16), // normal
		subImage(ss, 96, 88, 16, 16),  // hover
		subImage(ss, 128, 88, 16, 16), // clicked
	}

	// TODO: the background power ups need to have a darker background
	// to make the charge up animation more obvious

	// TODO: we nee a blank power up for when you don't have all 3
	// powerups in a match
	addMine = [2]*ebiten.Image{
		subImage(ss, 0, 0, 16, 16),
		subImage(ss, 112, 0, 16, 16),
	}
	scaredyCat = [2]*ebiten.Image{
		subImage(ss, 16, 0, 16, 16),
		subImage(ss, 128, 0, 16, 16),
	}
	tidalWave = [2]*ebiten.Image{
		subImage(ss, 32, 0, 16, 16),
		subImage(ss, 144, 0, 16, 16),
	}
	minusMine = [2]*ebiten.Image{
		subImage(ss, 48, 0, 16, 16),
		subImage(ss, 160, 0, 16, 16),
	}
	dogWistle = [2]*ebiten.Image{
		subImage(ss, 64, 0, 16, 16),
		subImage(ss, 176, 0, 16, 16),
	}
	shuffel = [2]*ebiten.Image{
		subImage(ss, 80, 0, 16, 16),
		subImage(ss, 192, 0, 16, 16),
	}
	dogABone = [2]*ebiten.Image{
		subImage(ss, 96, 0, 16, 16),
		subImage(ss, 208, 0, 16, 16),
	}

	l.miniMap = newMiniMap(l, v2f{0, 124}, l.board, l.settings.mineCount)
	l.levelTimer = newBoardTimer(v2f{})
	l.levelStarCounter = &starCounter{
		coord:         v2f{0, 20},
		timer:         l.levelTimer.timer,
		oneStarTime:   l.settings.starTimes[0],
		twoStarTime:   l.settings.starTimes[1],
		threeStarTime: l.settings.starTimes[2],
	}

	l.powerUps = [3]*powerUp{
		newPowerUp(l.powerUpTypes[0], ebiten.Key1, l.levelTimer.timer),
		newPowerUp(l.powerUpTypes[1], ebiten.Key2, l.levelTimer.timer),
		newPowerUp(l.powerUpTypes[2], ebiten.Key3, l.levelTimer.timer),
	}
	l.duckCharacterUI = newDuckFeedBack(
		v2f{138, 121},
		l.powerUps[0],
		l.powerUps[1],
		l.powerUps[2],
	)

	l.restart = newUIButton(v2f{82, 45}, restartBtn)
	l.quit = newUIButton(v2f{82, 65}, quitBtn)
	l.bestTime = &timer{timerAccumulator: l.settings.bestTime, coord: v2f{92, 96}}

	return nil
}

func (l *levelScean) unload() error {
	// we want to leave this stuff loaded since we'll be jumping in and
	// out of levels a bunch
	return nil
}

func (l *levelScean) update() error {
	// check for a win first thing so we get the lowest possible time
	var flagged int
	var flipped int
	for _, tile := range *l.board {
		if tile.flagged && tile.mine {
			flagged++
		}
		if tile.flipped {
			flipped++
		}
	}
	if flagged == l.settings.mineCount ||
		flipped == len(*l.board)-l.settings.mineCount {
		l.win = true
		l.levelTimer.timer.stop()
		l.duckCharacterUI.state = duckCool
	}

	// check the power ups
	if !l.usingPowerUp {
		if l.powerUps[0].wasSelected() {
			l.usingPowerUp = true
			l.usingPowerUpID = 0
		}
		if l.powerUps[1].wasSelected() {
			l.usingPowerUp = true
			l.usingPowerUpID = 1
		}
		if l.powerUps[2].wasSelected() {
			l.usingPowerUp = true
			l.usingPowerUpID = 2
		}
	}

	// check if were flipping a tile
	var flipCount int
	if mbtnr(ebiten.MouseButtonLeft) &&
		l.mouseAnchor.dist(mCoordsF()) < 5 &&
		l.clickCount < 30 &&
		!l.usingPowerUp &&
		!l.paused {
		minX := 9999999
		var selTile *n_tile
		for i, tile := range *l.board {
			if tile.hovered() && tile.index.x < minX {
				minX = tile.index.x
				selTile = &(*l.board)[i]
			}
		}

		if selTile != nil {
			if !l.filled {
				l.fillBoard(selTile)
				// make all the power ups available
				for i := 0; i < len(l.powerUps); i++ {
					l.powerUps[i].available = true
				}
				l.levelTimer.timer.start()
				l.filled = true
			}

			if !selTile.flipped {
				l.duckCharacterUI.state = duckSurprised
				l.duckCharacterUI.surprised = 30
				flipCount = selTile.flip()
				if flipCount == 0 {
					selTile.shake()
				}
				if selTile.mine && selTile.flipped {
					l.loose = true
					l.duckCharacterUI.state = duckDead
					l.levelTimer.timer.stop()
				}
			} else {
				var flags int
				for i := 0; i < 8; i++ {
					if selTile.adj[i] != nil && selTile.adj[i].flagged {
						flags++
					}
				}
				if flags == selTile.adjCount {
					l.duckCharacterUI.state = duckSurprised
					l.duckCharacterUI.surprised = 30
					for i := 0; i < 8; i++ {
						if selTile.adj[i] != nil && !selTile.adj[i].flagged {
							selTile.adj[i].flip()
							if selTile.adj[i].mine && selTile.adj[i].flipped {
								l.loose = true
								l.duckCharacterUI.state = duckDead
								l.levelTimer.timer.stop()
							}
						}
					}
				} else {
					for i := 0; i < 8; i++ {
						if selTile.adj[i] != nil && !selTile.adj[i].flipped {
							selTile.adj[i].shake()
						}
					}
				}
			}
		}
	}

	for i := 0; i < len(*l.board); i++ {
		(*l.board)[i].update(flipCount)
	}

	// finish checking for power up stuff
	if l.usingPowerUp {
		switch l.powerUps[l.usingPowerUpID].pType {
		case addMinePow:
			if doAddMinePow(l.board) {
				l.powerUps[l.usingPowerUpID].activte()
				l.usingPowerUp = false
			}
		case scaredyCatPow:
			if doScaredCat(l.board) {
				l.powerUps[l.usingPowerUpID].activte()
				l.usingPowerUp = false
			}
		case tidalWavePow:
			l.powerUps[l.usingPowerUpID].activte()
			l.usingPowerUp = false
			doTidalWave(l.board)
		case minusMinePow:
			if doMinusMinePow(l.board) {
				l.powerUps[l.usingPowerUpID].activte()
				l.usingPowerUp = false
			}
		case dogWistlePow:
			l.powerUps[l.usingPowerUpID].activte()
			l.usingPowerUp = false
			doDogWistle(l.board)
		case shuffelPow:
			l.powerUps[l.usingPowerUpID].activte()
			l.usingPowerUp = false
			doBoardShuffel(l.board)
		case dogABonePow:
			// this can only be activated when we loose
			l.usingPowerUp = false
		}
	}

	// flag a tile
	if mbtnp(ebiten.MouseButtonRight) && !l.paused {
		minX := 9999999
		var selTile *n_tile
		for i, tile := range *l.board {
			if tile.hovered() && tile.index.x < minX {
				minX = tile.index.x
				selTile = &(*l.board)[i]
			}
		}
		if selTile != nil {
			selTile.flag()
		}
	}

	// pan the board around
	cursorHold = false
	if mbtn(ebiten.MouseButtonLeft) {
		l.clickCount++
		cursorHold = true
	}

	if l.panning {
		mouse := mCoordsF()
		l.boardDXY.x = mouse.x - l.mouseAnchor.x
		l.boardDXY.y = mouse.y - l.mouseAnchor.y
	}
	if mbtnp(ebiten.MouseButtonLeft) {
		l.clickCount = 0
		l.panning = true
		l.mouseAnchor = mCoordsF()
	}
	if l.panning && !mbtn(ebiten.MouseButtonLeft) {
		l.panning = false
		l.boardXY.x += l.boardDXY.x
		l.boardXY.y += l.boardDXY.y
		l.boardDXY = v2f{}
	}

	// update the mini map and level timer + some other assets
	l.miniMap.update()
	l.levelTimer.update()
	if l.levelTimer.play.clicked {
		l.paused = false
	}
	if l.levelTimer.pause.clicked {
		l.paused = true
	}
	l.duckCharacterUI.update()
	n_mineDog.update()
	n_markerFlag.update()

	for _, tile := range *l.board {
		tile.update(0)
	}

	// the game is paused
	if l.paused {
		l.restart.update()
		l.quit.update()
		if l.restart.clicked {
			// restart the board
			currentScean = newLevelScean(l.settings, l.powerUpTypes)
			err := currentScean.load()
			if err != nil {
				return err
			}

			err = l.unload()
			if err != nil {
				return err
			}
		}
		if l.quit.clicked {
			// quit to the map
			currentScean = &levelSelect{}
			err := currentScean.load()
			if err != nil {
				return err
			}

			err = l.unload()
			if err != nil {
				return err
			}
		}
	}

	if l.win && btnp(ebiten.KeyEnter) {
		l.settings.beaten = true
		allLevels[l.settings.nextLevel].unlocked = true
		// quit to the map
		currentScean = &levelSelect{}
		err := currentScean.load()
		if err != nil {
			return err
		}

		err = l.unload()
		if err != nil {
			return err
		}
	}

	// dog a bone game saver
	for i := 0; i < 3; i++ {
		if l.loose && l.powerUps[i].ready && l.powerUps[i].pType == dogABonePow {
			l.loose = false
			l.powerUps[i].activte()
			l.duckCharacterUI.state = duckNormal
			l.levelTimer.timer.start()
			for i := 0; i < len(*l.board); i++ {
				if (*l.board)[i].mine && (*l.board)[i].flipped {
					(*l.board)[i].mine = false
					(*l.board)[i].flagged = true
					(*l.board)[i].flipped = false
					n_mineDog.pause()
					break
				}
			}
		}
	}

	return nil
}

func (l *levelScean) draw(screen *ebiten.Image) {
	if !l.paused || !l.filled {
		redraw := []int{}

		for i, tile := range *l.board {
			tile.draw(screen)
			if tile.mine && tile.flipped {
				redraw = append(redraw, i)
			}
			if tile.barkCounter > 0 {
				redraw = append(redraw, i)
			}
		}

		for _, tIndex := range redraw {
			(*l.board)[tIndex].draw(screen)
		}
	}

	if l.paused {
		// draw the pause menu
		pop := &ebiten.DrawImageOptions{}
		pop.GeoM.Translate(80, 42)
		screen.DrawImage(pauseMenu, pop)
		l.restart.draw(screen)
		l.quit.draw(screen)
		l.bestTime.draw(screen)
	}

	l.miniMap.draw(screen)
	l.levelTimer.draw(screen)
	l.levelStarCounter.draw(screen)
	l.duckCharacterUI.draw(screen)
}

func (l *levelScean) fillBoard(safe *n_tile) {
	mines := l.settings.mineCount
	for mines > 0 {
		i := rand.Intn(len(*l.board))
		if safe == &(*l.board)[i] || (*l.board)[i].mine {
			continue
		}
		valid := true
		for ii := 0; ii < 8; ii++ {
			if safe == (*l.board)[i].adj[ii] {
				valid = false
			}
		}
		if !valid {
			continue
		}

		(*l.board)[i].mine = true
		for ii := 0; ii < 8; ii++ {
			if (*l.board)[i].adj[ii] != nil {
				(*l.board)[i].adj[ii].adjCount++
			}
		}
		mines--
	}
}

func doTidalWave(board *[]n_tile) {
	soak := 4
	var saftey int
	for soak > 0 {
		saftey++
		if saftey > 100 {
			// if we've looked this hard for tiles to replace and
			// still havent found enough then just leave so we don't
			// get stuck here forever
			return
		}
		i := rand.Intn(len(*board))
		if (*board)[i].mine ||
			(*board)[i].flipped ||
			(*board)[i].flagged {
			continue
		}
		(*board)[i].water = true
		(*board)[i].gfx = n_newAniSprite(
			waterImg[:],
			[]uint{5, 5, 5, 5, 5, 5, 5},
			false,
		)
		soak--
	}
}

func doDogWistle(board *[]n_tile) {
	candidates := []int{}
	for i, tile := range *board {
		if tile.mine && !tile.flipped && !tile.flagged {
			candidates = append(candidates, i)
		}
	}

	mine := &(*board)[candidates[rand.Intn(len(candidates))]]
	mine.barkCounter = 240
	mine.shake()
}

func doBoardShuffel(board *[]n_tile) {
	for _, tile := range *board {
		if tile.mine {
			var done bool
			for !done {
				target := (*board)[rand.Intn(len(*board))]
				if !target.mine && !target.flipped && !target.flagged {
					tile.mine = false
					target.mine = true
					for _, adj := range tile.adj {
						if adj != nil {
							adj.adjCount--
						}
					}
					for _, adj := range target.adj {
						if adj != nil {
							adj.adjCount++
						}
					}
					done = true
					break
				}
			}
		}
		for _, tile := range *board {
			if tile.adjCount == 0 && tile.flipped {
				tile.flipped = false
				tile.flip()
			}
		}
	}
}

func doScaredCat(board *[]n_tile) bool {
	for i := 0; i < len(*board); i++ {
		if !(*board)[i].flipped {
			(*board)[i].bounce = true
		}
	}
	if mbtnr(ebiten.MouseButtonLeft) {
		minX := 999999
		var selTile *n_tile
		for i, tile := range *board {
			if tile.hovered() && tile.index.x < minX {
				minX = tile.index.x
				selTile = &(*board)[i]
			}
		}

		if selTile != nil && !selTile.flipped {
			if selTile.mine {
				for _, adj := range selTile.adj {
					if adj != nil {
						adj.adjCount--
					}
				}
			}
			selTile.mine = false
			selTile.flip()
			for _, adj := range selTile.adj {
				if adj != nil {
					if adj.mine {
						for _, adjadj := range adj.adj {
							if adjadj != nil {
								adjadj.adjCount--
							}
						}
					}
					adj.mine = false
					adj.flip()
				}
			}
			for _, adj := range selTile.adj {
				if adj != nil && adj.adjCount == 0 {
					adj.flipped = false
					adj.flip()
				}
			}
			for i := 0; i < len(*board); i++ {
				(*board)[i].bounce = false
			}
			return true
		}
	}
	return false
}

func doMinusMinePow(board *[]n_tile) bool {
	for i := 0; i < len(*board); i++ {
		if (*board)[i].flipped && (*board)[i].adjCount > 0 {
			(*board)[i].bounce = true
		}
	}
	if mbtnr(ebiten.MouseButtonLeft) {
		minX := 999999
		var selTile *n_tile
		for i, tile := range *board {
			if tile.hovered() && tile.index.x < minX {
				minX = tile.index.x
				selTile = &(*board)[i]
			}
		}

		if selTile != nil && selTile.adjCount > 0 {
			var done bool
			for !done {
				tile := selTile.adj[rand.Intn(8)]
				if tile != nil && tile.mine {
					tile.mine = false
					for _, adj := range tile.adj {
						if adj != nil {
							adj.adjCount--
							adj.flipped = false
							adj.flip()
							adj.bounce = false
						}
					}
					done = true
					break
				}
			}
			for i := 0; i < len(*board); i++ {
				if (*board)[i].flipped && (*board)[i].adjCount > 0 {
					(*board)[i].bounce = false
				}
			}
			return true
		}
	}
	return false
}

func doAddMinePow(board *[]n_tile) bool {
	for i := 0; i < len(*board); i++ {
		if (*board)[i].flipped && (*board)[i].adjCount > 0 {
			(*board)[i].bounce = true
		}
	}
	if mbtnr(ebiten.MouseButtonLeft) {
		minX := 9999999
		var selTile *n_tile
		for i, tile := range *board {
			if tile.hovered() && tile.index.x < minX {
				minX = tile.index.x
				selTile = &(*board)[i]
			}
		}

		if selTile != nil && selTile.adjCount > 0 {
			var candidates int
			for _, tile := range selTile.adj {
				if tile != nil && !tile.flipped && !tile.mine {
					candidates++
				}
			}
			if candidates > 0 {
				var done bool
				for !done {
					tile := selTile.adj[rand.Intn(8)]
					if tile != nil && !tile.mine && !tile.flipped {
						tile.mine = true
						for _, adj := range tile.adj {
							if adj != nil {
								adj.adjCount++
							}
						}
						done = true
						break
					}
				}
				for i := 0; i < len(*board); i++ {
					if (*board)[i].flipped && (*board)[i].adjCount > 0 {
						(*board)[i].bounce = false
					}
				}
				return true
			} else {
				selTile.shake()
			}
		}
	}
	return false
}
