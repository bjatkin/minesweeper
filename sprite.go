package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type sprite struct {
	op  *ebiten.DrawImageOptions
	spr *ebiten.Image
}

func newSprite(spr *ebiten.Image) sprite {
	op := &ebiten.DrawImageOptions{}
	return sprite{op: op, spr: spr}
}

func (s *sprite) draw(screen *ebiten.Image) {
	screen.DrawImage(s.spr, s.op)
}

type aniSprite struct {
	op       *ebiten.DrawImageOptions
	sprs     []*ebiten.Image
	timing   []uint16
	count    uint16
	frame    int
	loop     bool
	playFlag bool
}

func newAniSprite(sprs []*ebiten.Image, timing []uint16, loop bool) aniSprite {
	op := &ebiten.DrawImageOptions{}
	return aniSprite{
		op:     op,
		sprs:   sprs,
		timing: timing,
		loop:   loop,
	}
}

func (s *aniSprite) update() {
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
}

func (s *aniSprite) play() {
	s.playFlag = true
}

func (s *aniSprite) pause() {
	s.playFlag = false
}

func (s *aniSprite) reset() {
	s.frame = 0
	s.playFlag = false
}

func (s *aniSprite) draw(screen *ebiten.Image) {
	screen.DrawImage(s.sprs[s.frame], s.op)
}

func subImage(sheet *ebiten.Image, x, y, w, h uint) *ebiten.Image {
	return sheet.SubImage(image.Rect(int(x), int(y), int(x+w), int(y+h))).(*ebiten.Image)
}
