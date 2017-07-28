package util

import (
	"github.com/BurntSushi/toml"
	"io"
	"sync"
	"time"
)

// CookieConfig struct used for the cookies configuration options
type CookieConfig struct {
	Name     string
	MaxAge   int
	HashKey  string
	BlockKey string
}

// ShopConfig struct used for the shop configuration options
type ShopConfig struct {
	Enabled bool
}

// PluginConfig struct used for the plugin listener
type PluginConfig struct {
	Enabled  bool
	Username string
	Password string
	Origin   string
}

// RateLimiterConfig struct used for the rate limiting configuration options
type RateLimiterConfig struct {
	Number int64
	Time   StringDuration
}

// CacheConfig struct used for the cache configuration options
type CacheConfig struct {
	Default StringDuration
	Purge   StringDuration
}

// SSLConfig struct used for the ssl configuration options
type SSLConfig struct {
	Enabled bool
	Auto    bool
	Proxy   bool
	Cert    string
	Key     string
}

// MailConfig struct used for the mail configuration options
type MailConfig struct {
	Enabled  bool
	Server   string
	Port     int
	Username string
	Password string
}

// PaygolConfig struct used for the paygol configuration options
type PaygolConfig struct {
	Enabled  bool
	Service  int
	Currency string
	Language string
}

// PayPalConfig struct used for the paypal configuration options
type PayPalConfig struct {
	Enabled   bool
	PublicKey string
	SecretKey string
	Currency  string
	SandBox   bool
}

// FortumoConfig struct used for the fortumo configuration options
type FortumoConfig struct {
	Enabled bool
	Service string
	Secret  string
}

// ContentSecurityPolicyType struct used for CSP fields
type ContentSecurityPolicyType struct {
	Default []string
	SRC     []string
}

// ContentSecurityPolicyConfig struct used for CSP headers
type ContentSecurityPolicyConfig struct {
	Default []string
	Enabled bool
	Frame   ContentSecurityPolicyType
	Script  ContentSecurityPolicyType
	Font    ContentSecurityPolicyType
	Image   ContentSecurityPolicyType
	Connect ContentSecurityPolicyType
	Style   ContentSecurityPolicyType
}

// SecurityConfig struct used for the security of the application
type SecurityConfig struct {
	NonceEnabled      bool
	XSS               string
	STS               string
	Frame             string
	ContentType       string
	ReferrerPolicy    string
	CrossDomainPolicy string
	CSP               ContentSecurityPolicyConfig
}

// Configuration struct used for the main Castro config file TOML file
type Configuration struct {
	CheckUpdates bool
	Template     string
	Mode         string
	Port         int
	URL          string
	Datapack     string
	Security     SecurityConfig
	Plugin       PluginConfig
	Mail         MailConfig
	Captcha      CaptchaConfig
	SSL          SSLConfig
	PayPal       PayPalConfig
	PayGol       PaygolConfig
	Fortumo      FortumoConfig
	Shop         ShopConfig
	Cookies      CookieConfig
	Cache        CacheConfig
	RateLimit    RateLimiterConfig
	Custom       map[string]interface{}
}

// ConfigurationFile struct used to store a configuration pointer
type ConfigurationFile struct {
	rw            sync.RWMutex
	Configuration *Configuration
}

type StringDuration struct {
	Duration time.Duration
	String   string
}

var (
	// Config holds the main configuration file
	Config *ConfigurationFile

	// VERSION current version of the build
	VERSION string

	// BUILD_DATE date of the build
	BUILD_DATE string
)

func init() {
	Config = &ConfigurationFile{}
	Config.Configuration = &Configuration{}
}

// NewStringDuration returns a new string duration struct
func NewStringDuration(s string) StringDuration {
	return StringDuration{
		String: s,
	}
}

// MarshalText use toml interface to convert string durations to strings
func (s StringDuration) MarshalText() ([]byte, error) {
	return []byte(s.String), nil
}

// UnmarshalText use toml interface to convert strings to durations
func (s *StringDuration) UnmarshalText(text []byte) error {
	var err error
	s.Duration, err = time.ParseDuration(string(text))
	return err
}

// LoadConfig loads the configuration file to the given interface pointer
func LoadConfig(path string) error {
	// Lock mutex
	Config.rw.Lock()
	defer Config.rw.Unlock()

	// Decode the given file to the given interface
	if _, err := toml.DecodeFile(path, Config.Configuration); err != nil {
		return err
	}

	return nil
}

// IsDev checks if castro is running on development mode
func (c Configuration) IsDev() bool {
	return c.Mode == "dev"
}

// IsLog checks if castro is running on log mode
func (c Configuration) IsLog() bool {
	return c.Mode == "log"
}

// CSP returns a valid Content-Security-Policy header value
func (c Configuration) CSP() string {
	// Set default-src field
	buff := getCSPField("default-src", c.Security.CSP.Default, nil)

	// Set frame-src field
	buff += getCSPField("frame-src", c.Security.CSP.Frame.Default, c.Security.CSP.Frame.SRC)

	// Set script-src field
	buff += getCSPField("script-src", c.Security.CSP.Script.Default, c.Security.CSP.Script.SRC)

	// Set font-src field
	buff += getCSPField("font-src", c.Security.CSP.Font.Default, c.Security.CSP.Font.SRC)

	// Set connect-src field
	buff += getCSPField("connect-src", c.Security.CSP.Connect.Default, c.Security.CSP.Connect.SRC)

	// Set style-src field
	buff += getCSPField("style-src", c.Security.CSP.Style.Default, c.Security.CSP.Style.SRC)

	// Set img-src field
	buff += getCSPField("img-src", c.Security.CSP.Image.Default, c.Security.CSP.Image.SRC)

	return buff
}

func getCSPField(name string, def []string, src []string) string {
	// Data holder
	buff := name

	// Loop default values
	for _, d := range def {
		buff += " '" + d + "' "
	}

	// Loop src values
	for _, s := range src {
		buff += " " + s + " "
	}

	return buff + ";"
}

// IsSSL returns if the server is behind SSL
func (c Configuration) IsSSL() bool {
	if c.SSL.Enabled {
		return true
	}
	if c.SSL.Proxy {
		return true
	}

	return false
}

// EncodeConfig encodes the given io writer
func EncodeConfig(configFile io.Writer, c *Configuration) error {
	// Lock mutex
	Config.rw.Lock()
	defer Config.rw.Unlock()

	// Encode the given writer with the given interface
	return toml.NewEncoder(configFile).Encode(c)
}

// SetCustomValue sets a config custom value
func (c *ConfigurationFile) SetCustomValue(key string, v interface{}) {
	// Lock mutex
	c.rw.Lock()
	defer c.rw.Unlock()

	// Set custom value
	c.Configuration.Custom[key] = v
}
