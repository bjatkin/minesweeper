package main

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type levelData struct {
	rows, cols, mineCount uint
	maxTime               int
	shuffel, shuffelHard  bool
	shuffelTime           int
	flipLockCount         int
	adjLockCount          int
}

var levelList = []levelData{
	{8, 8, 10, 0, false, false, 0, 0, 0},     // normal 8 x 8
	{16, 16, 40, 0, false, false, 0, 0, 0},   // normal 16 x 16
	{16, 16, 40, 300, false, false, 0, 0, 0}, // timed 16 x 16
	{8, 8, 10, 0, true, false, 10, 0, 0},     // basic suffel 10s 8 x 8
	{16, 16, 40, 0, true, false, 30, 0, 0},   // basic shuffel 30s 16 x 16
	{8, 8, 10, 0, true, true, 60, 0, 0},      // had shuffel 16 x 16 60s
	{16, 16, 40, 0, false, false, 0, 60, 0},  // 20 flipLock tiles 16 x 16
	{16, 16, 40, 0, false, false, 0, 0, 20},  // 20 adjLock tiles 16 x 16
}

var levelCleared [100]bool

type levelSelectScean struct {
	sel int
}

func (l *levelSelectScean) load() error {
	return nil
}

func (l *levelSelectScean) unload() error {
	return nil
}

func (l *levelSelectScean) update() error {
	if btnp(ebiten.KeyEscape) {
		currentScean = &titleScreanScean{}

		var err error
		err = currentScean.load()
		if err != nil {
			return err
		}
		err = l.unload()
		if err != nil {
			return err
		}
	}
	if btnp(ebiten.KeyUp) || btnp(ebiten.KeyW) {
		l.sel--
	}
	if btnp(ebiten.KeyDown) || btnp(ebiten.KeyS) {
		l.sel++
	}
	if l.sel < 0 {
		l.sel = len(levelList) - 1
	}
	if l.sel >= len(levelList) {
		l.sel = 0
	}

	if btnp(ebiten.KeyEnter) {
		currentScean = &gameBoardScean{
			rows:          levelList[l.sel].rows,
			cols:          levelList[l.sel].cols,
			mineCount:     levelList[l.sel].mineCount,
			shuffel:       levelList[l.sel].shuffel,
			shuffelHard:   levelList[l.sel].shuffelHard,
			shuffelTime:   levelList[l.sel].shuffelTime,
			maxTime:       levelList[l.sel].maxTime,
			flipLockCount: levelList[l.sel].flipLockCount,
			adjLockCount:  levelList[l.sel].adjLockCount,
			levelID:       l.sel,
		}
		err := currentScean.load()
		if err != nil {
			return err
		}
		l.unload()
	}
	return nil
}

func (l *levelSelectScean) draw(screen *ebiten.Image) {
	for i, level := range levelList {
		if i == l.sel {
			pxlPrint(screen, mainFont, 5, float64(i*12), ">")
		}
		cleared := ""
		if levelCleared[i] {
			cleared = ", [cleared]"
		}

		maxTime := ""
		if level.maxTime > 0 {
			min := strconv.Itoa(level.maxTime / 60)
			if len(min) == 1 {
				min = "0" + min
			}
			sec := strconv.Itoa(level.maxTime % 60)
			if len(sec) == 1 {
				sec = "0" + sec
			}
			maxTime = ", Time Limit" + min + ":" + sec
		}

		shuffel := ""
		if level.shuffel {
			shuffel = ", Shuffel " + strconv.Itoa(level.shuffelTime) + " sec."
		}
		if level.shuffelHard {
			shuffel = ", Hard Shuffel " + strconv.Itoa(level.shuffelTime) + " sec."
		}

		extra := ""
		if level.flipLockCount > 0 {
			extra = strconv.Itoa(level.flipLockCount) + " tiles are locked"
		}
		if level.adjLockCount > 0 {
			extra = strconv.Itoa(level.adjLockCount) + " tiles are mysterious"
		}
		pxlPrint(screen, mainFont, 15, float64(i*12),
			fmt.Sprintf("%d x %d [%d]%s%s%s%s", level.rows, level.cols, level.mineCount, maxTime, shuffel, extra, cleared))
	}
}
