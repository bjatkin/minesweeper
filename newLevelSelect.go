package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type levelSelect struct {
	scrollPt      float64
	jeepCoord     v2f
	jeepIndex     int
	jeepFlip      bool
	jeepGoalIndex int
	levelPoints   []v2f
	connectPoints []v2f
	levelShake    []int
	selectLevel   bool
	currLevel     *n_levelData
	startMenu     *levelStartMenu
	levelNumber   int
	bestTime      *timer
}

// level select assets
var (
	levelSelectAssetsLoaded bool
	jeep                    *n_aniSprite
	mapBG                   *ebiten.Image
	water                   *n_aniSprite
	nessie                  *n_aniSprite
	whiteCloud              [4]*ebiten.Image
	blueCloud               [4]*ebiten.Image
	unlockedLvl             *ebiten.Image
	newLvl                  *ebiten.Image
	lvlConnect              *ebiten.Image
	mapUIHeader             *ebiten.Image
	redStar                 *ebiten.Image
)

func (l *levelSelect) load() error {
	cursorHold = false
	ss, err := getAsset("assets/sprite_sheet.png")
	if err != nil {
		return err
	}

	jeep = n_newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 77, 32, 20, 16),
			subImage(ss, 96, 32, 20, 16),
		},
		[]uint{15, 15},
		true,
	)
	jeep.play()

	mapBG = subImage(ss, 160, 168, 240, 240)
	water = n_newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 184, 152, 16, 16),
			subImage(ss, 200, 152, 16, 16),
		},
		[]uint{20, 20},
		true,
	)
	water.play()
	nessie = n_newAniSprite(
		[]*ebiten.Image{
			subImage(ss, 216, 152, 16, 16),
			subImage(ss, 232, 152, 16, 16),
		},
		[]uint{20, 20},
		true,
	)
	nessie.play()

	whiteCloud = [4]*ebiten.Image{
		subImage(ss, 400, 168, 104, 128), // NW
		subImage(ss, 520, 168, 120, 120), // NE
		subImage(ss, 400, 328, 48, 80),   // SW
		subImage(ss, 592, 360, 48, 48),   // SE
	}

	blueCloud = [4]*ebiten.Image{
		subImage(ss, 640, 168, 104, 128), // NW
		subImage(ss, 760, 168, 120, 120), // NE
		subImage(ss, 640, 328, 120, 120), // SW
		subImage(ss, 832, 360, 48, 48),   // SE
	}

	unlockedLvl = subImage(ss, 176, 16, 17, 16)
	newLvl = subImage(ss, 216, 16, 17, 16)
	lvlConnect = subImage(ss, 200, 24, 4, 4)

	lvlSelectMenu = subImage(ss, 0, 16, 80, 136)

	mapUIHeader = subImage(ss, 240, 16, 240, 10)

	levelHazard = [4]*ebiten.Image{
		subImage(ss, 224, 0, 16, 16),
		subImage(ss, 240, 0, 16, 16),
		subImage(ss, 256, 0, 16, 16),
		subImage(ss, 272, 0, 16, 16),
	}

	startBtn = [3]*ebiten.Image{
		subImage(ss, 80, 136, 66, 16),
		subImage(ss, 80, 120, 66, 16),
		subImage(ss, 80, 104, 66, 16),
	}

	powSelMenu = subImage(ss, 152, 32, 30, 132)
	redStar = subImage(ss, 208, 16, 8, 8)
	goldStar = subImage(ss, 200, 16, 8, 8)

	// These assets are defined in the newBoard file but we need it here as well
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
	// end foreign assets

	l.startMenu = newLevelStartMenu()
	l.bestTime = &timer{coord: v2f{194, 1}, timerAccumulator: allLevels[0].bestTime}
	l.currLevel = allLevels[0]

	l.levelPoints = []v2f{
		{40, 220},  // start
		{94, 217},  // palm grove
		{202, 206}, // secret forest
		{219, 152}, // far east
		{178, 139}, // transition lands
		{130, 155}, // on the savana
		{80, 149},  // ocean side
		{41, 112},  // deep ocean 1
		{99, 96},   // camping
		{178, 73},  // lakeside
		{196, 28},  // the canyon
		{138, 44},  // over the bridge
		{44, 50},   // deep ocean 1
		{86, 16},   // final battle
	}

	l.levelShake = make([]int, len(l.levelPoints))

	l.connectPoints = []v2f{
		{42, 222}, // start
		{60, 230},
		{75, 233},
		{85, 227},
		{96, 219}, // palm grov
		{122, 224},
		{140, 228},
		{166, 230},
		{190, 221},
		{204, 208}, // secret forest
		{222, 192},
		{226, 176},
		{221, 154}, // far east
		{217, 142},
		{209, 131},
		{191, 132},
		{180, 141}, // transition lands
		{167, 147},
		{155, 155},
		{132, 157}, // on the savana
		{121, 163},
		{104, 158},
		{82, 151}, // ocean side
		{71, 142},
		{59, 130},
		{43, 114}, // deep ocean 1
		{66, 108},
		{86, 103},
		{101, 98}, // camping
		{128, 99},
		{148, 94},
		{165, 85},
		{180, 75}, // lakeside
		{202, 74},
		{213, 61},
		{209, 47},
		{198, 30}, // the canyon
		{193, 19},
		{173, 18},
		{140, 46}, // over the bridge
		{155, 31},
		{127, 52},
		{111, 57},
		{96, 58},
		{79, 58},
		{63, 56},
		{88, 18}, // final battle
		{59, 43},
		{68, 33},
		{81, 23},
	}

	l.scrollPt = 80
	l.jeepCoord = l.connectPoints[l.jeepIndex]

	l.levelNumber = 1

	return nil
}

