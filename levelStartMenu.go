package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type levelStartMenu struct {
	powOne      *uiIcon
	powOneSel   bool
	powTwo      *uiIcon
	powTwoSel   bool
	powThree    *uiIcon
	powThreeSel bool
	startBtn    *uiButton
	powIcons    [7]*uiIcon
	levelData   *n_levelData
	loadTime    int
}

const (
	noHazard = iota
	icedTiles
	lockedTiles
	fadingFlags
	timeTrialTiles
)

// this should be loaded by the level select screen
var (
	lvlSelectMenu *ebiten.Image
	levelHazard   [4]*ebiten.Image
	startBtn      [3]*ebiten.Image
	powSelMenu    *ebiten.Image
	goldStar      *ebiten.Image
)

func newLevelStartMenu() *levelStartMenu {
	ret := &levelStartMenu{
		// TODO: we need to have a none power that draws a blank
		powOne:   &uiIcon{coord: v2f{88, 71}, size: v2i{16, 16}, img: minusMine[0]},
		powTwo:   &uiIcon{coord: v2f{109, 71}, size: v2i{16, 16}, img: minusMine[0]},
		powThree: &uiIcon{coord: v2f{130, 71}, size: v2i{16, 16}, img: minusMine[0]},
		startBtn: newUIButton(v2f{84, 126}, startBtn),
		powIcons: [7]*uiIcon{
			{coord: v2f{0, 17}, size: v2i{16, 16}, img: addMine[0], powType: addMinePow},
			{coord: v2f{0, 35}, size: v2i{16, 16}, img: scaredyCat[0], powType: scaredyCatPow},
			{coord: v2f{0, 53}, size: v2i{16, 16}, img: tidalWave[0], powType: tidalWavePow},
			{coord: v2f{0, 71}, size: v2i{16, 16}, img: minusMine[0], powType: minusMinePow},
			{coord: v2f{0, 89}, size: v2i{16, 16}, img: dogWistle[0], powType: dogWistlePow},
			{coord: v2f{0, 107}, size: v2i{16, 16}, img: shuffel[0], powType: shuffelPow},
			{coord: v2f{0, 125}, size: v2i{16, 16}, img: dogABone[0], powType: dogABonePow},
		},
	}

	return ret
}

func (l *levelStartMenu) update() {
	l.powOne.update()
	l.powTwo.update()
	l.powThree.update()

	l.loadTime--
	if !l.powOneSel && !l.powTwoSel && !l.powThreeSel && l.loadTime <= 0 {
		l.startBtn.update()
	}

	for _, i := range l.powIcons {
		i.update()
	}

	if l.powOneSel || l.powTwoSel || l.powThreeSel {
		if mbtnp(ebiten.MouseButtonLeft) {
			for _, i := range l.powIcons {
				if i.clicked {
					if l.powOneSel {
						l.powOne.img = i.img
						l.powOne.powType = i.powType
					}
					if l.powTwoSel {
						l.powTwo.img = i.img
						l.powTwo.powType = i.powType
					}
					if l.powThreeSel {
						l.powThree.img = i.img
						l.powThree.powType = i.powType
					}
				}
			}
			l.powOneSel = false
			l.powTwoSel = false
			l.powThreeSel = false
			l.loadTime = 15
			return
		}
	}

	var shift float64
	if l.powOne.clicked {
		l.powOneSel = true
		l.powTwoSel = false
		l.powThreeSel = false
		shift = 108
	}
	if l.powTwo.clicked {
		l.powTwoSel = true
		l.powOneSel = false
		l.powThreeSel = false
		shift = 128
	}
	if l.powThree.clicked {
		l.powThreeSel = true
		l.powOneSel = false
		l.powTwoSel = false
		shift = 148
	}

	if shift > 0 {
		for _, i := range l.powIcons {
			i.coord.x = shift
		}
	}

	if btnp(ebiten.KeyEscape) {
		l.powOneSel = false
		l.powTwoSel = false
		l.powThreeSel = false
	}
}

func (l *levelStartMenu) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 14)
	screen.DrawImage(lvlSelectMenu, op)

	// draw timers
	(&timer{coord: v2f{103, 20}, timerAccumulator: l.levelData.starTimes[0]}).draw(screen)
	(&timer{coord: v2f{103, 33}, timerAccumulator: l.levelData.starTimes[1]}).draw(screen)
	(&timer{coord: v2f{103, 46}, timerAccumulator: l.levelData.starTimes[2]}).draw(screen)

	// draw powerups
	l.powOne.draw(screen)
	l.powTwo.draw(screen)
	l.powThree.draw(screen)

	// draw hazards
	hop := &ebiten.DrawImageOptions{}
	hop.GeoM.Translate(88, 103)
	if l.levelData.frozenTileCount > 0 {
		screen.DrawImage(levelHazard[0], hop)
	}
	if l.levelData.lockedTileCount > 0 {
		screen.DrawImage(levelHazard[1], hop)
	}
	if l.levelData.fadeFlags {
		screen.DrawImage(levelHazard[2], hop)
	}
	if l.levelData.timeTrial {
		screen.DrawImage(levelHazard[3], hop)
	}

	// draw stars
	// fmt.Println("star count: ", l.levelData.stars)
	if l.levelData.stars > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(86, 22)
		screen.DrawImage(goldStar, op)
	}
	if l.levelData.stars > 1 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(86, 35)
		screen.DrawImage(goldStar, op)
		op.GeoM.Translate(7, 0)
		screen.DrawImage(goldStar, op)
	}
	if l.levelData.stars > 2 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(90, 45)
		screen.DrawImage(goldStar, op)
		op.GeoM.Translate(-4, 4)
		screen.DrawImage(goldStar, op)
		op.GeoM.Translate(8, 0)
		screen.DrawImage(goldStar, op)
	}

	// draw powerup menu
	l.startBtn.draw(screen)
	if l.powOneSel || l.powTwoSel || l.powThreeSel {
		sop := &ebiten.DrawImageOptions{}
		if l.powOneSel {
			sop.GeoM.Translate(98, 13)
		}
		if l.powTwoSel {
			sop.GeoM.Translate(118, 13)
		}
		if l.powThreeSel {
			sop.GeoM.Translate(138, 13)

		}
		screen.DrawImage(powSelMenu, sop)

		for _, i := range l.powIcons {
			i.draw(screen)
		}
	}
}

type uiIcon struct {
	coord   v2f
	size    v2i
	img     *ebiten.Image
	powType int
	hovered bool
	clicked bool
}

func newPowIcon(powType int) *uiIcon {
	var img *ebiten.Image
	switch powType {
	case addMinePow:
		img = addMine[0]
	case scaredyCatPow:
		img = scaredyCat[0]
	case tidalWavePow:
		img = tidalWave[0]
	case minusMinePow:
		img = minusMine[0]
	case dogWistlePow:
		img = dogWistle[0]
	case shuffelPow:
		img = shuffel[0]
	case dogABonePow:
		img = dogABone[0]
	}
	return &uiIcon{img: img, powType: powType}
}

func (i *uiIcon) update() {
	i.hovered = false
	i.clicked = false

	m := mCoordsF()
	if m.x > i.coord.x && m.x < i.coord.x+float64(i.size.x) &&
		m.y > i.coord.y && m.y < i.coord.y+float64(i.size.y) {
		i.hovered = true
	}

	if mbtnp(ebiten.MouseButtonLeft) && i.hovered {
		i.clicked = true
	}
}

func (i *uiIcon) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(i.coord.x, i.coord.y)
	screen.DrawImage(i.img, op)
}
