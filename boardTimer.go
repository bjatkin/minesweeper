package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type timer struct {
	coord            v2f
	timerAccumulator int64
	timerStart       int64
	running          bool
}

func (t *timer) start() {
	if !t.running {
		t.timerStart = time.Now().UnixNano()
		t.running = true
	}
}

func (t *timer) stop() {
	if t.running {
		t.timerAccumulator += time.Now().UnixNano() - t.timerStart
		t.running = false
	}
}

func (t *timer) time() int64 {
	total := t.timerAccumulator
	if t.running {
		total += time.Now().UnixNano() - t.timerStart
	}
	return total
}

func (t *timer) draw(screen *ebiten.Image) {
	total := t.timerAccumulator
	if t.running {
		total += time.Now().UnixNano() - t.timerStart
	}

	// nano seconds to ms
	total /= 1000000

	ms := fmt.Sprintf("%d", total%1000)
	if len(ms) == 1 {
		ms = "00" + ms
	}
	if len(ms) == 2 {
		ms = "0" + ms
	}
	sec := fmt.Sprintf("%d", (total/1000)%60)
	if len(sec) == 1 {
		sec = "0" + sec
	}
	min := fmt.Sprintf("%d", (total/1000)/60)
	if len(min) == 1 {
		min = "0" + min
	}

	time := min + ":" + sec + "." + ms
	bigRef := "0123456789:."
	smallRef := "0123456789."
	big := numberBigBlue
	small := numberSmallBlue
	if !t.running {
		big = numberBig
		small = numberSmall
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(t.coord.x, t.coord.y)
	for i, r := range time {
		if i < 5 {
			screen.DrawImage(big[strings.IndexRune(bigRef, r)], op)
			if i == 2 {
				op.GeoM.Translate(4, 0)
			} else {
				op.GeoM.Translate(7, 0)
			}
			continue
		}
		screen.DrawImage(small[strings.IndexRune(smallRef, r)], op)
		if i == 5 {
			op.GeoM.Translate(2, 0)
		} else {
			op.GeoM.Translate(4, 0)
		}
	}
}

type boardTimer struct {
	coord v2f
	timer *timer
	play  *uiButton
	pause *uiButton
}

func newBoardTimer(coord v2f) *boardTimer {
	return &boardTimer{
		coord: coord,
		timer: &timer{},
		play:  newUIButton(coord, timerPlayBtn),
		pause: newUIButton(coord, timerPauseBtn),
	}
}

func (b *boardTimer) update() {
	if b.timer.running {
		b.pause.update()
		if b.pause.wasClicked() {
			b.timer.stop()
		}
		return
	}

	if !b.timer.running {
		b.play.update()
		if b.play.wasClicked() {
			b.timer.start()
		}
	}
}

func (b *boardTimer) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.coord.x+16, b.coord.y)
	screen.DrawImage(timerBG, op)

	if b.timer.running {
		b.pause.draw(screen)
	} else {
		b.play.draw(screen)
	}

	b.timer.coord = b.coord
	b.timer.coord.y += 4
	b.timer.coord.x += 16
	b.timer.draw(screen)
}
