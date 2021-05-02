package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed assets/title_screen.png
	titleScreen []byte

	//go:embed assets/sprite_sheet.png
	spriteSheet []byte
)

var loadedAssets = make(map[string]*ebiten.Image, 100)

func getAsset(fileName string) (*ebiten.Image, error) {
	asset, found := loadedAssets[fileName]
	if found {
		return asset, nil
	}

	reader := bytes.NewReader(titleScreen)
	fmt.Println("FILE NAME: ", fileName)
	if fileName == "assets/sprite_sheet.png" {
		fmt.Println("Loading Sprite Sheet", len(spriteSheet))
		reader = bytes.NewReader(spriteSheet)
	} else {
		fmt.Println("Loading Title Screen", len(titleScreen))
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	ebtImg := ebiten.NewImageFromImage(img)
	loadedAssets[fileName] = ebtImg
	fmt.Println("size", ebtImg.Bounds().Max)
	fmt.Println("-----------------")

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
