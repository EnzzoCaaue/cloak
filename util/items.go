package util

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
)

const (
	itemWeightLabel      = "weight"
	itemDescriptionLabel = "description"
	itemSlotLabel        = "slotType"
    itemArmorLabel       = "armor"
)

// ServerItems holds all server items
type ServerItems struct {
	Items map[int]ItemDefinition
	rw    *sync.RWMutex
}

// ItemDefinition represents a server item
type ItemDefinition struct {
	ID         int             `xml:"id,attr"`
	Name       string          `xml:"name,attr"`
	Attributes []itemAttribute `xml:"attribute"`
}

type itemAttribute struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

type itemsDefinition struct {
	XMLName xml.Name         `xml:"items"`
	Items   []ItemDefinition `xml:"item"`
}

// ParseItems loads items xml into memory
func ParseItems(path string) {
	Items.rw.Lock()
	defer Items.rw.Unlock()
	Items.Items = make(map[int]ItemDefinition)
	b, err := ioutil.ReadFile(path + "/data/items/items.xml")
	if err != nil {
		log.Fatal(err)
	}
	itemList := &itemsDefinition{}
	buffer := bytes.NewBuffer(b)
	decoder := xml.NewDecoder(buffer)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(itemList)
	if err != nil {
		log.Fatal(err)
	}
	for _, k := range itemList.Items {
		Items.Items[k.ID] = k
	}
}

// Get gets an item by its ID
func (s *ServerItems) Get(id int) *ItemDefinition {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if v, ok := s.Items[id]; ok {
		return &v
	}
	return nil
}

// GetWeight returns an item weight
func (s *ItemDefinition) GetWeight() int {
	for _, v := range s.Attributes {
		if v.Key == itemWeightLabel {
			w, err := strconv.Atoi(v.Value)
			if err != nil {
				return 0
			}
			return w
		}
	}
	return 0
}

// GetDescription returns an item description
func (s *ItemDefinition) GetDescription() string {
	for _, v := range s.Attributes {
		if v.Key == itemDescriptionLabel {
			return v.Value
		}
	}
	return ""
}

// GetSlot returns an item slot type
func (s *ItemDefinition) GetSlot() string {
	for _, v := range s.Attributes {
		if v.Key == itemSlotLabel {
			return v.Value
		}
	}
	return ""
}

// GetArmor returns an item armor value
func (s *ItemDefinition) GetArmor() int {
    for _, v := range s.Attributes {
        if v.Key == itemArmorLabel {
            w, err := strconv.Atoi(v.Value)
			if err != nil {
				return 0
			}
			return w
        }
    }
    return 0
}