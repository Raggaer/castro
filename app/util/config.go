package util

import "github.com/BurntSushi/toml"

// database struct used for the config file
// database credentials
type database struct {
	Username string
	Password string
	Name     string
}

// cookie struct used for the cookies
// configuration options
type cookie struct {
	Name   string
	MaxAge int
}

// Configuration struct used for the main Castro
// config file TOML file
type Configuration struct {
	Mode     string
	Port     int
	Datapack string
	Database database
	Cookies  cookie
	Custom   map[string]interface{}
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

// IsDev checks if castro is running on
// development mode
func (c Configuration) IsDev() bool {
	return c.Mode == "dev"
}
