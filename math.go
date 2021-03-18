package main

import (
	"image"
	"math"
)

type v2f struct {
	x float64
	y float64
}

func (a v2f) equal(b v2f) bool {
	return a.x == b.x && a.y == b.y
}

func (a v2f) dist(b v2f) int {
	return int(math.Sqrt((a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y)))
}

func (a v2f) Int() v2i {
	return v2i{
		x: int(a.x),
		y: int(a.y),
	}
}

type v2i struct {
	x int
	y int
}

func (a v2i) equal(b v2i) bool {
	return a.x == b.x && a.y == b.y
}

func (a v2i) dist(b v2i) int {
	return int(math.Sqrt(float64(a.x-b.x)*float64(a.x-b.x) + float64(a.y-b.y)*float64(a.y-b.y)))
}

func (a v2i) Float64() v2f {
	return v2f{
		x: float64(a.x),
		y: float64(a.y),
	}
}

func toV2f(pt image.Point) v2f {
	return v2f{float64(pt.X), float64(pt.Y)}
}

func toV2i(pt image.Point) v2i {
	return v2i{pt.X, pt.Y}
}
