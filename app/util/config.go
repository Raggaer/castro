package util

import "github.com/BurntSushi/toml"

// Config struct used for the main Castro
// config file TOML file
type Config struct {
	Port int
}

// LoadConfig loads the configuration file to
// the given interface
func LoadConfig(data string, dest interface{}) error {
	if _, err := toml.Decode(data, dest); err != nil {
		return err
	}
	return nil
}
