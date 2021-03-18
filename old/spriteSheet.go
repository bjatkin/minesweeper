package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	spriteSheet    *ebiten.Image // the main sprite sheet
	gameArea       *ebiten.Image
	pauseMenu      *ebiten.Image
	grassSpr       [7]*ebiten.Image
	yellowGrassSpr [7]*ebiten.Image
	pinkGrassSpr   [7]*ebiten.Image
	blueGrassSpr   [7]*ebiten.Image
	redGrassSpr    [7]*ebiten.Image
	grassSpot      *ebiten.Image
	grassBG        *ebiten.Image
	bone           *ebiten.Image
	boneShadow     *ebiten.Image
	star           *ebiten.Image
	flagAni        aniSprite
	cursorAni      aniSprite
	dogAni         aniSprite
	barkBubble     *ebiten.Image
	boneAreaTop    *ebiten.Image
	boneAreaMid    *ebiten.Image
	boneAreaBot    *ebiten.Image
	primaryButton  [3]*ebiten.Image
	secondButton   [3]*ebiten.Image
)

func loadSprites() {
	var err error
	spriteSheet, _, err = ebitenutil.NewImageFromFile("assets/sprites_test_2.png")
	if err != nil {
		log.Fatal(err)
	}

	gameArea, _, err = ebitenutil.NewImageFromFile("assets/game_board.png")
	if err != nil {
		log.Fatal(err)
	}

	pauseMenu, _, err = ebitenutil.NewImageFromFile("assets/pause.png")
	if err != nil {
		log.Fatal(err)
	}

	grassSpr = [7]*ebiten.Image{
		subImage(0, 0, 16, 16),
		subImage(16, 0, 16, 16),
		subImage(32, 0, 16, 16),
		subImage(48, 0, 16, 16),
		subImage(64, 0, 16, 16),
		subImage(80, 0, 16, 16),
		subImage(96, 0, 16, 16),
	}
	yellowGrassSpr = [7]*ebiten.Image{
		subImage(0, 16, 16, 16),
		subImage(16, 16, 16, 16),
		subImage(32, 16, 16, 16),
		subImage(48, 16, 16, 16),
		subImage(64, 16, 16, 16),
		subImage(80, 16, 16, 16),
		subImage(96, 16, 16, 16),
	}
	pinkGrassSpr = [7]*ebiten.Image{
		subImage(0, 32, 16, 16),
		subImage(16, 32, 16, 16),
		subImage(32, 32, 16, 16),
		subImage(48, 32, 16, 16),
		subImage(64, 32, 16, 16),
		subImage(80, 32, 16, 16),
		subImage(96, 32, 16, 16),
	}
	// TODO: red grass and blue grass
	grassSpot = subImage(0, 80, 16, 16)
	grassBG = subImage(0, 96, 16, 16)
	bone = subImage(72, 80, 8, 8)
	boneShadow = subImage(72, 88, 8, 8)
	star = subImage(80, 80, 16, 16)

	flagAni = newAniSprite(
		[]*ebiten.Image{
			subImage(0, 112, 16, 16),
			subImage(16, 112, 16, 16),
			subImage(32, 112, 16, 16),
			subImage(48, 112, 16, 16),
		},
		[]uint16{8, 8, 8, 8},
		true,
	)
	flagAni.play()
	registerUp(&flagAni)

	cursorAni = newAniSprite(
		[]*ebiten.Image{
			subImage(64, 112, 16, 16),
			subImage(80, 112, 16, 16),
		},
		[]uint16{28, 28},
		true,
	)
	cursorAni.play()
	registerUp(&cursorAni)

	dogAni = newAniSprite(
		[]*ebiten.Image{
			subImage(112, 0, 32, 16),
			subImage(144, 0, 32, 16),
			subImage(176, 0, 32, 16),
			subImage(208, 0, 32, 16),
		},
		[]uint16{4, 4, 4, 4},
		false,
	)
	registerUp(&dogAni)

	barkBubble = subImage(112, 16, 16, 16)

	boneAreaTop = subImage(16, 80, 48, 8)
	boneAreaMid = subImage(16, 88, 48, 8)
	boneAreaBot = subImage(16, 96, 48, 8)

	primaryButton = [3]*ebiten.Image{
		subImage(184, 16, 56, 16),
		subImage(184, 32, 56, 16),
		subImage(184, 48, 56, 16),
	}

	secondButton = [3]*ebiten.Image{
		subImage(184, 64, 56, 16),
		subImage(184, 80, 56, 16),
		subImage(184, 96, 56, 16),
	}
}

func subImage(x, y, w, h uint) *ebiten.Image {
	return spriteSheet.SubImage(image.Rect(int(x), int(y), int(x+w), int(y+h))).(*ebiten.Image)
}

type transform struct {
	x, y float64
}

type sprite struct {
	flip  bool
	scale float64
	rot   float64
	spr   *ebiten.Image
}

func newSprite(spr *ebiten.Image) sprite {
	return sprite{scale: 1, spr: spr}
}

func (s *sprite) draw(screen *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	if s.flip {
		op.GeoM.Scale(-s.scale, s.scale)
	} else {
		op.GeoM.Scale(s.scale, s.scale)
	}
	op.GeoM.Rotate(s.rot)

	screen.DrawImage(s.spr, op)
}

type aniSprite struct {
	flip     bool
	scale    float64
	rot      float64
	loop     bool
	sprs     []*ebiten.Image
	timing   []uint16
	count    uint16
	frame    int
	playFlag bool
}

func newAniSprite(sprs []*ebiten.Image, timing []uint16, loop bool) aniSprite {
	return aniSprite{
		scale:  1,
		sprs:   sprs,
		timing: timing,
		loop:   loop,
	}
}

func (s *aniSprite) draw(screen *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	if s.flip {
		op.GeoM.Scale(-s.scale, s.scale)
	} else {
		op.GeoM.Scale(s.scale, s.scale)
	}
	op.GeoM.Rotate(s.rot)

	screen.DrawImage(s.sprs[s.frame], op)
}

func (s *aniSprite) update() error {
	if s.playFlag {
		s.count++
		if s.count >= s.timing[s.frame] {
			s.count = 0
			s.frame++
		}
		if s.frame >= len(s.sprs) {
			if s.loop {
				s.frame = 0
			} else {
				s.frame--
			}
		}
	}

	return nil
}

func (s *aniSprite) play() {
	s.playFlag = true
}

func (s *aniSprite) isPlaying() bool {
	return s.playFlag
}

func (s *aniSprite) pause() {
	s.playFlag = false
}

func (s *aniSprite) reset() {
	s.frame = 0
	s.playFlag = false
}
