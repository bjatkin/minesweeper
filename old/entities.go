package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var allEntities = []map[uint32]entity{}
var nextEID = []uint32{}

type entity interface {
	draw(screen *ebiten.Image)
	update() error
	setID(uint64)
	getID() uint64
}

func registerEnt(ent entity, layer uint32) uint64 {
	for i := len(allEntities); i <= int(layer); i++ {
		// TODO: is 500 the right number here?
		allEntities = append(allEntities, make(map[uint32]entity, 500))
		nextEID = append(nextEID, 0)
	}

	allEntities[layer][nextEID[layer]] = ent
	id := uint64(nextEID[layer])<<32 | uint64(layer)
	ent.setID(id)
	nextEID[layer]++
	return uint64(id)
}

func findEnt(eid uint64) (entity, bool) {
	layer := uint32(eid << 32 >> 32)
	smallEID := uint32(eid >> 32)
	ent, found := allEntities[layer][smallEID]
	return ent, found
}

func destroyEnt(eid uint64) bool {
	layer := uint32(eid << 32 >> 32)
	smallEID := uint32(eid >> 32)
	_, found := allEntities[layer][smallEID]
	if !found {
		return false
	}

	allEntities[layer][smallEID] = nil
	return true
}

var allUpdatables = map[uint]updatable{}
var nextUID = uint(0)

type updatable interface {
	update() error
}

func registerUp(up updatable) uint {
	allUpdatables[nextUID] = up
	nextUID++
	return nextUID - 1
}
