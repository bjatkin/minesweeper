package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type starCounter struct {
	coord         v2f
	threeStarTime int64
	twoStarTime   int64
	oneStarTime   int64
	starCount     int8
	timer         *timer
}

func (s *starCounter) getStarCount() int8 {
	var count int8
	if s.timer.time() < s.oneStarTime {
		count = 1
	}
	if s.timer.time() < s.twoStarTime {
		count = 2
	}
	if s.timer.time() < s.threeStarTime {
		count = 3
	}

	if count > s.starCount {
		s.starCount = count
		return count
	}

	return s.starCount
}

func (s *starCounter) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.coord.x, s.coord.y)
	bg := ebiten.NewImage(39, 11)
	bg.Fill(white)
	screen.DrawImage(bg, op)

	var offset float64
	if s.timer.time() < s.threeStarTime {
		offset = 40
		remain := s.threeStarTime - s.timer.time()
		prec := (1 - float64(remain)/float64(s.threeStarTime))

		// we convert to an int and then back to a float to floor the value
		offset -= float64(int(prec * 12))
	} else if s.timer.time() < s.twoStarTime {
		offset = 28
		remain := s.twoStarTime - s.timer.time()
		perc := (1 - float64(remain)/float64(s.twoStarTime-s.threeStarTime))

		// we convert to an int and then back to a float to floor the value
		offset -= float64(int(perc * 12))
	} else if s.timer.time() < s.oneStarTime {
		offset = 28
		remain := s.twoStarTime - s.timer.time()
		perc := (1 - float64(remain)/float64(s.twoStarTime-s.threeStarTime))

		// we convert to an int and then back to a float to floor the value
		offset -= float64(int(perc * 12))
	} else {
		offset = 4
	}

	x1, y1 := starCountHudBG.Bounds().Min.X, starCountHudBG.Bounds().Min.Y
	x2, y2 := starCountHudBG.Bounds().Max.X, starCountHudBG.Bounds().Max.Y
	slide := image.Rect(x1+int(offset), y1, x2, y2)
	op.GeoM.Translate(offset, 0)
	screen.DrawImage(starCountHudBG.SubImage(slide).(*ebiten.Image), op)
	op.GeoM.Translate(-offset, 0)

	screen.DrawImage(starCountHud, op)

	if s.starCount >= 1 {
		screen.DrawImage(star, op)
	}
	if s.starCount >= 2 {
		op.GeoM.Translate(12, 0)
		screen.DrawImage(star, op)
	}
	if s.starCount == 3 {
		op.GeoM.Translate(12, 0)
		screen.DrawImage(star, op)
	}
}
