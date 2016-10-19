package util

import "github.com/BurntSushi/toml"

// Database struct used for the config file
// database credentials
type Database struct {
	Username string
	Password string
	Name     string
}

// Config struct used for the main Castro
// config file TOML file
type Config struct {
	Port     int
	Database Database
}

// LoadConfig loads the configuration file to
// the given interface pointer
func LoadConfig(data string, dest interface{}) error {
	if _, err := toml.Decode(data, dest); err != nil {
		return err
	}
	return nil
}
