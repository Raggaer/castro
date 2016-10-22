package util

import (
	"time"

	"github.com/BurntSushi/toml"
)

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

// cache struct used for the cache
// configuration options
type cache struct {
	Default duration
	Purge   duration
}

// duration struct used for time.Duration conversion
type duration struct {
	time.Duration
}

// Configuration struct used for the main Castro
// config file TOML file
type Configuration struct {
	Mode     string
	Port     int
	Datapack string
	Database database
	Cookies  cookie
	Cache    cache
	Custom   map[string]interface{}
}

// Config holds the main configuration file
var Config *Configuration

func init() {
	Config = &Configuration{}
}

// UnmarshalText converts byte array to time.Duration
func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
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
