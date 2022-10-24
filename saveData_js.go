//go:build js
// +build js

package main

import (
	b64 "encoding/base64"
	"fmt"
	"syscall/js"
)

var CurrentSaveGame = &n_SaveGame{}

const SaveGameByteLen = 256

var SaveGameFile = "ready_set_duck_save_file"

type n_SaveGame struct {
	used           bool
	slot           int
	allLevels      [14]*n_levelData
	jeepIndex      int
	levelNumber    int
	currentPows    [3]int
	unlockedPowers [7]*uiIcon
	loaded         bool
}

func (s *n_SaveGame) saveGame(jeepIndex int, levelNumber int, pows [3]int) error {
	// save global state
	s.allLevels = allLevels
	s.unlockedPowers = unlockedPowers

	// save non global state
	s.jeepIndex = jeepIndex
	s.currentPows = pows
	s.levelNumber = levelNumber

	data := s.toBytes()
	saveFileData, err := readFromLocalStorage(SaveGameFile)
	if err != nil {
		return err
	}

	bl := SaveGameByteLen
	l := copy(saveFileData[s.slot*bl:s.slot*bl+bl], data)
	if l != bl {
		return fmt.Errorf("missing bytes, only %d/%d bytes coppied\n", l, bl)
	}

	writeToLocalStorage(SaveGameFile, saveFileData)

	return nil
}

func (s *n_SaveGame) loadGame(slot int) error {
	s.slot = slot
	createSaveGameFile()

	saveFileData, err := readFromLocalStorage(SaveGameFile)
	if err != nil {
		return err
	}

	bl := SaveGameByteLen
	s.fromBytes(saveFileData[slot*bl : slot*bl+bl])
	return nil
}

func createSaveGameFile() {
	existingData, _ := readFromLocalStorage(SaveGameFile)
	if len(existingData) > 0 {
		return
	}
	data := [SaveGameByteLen * 3]byte{}
	writeToLocalStorage(SaveGameFile, data[:])
}

func (s *n_SaveGame) toBytes() []byte {
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

	data = append(data, convBool(s.used))

	for len(data) < SaveGameByteLen {
		data = append(data, 0)
	}

	return data
}

func (s *n_SaveGame) fromBytes(data []byte) {
	for i := 0; i < len(s.allLevels); i++ {
		s.allLevels[i] = &n_levelData{}
		s.allLevels[i].loadLvl(data[i*11:])
	}

	s.jeepIndex = int(data[154])
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

	s.used = toBool(data[166])
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

func writeToLocalStorage(fileName string, data []byte) {
	stringData := b64.StdEncoding.EncodeToString(data)
	fmt.Println("writing to local storage", fileName, stringData)

	global := js.Global()
	localStorage := global.Get("localStorage")
	_ = localStorage.Call("setItem", fileName, stringData)
}

func readFromLocalStorage(fileName string) ([]byte, error) {
	fmt.Println("reading from local storage", fileName)
	global := js.Global()
	localStorage := global.Get("localStorage")
	fileContents := localStorage.Call("getItem", fileName)
	stringData := fileContents.String()
	return b64.StdEncoding.DecodeString(stringData)
}
