package main

import (
	"os"
	"time"

	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/dchest/uniuri"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	// configFileName the name of the application configuration file
	configFileName = "config.toml"

	// znoteTableName name of the znote table
	znoteTableName = "znote"
)

type znoteTable struct {
	Version   int
	Installed int64
}

// isInstalled check if application is installed
func isInstalled() bool {
	// Check if file exists
	_, err := os.Stat(configFileName)

	return err == nil
}

// isZnoteInstalled checks if znote_aac is already installed
func isZnoteInstalled(db *sqlx.DB) (bool, error) {
	// Check if table exists
	if _, err := db.Exec("DESCRIBE " + znoteTableName); err != nil {

		// Convert error to MySQL error type
		mErr, ok := err.(*mysql.MySQLError)

		// Check if table is installed
		if ok && mErr.Number == 1146 {

			return false, nil
		}

		return false, err
	}

	return true, nil
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

	// Get all sql files
	tables, err := ioutil.ReadDir(filepath.Join("install"))

	if err != nil {
		return err
	}

	// Loop files
	for _, table := range tables {

		// Check if table exists
		if _, err := db.Exec("DESCRIBE " + strings.TrimSuffix(table.Name(), ".sql")); err != nil {

			// Convert error to MySQL error type
			mErr, ok := err.(*mysql.MySQLError)

			// Check if table is installed
			if ok && mErr.Number == 1146 {

				// Read file
				buff, err := ioutil.ReadFile(filepath.Join("install", table.Name()))

				if err != nil {
					return err
				}

				// Execute query
				if _, err := db.Exec(string(buff)); err != nil {
					return err
				}

				continue
			}

			return err
		}
	}

	fmt.Println("Missing tables created")

	// Create configuration file
	return createConfigFile(configFileName, location)
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
