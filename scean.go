package main

import "github.com/hajimehoshi/ebiten/v2"

type scean interface {
	load() error
	unload() error
	update() error
	draw(*ebiten.Image)
}
