package util

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"
	"sync"
)

// SpriteFile contains all the information about the .spr file
type SpriteFile struct {
	Signature   uint32
	Amount      uint32
	spriteIndex []uint32
	rw          *sync.RWMutex
}

// ParseSpr parses tibia.spr file
func ParseSpr(path string) {
	Spr.rw.Lock()
	defer Spr.rw.Unlock()
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	if err := binary.Read(reader, binary.LittleEndian, &Spr.Signature); err != nil {
		log.Fatal(err)
	} else if err = binary.Read(reader, binary.LittleEndian, &Spr.Amount); err != nil {
		log.Fatal(err)
	}
	Spr.spriteIndex = make([]uint32, Spr.Amount)
	for i := uint32(0); i < Spr.Amount; i++ {
		if err := binary.Read(reader, binary.LittleEndian, &Spr.spriteIndex[i]); err != nil {
			log.Fatal(err)
		}
	}
}
