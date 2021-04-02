package main

import (
	"io/ioutil"
	"os"
)

type saveGame struct {
	allLevels      [14]*n_levelData
	jeepIndex      int
	levelNumber    int
	currentPows    [3]int
	unlockedPowers [7]*uiIcon
}

func (s *saveGame) updateSave(jeepIndex int, levelNumber int, pows [3]int) {
	// save global state
	s.allLevels = allLevels
	s.unlockedPowers = unlockedPowers

	// non global state
	s.jeepIndex = jeepIndex
	s.currentPows = pows
	s.levelNumber = levelNumber
}

func convInt(i int) []byte {
	return []byte{
		byte(i >> 56),
		byte(i >> 48),
		byte(i >> 40),
		byte(i >> 32),
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

func toInt(b []byte) int {
	return (int(b[0]) << 56) |
		(int(b[1]) << 48) |
		(int(b[2]) << 40) |
		(int(b[3]) << 32) |
		(int(b[4]) << 24) |
		(int(b[5]) << 16) |
		(int(b[6]) << 8) |
		int(b[7])
}

func convBool(b bool) byte {
	if b {
		return byte(1)
	}
	return byte(0)
}

func toBool(b byte) bool {
	if b == 0 {
		return false
	}
	return true
}

func (s *saveGame) saveData(fileName string) error {
	data := []byte{}
	for _, l := range s.allLevels {
		data = append(data, l.serializeLvl()...)
	}

	data = append(data, byte(s.jeepIndex))
	for _, p := range s.currentPows {
		data = append(data, byte(p))
	}

	data = append(data, byte(s.levelNumber))

	data = append(data, byte(s.unlockedPowers[0].powType))
	data = append(data, byte(s.unlockedPowers[1].powType))
	data = append(data, byte(s.unlockedPowers[2].powType))
	data = append(data, byte(s.unlockedPowers[3].powType))
	data = append(data, byte(s.unlockedPowers[4].powType))
	data = append(data, byte(s.unlockedPowers[5].powType))
	data = append(data, byte(s.unlockedPowers[6].powType))

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *saveGame) loadData(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	for i := 0; i < len(s.allLevels); i++ {
		s.allLevels[i] = &n_levelData{}
		s.allLevels[i].loadLvl(data[i*11:])
	}

	s.jeepIndex = int(data[98])
	s.currentPows = [3]int{
		int(data[155]),
		int(data[156]),
		int(data[157]),
	}

	s.levelNumber = int(data[158])

	s.unlockedPowers[0] = newPowIcon(int(data[159]), v2f{0, 17})
	s.unlockedPowers[1] = newPowIcon(int(data[160]), v2f{0, 35})
	s.unlockedPowers[2] = newPowIcon(int(data[161]), v2f{0, 53})
	s.unlockedPowers[3] = newPowIcon(int(data[162]), v2f{0, 71})
	s.unlockedPowers[4] = newPowIcon(int(data[163]), v2f{0, 89})
	s.unlockedPowers[5] = newPowIcon(int(data[164]), v2f{0, 107})
	s.unlockedPowers[6] = newPowIcon(int(data[165]), v2f{0, 125})

	return nil
}
