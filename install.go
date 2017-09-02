package main

import (
	"html/template"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"

	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
)

const (
	// configFileName the name of the application configuration file
	configFileName = "config.toml"

	// znoteTableName name of the znote table
	znoteTableName = "znote"
)

var (
	installationTemplate = template.New("install")

	installationConfigFile = &util.Configuration{
		CheckUpdates: true,
		Template:     "views/default",
		Mode:         "dev",
		Port:         80,
		URL:          "localhost",
		Datapack:     "",
		Cookies: util.CookieConfig{
			Name:     fmt.Sprintf("castro-%v", uniuri.NewLen(5)),
			MaxAge:   1000000,
			HashKey:  uniuri.NewLen(32),
			BlockKey: uniuri.NewLen(32),
		},
		Cache: util.CacheConfig{
			Default: util.NewStringDuration("5m"),
			Purge:   util.NewStringDuration("1m"),
		},
		RateLimit: util.RateLimiterConfig{
			Number: 100,
			Time:   util.NewStringDuration("1m"),
		},
		Security: util.SecurityConfig{
			NonceEnabled:      true,
			STS:               "max-age=10000",
			XSS:               "1; mode=block",
			Frame:             "DENY",
			ContentType:       "nosniff",
			ReferrerPolicy:    "origin",
			CrossDomainPolicy: "none",
			CSP: util.ContentSecurityPolicyConfig{
				Default: []string{"none"},
				Frame: util.ContentSecurityPolicyType{
					SRC: []string{"http://pay.fortumo.com", "https://www.google.com"},
				},
				Script: util.ContentSecurityPolicyType{
					Default: []string{"self"},
					SRC:     []string{"https://ajax.googleapis.com", "https://assets.fortumo.com", "https://www.google.com", "https://code.jquery.com", "https://cdn.datatables.net", "https://www.gstatic.com"},
				},
				Font: util.ContentSecurityPolicyType{
					Default: []string{"self"},
					SRC:     []string{"http://fonts.gstatic.com", "http://fonts.googleapis.com"},
				},
				Connect: util.ContentSecurityPolicyType{
					Default: []string{"self"},
				},
				Style: util.ContentSecurityPolicyType{
					Default: []string{"unsafe-inline", "self"},
					SRC:     []string{"https://assets.fortumo.com", "http://fonts.googleapis.com", "https://cdn.datatables.net"},
				},
				Image: util.ContentSecurityPolicyType{
					Default: []string{"self"},
					SRC:     []string{"https://assets.fortumo.com", "https://*.githubusercontent.com", "data:"},
				},
			},
		},
	}
)

// znoteTable main znote installation table to look for
type znoteTable struct {
	Version   int
	Installed int64
}

// znoteAccount main znote accounts table
type znoteAccount struct {
	ID         uint64
	Account_id int64
	Points     uint
}

type installError struct {
	Error string
}

// isInstalled check if application is installed
func isInstalled() bool {
	// Check if file exists
	_, err := os.Stat(configFileName)

	return err == nil
}

// isZnoteInstalled checks if znote_aac is already installed
func isZnoteInstalled(db *sqlx.Tx) (bool, error) {
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

func accountExists(id int64, db *sqlx.Tx) bool {
	// Data holder
	exists := false

	// Check if account exists
	if err := db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM castro_accounts WHERE account_id = ?)", id); err != nil {
		return false
	}

	return exists
}

// startInstallerApplication starts the installer server and the installer template
func startInstallerApplication() error {
	// Register html files
	if err := registerInstallationTemplate(); err != nil {
		return err
	}

	// Create installer router
	router := httprouter.New()

	// Declare installation routes
	router.GET("/", showInstallationPage)
	router.POST("/install/path", installationServerPath)
	router.GET("/install/captcha", showInstallationServerCaptcha)
	router.POST("/install/captcha", installationServerCaptcha)
	router.GET("/install/encode", showInstallationServerEncode)
	router.POST("/install/encode", installationServerEncode)
	router.GET("/install/finish", showInstallationServerFinish)

	// Create installer listener
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		return err
	}

	// Create negroni middleware
	n := negroni.New(
		negroni.NewStatic(http.Dir("public/")),
	)

	// Use httprouter router
	n.UseHandler(router)

	fmt.Println("Castro is not installed. Installer will listen on " + listener.Addr().String())

	// Start installer http server
	return http.Serve(listener, n)
}

// registerInstallationTemplate register the html template files
func registerInstallationTemplate() error {
	// Parse installation html files
	_, err := installationTemplate.ParseGlob(filepath.Join("views", "*.html"))

	return err
}

func showInstallationServerFinish(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Execute finish template
	installationTemplate.ExecuteTemplate(res, "install_finish.html", nil)
}

// showInstallationPage handles the installation GET request
func showInstallationPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Execute install template
	installationTemplate.ExecuteTemplate(res, "install_path.html", nil)
}

func installationServerEncode(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Encode config struct to file
	if err := createConfigFile("config.toml", installationConfigFile); err != nil {
		installationTemplate.ExecuteTemplate(res, "install_encode.html", installError{
			Error: err.Error(),
		})
		return
	}

	// Redirect to finish page
	http.Redirect(res, req, "/install/finish", 302)
}

func showInstallationServerEncode(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Execute encode done template
	installationTemplate.ExecuteTemplate(res, "install_encode.html", nil)
}

