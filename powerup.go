package main

import (
	"fmt"
	"image"
	"log"
	"math"

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
	shake          int
}

var (
	addMine    [2]*ebiten.Image
	minusMine  [2]*ebiten.Image
	tidalWave  [2]*ebiten.Image
	scaredyCat [2]*ebiten.Image
	dogWistle  [2]*ebiten.Image
	shuffel    [2]*ebiten.Image
	dogABone   [2]*ebiten.Image
	locked     [2]*ebiten.Image
)

const (
	lockedPow     = 0
	addMinePow    = 1 // select any (flipped tile) next to one or more (grass tiles). If possible, a (dog) will be added under a random (grass tile)
	minusMinePow  = 2 // select any (flipped tile) next to one or more (grass tiles). If possible, a (dog) will be removed from a random (grass tile)
	tidalWavePow  = 3 // convert up to 4 (grass tile) into (water tile). (water tile) will never contain a (dog)
	scaredyCatPow = 4 // select any (grass tile). That tile and it's neighbors will be flipped. Any (dog) will be removed before being flipping.
	dogWistlePow  = 5 // a random (grass tile) will be hilighted, inidicating that (grass tile) is hiding a (dog)
	shuffelPow    = 6 // shuffle all unflagged mines to new (grass tile)
	dogABonePow   = 7 // prevent a game over once per level
)

func newPowerUp(powType int, boundKey ebiten.Key, timer *timer) *powerUp {
	pow := powerUp{pType: powType, boundKey: boundKey, timer: timer, ready: true}
	nSec := int64(1000000000)
	switch powType {
	case lockedPow:
		pow.icons = locked
	case addMinePow:
		pow.coolDown = 60 * nSec
		pow.icons = addMine
	case minusMinePow:
		pow.coolDown = 60 * nSec
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

func (p *powerUp) update() {
	p.shake--
}

func (p *powerUp) wasSelected() bool {
	if !p.available || p.pType == lockedPow {
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

	return hover && mbtnp(ebiten.MouseButtonLeft)
}

func (p *powerUp) activte() {
	p.ready = false
	p.countDownTimer = p.timer.time() + p.coolDown
	if p.pType == dogABonePow { // this is a single use powerup
		p.countDownTimer = 9223372036854775807
	}
}

func (p *powerUp) draw(screen *ebiten.Image) {
	var offset float64
	if p.shake > 0 {
		offset = math.Sin(float64(tickCounter) * 2)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.coord.x+offset, p.coord.y)
	screen.DrawImage(p.icons[1], op) // background

	if !p.available {
		return
	}

	var fill int
	now := p.timer.time()
	if now < p.countDownTimer {
		fill = int(16*(float64(p.countDownTimer-now)/float64(p.coolDown))) + 2
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
