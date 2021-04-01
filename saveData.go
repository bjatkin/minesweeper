package main

import (
	"io/ioutil"
	"os"
)

type saveGame struct {
	allLevels   [14]*n_levelData
	jeepIndex   int
	levelNumber int
	currentPows [3]int
}

func (s *saveGame) updateSave(jeepIndex int, levelNumber int, pows [3]int) {
	s.allLevels = allLevels
	s.jeepIndex = jeepIndex
	s.currentPows = pows
	s.levelNumber = levelNumber
}

func convInt(i int) []byte {
	return []byte{
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

func toInt(b []byte) int {
	return (int(b[0]) << 24) |
		(int(b[1]) << 16) |
		(int(b[2]) << 8) |
		int(b[3])
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

	data = append(data, byte(unlockedPowers[0].powType))
	data = append(data, byte(unlockedPowers[1].powType))
	data = append(data, byte(unlockedPowers[2].powType))
	data = append(data, byte(unlockedPowers[3].powType))
	data = append(data, byte(unlockedPowers[4].powType))
	data = append(data, byte(unlockedPowers[5].powType))
	data = append(data, byte(unlockedPowers[6].powType))

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

	for i := 0; i < len(allLevels); i++ {
		allLevels[i].loadLvl(data[i*7:])
	}

	s.jeepIndex = int(data[98])
	s.currentPows = [3]int{
		int(data[99]),
		int(data[100]),
		int(data[101]),
	}

	s.levelNumber = int(data[102])

	unlockedPowers[0] = newPowIcon(int(data[103]), unlockedPowers[0].coord)
	unlockedPowers[1] = newPowIcon(int(data[104]), unlockedPowers[1].coord)
	unlockedPowers[2] = newPowIcon(int(data[105]), unlockedPowers[2].coord)
	unlockedPowers[3] = newPowIcon(int(data[106]), unlockedPowers[3].coord)
	unlockedPowers[4] = newPowIcon(int(data[107]), unlockedPowers[4].coord)
	unlockedPowers[5] = newPowIcon(int(data[108]), unlockedPowers[5].coord)
	unlockedPowers[6] = newPowIcon(int(data[109]), unlockedPowers[6].coord)

	return nil
}
