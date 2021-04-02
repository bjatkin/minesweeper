package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type titleScreanScean struct {
	title       *ebiten.Image
	newGame     bool
	loadGame    bool
	loadedGames [3]*saveGame
	selSlot     int
	cursorIndex int
}

var gameSlot int
var eggCursor *n_aniSprite
var gameSlots *ebiten.Image
var newGame *ebiten.Image
var loadGame *ebiten.Image
var emptySlot *ebiten.Image
var lvlSlot *ebiten.Image
var powSlot *ebiten.Image
var hilightSlot *ebiten.Image
var slot1, slot2, slot3 *ebiten.Image
var starsSlot *ebiten.Image
var greenTick *ebiten.Image
var pinkTick *ebiten.Image
var blueTick *ebiten.Image

func (t *titleScreanScean) load() error {
	var err error
	t.title, err = getAsset("assets/title_screen.png")
	if err != nil {
		return err
	}

	ss, err := getAsset("assets/sprite_sheet.png")
	if err != nil {
		return err
	}

	eggCursor = n_newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 64, 328, 16, 16), // normal
			subImage(ss, 80, 328, 16, 16), // left
			subImage(ss, 64, 328, 16, 16), // normal
			subImage(ss, 96, 328, 16, 16), // right
			subImage(ss, 64, 328, 16, 16), // normal
			subImage(ss, 80, 328, 16, 16), // left
			subImage(ss, 64, 328, 16, 16), // normal
			subImage(ss, 96, 328, 16, 16), // right
		},
		[]uint{80, 8, 8, 8, 8, 8, 8, 8},
		true,
	)
	eggCursor.play()

	gameSlots = subImage(ss, 680, 0, 218, 132)
	newGame = subImage(ss, 680, 136, 53, 11)
	loadGame = subImage(ss, 672, 152, 58, 11)
	emptySlot = subImage(ss, 632, 144, 38, 13)
	lvlSlot = subImage(ss, 736, 133, 38, 35)
	powSlot = subImage(ss, 776, 133, 38, 35)
	starsSlot = subImage(ss, 816, 133, 88, 35)
	greenTick = subImage(ss, 608, 144, 3, 3)
	pinkTick = subImage(ss, 616, 144, 2, 2)
	blueTick = subImage(ss, 624, 144, 2, 2)
	hilightSlot = subImage(ss, 198, 408, 174, 40)
	slot1 = subImage(ss, 160, 408, 36, 15)
	slot2 = subImage(ss, 160, 424, 36, 15)
	slot3 = subImage(ss, 160, 440, 36, 15)

	// foreign assets
	redStar = subImage(ss, 208, 16, 8, 8)
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
	// end foreign assets

	s := saveGame{}
	err = s.loadData("test.save")
	if err != nil {
		return err
	}
	t.loadedGames[0] = &s

	return nil
}

func (t *titleScreanScean) unload() error {
	t.title = nil
	unloadAsset("assets/title_Screen.png")
	return nil
}

func (t *titleScreanScean) update() error {
	var hovered bool
	if !t.loadGame && !t.newGame {
		if btnp(ebiten.KeyDown) || btnp(ebiten.KeyS) {
			t.cursorIndex = 1
		}
		if btnp(ebiten.KeyUp) || btnp(ebiten.KeyW) {
			t.cursorIndex = 0
		}

		m := mCoords()
		if m.x > 34 && m.y > 124 &&
			m.x < 99 && m.y < 135 {
			t.cursorIndex = 0
			hovered = true
		}
		if m.x > 37 && m.y > 139 &&
			m.x < 107 && m.y < 149 {
			t.cursorIndex = 1
			hovered = true
		}
	}
	if t.loadGame || t.newGame {
		t.selSlot = 0
		m := mCoords()
		if m.x > 16 && m.y > 16 &&
			m.x < 227 && m.y < 55 {
			t.selSlot = 1
		}
		if m.x > 16 && m.y > 59 &&
			m.x < 227 && m.y < 98 {
			t.selSlot = 2
		}
		if m.x > 16 && m.y > 102 &&
			m.x < 227 && m.y < 141 {
			t.selSlot = 3
		}

		if mbtnp(ebiten.MouseButtonLeft) {
			if t.selSlot == 0 {
				t.loadGame = false
				t.newGame = false
			}
		}
	}
	if btnp(ebiten.KeyEnter) || mbtnp(ebiten.MouseButtonLeft) {
		// Load into a new game
		if t.loadGame {
			gameSlot = t.selSlot

			lvlMap := &levelSelect{
				jeepIndex:   t.loadedGames[0].jeepIndex,
				levelNumber: t.loadedGames[0].levelNumber,
				startMenu:   newLevelStartMenu(t.loadedGames[0].currentPows),
			}

			// unlock levels
			for i, lvl := range t.loadedGames[0].allLevels {
				allLevels[i].beaten = lvl.beaten
				allLevels[i].stars = lvl.stars
				allLevels[i].unlocked = lvl.unlocked
				allLevels[i].bestTime = lvl.bestTime
			}

			// unlock powerups
			for i, pow := range t.loadedGames[0].unlockedPowers {
				unlockedPowers[i] = newPowIcon(pow.powType, unlockedPowers[i].coord)
			}

			currentScean = lvlMap
			err := currentScean.load()
			if err != nil {
				return err
			}

			err = t.unload()
			if err != nil {
				return err
			}
		}
		if t.newGame {
			gameSlot = t.selSlot

			currentScean = newLevelScean(
				allLevels[0],
				[3]int{lockedPow, lockedPow, lockedPow},
				0,
				1,
			)

			var err error
			err = currentScean.load()
			if err != nil {
				return err
			}

			err = t.unload()
			if err != nil {
				return err
			}
		}

		if mbtnp(ebiten.MouseButtonLeft) && hovered {
			if t.cursorIndex == 0 {
				t.newGame = true
				t.loadGame = false
			}
			if t.cursorIndex == 1 {
				t.loadGame = true
				t.newGame = false
			}
		}
	}

	eggCursor.update()
	return nil
}

