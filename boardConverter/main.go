package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage converter convert.bmp layout_name")
		return
	}
	layoutName := os.Args[2]

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Min.Y
	if width > 64 || height > 64 {
		fmt.Printf("Error: image is too large (%dx%d)\n", width, height)
		return
	}

	var layout []*layoutTile
	white := color.RGBA{255, 255, 255, 255}
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			if !colEqual(img.At(x, y), white) {
				layout = append(layout, &layoutTile{index: v2i{x, y}, adj: [8]int{-1, -1, -1, -1, -1, -1, -1, -1}})
			}
		}
	}

	for _, tile := range layout {
		search := [8]v2i{
			{tile.index.x + 1, tile.index.y + 1},
			{tile.index.x + 1, tile.index.y},
			{tile.index.x + 1, tile.index.y - 1},
			{tile.index.x, tile.index.y + 1},
			{tile.index.x, tile.index.y - 1},
			{tile.index.x - 1, tile.index.y + 1},
			{tile.index.x - 1, tile.index.y},
			{tile.index.x - 1, tile.index.y - 1},
		}
		for i, adj := range layout {
			for ii, find := range search {
				if adj.index.equal(find) {
					tile.adj[ii] = i
					tile.adjCount++
				}
			}
		}
	}

	dest := "./../boardLayouts.go"
	preface, err := ioutil.ReadFile(dest)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	layoutFile := string(preface) + "\nvar " + layoutName + "BoardLayout = []*layoutTile{\n"
	for _, tile := range layout {
		layoutFile += tile.String() + ","
	}
	layoutFile += "\n}"
	ioutil.WriteFile(dest, []byte(layoutFile), os.ModeAppend)
}

type v2i struct {
	x, y int
}

func (a *v2i) equal(b v2i) bool {
	return a.x == b.x && a.y == b.y
}

type layoutTile struct {
	index    v2i
	adj      [8]int
	adjCount int
}

func (t *layoutTile) String() string {
	return fmt.Sprintf(
		"{v2i{%d,%d},[8]int{%d,%d,%d,%d,%d,%d,%d,%d},%d}",
		t.index.x, t.index.y,
		t.adj[0], t.adj[1], t.adj[2], t.adj[3],
		t.adj[4], t.adj[5], t.adj[6], t.adj[7],
		t.adjCount,
	)
}

func colEqual(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}
