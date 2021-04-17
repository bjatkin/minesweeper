package main

import "github.com/hajimehoshi/ebiten/v2"

// alert types
const (
	quitGameAlertType = iota
	exitMapAlertType
	deleteGameAlertType
)

type alertDialogue struct {
	coord     v2f
	alertType int
	yesBtn    *uiButton
	noBtn     *uiButton
	yes       bool
	no        bool
}

var (
	quitGameAlert *ebiten.Image
	exitMapAlert  *ebiten.Image
	deleteAlert   *ebiten.Image
	yesBtn        [3]*ebiten.Image
	noBtn         [3]*ebiten.Image
)

// TODO: fixup these coord positions
func newAlertDialogue(coord v2f, alertType int) *alertDialogue {
	yesCoord := coord
	yesCoord.x += 15
	yesCoord.y += 10

	noCoord := coord
	yesCoord.x += 5
	yesCoord.y += 10

	return &alertDialogue{
		coord:     coord,
		alertType: alertType,
		yesBtn:    newUIButton(yesCoord, yesBtn),
		noBtn:     newUIButton(noCoord, noBtn),
	}
}

func loadAlertDialogue(ss *ebiten.Image) {
	quitGameAlert = subImage(ss, 904, 0, 96, 57)
	exitMapAlert = subImage(ss, 904, 64, 96, 64)
	deleteAlert = subImage(ss, 904, 128, 96, 64)
	yesBtn = [3]*ebiten.Image{
		subImage(ss, 904, 192, 34, 16),
		subImage(ss, 904, 208, 34, 16),
		subImage(ss, 904, 225, 34, 16),
	}
	noBtn = [3]*ebiten.Image{
		subImage(ss, 952, 192, 34, 16),
		subImage(ss, 952, 208, 34, 16),
		subImage(ss, 952, 225, 34, 16),
	}
}

func (a *alertDialogue) update() {
	a.yesBtn.update()
	a.noBtn.update()

	if a.yesBtn.clickDone {
		a.yes = true
	}
	if a.noBtn.clickDone {
		a.no = true
	}
}

// TODO: fixup these coord positions
func (a *alertDialogue) draw(screen *ebiten.Image) {
	aop := ebiten.DrawImageOptions{}
	aop.GeoM.Translate(a.coord.x, a.coord.y)
	switch a.alertType {
	case quitGameAlertType:
		screen.DrawImage(quitGameAlert, &aop)
	case exitMapAlertType:
		screen.DrawImage(exitMapAlert, &aop)
	case deleteGameAlertType:
		screen.DrawImage(deleteAlert, &aop)
	}

	a.yesBtn.draw(screen)
	a.noBtn.draw(screen)
}
