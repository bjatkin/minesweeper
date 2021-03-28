package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type pxlFont struct {
	fontImage           *ebiten.Image
	sprWidth, sprHeight uint8
	fontwidth           [][]int8
	lineHeight          uint8
	baseLine            []uint8
	stringRef           []string
}

func load6x8Font() *pxlFont {
	img, _, err := ebitenutil.NewImageFromFile("assets/pixel_font.png")
	if err != nil {
		log.Fatal(err)
	}
	font := &pxlFont{
		fontImage: img,
		sprWidth:  6,
		sprHeight: 8,
		fontwidth: [][]int8{
			{6, 6, 6, 6, 6, 6, 6, 6, 5, 6, 6, 6, 7, 6, 6, 6, 6, 6, 6, 6, 6, 6, 7, 6, 7, 7},
			{5, 5, 4, 5, 4, 4, 4, 4, 2, 4, 4, 3, 6, 4, 4, 4, 4, 3, 3, 4, 5, 4, 6, 4, 5, 4},
			{6, 3, 6, 6, 6, 4, 4, 6, 4, 4, 4, 5, 5, 4, 4, 4, 4, 4, 4, 2, 4, 4, 3, 3, 4, 4, 4, 4, 3, 5},
			{4, 4, 4, 4, 4, 4, 4, 4, 4, 5},
		},
		lineHeight: 10,
		baseLine:   []uint8{0, 2, 1, 1},
		stringRef: []string{
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"abcdefghijklmnopqrstuvwxyz",
			"~!@#$%^&*()-_+={}[]|\\/:;\",<>.?",
			"1234567890",
		},
	}
	return font
}

func pxlPrint(dest *ebiten.Image, font *pxlFont, x, y float64, str string) (*ebiten.Image, error) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	for _, r := range str {
		if r == ' ' {
			op.GeoM.Translate(float64(font.sprWidth), 0)
			continue
		}
		i := -1
		line := 0
		for _, ref := range font.stringRef {
			i = strings.IndexRune(ref, r)
			if i >= 0 {
				break
			}
			line++
		}
		if i == -1 {
			return dest, fmt.Errorf("Unknow character %s", string(r))
		}
		width := font.fontwidth[line][i]
		height := font.baseLine[line]
		op.GeoM.Translate(0, float64(height))

		dest.DrawImage(font.fontImage.SubImage(
			image.Rect(
				i*int(font.sprWidth),
				line*int(font.sprHeight),
				i*int(font.sprWidth)+int(font.sprWidth),
				line*int(font.sprHeight)+int(font.sprHeight),
			),
		).(*ebiten.Image), op)
		op.GeoM.Translate(float64(width), -float64(height))
	}

	return dest, nil
}

func pxlLen(font *pxlFont, str string) int {
	var len int
	for _, r := range str {
		if r == ' ' {
			len += int(font.sprWidth)
			continue
		}
		i := -1
		line := 0
		for _, ref := range font.stringRef {
			i = strings.IndexRune(ref, r)
			if i >= 0 {
				break
			}
			line++
		}
		if i == -1 {
			len += int(font.sprWidth)
		}
		len += int(font.fontwidth[line][i])
	}

	return len
}
