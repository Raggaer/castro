package util

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
	"github.com/raggaer/otmap"
	"io/ioutil"
)

var (
	// OTBMap holds the main server map parsed using otmap library
	OTBMap *CastroMap

	// ServerHouseList holds the main server house list XML list
	ServerHouseList = ServerHouses{
		List: &HouseList{},
	}
)

// CastroMap struct used to decode and encode tibia maps
type CastroMap struct {
	Towns     []otmap.Town
	HouseFile string
}

// House holds all information about a game house
type House struct {
	ID     uint32 `xml:"houseid,attr"`
	Name   string `xml:"name,attr"`
	EntryX uint16 `xml:"entryx,attr"`
	EntryY uint16 `xml:"entryy,attr"`
	EntryZ uint16 `xml:"entryz,attr"`
	Size   int    `xml:"size,attr"`
	TownID uint32 `xml:"townid,attr"`
}

// HouseList holds the house array
type HouseList struct {
	XMLName xml.Name `xml:"houses"`
	Houses  []*House `xml:"house"`
}

// ServerHouses contains the whole house list of the server
type ServerHouses struct {
	List *HouseList
}

// LoadHouses parses the server map houses
func LoadHouses(file string, list ServerHouses) error {
	// Load houses file
	f, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	// Unmarshal houses file
	return xml.Unmarshal(f, &list.List)
}

// EncodeMap encodes the server map to the given destination
func EncodeMap(path string) ([]byte, error) {
	// Parse server map
	m, err := otmap.Parse(path, true)

	if err != nil {
		return nil, err
	}

	// Create castro map holder
	c := CastroMap{
		Towns:     m.Towns,
		HouseFile: m.HouseFile,
	}

	// Map buffer
	buff := bytes.NewBuffer([]byte{})

	// Create map encoder
	encoder := gob.NewEncoder(buff)

	// Encode map
	if err := encoder.Encode(&c); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// DecodeMap decodes the server map to the given destination
func DecodeMap(mapData []byte) (*CastroMap, error) {
	// Map buffer
	buff := bytes.NewBuffer(mapData)

	// Create decoder
	decoder := gob.NewDecoder(buff)

	// Map holder
	m := CastroMap{}

	// Decode map
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
