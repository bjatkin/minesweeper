package main

import "github.com/hajimehoshi/ebiten/v2"

type n_aniSprite struct {
	sprs     []*ebiten.Image
	timing   []uint
	count    uint
	frame    int
	loop     bool
	done     bool
	playFlag bool
}

func n_newAniSprite(sprs []*ebiten.Image, timing []uint, loop bool) *n_aniSprite {
	return &n_aniSprite{
		sprs:   sprs,
		timing: timing,
		loop:   loop,
	}
}

func (s *n_aniSprite) update() {
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
				s.done = true
			}
		}
	}
}

func (s *n_aniSprite) draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	screen.DrawImage(s.sprs[s.frame], op)
}

func (s *n_aniSprite) next() {
	if s.frame == len(s.sprs)-1 {
		return
	}
	s.frame++
}

func (s *n_aniSprite) prev() {
	if s.frame == 0 {
		return
	}
	s.frame--
}

func (s *n_aniSprite) play() {
	s.playFlag = true
}

func (s *n_aniSprite) pause() {
	s.playFlag = false
}

func (s *n_aniSprite) reset() {
	s.frame = 0
	s.playFlag = false
}
