package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type levelScean struct {
	boardXY                 v2f
	boardDXY                v2f
	boardWidth, boardHeight int
	mouseAnchor             v2f
	clickCount              int
	panning                 bool
	mouseXY                 v2i
	board                   *[]n_tile
	filled                  bool
	win, loose              bool
	paused                  bool
	bestTime                *timer
	quit                    *uiButton
	restart                 *uiButton
	continueGame            *uiButton
	timerAccumulator        int64
	start                   int64
	flagCount               int
	settings                *n_levelData
	jeepIndexReturn         int
	levelIndexReturn        int
	usingPowerUp            bool
	usingPowerUpID          int
	powerUps                [3]*powerUp
	powerUpTypes            [3]int
	powSelDone              bool

	mineCount        int
	miniMap          *miniMap
	duckCharacterUI  *duckFeedBack
	levelTimer       *boardTimer
	levelStarCounter *starCounter
}

func newLevelScean(data *n_levelData, powerUpTypes [3]int, jeepIndexReturn int, levelIndexReturn int) *levelScean {
	ret := &levelScean{
		settings:         data,
		powerUpTypes:     powerUpTypes,
		jeepIndexReturn:  jeepIndexReturn,
		mineCount:        data.mineCount,
		levelIndexReturn: levelIndexReturn,
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
	ret.boardWidth = width
	ret.boardHeight = height
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
	grassImg         [7]*ebiten.Image
	yellowGrassImg   [7]*ebiten.Image
	pinkGrassImg     [7]*ebiten.Image
	blueGrassImg     [7]*ebiten.Image
	waterImg         [7]*ebiten.Image
	iceImg           [3]*ebiten.Image
	lock             *ebiten.Image
	n_mineDog        *n_aniSprite
	n_dogBark        *ebiten.Image
	sideDuck         [4]*ebiten.Image
	activePowerUp    [7]*ebiten.Image
	inactivePowerUp  [7]*ebiten.Image
	n_markerFlag     *n_aniSprite
	powerUpHud       *ebiten.Image
	powerUpHudBG     *ebiten.Image
	numberBig        [12]*ebiten.Image
	numberSmall      [11]*ebiten.Image
	numberBigBlue    [12]*ebiten.Image
	numberSmallBlue  [11]*ebiten.Image
	numberBigGray    [12]*ebiten.Image
	numberSmallGray  [11]*ebiten.Image
	numberBigWhite   [12]*ebiten.Image
	numberSmallWhite [11]*ebiten.Image
	miniMapHud       *ebiten.Image
	starCountHud     *ebiten.Image
	starCountHudBG   *ebiten.Image
	star             *ebiten.Image
	pauseMenu        *ebiten.Image
	restartBtn       [3]*ebiten.Image
	continueBtn      [3]*ebiten.Image
	quitBtn          [3]*ebiten.Image
	timerBG          *ebiten.Image
	timerPlayBtn     [3]*ebiten.Image
	timerPauseBtn    [3]*ebiten.Image
	timeToken        *ebiten.Image
	scrollArrowLeft  *ebiten.Image
	scrollArrowUp    *ebiten.Image
	winMenu          *ebiten.Image
	looseMenu        *ebiten.Image
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
		subImage(ss, 288, 32, 8, 8),
		subImage(ss, 296, 32, 8, 8),
		subImage(ss, 304, 32, 8, 8),
		subImage(ss, 312, 32, 8, 8),
		subImage(ss, 320, 32, 8, 8),
		subImage(ss, 328, 32, 8, 8),
		subImage(ss, 336, 32, 8, 8),
		subImage(ss, 344, 32, 8, 8),
		subImage(ss, 352, 32, 8, 8),
		subImage(ss, 360, 32, 8, 8),
		subImage(ss, 368, 32, 8, 8),
		subImage(ss, 376, 32, 8, 8),
	}
	numberSmallBlue = [11]*ebiten.Image{
		subImage(ss, 288, 40, 8, 8),
		subImage(ss, 296, 40, 8, 8),
		subImage(ss, 304, 40, 8, 8),
		subImage(ss, 312, 40, 8, 8),
		subImage(ss, 320, 40, 8, 8),
		subImage(ss, 328, 40, 8, 8),
		subImage(ss, 336, 40, 8, 8),
		subImage(ss, 344, 40, 8, 8),
		subImage(ss, 352, 40, 8, 8),
		subImage(ss, 360, 40, 8, 8),
		subImage(ss, 368, 40, 8, 8),
	}

	numberBigGray = [12]*ebiten.Image{
		subImage(ss, 288, 48, 8, 8),
		subImage(ss, 296, 48, 8, 8),
		subImage(ss, 304, 48, 8, 8),
		subImage(ss, 312, 48, 8, 8),
		subImage(ss, 320, 48, 8, 8),
		subImage(ss, 328, 48, 8, 8),
		subImage(ss, 336, 48, 8, 8),
		subImage(ss, 344, 48, 8, 8),
		subImage(ss, 352, 48, 8, 8),
		subImage(ss, 360, 48, 8, 8),
		subImage(ss, 368, 48, 8, 8),
		subImage(ss, 376, 48, 8, 8),
	}

	numberSmallGray = [11]*ebiten.Image{
		subImage(ss, 288, 56, 8, 8),
		subImage(ss, 296, 56, 8, 8),
		subImage(ss, 304, 56, 8, 8),
		subImage(ss, 312, 56, 8, 8),
		subImage(ss, 320, 56, 8, 8),
		subImage(ss, 328, 56, 8, 8),
		subImage(ss, 336, 56, 8, 8),
		subImage(ss, 344, 56, 8, 8),
		subImage(ss, 352, 56, 8, 8),
		subImage(ss, 360, 56, 8, 8),
		subImage(ss, 368, 56, 8, 8),
	}

	numberBigWhite = [12]*ebiten.Image{
		subImage(ss, 304, 64, 8, 8),
		subImage(ss, 312, 64, 8, 8),
		subImage(ss, 320, 64, 8, 8),
		subImage(ss, 328, 64, 8, 8),
		subImage(ss, 336, 64, 8, 8),
		subImage(ss, 344, 64, 8, 8),
		subImage(ss, 352, 64, 8, 8),
		subImage(ss, 360, 64, 8, 8),
		subImage(ss, 368, 64, 8, 8),
		subImage(ss, 376, 64, 8, 8),
		subImage(ss, 384, 64, 8, 8),
		subImage(ss, 392, 64, 8, 8),
	}

	numberSmallWhite = [11]*ebiten.Image{
		subImage(ss, 304, 72, 8, 8),
		subImage(ss, 312, 72, 8, 8),
		subImage(ss, 320, 72, 8, 8),
		subImage(ss, 328, 72, 8, 8),
		subImage(ss, 336, 72, 8, 8),
		subImage(ss, 344, 72, 8, 8),
		subImage(ss, 352, 72, 8, 8),
		subImage(ss, 360, 72, 8, 8),
		subImage(ss, 368, 72, 8, 8),
		subImage(ss, 376, 72, 8, 8),
		subImage(ss, 384, 72, 8, 8),
	}

	startBtn = [3]*ebiten.Image{
		subImage(ss, 80, 136, 66, 16),
		subImage(ss, 80, 120, 66, 16),
		subImage(ss, 80, 104, 66, 16),
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
	locked = [2]*ebiten.Image{
		subImage(ss, 288, 0, 16, 16),
		subImage(ss, 288, 0, 16, 16),
	}
	timeToken = subImage(ss, 0, 352, 16, 16)

	miniMapHud = subImage(ss, 184, 32, 94, 40)
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

	continueBtn = [3]*ebiten.Image{
		subImage(ss, 304, 152, 66, 16), // normal
		subImage(ss, 304, 136, 66, 16), // hover
		subImage(ss, 304, 120, 66, 16), // clicked
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

	winMenu = subImage(ss, 0, 368, 122, 96)
	looseMenu = subImage(ss, 304, 80, 122, 40)

	l.miniMap = newMiniMap(l, v2f{0, 124}, l.board, l.mineCount)
	l.levelTimer = newBoardTimer(v2f{})
	if l.settings.timeTrial {
		l.levelTimer.timer.countDown = true
		l.levelTimer.timer.maxTime = l.settings.timeTrialMaxTime
	}
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

	scrollArrowLeft = subImage(ss, 16, 352, 8, 8)
	scrollArrowUp = subImage(ss, 24, 352, 8, 8)

	l.restart = newUIButton(v2f{82, 45}, restartBtn)
	l.quit = newUIButton(v2f{82, 65}, quitBtn)
	l.continueGame = newUIButton(v2f{87, 72}, continueBtn)
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
	if !l.loose {
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
		if flagged == l.mineCount ||
			flipped == len(*l.board)-l.mineCount {
			l.win = true
			l.levelTimer.timer.stop()
			l.duckCharacterUI.state = duckCool
		}
	}

	// check if we are out of time
	if l.levelTimer.timer.overTime() {
		l.loose = true
		l.levelTimer.timer.stop()
	}

	// check the power ups
	if !l.usingPowerUp {
		if l.powerUps[0].wasSelected() {
			l.usingPowerUp = true
			l.powSelDone = true
			if mbtn(ebiten.MouseButtonLeft) {
				l.powSelDone = false
			}
			l.usingPowerUpID = 0
		}
		if l.powerUps[1].wasSelected() {
			l.usingPowerUp = true
			l.powSelDone = true
			if mbtn(ebiten.MouseButtonLeft) {
				l.powSelDone = false
			}
			l.usingPowerUpID = 1
		}
		if l.powerUps[2].wasSelected() {
			l.usingPowerUp = true
			l.powSelDone = true
			if mbtn(ebiten.MouseButtonLeft) {
				l.powSelDone = false
			}
			l.usingPowerUpID = 2
		}
	}

	// check if were flipping a tile
	var flipCount int
	if !l.win && !l.loose &&
		mbtnr(ebiten.MouseButtonLeft) &&
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

			// Initalize the board on the first flip
			if !l.filled {
				l.fillBoard(selTile)
				// make all the power ups available
				for i := 0; i < len(l.powerUps); i++ {
					l.powerUps[i].available = true
				}
				l.levelTimer.timer.start()
				l.filled = true

				if l.settings.frozenTileCount > 0 {
					freezeTiles(l.board, l.settings.frozenTileCount)
				}
				if l.settings.lockedTileCount > 0 {
					lockTiles(l.board, l.settings.lockedTileCount)
				}
				if l.settings.timeTrial {
					addTimeTiles(l.board, l.settings.timeTrialCount)
				}
			}

			// this is a virgin tile we are looking to flip
			if !selTile.flipped {
				flipCount = selTile.flip()
				if flipCount > 0 {
					l.duckCharacterUI.state = duckSurprised
					l.duckCharacterUI.surprised = 30
				}
				if flipCount == 0 {
					selTile.shake()
				}
				if selTile.mine && selTile.flipped {
					l.loose = true
					l.duckCharacterUI.state = duckDead
					l.levelTimer.timer.stop()
				}
			} else { // here we are trying to flip adjcent tiles
				var flags int
				for i := 0; i < 8; i++ {
					if selTile.adj[i] != nil && selTile.adj[i].flagged {
						flags++
					}
				}
				if flags == selTile.adjCount {
					for i := 0; i < 8; i++ {
						if selTile.adj[i] != nil && !selTile.adj[i].flagged {
							flipCount += selTile.adj[i].flip()
							if selTile.adj[i].mine && selTile.adj[i].flipped {
								l.loose = true
								l.duckCharacterUI.state = duckDead
								l.levelTimer.timer.stop()
							}
						}
					}
					if flipCount > 0 {
						l.duckCharacterUI.state = duckSurprised
						l.duckCharacterUI.surprised = 30
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

	speed := 30.0
	if btn(ebiten.KeyW) {
		l.boardXY.y -= speed
		if l.boardXY.y < 1 {
			l.boardXY.y = 1
		}
	}
	if btn(ebiten.KeyA) {
		l.boardXY.x -= speed
		if l.boardXY.x < 1 {
			l.boardXY.x = 1
		}
	}
	if btn(ebiten.KeyS) {
		l.boardXY.y += speed
		if l.boardXY.y+float64(l.boardHeight+1)*11 > 160 {
			l.boardXY.y = 160 - float64(l.boardHeight+1)*11
		}
	}
	if btn(ebiten.KeyD) {
		l.boardXY.x += speed
		if l.boardXY.x+float64(l.boardWidth+1)*17 > 239 {
			l.boardXY.x = 239 - float64(l.boardWidth+1)*17
		}
	}

	// finish checking for power up stuff
	if !l.win && !l.loose &&
		l.usingPowerUp &&
		l.powSelDone &&
		!l.paused {
		switch l.powerUps[l.usingPowerUpID].pType {
		case addMinePow:
			if doAddMinePow(l.board, l.mouseAnchor, l.clickCount) {
				l.mineCount++
				l.miniMap.mineCount = l.mineCount
				l.powerUps[l.usingPowerUpID].activte()
				l.usingPowerUp = false
			}
		case scaredyCatPow:
			if doScaredCat(l.board, l.mouseAnchor, l.clickCount) {
				var mineCount int
				for _, tile := range *l.board {
					if tile.mine {
						mineCount++
					}
				}
				l.mineCount = mineCount
				l.miniMap.mineCount = l.mineCount
				l.powerUps[l.usingPowerUpID].activte()
				l.usingPowerUp = false
			}
		case tidalWavePow:
			l.powerUps[l.usingPowerUpID].activte()
			l.usingPowerUp = false
			doTidalWave(l.board)
		case minusMinePow:
			if doMinusMinePow(l.board, l.mouseAnchor, l.clickCount) {
				l.mineCount--
				l.miniMap.mineCount = l.mineCount
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

	if mbtnr(ebiten.MouseButtonLeft) {
		l.powSelDone = true
	}

	// flag a tile
	if !l.win && !l.loose &&
		mbtnp(ebiten.MouseButtonRight) &&
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

	// update the level timer + some other assets
	if !l.win {
		l.levelTimer.update()
		if l.levelTimer.play.clicked {
			l.paused = false
		}
		if l.levelTimer.pause.clicked {
			l.paused = true
		}
		if btnp(ebiten.KeyEscape) {
			if l.usingPowerUp {
				l.usingPowerUp = false
				for i := 0; i < len(*l.board); i++ {
					(*l.board)[i].bounce = false
				}
			} else {
				if !l.paused {
					l.paused = true
					l.levelTimer.timer.stop()
				} else {
					l.paused = false
					l.levelTimer.timer.start()
				}
			}
		}

	}

	l.miniMap.update()
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
			currentScean = newLevelScean(l.settings, l.powerUpTypes, l.jeepIndexReturn, l.levelIndexReturn)
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
			currentScean = &levelSelect{
				startMenu:   newLevelStartMenu([3]int{l.powerUps[0].pType, l.powerUps[1].pType, l.powerUps[2].pType}),
				jeepIndex:   l.jeepIndexReturn,
				levelNumber: l.levelIndexReturn,
			}
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

	if l.win {
		for i := 0; i < len(*l.board); i++ {
			if !(*l.board)[i].mine {
				(*l.board)[i].flip()
			}
			if (*l.board)[i].mine && !(*l.board)[i].flagged {
				(*l.board)[i].flag()
			}
		}
		l.powerUps[0].available = false
		l.powerUps[1].available = false
		l.powerUps[2].available = false
		l.continueGame.update()

		if btnp(ebiten.KeyEnter) || l.continueGame.clicked {
			l.settings.beaten = true
			allLevels[l.settings.nextLevel].unlocked = true
			// quit to the map
			currentScean = &levelSelect{
				startMenu:   newLevelStartMenu([3]int{l.powerUps[0].pType, l.powerUps[1].pType, l.powerUps[2].pType}),
				jeepIndex:   l.jeepIndexReturn,
				levelNumber: l.levelIndexReturn,
			}
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

	// dog a bone game saver
	if !l.levelTimer.timer.overTime() {
		for i := 0; i < 3; i++ {
			if l.loose && l.powerUps[i].ready && l.powerUps[i].pType == dogABonePow {
				l.loose = false
				l.powerUps[i].activte()
				l.duckCharacterUI.state = duckNormal
				l.levelTimer.timer.start()
				for i := 0; i < len(*l.board); i++ {
					if (*l.board)[i].mine && (*l.board)[i].flipped {
						(*l.board)[i].flagged = true
						(*l.board)[i].flipped = false
						n_mineDog.pause()
						n_mineDog.reset()
						break
					}
				}
			}
		}
	}

	// run the you lost code
	if l.loose {
		for i := 0; i < len(*l.board); i++ {
			if (*l.board)[i].mine && !(*l.board)[i].flagged {
				(*l.board)[i].flag()
			}
			if !(*l.board)[i].mine && (*l.board)[i].flagged {
				(*l.board)[i].flag()
			}
		}
		l.powerUps[0].available = false
		l.powerUps[1].available = false
		l.powerUps[2].available = false
		l.restart.coord = v2f{64, 19}
		l.quit.coord = v2f{122, 19}

		l.restart.update()
		l.quit.update()

		if l.restart.clicked {
			// restart the board
			currentScean = newLevelScean(l.settings, l.powerUpTypes, l.jeepIndexReturn, l.levelIndexReturn)
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
			currentScean = &levelSelect{
				startMenu:   newLevelStartMenu([3]int{l.powerUps[0].pType, l.powerUps[1].pType, l.powerUps[2].pType}),
				jeepIndex:   l.jeepIndexReturn,
				levelNumber: l.levelIndexReturn,
			}
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

	// draw scroll arrows
	if !l.paused {
		offset := math.Abs(math.Sin(float64(tickCounter) / 10))
		x, y := l.boardXY.x+l.boardDXY.x, l.boardXY.y+l.boardDXY.y

		// left scroll arrow
		if x < 0 {
			left := &ebiten.DrawImageOptions{}
			left.GeoM.Translate(3-offset, 80)
			screen.DrawImage(scrollArrowLeft, left)
		}

		// right scroll arrow
		if x+float64(l.boardWidth+1)*17 > 240 {
			right := &ebiten.DrawImageOptions{}
			right.GeoM.Scale(-1, 1)
			right.GeoM.Translate(237+offset, 80)
			screen.DrawImage(scrollArrowLeft, right)
		}

		// up scroll arrow
		if y < 0 {
			up := &ebiten.DrawImageOptions{}
			up.GeoM.Translate(120, 3-offset)
			screen.DrawImage(scrollArrowUp, up)
		}

		// down scroll arrow
		if y+float64(l.boardHeight+1)*11 > 160 {
			down := &ebiten.DrawImageOptions{}
			down.GeoM.Scale(1, -1)
			down.GeoM.Translate(120, 157+offset)
			screen.DrawImage(scrollArrowUp, down)
		}
	}

	// draw the win screen
	if l.win {
		if l.bestTime.timerAccumulator > l.levelTimer.timer.timerAccumulator {
			l.bestTime.timerAccumulator = l.levelTimer.timer.timerAccumulator
			l.settings.bestTime = l.bestTime.timerAccumulator
		}
		winOP := &ebiten.DrawImageOptions{}
		winOP.GeoM.Translate(59, 0)
		screen.DrawImage(winMenu, winOP)

		// steal the best time timer
		old := l.bestTime.coord
		l.bestTime.coord = v2f{130, 57}
		l.bestTime.draw(screen)
		l.bestTime.coord = old

		// steal the level timer
		yourTime := timer{}
		yourTime.coord = v2f{65, 57}
		yourTime.timerAccumulator = l.levelTimer.timer.timerAccumulator
		yourTime.draw(screen)

		// draw the continue button
		l.continueGame.draw(screen)

		// draw your reward stars
		l.settings.stars = l.levelStarCounter.getStarCount()
		starOP := &ebiten.DrawImageOptions{}
		starOP.GeoM.Translate(96, 20)
		if l.settings.stars > 0 {
			screen.DrawImage(star, starOP)
			starOP.GeoM.Translate(16, 0)
		}
		if l.settings.stars > 1 {
			screen.DrawImage(star, starOP)
			starOP.GeoM.Translate(16, 0)
		}
		if l.settings.stars > 2 {
			screen.DrawImage(star, starOP)
		}
	}

	if l.loose {
		looseOP := &ebiten.DrawImageOptions{}
		looseOP.GeoM.Translate(59, 0)
		screen.DrawImage(looseMenu, looseOP)

		l.restart.draw(screen)
		l.quit.draw(screen)
	}
}

func (l *levelScean) fillBoard(safe *n_tile) {
	mines := l.mineCount
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
	for i := 0; i < len(*board); i++ {
		tile := &(*board)[i]
		if tile.mine && !tile.flagged && rand.Float64() > 0.5 {
			var done bool
			var safe int
			for !done {
				safe++
				if safe > 100 {
					break
				}
				ii := rand.Intn(len(*board))
				if ii == i {
					continue
				}
				target := &(*board)[ii]
				if !target.mine && !target.flipped && !target.flagged {
					tile.mine = false
					target.mine = true
					done = true
				}
			}
		}

		// recalculate the adj counts
		for i := 0; i < len(*board); i++ {
			tile := &(*board)[i]
			tile.adjCount = 0
			for _, adj := range tile.adj {
				if adj != nil && adj.mine {
					tile.adjCount++
				}
			}
		}

		// filp any new open spaces
		for _, tile := range *board {
			if tile.adjCount == 0 && tile.flipped {
				tile.flipped = false
				tile.flip()
			}
		}
	}
}

func doScaredCat(board *[]n_tile, mouseAnchor v2f, clickCount int) bool {
	for i := 0; i < len(*board); i++ {
		if !(*board)[i].flipped {
			(*board)[i].bounce = true
		}
	}
	if mouseAnchor.dist(mCoordsF()) < 5 &&
		clickCount < 30 &&
		mbtnr(ebiten.MouseButtonLeft) {
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

func doMinusMinePow(board *[]n_tile, mouseAnchor v2f, clickCount int) bool {
	for i := 0; i < len(*board); i++ {
		if (*board)[i].flipped && (*board)[i].adjCount > 0 {
			(*board)[i].bounce = true
		}
	}
	if mouseAnchor.dist(mCoordsF()) < 5 &&
		clickCount < 30 &&
		mbtnr(ebiten.MouseButtonLeft) {
		minX := 999999
		var selTile *n_tile
		for i, tile := range *board {
			if tile.hovered() && tile.index.x < minX {
				minX = tile.index.x
				selTile = &(*board)[i]
			}
		}

		var flagCount int
		if selTile != nil {
			for _, adj := range selTile.adj {
				if adj != nil && adj.flagged {
					flagCount++
				}
			}
		}
		if selTile != nil && selTile.adjCount > flagCount {
			var done bool
			for !done {
				tile := selTile.adj[rand.Intn(8)]
				if tile != nil && tile.mine && !tile.flagged {
					tile.mine = false
					for _, adj := range tile.adj {
						if adj != nil {
							adj.adjCount--
						}
					}
					done = true
					break
				}
			}
			for i := 0; i < len(*board); i++ {
				(*board)[i].bounce = false
			}

			if selTile.adjCount == 0 {
				selTile.flipped = false
				selTile.flip()
			}

			for _, adj := range selTile.adj {
				if adj != nil && adj.adjCount == 0 && adj.flipped {
					adj.flipped = false
					adj.flip()
				}
			}
			return true
		}
		if selTile != nil && selTile.adjCount <= flagCount {
			selTile.shake()
			for _, adj := range selTile.adj {
				if adj != nil && adj.flagged {
					adj.shake()
				}
			}
		}
	}
	return false
}

func doAddMinePow(board *[]n_tile, mouseAnchor v2f, clickCount int) bool {
	for i := 0; i < len(*board); i++ {
		if (*board)[i].flipped && (*board)[i].adjCount > 0 {
			(*board)[i].bounce = true
		}
	}
	if mouseAnchor.dist(mCoordsF()) < 5 &&
		clickCount < 30 &&
		mbtnr(ebiten.MouseButtonLeft) {
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
				for _, adj := range selTile.adj {
					if adj != nil && !adj.flipped {
						adj.shake()
					}
				}
			}
		}
	}
	return false
}

func lockTiles(board *[]n_tile, tileCount int) {
	count := tileCount
	for count > 0 {
		tile := &(*board)[rand.Intn(len(*board))]
		if tile.lockedCount == 0 && tile.adjCount > 0 && !tile.mine {
			tile.lockedCount = 10
			count--
		}
	}
}

func freezeTiles(board *[]n_tile, tileCount int) {
	count := tileCount
	for count > 0 {
		tile := &(*board)[rand.Intn(len(*board))]
		if !tile.iced {
			tile.iced = true
			tile.gfx = n_newAniSprite(
				iceImg[:],
				[]uint{6, 6, 6},
				false,
			)
			count--
		}
	}
}

func addTimeTiles(board *[]n_tile, tileCount int) {
	count := tileCount
	for count > 0 {
		tile := &(*board)[rand.Intn(len(*board))]
		if !tile.timeTile {
			tile.timeTile = true
			count--
		}
	}
}
