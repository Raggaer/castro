package main

import (
	"errors"
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
	installationSteps = []installationStep{
		{
			Name: "Path",
			URL:  "/",
			Description: template.HTML(`
			<p>Welcome to the installation wizard. The wizard will guide you through a very simple process where you will be able to setup some Castro features like the SMTP server and the Google reCAPTCHA service</p>
			<p>These features are optional and you can activate them later</p>
			<p>Just fill the information below to start using one of the most powerful and extensible Open Tibia content management system</p>
			`),
			Optional: false,
			Form: []installationFormField{
				{
					Name:        "path",
					Type:        "text",
					Placeholder: "Server full path",
					HelperText:  "We need your full server folder path. Castro will load all the necessary information from your config.lua file",
				},
				{
					Name:        "port",
					Type:        "number",
					Placeholder: "Website port",
					HelperText:  "Port where Castro will listen. You should use port 80 unless you run Castro behind a proxy server like NGINX or Caddy",
				},
				{
					Name:        "url",
					Type:        "text",
					Placeholder: "Website url",
					HelperText:  "Your absolute website URL. This URL will be used to create links. URL/destination",
				},
			},
			Post: func(res http.ResponseWriter, req *http.Request, s installationStep) error {
				// Get server port
				port, err := strconv.Atoi(req.FormValue("port"))

				if err != nil {
					return errors.New("Invalid port number")
				}

				// Install database tables
				if err := installApplication(req.FormValue("path")); err != nil {
					return err
				}

				// Set application settings
				installationConfigFile.Datapack = req.FormValue("path")
				installationConfigFile.Port = port
				installationConfigFile.URL = req.FormValue("url")

				return nil
			},
		},
		{
			Name:     "Captcha",
			URL:      "/install/captcha",
			Optional: true,
			Description: template.HTML(`
			<p>You can configure your Google reCAPTCHA credentials. reCAPTCHA offers an easy way to stop bots at saturaing your database</p>
			<p>By default if captcha is enabled it will appear on the registration form, you can use lua bindings to add captcha security to any other form of the website</p>
			<p>To setup your captcha service head to <a href="https://www.google.com/recaptcha/admin#list">Google reCAPTCHA</a> and create a new application, make sure to select <b>reCAPTCHA v2</b> as your application type. You can also learn how to integreate captchas on Castro by heading to the <a href="https://docs.castroaac.org/docs/lua/captcha">documentation page</a></p>
			`),
			Form: []installationFormField{
				{
					Name:        "public",
					Type:        "text",
					Placeholder: "Captcha public key",
					HelperText:  "Google reCAPTCHA public key",
				},
				{
					Name:        "private",
					Type:        "text",
					Placeholder: "Captcha private key",
					HelperText:  "Google reCAPTCHA private key",
				},
			},
			Post: func(res http.ResponseWriter, req *http.Request, s installationStep) error {
				// Update fields
				installationConfigFile.Captcha = util.CaptchaConfig{
					Public:  req.FormValue("public"),
					Secret:  req.FormValue("private"),
					Enabled: true,
				}

				return nil
			},
		},
		{
			Name:     "Mail",
			URL:      "/install/mail",
			Optional: true,
			Description: template.HTML(`
			<p>You can configure an SMTP server to send emails within Castro. Please fill the form below</p>
			<p>If you want to read more about sending emails using Castro lua bindings head to the <a href="https://docs.castroaac.org/docs/lua/mail">documentation page</a></p>
			`),
			Form: []installationFormField{
				{
					Name:        "server",
					Type:        "text",
					Placeholder: "SMTP server address",
					HelperText:  "Address of your SMTP server",
				},
				{
					Name:        "port",
					Type:        "number",
					Placeholder: "SMTP server port",
					HelperText:  "Port where your SMTP server listens on",
				},
				{
					Name:        "username",
					Type:        "text",
					Placeholder: "SMTP server username",
					HelperText:  "Login username for your SMTP server",
				},
				{
					Name:        "password",
					Type:        "password",
					Placeholder: "SMTP server password",
					HelperText:  "Login password for your SMTP server",
				},
			},
			Post: func(res http.ResponseWriter, req *http.Request, s installationStep) error {
				// Get server port
				port, err := strconv.Atoi(req.FormValue("port"))

				if err != nil {
					return err
				}

				// Update fields
				installationConfigFile.Mail = util.MailConfig{
					Server:   req.FormValue("server"),
					Port:     port,
					Username: req.FormValue("username"),
					Password: req.FormValue("password"),
				}

				return nil
			},
		},
	}

	// Installation template holder
	installationTemplate = template.New("install")

	// Installation config file holder
	installationConfigFile = &util.Configuration{
		CheckUpdates: true,
		Template:     "views/default",
		Mode:         "dev",
		Port:         80,
		URL:          "localhost",
		Datapack:     "",
		Plugin: util.PluginConfig{
			Enabled: true,
			Origin:  "https://plugins.castroaac.org",
		},
		MapWatch: util.MapWatchConfig{
			Enabled: true,
			Check:   util.NewStringDuration("1h"),
		},
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
			Number:  100,
			Enabled: false,
			Time:    util.NewStringDuration("1m"),
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

type installationStep struct {
	URL         string
	Optional    bool
	Description template.HTML
	Name        string
	Form        []installationFormField
	Post        installationFormHandle
}

type installationFormHandle func(http.ResponseWriter, *http.Request, installationStep) error

type installationTemplateData struct {
	Step    installationStep
	Success string
	Error   string
	Next    string
	Last    bool
	Sidebar []installationStep
}

type installationFormField struct {
	Name        string
	Type        string
	Placeholder string
	HelperText  string
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
	router.GET("/install/finish", showInstallationFinish)
	router.POST("/install/finish", installationFinish)

	// Loop installation steps
	for i, step := range installationSteps {

		// Register get route
		router.GET(step.URL, installationStepGet(i, step))

		// Register post route
		router.POST(step.URL, installationStepPost(i, step))
	}

	// Create installer listener
	listener, err := net.Listen("tcp", ":8080")

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

func installationStepPost(i int, step installationStep) httprouter.Handle {
	// Set template data
	d := &installationTemplateData{
		Error:   "",
		Success: "",
		Step:    step,
		Next:    "",
		Sidebar: installationSteps,
	}

	// Return httprouter handle
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

		// Parse form
		if err := req.ParseForm(); err != nil {
			d.Error = err.Error()
			installationTemplate.ExecuteTemplate(res, "install.html", d)
			return
		}

		// Call form parser
		if err := step.Post(res, req, step); err != nil {
			d.Error = err.Error()
			installationTemplate.ExecuteTemplate(res, "install.html", d)
			return
		}

		// Check if there is next step
		if i+1 >= len(installationSteps) {

			// Redirect to final page
			http.Redirect(res, req, "/install/finish", 302)

			return
		}

		// Redirect to next step
		http.Redirect(res, req, installationSteps[i+1].URL, 302)
	}
}

func installationStepGet(i int, step installationStep) httprouter.Handle {
	// Set template data
	d := &installationTemplateData{
		Error:   "",
		Success: "",
		Step:    step,
		Next:    "",
		Last:    false,
		Sidebar: installationSteps,
	}

	// Check if step is optional
	if step.Optional {

		// Check if there are more steps
		if i+1 < len(installationSteps) {
			d.Next = installationSteps[i+1].URL
		}
	}

	if i+1 >= len(installationSteps) {
		d.Last = true
	}

	// Return httprouter handle
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// Execute step template
		installationTemplate.ExecuteTemplate(res, "install.html", d)
	}
}

func installationFinish(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Create config file
	if err := createConfigFile("config.toml", installationConfigFile); err != nil {
		installationTemplate.ExecuteTemplate(res, "install_finish.html", installationTemplateData{
			Error: err.Error(),
		})
		return
	}

	// Render finish template
	installationTemplate.ExecuteTemplate(res, "install_finish.html", nil)
}

func showInstallationFinish(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Execute encode layout
	installationTemplate.ExecuteTemplate(res, "install_encode.html", nil)
}

// registerInstallationTemplate register the html template files
func registerInstallationTemplate() error {
	// Parse installation html files
	_, err := installationTemplate.ParseGlob(filepath.Join("views", "install", "*.html"))

	return err
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
		"&multiStatements=true",
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
