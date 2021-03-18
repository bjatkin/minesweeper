package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type levelScean struct {
	boardXY          v2f
	mouseXY          v2i
	board            *[]n_tile
	filled           bool
	win, loose       bool
	paused           bool
	timerAccumulator int64
	start            int64
	flagCount        int
	settings         n_levelData
	usingPowerUp     bool
	powerUp          powerUp
}

// Level Assets, all of these should be loaded from the main sprite sheet
var (
	levelAssetsLoaded bool
	grassImg          [7]*ebiten.Image
	yellowGrassImg    [7]*ebiten.Image
	pinkGrassImg      [7]*ebiten.Image
	blueGrassImg      [7]*ebiten.Image
	waterImg          [7]*ebiten.Image
	iceImg            [3]*ebiten.Image
	lock              *ebiten.Image
	n_mineDog         *n_aniSprite
	n_dogBark         *ebiten.Image
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
	// we only need to run this function once and then these assest are left in memory going forward
	if levelAssetsLoaded {
		return nil
	}

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
		[]uint{5, 5, 5, 5},
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

	levelAssetsLoaded = true
	return nil
}

func (l *levelScean) unload() error {
	// we want to leave this stuff loaded since we'll be jumping in and
	// out of levels a bunch
	return nil
}

func (l *levelScean) update() error {
	return nil
}

func (l *levelScean) draw(screen *ebiten.Image) {
}
