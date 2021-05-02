package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed assets/title_screen.png
	titleScreen []byte

	//go:embed assets/sprite_sheet.png
	spriteSheet []byte

	//go:embed assets/pixel_font.png
	pixelFont []byte
)

var loadedAssets = make(map[string]*ebiten.Image, 100)

func getAsset(fileName string) (*ebiten.Image, error) {
	asset, found := loadedAssets[fileName]
	if found {
		return asset, nil
	}

	reader := bytes.NewReader(titleScreen)
	if fileName == "assets/sprite_sheet.png" {
		reader = bytes.NewReader(spriteSheet)
	}
	if fileName == "assets/pixel_font.png" {
		reader = bytes.NewReader(pixelFont)
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	ebtImg := ebiten.NewImageFromImage(img)
	loadedAssets[fileName] = ebtImg

	return ebtImg, nil
}

// func getAsset(fileName string) (*ebiten.Image, error) {
// 	asset, found := loadedAssets[fileName]
// 	if found {
// 		return asset, nil
// 	}

// 	asset, _, err := ebitenutil.NewImageFromFile(fileName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	loadedAssets[fileName] = asset
// 	return asset, nil
// }

// func preloadAsset(fileName string) error {
// 	_, found := loadedAssets[fileName]
// 	if found {
// 		return nil
// 	}

// 	asset, _, err := ebitenutil.NewImageFromFile(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	loadedAssets[fileName] = asset
// 	return nil
// }

// func unloadAsset(fileName string) {
// 	loadedAssets[fileName] = nil
// }
