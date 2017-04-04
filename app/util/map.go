package util

import (
	"encoding/gob"
	"encoding/xml"
	"github.com/raggaer/otmap"
	"io/ioutil"
	"os"
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
func EncodeMap(path, dest string) error {
	// Parse server map
	m, err := otmap.Parse(path, true)

	if err != nil {
		return err
	}

	// Create castro map holder
	c := CastroMap{
		Towns:     m.Towns,
		HouseFile: m.HouseFile,
	}

	// Create map file
	f, err := os.Create(dest)

	if err != nil {
		return err
	}

	// Close map file
	defer f.Close()

	// Create map encoder
	encoder := gob.NewEncoder(f)

	// Encode map
	return encoder.Encode(&c)
}

// DecodeMap decodes the server map to the given destination
func DecodeMap(path string) (*CastroMap, error) {
	// Open map file
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	// Close map file
	defer f.Close()

	// Create decoder
	decoder := gob.NewDecoder(f)

	// Map holder
	m := CastroMap{}

	// Decode map
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
