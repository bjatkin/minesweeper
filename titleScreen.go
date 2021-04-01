package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type titleScreanScean struct {
	title       *ebiten.Image
	newGame     bool
	loadGame    bool
	loadedGames [3]*saveGame
}

var eggCursor *n_aniSprite
var gameSlots *ebiten.Image
var newGame *ebiten.Image
var loadGame *ebiten.Image
var emptySlot *ebiten.Image
var lvlSlot *ebiten.Image
var powSlot *ebiten.Image
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

	// foreign assets
	redStar = subImage(ss, 208, 16, 8, 8)

	return nil
}

func (t *titleScreanScean) unload() error {
	t.title = nil
	unloadAsset("assets/title_Screen.png")
	return nil
}

func (t *titleScreanScean) update() error {
	if btnp(ebiten.KeyEnter) || mbtnp(ebiten.MouseButtonLeft) {
		if t.loadGame {
			lvlMap := &levelSelect{
				jeepIndex:   t.loadedGames[0].jeepIndex,
				levelNumber: t.loadedGames[0].levelNumber,
			}

			// unlock levels
			for i, lvl := range t.loadedGames[0].allLevels {
				allLevels[i].beaten = lvl.beaten
				allLevels[i].stars = lvl.stars
				allLevels[i].unlocked = lvl.unlocked
				allLevels[i].bestTime = lvl.bestTime
			}

			currentScean = lvlMap
			err := currentScean.load()
			if err != nil {
				return err
			}

			// unlock powerups
			for i, pow := range t.loadedGames[0].unlockedPowers {
				unlockedPowers[i] = newPowIcon(pow.powType, unlockedPowers[i].coord)
			}

			// set powerups
			lvlMap.startMenu = newLevelStartMenu(t.loadedGames[0].currentPows)

			err = t.unload()
			if err != nil {
				return err
			}
		}

		t.loadGame = true
		s := saveGame{}
		s.loadData("test.save")
		t.loadedGames[0] = &s
		t.loadedGames[2] = &s

		// THIS IS THE CODE THAT WILL LOAD US INTO A NEW SCEAN
		// currentScean = newLevelScean(
		// 	allLevels[0],
		// 	[3]int{lockedPow, lockedPow, lockedPow},
		// 	0,
		// 	1,
		// )

		// var err error
		// err = currentScean.load()
		// if err != nil {
		// 	return err
		// }

		// err = t.unload()
		// if err != nil {
		// 	return err
		// }
	}

	eggCursor.update()
	return nil
}

func (t *titleScreanScean) draw(screen *ebiten.Image) {
	screen.DrawImage(t.title, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(33, 121)
	eggCursor.draw(screen, op)

	if t.newGame || t.loadGame {
		// draw New Game
		newOp := &ebiten.DrawImageOptions{}
		newOp.GeoM.Translate(13, 3)
		if t.newGame {
			screen.DrawImage(newGame, newOp)
		} else {
			screen.DrawImage(loadGame, newOp)
		}
		newOp.GeoM.Translate(0, 10)
		screen.DrawImage(gameSlots, newOp)

		for i := 0; i < 3; i++ {
			newOp.GeoM.Reset()
			newOp.GeoM.Translate(120, 29+float64(i*43))
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
