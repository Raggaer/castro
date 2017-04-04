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
	glua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"path/filepath"
	"strings"
	"github.com/raggaer/castro/app/models"
)

const (
	// configFileName the name of the application configuration file
	configFileName = "config.toml"

	// znoteTableName name of the znote table
	znoteTableName = "znote"
)

// znoteTable main znote installation table to look for
type znoteTable struct {
	Version   int
	Installed int64
}

// znoteAccount main znote accounts table
type znoteAccount struct {
	ID         uint64
	Account_id uint64
	Points     uint
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

	// Load config.lua file
	if err := lua.LoadConfig(
		filepath.Join(filepath.Join(location, "config.lua")),
	); err != nil {
		return err
	}

	// Connect to database
	db, err := database.Open(lua.Config.GetGlobal("mysqlUser").String(), lua.Config.GetGlobal("mysqlPass").String(), lua.Config.GetGlobal("mysqlDatabase").String())

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

	// Check if znote is installed
	z, err := isZnoteInstalled(db)

	if err != nil {
		return err
	}

	if z {

		// Znote accounts placeholder
		znoteAccounts := []znoteAccount{}

		// Get znote accounts
		if err := db.Select(&znoteAccounts, "SELECT id, account_id, points FROM znote_accounts ORDER BY id"); err != nil {
			return err
		}

		// Loop znote accounts
		for _, acc := range znoteAccounts {

			// Insert castro account from znote account
			if _, err := db.Exec("INSERT INTO castro_accounts (account_id, points) VALUES (?, ?)", acc.Account_id, acc.Points); err != nil {
				return err
			}
		}
	} else {

		// Normal accounts placeholder
		accountList := []models.Account{}

		// Get all accounts
		if err := db.Select(&accountList, "SELECT id FROM accounts ORDER BY id"); err != nil {
			return err
		}

		// Loop accounts
		for _, acc := range accountList {

			// Insert castro account from normal account
			if _, err := db.Exec("INSERT INTO castro_accounts (account_id) VALUES (?)", acc.ID); err != nil {
				return err
			}
		}
	}

	// Create configuration file
	return createConfigFile(configFileName, location)
}

// createConfigFile encodes a configuration file with the given name and location
func createConfigFile(name, location string) error {
	// Get lua state
	luaState := glua.NewState()

	// Close state
	defer luaState.Close()

	// Create events metatable
	lua.SetEventsMetaTable(luaState)

	// Create storage metatable
	lua.SetStorageMetaTable(luaState)

	// Create time metatable
	lua.SetTimeMetaTable(luaState)

	// Create url metatable
	lua.SetURLMetaTable(luaState)

	// Create debug metatable
	lua.SetDebugMetaTable(luaState)

	// Create XML metatable
	lua.SetXMLMetaTable(luaState)

	// Create captcha metatable
	lua.SetCaptchaMetaTable(luaState)

	// Create crypto metatable
	lua.SetCryptoMetaTable(luaState)

	// Create validator metatable
	lua.SetValidatorMetaTable(luaState)

	// Create database metatable
	lua.SetDatabaseMetaTable(luaState)

	// Create config metatable
	lua.SetConfigMetaTable(luaState)

	// Create map metatable
	lua.SetMapMetaTable(luaState)

	// Create mail metatable
	lua.SetMailMetaTable(luaState)

	// Create cache metatable
	lua.SetCacheMetaTable(luaState)

	// Create reflect metatable
	lua.SetReflectMetaTable(luaState)

	// Create json metatable
	lua.SetJSONMetaTable(luaState)

	// Set config metatable
	lua.SetConfigGlobal(luaState)

	// Execute init file
	if err := luaState.DoFile(filepath.Join("engine", "install.lua")); err != nil {
		return err
	}

	// Create configuration file handle
	configFile, err := os.Create(name)
	if err != nil {
		return err
	}

	// Close file handle
	defer configFile.Close()

	// Get lua file table
	tbl := luaState.ToTable(-1)

	// Encode the given configuration struct into the file
	return toml.NewEncoder(configFile).Encode(util.Configuration{
		Version:      VERSION,
		CheckUpdates: true,
		Mode:         "dev",
		Port:         8080,
		URL:          "localhost",
		Datapack:     location,
		Cookies: util.CookieConfig{
			Name:     fmt.Sprintf("castro-%v", uniuri.NewLen(5)),
			MaxAge:   1000000,
			HashKey:  uniuri.NewLen(32),
			BlockKey: uniuri.NewLen(32),
		},
		Cache: util.CacheConfig{
			Default: time.Minute * 5,
			Purge:   time.Minute,
		},
		RateLimit: util.RateLimiterConfig{
			Number: 100,
			Time:   time.Minute,
		},
		Custom: lua.TableToMap(tbl),
	})
}
