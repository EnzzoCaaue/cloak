package util

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"log"
	"sync"
)

// SpriteFile contains all the information about the .spr file
type SpriteFile struct {
	Signature   uint32
	Amount      uint32
	spriteIndex []uint32
	data        []byte
	rw          *sync.RWMutex
}

// ParseSpr parses tibia.spr file
func ParseSpr(path string) {
	//Spr.rw.Lock()
	//defer Spr.rw.Unlock()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	buffer := bytes.NewBuffer(data)
	reader := bufio.NewReader(buffer)
	if err := binary.Read(reader, binary.LittleEndian, &Spr.Signature); err != nil {
		log.Fatal(err)
	} else if err = binary.Read(reader, binary.LittleEndian, &Spr.Amount); err != nil {
		log.Fatal(err)
	}
	offset := (int64(Spr.Amount) * 4) - 8
	Spr.data = data[offset:]
	Spr.spriteIndex = make([]uint32, Spr.Amount+1)
	for i := uint32(1); i < Spr.Amount; i++ {
		if err := binary.Read(reader, binary.LittleEndian, &Spr.spriteIndex[i]); err != nil {
			log.Fatal(err)
		}
	}
}
