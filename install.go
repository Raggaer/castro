package main

import (
	"os"
	"time"

	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/dchest/uniuri"
	"github.com/go-sql-driver/mysql"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
	"path/filepath"
)

const (
	// configFileName the name of the application configuration file
	configFileName = "config.toml"
)

var (
	// tables list of application required tables
	tables = []installationTable{
		{
			name: "castro_accounts",
			source: `CREATE TABLE castro_accounts (
				id INT(11) NOT NULL AUTO_INCREMENT,
				account_id INT(11) NOT NULL,
				points INT(11) DEFAULT 0,
				admin TINYINT(1) DEFAULT 0,
				PRIMARY KEY (id)
			) ENGINE=InnoDB`,
		},
	}
)

// installationTable struct used for application tables
type installationTable struct {
	name   string
	source string
}

// isInstalled check if application is installed
func isInstalled() bool {
	// Check if file exists
	_, err := os.Stat(configFileName)

	return err == nil
}

// installApplication runs the installation process
func installApplication() error {

	// Ask user for server directory location
	fmt.Print("Insert your server location: ")

	// Location holder
	var location string

	// Get user response
	_, err := fmt.Scanln(&location)

	if err != nil {
		return err
	}

	// Configuration holder
	cfg := lua.Configuration{}

	// Load config.lua file
	if err := lua.LoadConfig(
		filepath.Join(location),
		&cfg,
	); err != nil {
		return err
	}

	// Connect to database
	db, err := database.Open(cfg.MySQLUser, cfg.MySQLPass, cfg.MySQLDatabase)

	if err != nil {
		return err
	}

	// Ping database
	if err := db.Ping(); err != nil {
		return err
	}

	// Close database handle
	defer db.Close()

	// Loop installation tables
	for _, table := range tables {

		// Check if table exists
		if _, err := db.Exec("DESCRIBE " + table.name); err != nil {

			// Convert error to MySQL error type
			mErr, ok := err.(*mysql.MySQLError)

			// Check if table is installed
			if ok && mErr.Number == 1146 {

				// Create missing table
				if _, err := db.Exec(table.source); err != nil {
					return err
				}

				continue
			}

			return err
		}

	}

	return nil
}

// createConfigFile encodes a configuration file with the given name and location
func createConfigFile(name, location string) error {
	// Create configuration file handle
	configFile, err := os.Create(name)
	if err != nil {
		return err
	}

	// Close file handle
	defer configFile.Close()

	// Encode the given configuration struct into the file
	return toml.NewEncoder(configFile).Encode(util.Configuration{
		Mode:     "dev",
		Port:     8080,
		URL:      "localhost",
		Datapack: location,
		Secret:   uniuri.NewLen(35),
		Captcha: util.CaptchaConfig{
			Enabled: false,
		},
		Cookies: util.CookieConfig{
			Name:   "castro",
			MaxAge: 1000000,
		},
		Cache: util.CacheConfig{
			Default: time.Minute * 5,
			Purge:   time.Minute,
		},
		SSL: util.SSLConfig{
			Enabled: false,
		},
		RateLimit: util.RateLimiterConfig{
			Number: 100,
			Time:   time.Minute,
		},
		Custom: make(map[string]interface{}),
	})
}