func (t *titleScreanScean) draw(screen *ebiten.Image) {
	screen.DrawImage(t.title, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	if t.cursorIndex == 0 {
		op.GeoM.Translate(33, 121)
	} else {
		op.GeoM.Translate(36, 135)
	}
	eggCursor.draw(screen, op)

	// draw the new/ load game ui
	if t.newGame || t.loadGame {
		if btnp(ebiten.KeyEscape) {
			t.newGame = false
			t.loadGame = false
		}
		newOp := &ebiten.DrawImageOptions{}
		newOp.GeoM.Translate(13, 3)
		if t.newGame {
			screen.DrawImage(newGame, newOp)
		} else {
			screen.DrawImage(loadGame, newOp)
		}
		newOp.GeoM.Translate(0, 10)
		screen.DrawImage(gameSlots, newOp)

		selOp := &ebiten.DrawImageOptions{}
		if t.selSlot == 1 {
			selOp.GeoM.Translate(16, 16)
			screen.DrawImage(slot1, selOp)
			selOp.GeoM.Translate(38, 0)
			screen.DrawImage(hilightSlot, selOp)
		}
		if t.selSlot == 2 {
			selOp.GeoM.Translate(16, 59)
			screen.DrawImage(slot2, selOp)
			selOp.GeoM.Translate(38, 0)
			screen.DrawImage(hilightSlot, selOp)
		}
		if t.selSlot == 3 {
			selOp.GeoM.Translate(16, 102)
			screen.DrawImage(slot3, selOp)
			selOp.GeoM.Translate(38, 0)
			screen.DrawImage(hilightSlot, selOp)
		}

		for i := 0; i < 3; i++ {
			newOp.GeoM.Reset()
			newOp.GeoM.Translate(118, 29+float64(i*43))
			if t.loadedGames[i] == nil {
				screen.DrawImage(emptySlot, newOp)
			} else {
				// draw slot data
				dataOp := &ebiten.DrawImageOptions{}
				dataOp.GeoM.Translate(57, 18+float64(i*43))
				screen.DrawImage(lvlSlot, dataOp)
				dataOp.GeoM.Translate(40, 0)
				screen.DrawImage(powSlot, dataOp)
				dataOp.GeoM.Translate(40, 0)
				screen.DrawImage(starsSlot, dataOp)

				// draw all the icons here
				// lvl icons
				iconOp := &ebiten.DrawImageOptions{}
				iconOp.GeoM.Translate(60, 40+float64(i*43))
				for i, lvl := range t.loadedGames[i].allLevels {
					if lvl.beaten {
						screen.DrawImage(greenTick, iconOp)
					}
					iconOp.GeoM.Translate(5, 0)
					if i == 6 {
						// start row 2
						iconOp.GeoM.Translate(-35, 7)
					}
				}

				// pow icons
				iconOp.GeoM.Reset()
				iconOp.GeoM.Translate(100, 40+float64(i*43))
				for _, pow := range t.loadedGames[i].unlockedPowers {
					if pow.powType != lockedPow {
						screen.DrawImage(pinkTick, iconOp)
						iconOp.GeoM.Translate(5, 0)
					}
				}

				// slot icons
				iconOp.GeoM.Reset()
				iconOp.GeoM.Translate(100, 47+float64(i*43))
				for _, slot := range t.loadedGames[i].currentPows {
					if slot != lockedPow {
						screen.DrawImage(blueTick, iconOp)
						iconOp.GeoM.Translate(5, 0)
					}
				}

				// stars
				var totalStars int
				for _, lvl := range t.loadedGames[i].allLevels {
					totalStars += int(lvl.stars)
				}

				starOp := &ebiten.DrawImageOptions{}
				starOp.GeoM.Translate(139, 38+float64(i*43))
				var counter int
				for i := 0; i < 2; i++ {
					for ii := 0; ii < 14; ii++ {
						if counter < totalStars {
							screen.DrawImage(redStar, starOp)
							starOp.GeoM.Translate(6, 0)
						}
						counter++
					}
					starOp.GeoM.Translate(-84, 5)
				}
			}
		}
	}
}
