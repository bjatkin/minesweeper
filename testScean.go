package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type testScean struct {
	baseBoard    *levelScean
	tile         *n_tile
	zeroTile     *n_tile
	flaggedTile  *n_tile
	blownTile    *n_tile
	waterTile    *n_tile
	icedTile     *n_tile
	iceFlipIndex int
	lockedTile   *n_tile
	powerUps     [8]*powerUp

	btnRestart *uiButton
	btnQuit    *uiButton

	timer   *timer
	uiTimer *boardTimer

	duckUI *duckFeedBack

	miniMap *miniMap

	stars *starCounter
}

func (t *testScean) load() error {

	// a quick hack to get the level assets loaded
	t.baseBoard = newLevelScean(allLevels[0], [3]int{}, 0, 1)
	t.baseBoard.boardXY.x = 5
	err := t.baseBoard.load()
	if err != nil {
		return err
	}

	// load in some assets to test
	t.tile = n_newTile(t.baseBoard, v2i{0, 5}, false, false, false)
	t.tile.adjCount = 4

	t.zeroTile = n_newTile(t.baseBoard, v2i{0, 7}, false, false, false)

	t.flaggedTile = n_newTile(t.baseBoard, v2i{0, 9}, false, false, false)
	t.flaggedTile.flag()

	t.blownTile = n_newTile(t.baseBoard, v2i{1, 5}, false, false, false)
	t.blownTile.mine = true

	t.icedTile = n_newTile(t.baseBoard, v2i{1, 7}, true, false, false)
	t.icedTile.adjCount = 1
	for i := 0; i < 8; i++ {
		t.icedTile.adj[i] = n_newTile(t.baseBoard, v2i{}, false, false, false)
	}

	t.lockedTile = n_newTile(t.baseBoard, v2i{1, 9}, false, true, false)
	t.lockedTile.adjCount = 8

	t.waterTile = n_newTile(t.baseBoard, v2i{2, 5}, false, false, true)
	t.tile.adjCount = 5

	t.btnRestart = newUIButton(v2f{5, 10}, restartBtn)

	t.btnQuit = newUIButton(v2f{5, 30}, quitBtn)

	t.timer = &timer{coord: v2f{110, 30}}
	t.timer.start()

	t.uiTimer = newBoardTimer(v2f{110, 50})

	t.powerUps[0] = newPowerUp(addMinePow, ebiten.Key1, t.uiTimer.timer)
	t.powerUps[0].coord = v2f{90, 5}
	t.powerUps[0].available = true

	t.powerUps[1] = newPowerUp(minusMinePow, ebiten.Key2, t.uiTimer.timer)
	t.powerUps[1].coord = v2f{90, 25}
	t.powerUps[1].available = true

	t.powerUps[2] = newPowerUp(tidalWavePow, ebiten.Key3, t.uiTimer.timer)
	t.powerUps[2].coord = v2f{90, 45}
	t.powerUps[2].available = true

	t.powerUps[3] = newPowerUp(scaredyCatPow, ebiten.Key0, t.uiTimer.timer)
	t.powerUps[3].coord = v2f{90, 65}
	t.powerUps[3].available = true

	t.powerUps[4] = newPowerUp(dogWistlePow, ebiten.Key0, t.uiTimer.timer)
	t.powerUps[4].coord = v2f{90, 85}
	t.powerUps[4].available = true

	t.powerUps[5] = newPowerUp(shuffelPow, ebiten.Key0, t.uiTimer.timer)
	t.powerUps[5].coord = v2f{90, 105}
	t.powerUps[5].available = true

	t.powerUps[6] = newPowerUp(dogABonePow, ebiten.Key0, t.uiTimer.timer)
	t.powerUps[6].coord = v2f{90, 125}
	t.powerUps[6].available = true

	t.powerUps[7] = newPowerUp(tidalWavePow, ebiten.Key0, t.uiTimer.timer)
	t.powerUps[7].coord = v2f{110, 5}

	t.stars = &starCounter{
		coord:         v2f{0, 1},
		timer:         t.uiTimer.timer,
		oneStarTime:   15 * 1000000000,
		twoStarTime:   10 * 1000000000,
		threeStarTime: 5 * 1000000000,
	}

	t.duckUI = newDuckFeedBack(
		v2f{138, 121},
		newPowerUp(addMinePow, ebiten.Key1, t.uiTimer.timer),
		newPowerUp(minusMinePow, ebiten.Key2, t.uiTimer.timer),
		newPowerUp(dogABonePow, ebiten.Key3, t.uiTimer.timer),
	)
	t.duckUI.powOne.available = true
	t.duckUI.powTwo.available = true
	t.duckUI.powThree.available = true

	tiles := newSquareLayout(8, 8)
	var miniTiles []n_tile
	for i, tile := range tiles {
		miniTiles = append(miniTiles, *n_newTile(t.baseBoard, tile.index, false, false, false))
		if tile.index.x%2 == 0 {
			miniTiles[i].flipped = true
		}
	}
	for i, tile := range miniTiles {
		for ii, adj := range tiles[i].adj {
			if adj > -1 {
				tile.adj[ii] = &miniTiles[adj]
			}
		}
	}

	t.miniMap = newMiniMap(t.baseBoard, v2f{0, 124}, &miniTiles, 10)

	return nil
}

