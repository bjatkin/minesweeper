package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type powerUp struct {
	coord          v2f
	coolDown       int64
	countDownTimer int64
	icons          [2]*ebiten.Image
	pType          int
	available      bool
	ready          bool
	boundKey       ebiten.Key
	timer          *timer
}

var (
	addMine    [2]*ebiten.Image
	minusMine  [2]*ebiten.Image
	tidalWave  [2]*ebiten.Image
	scaredyCat [2]*ebiten.Image
	dogWistle  [2]*ebiten.Image
	shuffel    [2]*ebiten.Image
	dogABone   [2]*ebiten.Image
)

const (
	addMinePow = iota
	minusMinePow
	tidalWavePow
	scaredyCatPow
	dogWistlePow
	shuffelPow
	dogABonePow
)

func newPowerUp(powType int, boundKey ebiten.Key, timer *timer) *powerUp {
	pow := powerUp{pType: powType, boundKey: boundKey, timer: timer, ready: true}
	nSec := int64(1000000000)
	switch powType {
	case addMinePow:
		pow.coolDown = 120 * nSec
		pow.icons = addMine
	case minusMinePow:
		pow.coolDown = 120 * nSec
		pow.icons = minusMine
	case tidalWavePow:
		pow.coolDown = 180 * nSec
		pow.icons = tidalWave
	case scaredyCatPow:
		pow.coolDown = 180 * nSec
		pow.icons = scaredyCat
	case dogWistlePow:
		pow.coolDown = 180 * nSec
		pow.icons = dogWistle
	case shuffelPow:
		pow.coolDown = 180 * nSec
		pow.icons = shuffel
	case dogABonePow:
		pow.icons = dogABone
	default:
		log.Fatal(fmt.Errorf("unknown powerup id: %d", powType))
	}
	return &pow
}

func (p *powerUp) wasSelected() bool {
	if !p.available || !p.ready {
		return false
	}

	now := p.timer.time()
	if now < p.countDownTimer {
		return false
	}
	p.ready = true

	if btnp(p.boundKey) {
		return true
	}

	mouse := mCoordsF()
	var hover bool
	if mouse.x > p.coord.x &&
		mouse.x < p.coord.x+16 &&
		mouse.y > p.coord.y &&
		mouse.y < p.coord.y+16 {
		hover = true
	}

	if hover && mbtnp(ebiten.MouseButtonLeft) {
	}

	return hover && mbtnp(ebiten.MouseButtonLeft)
}

func (p *powerUp) activte() {
	p.ready = false
	p.countDownTimer = p.timer.time() + p.coolDown
	// p.countDownTimer = time.Now().UnixNano() + p.coolDown
	if p.pType == dogABonePow { // this is a single use powerup
		p.countDownTimer = 9223372036854775807
	}
}

func (p *powerUp) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.coord.x, p.coord.y)
	screen.DrawImage(p.icons[1], op) // background

	if !p.available {
		return
	}

	var fill int
	now := p.timer.time()
	// now := time.Now().UnixNano()
	if now < p.countDownTimer {
		fill = int(16 * (float64(p.countDownTimer-now) / float64(p.coolDown)))
	}

	rect := image.Rect(
		p.icons[0].Bounds().Min.X,
		p.icons[0].Bounds().Min.Y+fill,
		p.icons[0].Bounds().Min.X+16,
		p.icons[0].Bounds().Min.Y+16,
	)
	op.GeoM.Translate(0, float64(fill))
	screen.DrawImage( // foreground
		p.icons[0].SubImage(rect).(*ebiten.Image),
		op,
	)
}
