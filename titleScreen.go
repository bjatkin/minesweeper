package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type titleScreanScean struct {
	title *ebiten.Image
}

func (t *titleScreanScean) load() error {
	var err error
	t.title, err = getAsset("assets/title_screen.png")
	if err != nil {
		return err
	}

	return nil
}

func (t *titleScreanScean) unload() error {
	t.title = nil
	unloadAsset("assets/title_Screen.png")
	return nil
}

func (t *titleScreanScean) update() error {
	if btnp(ebiten.KeyEnter) {
		currentScean = &levelSelectScean{}

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
	return nil
}

func (t *titleScreanScean) draw(screen *ebiten.Image) {
	screen.DrawImage(t.title, &ebiten.DrawImageOptions{})
}