func (l *levelSelect) unload() error {
	// these assets are part of the core gameloop and are
	// pretty small so we don't really want to unload and
	// reload them a lot
	return nil
}

func (l *levelSelect) update() error {
	for i := 0; i < len(l.levelShake); i++ {
		l.levelShake[i]--
	}
	jeep.update()

	for l.jeepCoord.y-l.scrollPt > 100 {
		l.scrollPt++
		if l.scrollPt > 80 {
			l.scrollPt = 80
			break
		}
	}
	for l.jeepCoord.y-l.scrollPt < 60 {
		l.scrollPt--
		if l.scrollPt < 0 {
			l.scrollPt = 0
			break
		}
	}

	if !l.selectLevel {
		var clickedLvl bool
		if mbtnp(ebiten.MouseButtonLeft) {
			mouse := mCoordsF()
			mouse.y += l.scrollPt
			for i, lvl := range l.levelPoints {
				if lvl.dist(mouse) < 20 {
					if allLevels[i].unlocked {
						l.levelNumber = i + 1
						for ii, c := range l.connectPoints {
							if c.dist(lvl) < 5 {
								l.currLevel = allLevels[i]
								l.startMenu.levelData = allLevels[i]
								l.bestTime.timerAccumulator = allLevels[i].bestTime
								l.jeepGoalIndex = ii
								if l.jeepGoalIndex == l.jeepIndex {
									clickedLvl = true
								}
								break
							}
						}
					} else {
						l.levelShake[i] = 30
					}
					break
				}
			}
		}

		if btnp(ebiten.KeyEnter) {
			clickedLvl = true
		}

		goal := l.connectPoints[l.jeepIndex]
		if l.jeepCoord.dist(goal) < 2 {
			if l.jeepIndex < l.jeepGoalIndex {
				l.jeepIndex++
			}
			if l.jeepIndex > l.jeepGoalIndex {
				l.jeepIndex--
			}
		} else {
			speed := 0.75
			if l.jeepCoord.x < goal.x {
				l.jeepFlip = false
				l.jeepCoord.x += speed
			}
			if l.jeepCoord.x > goal.x+speed {
				l.jeepFlip = true
				l.jeepCoord.x -= speed
			}
			if l.jeepCoord.y < goal.y {
				l.jeepCoord.y += speed
			}
			if l.jeepCoord.y > goal.y+speed {
				l.jeepCoord.y -= speed
			}
		}

		if l.jeepCoord.dist(goal) < 2 &&
			l.jeepIndex == l.jeepGoalIndex &&
			clickedLvl {
			// pull up the level select menu
			l.selectLevel = true
			l.startMenu.loadTime = 30
		}

		return nil
	}

	if l.selectLevel {
		if btnp(ebiten.KeyEscape) {
			l.selectLevel = false
		}

		m := mCoordsF()
		if mbtnp(ebiten.MouseButtonLeft) &&
			(m.x < 80 || m.x > 160 ||
				m.y < 14 || m.y > 160) {
			l.selectLevel = false
		}
		l.startMenu.update()
	}

	if l.startMenu.startBtn.clicked {
		currentScean = newLevelScean(
			l.startMenu.levelData,
			[3]int{
				l.startMenu.powOne.powType,
				l.startMenu.powTwo.powType,
				l.startMenu.powThree.powType,
			},
		)

		err := currentScean.load()
		if err != nil {
			return err
		}

		err = l.unload()
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *levelSelect) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, -l.scrollPt)
	screen.DrawImage(mapBG, op)

	for i, pt := range l.connectPoints {
		offX := math.Sin(float64((int(tickCounter) + i*15) / 25))
		offY := math.Sin(float64((int(tickCounter) + i*10) / 20))
		cop := &ebiten.DrawImageOptions{}
		cop.GeoM.Translate(pt.x+offX, pt.y-l.scrollPt+offY)
		screen.DrawImage(lvlConnect, cop)
	}

	for i, pt := range l.levelPoints {
		lop := &ebiten.DrawImageOptions{}
		lop.GeoM.Translate(pt.x, pt.y-l.scrollPt)
		if l.levelShake[i] > 0 {
			lop.GeoM.Translate(math.Sin(float64(tickCounter)), 0)
		}
		if allLevels[i].unlocked {
			// draw the stars on the level icon
			screen.DrawImage(unlockedLvl, lop)
			if allLevels[i].stars > 0 {
				lop.GeoM.Translate(6, 2)
				screen.DrawImage(goldStar, lop)
				lop.GeoM.Translate(-6, -2)
			}
			if allLevels[i].stars > 1 {
				lop.GeoM.Translate(2, 6)
				screen.DrawImage(goldStar, lop)
				lop.GeoM.Translate(-2, -6)
			}
			if allLevels[i].stars > 2 {
				lop.GeoM.Translate(10, 6)
				screen.DrawImage(goldStar, lop)
				lop.GeoM.Translate(-10, -6)
			}
		} else {
			screen.DrawImage(newLvl, lop)
		}
	}

	jop := &ebiten.DrawImageOptions{}
	if l.jeepFlip {
		jop.GeoM.Scale(-1, 1)
		jop.GeoM.Translate(28, 0)
	}
	jop.GeoM.Translate(l.jeepCoord.x-8, l.jeepCoord.y-l.scrollPt-4)

	jeep.draw(screen, jop)

	// draw the clouds
	cop := &ebiten.DrawImageOptions{}
	cop.GeoM.Translate(0, 8-l.scrollPt*2)
	screen.DrawImage(whiteCloud[0], cop)
	cop.GeoM.Reset()
	cop.GeoM.Translate(120, 8-l.scrollPt*2)
	screen.DrawImage(whiteCloud[1], cop)
	cop.GeoM.Reset()
	cop.GeoM.Translate(0, 240-l.scrollPt*2)
	screen.DrawImage(whiteCloud[2], cop)
	cop.GeoM.Reset()
	cop.GeoM.Translate(192, 272-l.scrollPt*2)
	screen.DrawImage(whiteCloud[3], cop)

	// blue clouds
	cop.GeoM.Reset()
	cop.GeoM.Translate(0, 8-l.scrollPt*2)
	screen.DrawImage(blueCloud[0], cop)
	cop.GeoM.Reset()
	cop.GeoM.Translate(120, 8-l.scrollPt*2)
	screen.DrawImage(blueCloud[1], cop)
	cop.GeoM.Reset()
	cop.GeoM.Translate(0, 240-l.scrollPt*2)
	screen.DrawImage(blueCloud[2], cop)
	cop.GeoM.Reset()
	cop.GeoM.Translate(192, 272-l.scrollPt*2)
	screen.DrawImage(blueCloud[3], cop)

	// draw the ui
	if l.selectLevel {
		l.startMenu.draw(screen)
	}

	// map ui header
	op.GeoM.Reset()
	screen.DrawImage(mapUIHeader, op)
	op.GeoM.Translate(38, 1)
	if l.levelNumber > 9 {
		screen.DrawImage(numberBig[l.levelNumber/10], op)
		op.GeoM.Translate(6, 0)
		screen.DrawImage(numberBig[l.levelNumber%10], op)
	} else {
		screen.DrawImage(numberBig[l.levelNumber], op)
	}

	// best timer
	l.bestTime.draw(screen)

	// draw the map ui header star
	sop := &ebiten.DrawImageOptions{}
	if l.currLevel.stars > 0 {
		sop.GeoM.Translate(54, 2)
		screen.DrawImage(redStar, sop)
	}
	if l.currLevel.stars > 1 {
		sop.GeoM.Translate(7, 0)
		screen.DrawImage(redStar, sop)
	}
	if l.currLevel.stars > 2 {
		sop.GeoM.Translate(7, 0)
		screen.DrawImage(redStar, sop)
	}
}
