package main

type n_levelData struct {
	frozenTileCount  int
	lockedTileCount  int
	timeTrial        bool
	timeTrialCount   int
	timeTrialMaxTime int64
	fadeFlags        bool
	layout           []*layoutTile
	mineCount        int
	bestTime         int64
	stars            int8
	starTimes        [3]int64
	unlocked         bool
	unlockedPow      int
	unlockSlot       int
	beaten           bool
	nextLevel        int
}

func (lvl *n_levelData) serializeLvl() []byte {
	data := []byte{
		convBool(lvl.unlocked),
		convBool(lvl.beaten),
		byte(lvl.stars),
	}

	data = append(data, convInt(int(lvl.bestTime))...)

	return data
}

func (lvl *n_levelData) loadLvl(data []byte) {
	lvl.unlocked = toBool(data[0])
	lvl.beaten = toBool(data[1])
	lvl.stars = int8(data[2])
	lvl.bestTime = int64(toInt(data[3:11]))
}

// one second in nano seconds
var sec = int64(1000000000)
var min = 60 * sec

// FUN LEVEL LIST
/*
	8x8, 10 very easy
	16x16, 40 easy but long
	heart, 40 med easy
	heartBoardLayout, 40, frozenTileCount: 60 Med
	ditherBoardLayout, 25 Hard
	ditherBoardLayout, 25 Hard, lockedTileCount: 25 Hard
	croseeBoardLayout, 25 Hard
	mineBoardLayout, 50 Med
	dogBoardLayout, 50 Med
	squaresBoardLayout, 70 Med + Long
	squaresBoardLayout, 99 Hard + Long
	squares2BoardLayout, 70 Easy + Long
	mineBoardLayout, 60, fadeFlags: true, Hard, final level
*/

// all the levels
var allLevels = [14]*n_levelData{
	// TEST LEVELS
	// level 1
	{
		layout:      newSquareLayout(8, 8),
		mineCount:   10,
		starTimes:   [3]int64{5 * min, 2 * min, 60 * sec},
		bestTime:    4 * min,
		unlocked:    true,
		unlockedPow: minusMinePow,
		unlockSlot:  1,
		nextLevel:   1,
	},
	// level 2
	{
		layout:      newSquareLayout(16, 16),
		mineCount:   40,
		starTimes:   [3]int64{6*min + 35*sec, 4*min + 15*sec, 3 * min},
		bestTime:    3*min + 45*sec,
		unlockedPow: tidalWavePow,
		nextLevel:   2,
	},
	// level 3
	{
		layout:      heartBoardLayout,
		mineCount:   40,
		starTimes:   [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		bestTime:    5*min + 45*sec,
		unlockedPow: addMinePow,
		nextLevel:   3,
	},
	// level 4
	{
		layout:          newSquareLayout(16, 16),
		mineCount:       40,
		frozenTileCount: 50,
		starTimes:       [3]int64{5*min + 35*sec, 4*min + 15*sec, 3 * min},
		bestTime:        3*min + 10*sec,
		unlockSlot:      2,
		nextLevel:       4,
	},
	// level 5
	{
		layout:      newSquareLayout(16, 30),
		mineCount:   99,
		starTimes:   [3]int64{4*min + 30*sec, 4 * min, 3*min + 40*sec},
		unlockedPow: dogWistlePow,
		bestTime:    4*min + 45*sec,
		nextLevel:   5,
	},
	// level 6
	{
		layout:    dogBoardLayout,
		mineCount: 50,
		starTimes: [3]int64{4 * min, 3*min + 30*sec, 3 * min},
		bestTime:  3*min + 40*sec,
		nextLevel: 6,
	},
	// level 7
	{
		layout:          mineBoardLayout,
		mineCount:       60,
		unlockedPow:     scaredyCatPow,
		lockedTileCount: 30,
		starTimes:       [3]int64{5 * min, 4 * min, 3*min + 30*sec},
		bestTime:        4 * min,
		nextLevel:       7,
	},
	// level 8
	{
		layout:     squares2BoardLayout,
		mineCount:  80,
		unlockSlot: 3,
		starTimes:  [3]int64{6 * min, 5 * min, 4*min + 30*sec},
		bestTime:   5 * min,
		nextLevel:  8,
	},
	// level 9
	{
		layout:    crossBoardLayout,
		mineCount: 25,
		nextLevel: 9,
		starTimes: [3]int64{5 * min, 4*min + 30*sec, 3*min + 30*sec},
		bestTime:  5 * min,
	},
	// level 10
	{
		layout:      ditherBoardLayout,
		mineCount:   30,
		unlockedPow: shuffelPow,
		nextLevel:   10,
		starTimes:   [3]int64{3*min + 30*sec, 3 * min, 2 * min},
		bestTime:    3 * min,
	},
	// level 11
	{
		layout:          wheelBoardLayout,
		mineCount:       80,
		frozenTileCount: 50,
		nextLevel:       11,
		starTimes:       [3]int64{5*min + 45*sec, 5 * min, 4 * min},
		bestTime:        5 * min,
	},
	// level 12
	{
		layout:          ditherBoardLayout,
		mineCount:       30,
		lockedTileCount: 30,
		nextLevel:       12,
		starTimes:       [3]int64{3 * min, 2 * min, 1*min + 45*sec},
		bestTime:        3 * min,
	},
	// level 13
	{
		layout:          mineBoardLayout,
		mineCount:       70,
		frozenTileCount: 50,
		nextLevel:       13,
		starTimes:       [3]int64{4 * min, 3*min + 45*sec, 3 * min},
		bestTime:        3*min + 50*sec,
	},
	// level 14
	{
		layout:      newSquareLayout(16, 30),
		mineCount:   99,
		unlockedPow: dogABonePow,
		fadeFlags:   false,
		nextLevel:   14,
		starTimes:   [3]int64{4 * min, 3*min + 30*sec, 3 * min},
		bestTime:    3*min + 45*sec,
	},
}

// Note we may not need the layout tile (we can maybe use use v2i here and chuck the adj stuff)
type layoutTile struct {
	index    v2i
	adj      [8]int
	adjCount int
}

func newSquareLayout(rows, cols int) []*layoutTile {
	var layout []*layoutTile

	for i := 0; i < cols; i++ {
		for ii := 0; ii < rows; ii++ {
			tile := layoutTile{
				index: v2i{i, ii},
				adj:   [8]int{-1, -1, -1, -1, -1, -1, -1, -1},
			}
			layout = append(layout, &tile)
		}
	}

	// OPTIM: with a little hash map magic here we could make this O(n)
	// this sets up all the adj tiles
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

	return layout
}
