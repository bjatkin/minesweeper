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
	powerups         [3]powerUp
	unlocked         bool
	beaten           bool
}

// al the levels
var allLevels = [14]*n_levelData{
	// level 1
	{
		layout:    newSquareLayout(8, 8),
		mineCount: 10,
		starTimes: [3]int64{50000000, 9000000, 10000000},
		bestTime:  1018402938475,
		unlocked:  true,
	},
	// level 2
	{
		layout:          newSquareLayout(16, 16),
		mineCount:       25,
		starTimes:       [3]int64{},
		unlocked:        true,
		lockedTileCount: 10,
		stars:           3,
	},
	// level 3
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 4
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 5
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 6
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 7
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 8
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 9
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 10
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 11
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 12
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 13
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
	},
	// level 14
	{
		layout:    newSquareLayout(16, 16),
		mineCount: 25,
		starTimes: [3]int64{},
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
