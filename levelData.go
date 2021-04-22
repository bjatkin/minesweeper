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
	ditherBoardLayout, 25 Hard
	croseeBoardLayout, 25 Hard
	mineBoardLayout, 50 Med
*/

// al the levels
var allLevels = [14]*n_levelData{
	// TEST LEVELS
	{layout: mineBoardLayout, mineCount: 40},
	// level 1
	// {
	// 	layout:      newSquareLayout(8, 8),
	// 	mineCount:   10,
	// 	starTimes:   [3]int64{60 * sec, 45 * sec, 35 * sec},
	// 	bestTime:    65 * sec,
	// 	unlocked:    true,
	// 	unlockedPow: minusMinePow,
	// 	unlockSlot:  1,
	// 	nextLevel:   1,
	// },
	// level 2
	{
		layout:      newSquareLayout(16, 16),
		mineCount:   40,
		starTimes:   [3]int64{3*min + 35*sec, 3*min + 15*sec, 3 * min},
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
		starTimes:       [3]int64{3*min + 35*sec, 3*min + 15*sec, 3 * min},
		bestTime:        3*min + 10*sec,
		unlockSlot:      2,
		nextLevel:       4,
	},
	// level 5
	{
		layout:      newSquareLayout(16, 30),
		mineCount:   99,
		starTimes:   [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		unlockedPow: dogWistlePow,
		bestTime:    5*min + 45*sec,
		nextLevel:   5,
	},
	// level 6
	{
		layout:    newSquareLayout(16, 30),
		mineCount: 99,
		starTimes: [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		bestTime:  5*min + 45*sec,
		nextLevel: 6,
	},
	// level 7
	{
		layout:      newSquareLayout(16, 30),
		mineCount:   99,
		starTimes:   [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		unlockedPow: scaredyCatPow,
		bestTime:    5*min + 45*sec,
		nextLevel:   7,
	},
	// level 8
	{
		layout:     newSquareLayout(16, 30),
		mineCount:  99,
		starTimes:  [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		bestTime:   5*min + 45*sec,
		unlockSlot: 3,
		nextLevel:  8,
	},
	// level 9
	{
		layout:    newSquareLayout(16, 30),
		mineCount: 99,
		starTimes: [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		bestTime:  5*min + 45*sec,
		nextLevel: 9,
	},
	// level 10
	{
		layout:      newSquareLayout(16, 30),
		mineCount:   99,
		starTimes:   [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		unlockedPow: shuffelPow,
		bestTime:    5*min + 45*sec,
		nextLevel:   10,
	},
	// level 11
	{
		layout:    newSquareLayout(16, 30),
		mineCount: 99,
		starTimes: [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		bestTime:  5*min + 45*sec,
		nextLevel: 11,
	},
	// level 12
	{
		layout:    newSquareLayout(16, 30),
		mineCount: 99,
		starTimes: [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		bestTime:  5*min + 45*sec,
		nextLevel: 12,
	},
	// level 13
	{
		layout:    newSquareLayout(16, 30),
		mineCount: 99,
		starTimes: [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		bestTime:  5*min + 45*sec,
		nextLevel: 13,
	},
	// level 14
	{
		layout:      newSquareLayout(16, 30),
		mineCount:   99,
		starTimes:   [3]int64{5*min + 30*sec, 4*min + 15*sec, 3*min + 25*sec},
		bestTime:    5*min + 45*sec,
		unlockedPow: dogABonePow,
		nextLevel:   14,
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
