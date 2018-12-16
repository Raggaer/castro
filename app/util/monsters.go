package util

import (
	"encoding/xml"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
)

var MonstersList = []*Monster{}

// MonsterList defines the monsters.xml file
type MonsterList struct {
	XMLName  xml.Name             `xml:"monsters"`
	Monsters []monsterListElement `xml:"monster"`
}

type monsterListElement struct {
	XMLName  xml.Name `xml:"monster"`
	Name     string   `xml:"name,attr"`
	File     string   `xml:"file,attr"`
	Disabled bool     `xml:"disablewebsite,attr"`
}

// Monster defines a server cerature
type Monster struct {
	XMLName      xml.Name            `xml:"monster"`
	Name         string              `xml:"name,attr"`
	Description  string              `xml:"nameDescription,attr"`
	Race         string              `xml:"race,attr"`
	Experience   int                 `xml:"experience,attr"`
	Speed        int                 `xml:"speed,attr"`
	Health       MonsterHealth       `xml:"health"`
	Look         MonsterLook         `xml:"look"`
	TargetChange MonsterTargetChange `xml:"targetchange"`
	Attacks      MonsterAttackList   `xml:"attacks"`
	Defenses     MonsterDefenseList  `xml:"defenses"`
	Voices       MonsterVoiceList    `xml:"voices"`
	Loot         MonsterLootList     `xml:"loot"`
	Elements     MonsterElements     `xml:"elements>element"`
}

type MonsterElements struct {
	Ice      int `xml:"icePercent,attr"`
	Earth    int `xml:"earthPercent,attr"`
	Energy   int `xml:"energyPercent,attr"`
	Fire     int `xml:"firePercent,attr"`
	Holy     int `xml:"holyPercent,attr"`
	Physical int `xml:"physicalPercent,attr"`
	Death    int `xml:"deathPercent,attr"`
}

// MonsterHealth defines the monster health values
type MonsterHealth struct {
	XMLName xml.Name `xml:"health"`
	Now     int      `xml:"now,attr"`
	Max     int      `xml:"max,attr"`
}

// MonsterLook defines the monster looktype values
type MonsterLook struct {
	XMLName xml.Name `xml:"look"`
	Type    int      `xml:"type,attr"`
	Corpse  int      `xml:"corpse,attr"`
}

// MonsterTargetChange defines the monster targetting change values
type MonsterTargetChange struct {
	XMLName  xml.Name `xml:"targetchange"`
	Interval int      `xml:"interval,attr"`
	Chance   int      `xml:"chance,attr"`
}

// MonsterAttackList defines a list of monster attacks
type MonsterAttackList struct {
	Attacks []MonsterAttack `xml:"attack"`
}

// MonsterAttack defines a monster attack
type MonsterAttack struct {
	XMLName    xml.Name           `xml:"attack"`
	Name       string             `xml:"name,attr"`
	Interval   int                `xml:"interval,attr"`
	Range      int                `xml:"range,attr"`
	Min        int                `xml:"min,attr"`
	Max        int                `xml:"max,attr"`
	Target     int                `xml:"target,attr"`
	Attributes []MonsterAttribute `xml:"attribute"`
}

// MonsterAttribute defines a monster attribute
type MonsterAttribute struct {
	XMLName xml.Name `xml:"attribute"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

// MonsterDefenseList defines a monster defense list
type MonsterDefenseList struct {
	Armor    int              `xml:"armor,attr"`
	Defense  int              `xml:"defense,attr"`
	Defenses []MonsterDefense `xml:"defense"`
}

// MonsterDefense defines a monster defense value
type MonsterDefense struct {
	Name       string             `xml:"name,attr"`
	Interval   int                `xml:"interval,attr"`
	Chance     int                `xml:"chance,attr"`
	Min        int                `xml:"min,attr"`
	Max        int                `xml:"max,attr"`
	Attributes []MonsterAttribute `xml:"attribute"`
}

// MonsterVoiceList defines a list of monster voices
type MonsterVoiceList struct {
	Interval int            `xml:"interval,attr"`
	Chance   int            `xml:"chance,attr"`
	Voices   []MonsterVoice `xml:"voice"`
}

// MonsterVoice defines a monster sentence
type MonsterVoice struct {
	XMLName  xml.Name `xml:"voice"`
	Sentence string   `xml:"sentence,attr"`
}

// MonsterLootList defines a list of monster lootable items
type MonsterLootList struct {
	Loot []MonsterItem `xml:"item"`
}

// MonsterItem defines a monster lootable item
type MonsterItem struct {
	XMLName  xml.Name `xml:"item"`
	ID       int      `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	CountMax int      `xml:"countmax,attr"`
	Chance   int      `xml:"chance,attr"`
}

// LoadMonsterList loads the monsters.xml file
func LoadMonsterList(path string) (*MonsterList, error) {
	// Open monsters.xml file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create xml decoder
	list := MonsterList{}
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&list); err != nil {
		return nil, err
	}

	return &list, nil
}

// LoadMonster loads the given monster xml file
func LoadMonster(path string) (*Monster, error) {
	// Open monster .xml file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create xml decoder
	monster := Monster{}
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&monster); err != nil {
		return nil, err
	}

	return &monster, nil
}

// LoadServerMonsters loads the server monsters and sets the variable
func LoadServerMonsters(path string) error {
	// Load monsters.xml first
	list, err := LoadMonsterList(filepath.Join(path, "data", "monster", "monsters.xml"))
	if err != nil {
		return err
	}

	// Start loading each monster
	for _, m := range list.Monsters {
		if m.Disabled {
			continue
		}
		mst, err := LoadMonster(filepath.Join(path, "data", "monster", m.File))
		if err != nil {
			return err
		}
		MonstersList = append(MonstersList, mst)
	}
	return nil
}