func (t *testScean) unload() error {
	return nil
}

func (t *testScean) update() error {
	flipCount := 0
	if btnp(ebiten.KeyA) {
		flipCount = 5
		if t.iceFlipIndex < 8 {
			t.icedTile.adj[t.iceFlipIndex].flip()
			t.iceFlipIndex++
		}

		if t.timer.running {
			t.timer.stop()
		} else {
			t.timer.start()
		}

		if t.stars.starCount > 0 {
			t.stars.starCount--
		}
	}

	if btnp(ebiten.KeyD) {
		for _, adj := range t.icedTile.adj {
			adj.flipped = true
		}
		t.iceFlipIndex = 10
		if t.stars.starCount < 3 {
			t.stars.starCount++
		}
	}

	if mbtn(ebiten.MouseButtonLeft) {
		if t.duckUI.state == duckNormal {
			t.duckUI.state = duckSurprised
		}
		if t.btnQuit.clicked {
			t.tile.shake()
		}
	}

	if mbtn(ebiten.MouseButtonLeft) && t.btnRestart.clicked {
		t.tile.gfx.reset()
		t.tile.flip()

		t.zeroTile.gfx.reset()
		t.zeroTile.flip()

		t.blownTile.flip()
		t.lockedTile.flip()

		t.waterTile.gfx.reset()
		t.waterTile.flip()

		t.flaggedTile.flip()
	}

	t.btnQuit.update()
	t.btnRestart.update()

	t.tile.update(flipCount)
	t.zeroTile.update(flipCount)
	t.flaggedTile.update(flipCount)
	t.blownTile.update(flipCount)
	t.icedTile.update(flipCount)
	t.lockedTile.update(flipCount)
	t.waterTile.update(flipCount)

	for _, pow := range t.powerUps {
		if pow.wasSelected() {
			pow.activte()
		}
	}

	if t.duckUI.powOne.wasSelected() {
		t.duckUI.powOne.activte()
		t.duckUI.state = duckCool
	}
	if t.duckUI.powTwo.wasSelected() {
		t.duckUI.powTwo.activte()
		t.duckUI.state = duckDead
	}
	if t.duckUI.powThree.wasSelected() {
		t.duckUI.powThree.activte()
		t.duckUI.state = duckNormal
	}

	n_mineDog.update()
	n_markerFlag.update()

	t.uiTimer.update()

	t.miniMap.update()

	if btnp(ebiten.KeyW) {
		t.tile.bounce = !t.tile.bounce
	}

	return nil
}

func (t *testScean) draw(screen *ebiten.Image) {
	t.btnQuit.draw(screen)
	t.btnRestart.draw(screen)

	t.tile.draw(screen)
	t.zeroTile.draw(screen)
	t.flaggedTile.draw(screen)
	t.blownTile.draw(screen)
	t.icedTile.draw(screen)
	t.lockedTile.draw(screen)
	t.waterTile.draw(screen)

	for _, pow := range t.powerUps {
		pow.draw(screen)
	}

	t.timer.draw(screen)
	t.uiTimer.draw(screen)

	t.stars.draw(screen)

	t.duckUI.draw(screen)

	t.miniMap.draw(screen)
}
