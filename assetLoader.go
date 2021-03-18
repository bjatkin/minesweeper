package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var loadedAssets = make(map[string]*ebiten.Image, 100)

func getAsset(fileName string) (*ebiten.Image, error) {
	asset, found := loadedAssets[fileName]
	if found {
		return asset, nil
	}

	asset, _, err := ebitenutil.NewImageFromFile(fileName)
	if err != nil {
		return nil, err
	}

	loadedAssets[fileName] = asset
	return asset, nil
}

func preloadAsset(fileName string) error {
	_, found := loadedAssets[fileName]
	if found {
		return nil
	}

	asset, _, err := ebitenutil.NewImageFromFile(fileName)
	if err != nil {
		return err
	}
	loadedAssets[fileName] = asset
	return nil
}

func unloadAsset(fileName string) {
	loadedAssets[fileName] = nil
}
