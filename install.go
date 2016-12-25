package main

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/raggaer/castro/app/util"
	"github.com/dchest/uniuri"
)

func isInstalled() bool {
	_, err := os.Stat("config.toml")
	return err == nil
}

func createConfigFile(name string) error {
	// Create configuration file handle
	configFile, err := os.Create("config.toml")
	if err != nil {
		return err
	}
	defer configFile.Close()

	// Encode the given configuration struct into the file
	return toml.NewEncoder(configFile).Encode(util.Configuration{
		Mode:     "dev",
		Port:     8080,
		URL:      "localhost",
		Datapack: "/",
		Secret: uniuri.New(),
		Cookies: util.CookieConfig{
			Name:   "castro",
			MaxAge: 1000000,
		},
		Cache: util.CacheConfig{
			Default: int(time.Minute) * 5,
			Purge:   int(time.Minute),
		},
		SSL: util.SSLConfig{
			Enabled: false,
		},
		Custom: make(map[string]interface{}),
	})
}
