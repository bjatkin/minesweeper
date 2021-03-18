package main

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 240
	screenHeight = 160
)

const (
	gameplay = iota
	pause
)

// Game is an ebiten game
type Game struct {
	frames uint64
}

// Update is the game update function
func (g *Game) Update() error {
	g.frames++
	pauseCoolDown--

	if ebiten.IsKeyPressed(ebiten.KeyEscape) && pauseCoolDown < 0 {
		pauseCoolDown = 15
		if gameState != pause {
			gameState = pause
		} else {
			gameState = gameplay
		}
	}

	test.hover(ebiten.CursorPosition())

	if gameState == pause {
		cursorAni.update()
		return nil
	}

	leftBtn := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	rightBtn := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	if leftBtn && rightBtn {
		row, col := currentBoard.click(ebiten.CursorPosition())
		if row >= 0 && col >= 0 {
			currentBoard.clearAround(row, col)
		}
	}

	if leftBtn && !rightBtn {
		row, col := currentBoard.click(ebiten.CursorPosition())
		if row >= 0 && col >= 0 {
			currentBoard.flipTiles(row, col)
		}
	}

	if rightBtn && !leftBtn {
		row, col := currentBoard.click(ebiten.CursorPosition())
		if row >= 0 && col >= 0 {
			tile := currentBoard.field[row][col]
			if tile != nil && !tile.flipped && tile.flaggedCoolDown <= 0 {
				tile.flaggedCoolDown = 15
				tile.flagged = !tile.flagged
				if tile.flagged {
					currentBoneYard.filled++
				} else {
					currentBoneYard.filled--

				}
			}
		}
	}

	var err error
	for _, up := range allUpdatables {
		err = up.update()
		if err != nil {
			return err
		}
	}

	for _, layer := range allEntities {
		for _, ent := range layer {
			err = ent.update()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Draw is the draw funciton for the game
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 241, 232, 255})

	// draw background
	for x := 0.0; x < screenWidth; x += 16 {
		for y := 0.0; y < screenHeight; y += 16 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y)
			screen.DrawImage(grassBG, op)
		}
	}

	screen.DrawImage(gameArea, nil)
	diff := time.Now().Unix() - startTime
	sec := strconv.Itoa(int(diff % 60))
	if len(sec) == 1 {
		sec = "0" + sec
	}
	min := strconv.Itoa(int(diff / 60))
	if len(min) == 1 {
		min = "0" + min
	}
	pxlPrint(screen, mainFont, 31, 8, min+":"+sec)

	for _, layer := range allEntities {
		for _, ent := range layer {
			ent.draw(screen)
		}
	}

	if gameState == pause {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(pauseMenu, op)
	}

	// draw the new cursor
	x, y := ebiten.CursorPosition()
	cursorAni.draw(screen, float64(x)-2, float64(y)-2)
}

// Layout determins the layout of the game
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	mainFont        *pxlFont
	startTime       = time.Now().Unix()
	currentBoard    *gameBoard
	currentBoneYard *boneyard
	gameState       = gameplay

	// TODO: remove this
	pauseCoolDown = 0
	test          *button
)

func main() {
	test = &button{}
	registerEnt(test, 10)

	rand.Seed(time.Now().Unix())

	loadSprites()
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	// load the game assets here
	mainFont = load6x8Font()

	currentBoard = newGameBoard(8, 8, 10)
	currentBoard.translate(76, 50)

	currentBoneYard = &boneyard{total: 10}
	registerEnt(currentBoneYard, 10)

	// start the game
	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Mine Sweeper")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
