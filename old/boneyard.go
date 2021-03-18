package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type boneyard struct {
	total  uint
	filled uint
	eid    uint64
	frames int
}

func (b *boneyard) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(5, 48)
	screen.DrawImage(boneAreaTop, op)

	op.GeoM.Translate(3, 2)
	count := int(b.total)
	fill := int(b.filled)
	for count > 0 {
		for i := 0; i < 4; i++ {
			screen.DrawImage(boneShadow, op)

			if fill > 0 {
				up := 3 + math.Sin(float64(b.frames)/20)
				op.GeoM.Translate(0, -up)
				screen.DrawImage(bone, op)
				fill--
				op.GeoM.Translate(0, up)
			}

			op.GeoM.Translate(9, 0)

			count--
			if count <= 0 {
				break
			}
		}
		if count > 0 {
			op.GeoM.Translate(-39, 6)
			screen.DrawImage(boneAreaMid, op)

			op.GeoM.Translate(0, 3)
			screen.DrawImage(boneAreaMid, op)

			op.GeoM.Translate(3, 0)
		}
	}

	op.GeoM.Translate(-21, 8)
	screen.DrawImage(boneAreaBot, op)
}

func (b *boneyard) update() error {
	b.frames++
	return nil
}

func (b *boneyard) setID(eid uint64) {
	b.eid = eid
}

func (b *boneyard) getID() uint64 {
	return b.eid
}
