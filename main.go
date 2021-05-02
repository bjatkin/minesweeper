package main

import (
	"errors"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 240
	screenHeight = 160
)

var quitGame = errors.New("GameOver")
var tickCounter uint64

// Game the main game object
type Game struct{}

// Update the update func that get's called every tick
func (g *Game) Update() error {
	tickCounter++
	updateKeys()
	updateMouse()

	err := currentScean.update()
	if err != nil {
		return err
	}

	gameCursor.update()
	gameCursor.op.GeoM.Reset()
	x, y := ebiten.CursorPosition()
	gameCursor.op.GeoM.Translate(float64(x)-2, float64(y)-2)
	gameCursorHold.op.GeoM.Reset()
	gameCursorHold.op.GeoM.Translate(float64(x), float64(y))

	return nil
}

// Draw the code to draw the screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 241, 232, 255}) // pico-8 white

	currentScean.draw(screen)
	if cursorHold {
		gameCursorHold.draw(screen)
	} else {
		gameCursor.draw(screen)
	}
}

// Layout determines the layout of the game window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	mainFont       *pxlFont
	currentScean   scean
	gameCursor     aniSprite
	gameCursorHold sprite
	cursorHold     bool
)

func main() {
	rand.Seed(time.Now().Unix())

	mainFont = load6x8Font()

	var err error
	ss, err := getAsset("assets/sprite_sheet.png")
	if err != nil {
		log.Fatal(err)
	}
	gameCursor = newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 120, 32, 16, 16),
			subImage(ss, 136, 32, 16, 16),
		},
		[]uint16{28, 28},
		true,
	)
	gameCursor.play()

	gameCursorHold = newSprite(subImage(ss, 120, 48, 16, 16))

	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	currentScean = &titleScreanScean{}
	err = currentScean.load()
	if err != nil {
		log.Fatal(err)
	}

	// TEST
	// currentScean = &testScean{}
	// err = currentScean.load()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("MineSweeper")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(&Game{}); err != nil && err != quitGame {
		log.Fatal(err)
	}
}
