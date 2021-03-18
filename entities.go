package main

var allEntities map[int]*entity
var nextEnt int

type entity interface {
	update() error
	draw()
	setID(int)
	getID() int
}

func registerEnt(ent *entity) int {
	allEntities[nextEnt] = ent
	nextEnt++
	return nextEnt - 1
}

func getEnt(id int) (*entity, bool) {
	ent, found := allEntities[id]
	return ent, found
}

func deleteEnt(id int) bool {
	_, found := allEntities[id]
	if !found {
		return false
	}
	allEntities[id] = nil
	return true
}
