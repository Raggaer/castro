package util

import (
	"encoding/xml"
	"github.com/raggaer/otmap"
	"io/ioutil"
)

var (
	// OTBMap holds the main server map parsed
	// using otmap library
	OTBMap *otmap.Map

	// HouseList holds the main server house list
	// XML list
	ServerHouseList = ServerHouses{
		List: &HouseList{},
	}
)

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
