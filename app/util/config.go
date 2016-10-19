package util

import "github.com/BurntSushi/toml"

// Database struct used for the config file
// database credentials
type Database struct {
	Username string
	Password string
	Name     string
}

// Configuration struct used for the main Castro
// config file TOML file
type Configuration struct {
	Port     int
	Datapack string
	Database Database
}

// Config holds the main configuration file
var Config *Configuration

func init() {
	Config = &Configuration{}
}

// LoadConfig loads the configuration file to
// the given interface pointer
func LoadConfig(data string, dest interface{}) error {
	if _, err := toml.Decode(data, dest); err != nil {
		return err
	}
	return nil
}