func showInstallationServerCaptcha(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Execute install captcha template
	installationTemplate.ExecuteTemplate(res, "install_captcha.html", nil)
}

func installationServerCaptcha(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Parse form
	if err := req.ParseForm(); err != nil {
		installationTemplate.ExecuteTemplate(res, "install_captcha.html", installError{
			Error: err.Error(),
		})
		return
	}

	// Set configuration captcha options
	installationConfigFile.Captcha = util.CaptchaConfig{
		Enabled: true,
		Public:  req.FormValue("public"),
		Secret:  req.FormValue("public"),
	}

	// Redirect user to next step
	http.Redirect(res, req, "/install/encode", 302)
}

func installationServerPath(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Parse form
	if err := req.ParseForm(); err != nil {
		installationTemplate.ExecuteTemplate(res, "install_path.html", installError{
			Error: err.Error(),
		})
		return
	}

	// Get server path
	location := req.FormValue("path")

	// Install needed database tables
	if err := installApplication(location); err != nil {
		installationTemplate.ExecuteTemplate(res, "install_path.html", installError{
			Error: err.Error(),
		})
		return
	}

	// Get server port
	port, err := strconv.Atoi(req.FormValue("port"))

	if err != nil {
		installationTemplate.ExecuteTemplate(res, "install_path.html", installError{
			Error: err.Error(),
		})
		return
	}

	// Update config values
	installationConfigFile.Datapack = location
	installationConfigFile.Port = port

	// Redirect to next step
	http.Redirect(res, req, "/install/captcha", 302)
}

// installApplication runs the installation process
func installApplication(location string) error {
	// Load config.lua file
	if err := lua.LoadConfig(
		filepath.Join(filepath.Join(location, "config.lua")),
	); err != nil {
		return err
	}

	// Connect to database
	conn, err := database.Open(
		lua.Config.GetGlobal("mysqlUser").String(),
		lua.Config.GetGlobal("mysqlPass").String(),
		lua.Config.GetGlobal("mysqlDatabase").String(),
	)

	if err != nil {
		return err
	}

	// Set global handler for lua states
	database.DB = conn

	// Ping database
	if err := database.DB.Ping(); err != nil {
		return err
	}

	// Close database handle
	defer database.DB.Close()

	// Begin transaction
	db, err := database.DB.Beginx()

	if err != nil {
		db.Rollback()
		return err
	}

	// Get all sql files
	tables, err := ioutil.ReadDir(filepath.Join("install"))

	if err != nil {
		db.Rollback()
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
					db.Rollback()
					return err
				}

				// Execute query
				if _, err := db.Exec(string(buff)); err != nil {
					db.Rollback()
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
		db.Rollback()
		return err
	}

	if z {
		// Znote accounts placeholder
		znoteAccounts := []znoteAccount{}

		// Get znote accounts
		if err := db.Select(&znoteAccounts, "SELECT id, account_id, points FROM znote_accounts ORDER BY id"); err != nil {
			db.Rollback()
			return err
		}

		// Loop znote accounts
		for _, acc := range znoteAccounts {

			// Check if account exists
			if accountExists(acc.Account_id, db) {
				continue
			}

			// Insert castro account from znote account
			if _, err := db.Exec("INSERT INTO castro_accounts (account_id, points) VALUES (?, ?)", acc.Account_id, acc.Points); err != nil {
				db.Rollback()
				return err
			}
		}
	}

	// Normal accounts placeholder
	accountList := []models.Account{}

	// Get all accounts
	if err := db.Select(&accountList, "SELECT id FROM accounts ORDER BY id"); err != nil {
		db.Rollback()
		return err
	}

	// Loop accounts
	for _, acc := range accountList {

		// Check if account exists
		if accountExists(acc.ID, db) {
			continue
		}

		// Insert castro account from normal account
		if _, err := db.Exec("INSERT INTO castro_accounts (account_id) VALUES (?)", acc.ID); err != nil {
			db.Rollback()
			return err
		}
	}

	fmt.Println(">> Encoding map file. This process can take several minutes")

	// Encode server map
	mapData, err := util.EncodeMap(
		filepath.Join(location, "data", "world", lua.Config.GetGlobal("mapName").String()+".otbm"),
	)

	if err != nil {
		db.Rollback()
		return err
	}

	// Save encoded map
	if _, err := db.Exec("INSERT INTO castro_map (name, data, created_at, updated_at) VALUES (?, ?, ?, ?)", lua.Config.GetGlobal("mapName").String(), mapData, time.Now(), time.Now()); err != nil {
		db.Rollback()
		return err
	}

	// Commit changes
	if err := db.Commit(); err != nil {
		db.Rollback()
		return err
	}

	return nil
}

// createConfigFile encodes a configuration file with the given name and location
func createConfigFile(name string, cfg *util.Configuration) error {
	// Get lua state
	luaState := glua.NewState()

	// Close state
	defer luaState.Close()

	// Get application state ready
	lua.GetApplicationState(luaState)

	// Execute init file
	if err := lua.ExecuteFile(luaState, filepath.Join("engine", "install.lua")); err != nil {
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
	cfg.Custom = lua.TableToMap(luaState.ToTable(-1))

	// Encode the given configuration struct into the file
	return toml.NewEncoder(configFile).Encode(cfg)
}
