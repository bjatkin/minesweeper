package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type titleScreanScean struct {
	title *ebiten.Image
}

var eggCursor *n_aniSprite

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

	return nil
}

func (t *titleScreanScean) unload() error {
	t.title = nil
	unloadAsset("assets/title_Screen.png")
	return nil
}

func (t *titleScreanScean) update() error {
	if btnp(ebiten.KeyEnter) || mbtnp(ebiten.MouseButtonLeft) {
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

	eggCursor.update()
	return nil
}

func (t *titleScreanScean) draw(screen *ebiten.Image) {
	screen.DrawImage(t.title, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(33, 121)
	eggCursor.draw(screen, op)
}
